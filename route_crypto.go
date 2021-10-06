package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"
)

func SelectPairs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectPairs..."))

	_ = SelectSession(r)

	type PairInfo struct {
		CoinId int
		Symbol string
		Slug   string
	}

	pairs := make(map[string]PairInfo)

	pairs_sql := `
		SELECT DISTINCT
			name,
			coinid,
			symbol,
			slug
		FROM coins;`
	pairs_rows, err := DbWebApp.Query(pairs_sql)
	defer pairs_rows.Close()
	if err != nil {
		log.Error(err)
	}
	for pairs_rows.Next() {
		var name string
		pair_info := PairInfo{}
		if err = pairs_rows.Scan(
			&name,
			&pair_info.CoinId,
			&pair_info.Symbol,
			&pair_info.Slug,
		); err != nil {
			log.Error(err)
		}
		pairs[name] = pair_info
	}

	json.NewEncoder(w).Encode(pairs)
}

func SelectStellarPrice(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectStellarPrice..."))

	_ = SelectSession(r)

	stellar_price := struct {
		Price string
	}{}

	price_sql := `
		SELECT
			p.price::text
		FROM prices p
		LEFT JOIN coins c ON(p.coinid = c.coinid)
		WHERE c.symbol = 'XLM'
		ORDER BY p.createdat DESC
		LIMIT 1;`
	err := DbWebApp.QueryRow(
		price_sql).Scan(
		&stellar_price.Price,
	)
	if err != nil {
		log.Error(err)
	}

	json.NewEncoder(w).Encode(stellar_price)
}

func SelectTransactionCredentials(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting SelectTransactionCredentials..."))

	user_session := SelectSession(r)
	user := UserByEmail(user_session.Email)

	wallet_sql := `
		SELECT
			blockchain,
			currency,
			address
		FROM internalwallets
		WHERE blockchain = 'Stellar'`

	var blockchain string
	var currency string
	var deposit_address string
	err := DbWebApp.QueryRow(
		wallet_sql).Scan(
		&blockchain,
		&currency,
		&deposit_address,
	)
	if err != nil {
		log.Error(err)
	}

	var memo string
	statement := `
		INSERT INTO memos (userid, blockchain, currency, depositaddress, status, memo, createdat)
		VALUES ($1, $2, $3, $4, $5, SUBSTR(MD5(RANDOM()::TEXT), 0, 20), current_timestamp)
		RETURNING memo;`
	err = DbWebApp.QueryRow(
		statement,
		user.Id,
		blockchain,
		currency,
		deposit_address,
		"pending").Scan(&memo)
	if err != nil {
		log.Error(err)
	}

	credentials := struct {
		DepositAddress string
		Memo           string
	}{
		deposit_address,
		memo,
	}

	json.NewEncoder(w).Encode(credentials)
}

func ValidateStellarTransaction(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ValidateTransaction..."))

	transaction_detail := struct {
		Id     string  `json:"Id"`
		Memo   string  `json:"Memo"`
		Amount float64 `json:"Amount"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&transaction_detail)
	if err != nil {
		log.Error(err)
	}

	fmt.Println("tx_id", transaction_detail.Id)
	fmt.Println("tx_memo", transaction_detail.Memo)
	fmt.Println("tx_amount", transaction_detail.Amount)

	time.Sleep(2 * time.Second)

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	type StellarTx struct {
		Memo           string    `json:"memo"`
		ID             string    `json:"id"`
		Successful     bool      `json:"successful"`
		CreatedAt      time.Time `json:"created_at"`
		SourceAccount  string    `json:"source_account"`
		FeeAccount     string    `json:"fee_account"`
		FeeCharged     string    `json:"fee_charged"`
		MaxFee         string    `json:"max_fee"`
		OperationCount int       `json:"operation_count"`
		EnvelopeXdr    string    `json:"envelope_xdr"`
		ResultXdr      string    `json:"result_xdr"`
		ResultMetaXdr  string    `json:"result_meta_xdr"`
		FeeMetaXdr     string    `json:"fee_meta_xdr"`
		MemoType       string    `json:"memo_type"`
		Signatures     []string  `json:"signatures"`
	}

	stellar_tx_url := fmt.Sprintf("https://horizon.stellar.org/transactions/%s", transaction_detail.Id)

	res, err := http.Get(stellar_tx_url)
	if err != nil {
		log.WithFields(log.Fields{
			"session": session.Id,
			"user":    user.Id,
		}).Error("Failed fetching TX from Horizon API")
		json.NewEncoder(w).Encode("KO")
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"session": session.Id,
			"user":    user.Id,
		}).Error("Failed converting TX into struct")
		json.NewEncoder(w).Encode("KO")
		return
	}

	var stellar_tx StellarTx
	json.Unmarshal(body, &stellar_tx)

	if !stellar_tx.Successful {
		log.WithFields(log.Fields{
			"session": session.Id,
			"user":    user.Id,
		}).Error("Unsucsessfull TX")
		json.NewEncoder(w).Encode("KO")
		return
	}

	envelopeXDR := stellar_tx.EnvelopeXdr

	envelope := xdr.TransactionEnvelope{}
	err = xdr.SafeUnmarshalBase64(envelopeXDR, &envelope)
	if err != nil {
		log.WithFields(log.Fields{
			"session": session.Id,
			"user":    user.Id,
		}).Error("Corrupted EnvelopeXDR")
		json.NewEncoder(w).Encode("KO")
		return
	}

	if envelope.V1.Tx.Operations[0].Body.PaymentOp == nil {
		log.WithFields(log.Fields{
			"session": session.Id,
			"user":    user.Id,
		}).Error("Not a tx payment type")
		json.NewEncoder(w).Encode("KO")
		return
	}

	// VERIFICIARE CHE ANCHE L IMPORTO DELLA TRANSAZIONE SIA CORRETTO, QUINDI POST REQUEST
	amount := envelope.V1.Tx.Operations[0].Body.PaymentOp.Amount
	asset := envelope.V1.Tx.Operations[0].Body.PaymentOp.Asset.String()
	fmt.Println("Memo", stellar_tx.Memo)
	fmt.Println("Asset", asset)
	fmt.Println("Amount", amount)

	if asset != "native" {
		log.WithFields(log.Fields{
			"session": session.Id,
			"user":    user.Id,
		}).Error("No XLM transaction")
		json.NewEncoder(w).Encode("KO")
		return
	}

	json.NewEncoder(w).Encode("OK")
}

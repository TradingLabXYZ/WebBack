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

	credentials.Memo = "141f965131cf9f89f97" // TEMP TO TEST

	json.NewEncoder(w).Encode(credentials)
}

func ValidateStellarTransaction(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting ValidateTransaction..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	time.Sleep(2 * time.Second)

	input_tx := struct {
		Id         string  `json:"Id"`
		Memo       string  `json:"Memo"`
		Months     int     `json:"Months"`
		Amount     float64 `json:"Amount"`
		Blockchain string  `json:"Blockchain"`
		Currency   string  `json:"Currency"`
		AmountXdr  xdr.Int64
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input_tx)
	if err != nil {
		fmt.Println(err)
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"custom_msg": "Failed decoding data from user",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	input_tx.AmountXdr = xdr.Int64(input_tx.Amount)

	type StellarTx struct {
		Memo          string    `json:"memo"`
		ID            string    `json:"id"`
		Successful    bool      `json:"successful"`
		CreatedAt     time.Time `json:"created_at"`
		SourceAccount string    `json:"source_account"`
		EnvelopeXdr   string    `json:"envelope_xdr"`
	}

	stellar_tx_url := fmt.Sprintf(
		"https://horizon.stellar.org/transactions/%s",
		input_tx.Id,
	)

	res, err := http.Get(stellar_tx_url)
	if err != nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Failed fetching TX from Horizon API",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Failed converting TX into struct",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	var stellar_tx StellarTx
	json.Unmarshal(body, &stellar_tx)

	if !stellar_tx.Successful {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Unsucsessfull TX",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	envelopeXDR := stellar_tx.EnvelopeXdr

	envelope := xdr.TransactionEnvelope{}
	err = xdr.SafeUnmarshalBase64(envelopeXDR, &envelope)
	if err != nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Corrupted EnvelopeXDR",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	paymentOp := envelope.V1.Tx.Operations[0].Body.PaymentOp
	createAccountOp := envelope.V1.Tx.Operations[0].Body.CreateAccountOp
	if paymentOp == nil && createAccountOp == nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Not a tx payment type",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	var amount xdr.Int64
	var asset string
	if paymentOp != nil {
		amount = envelope.V1.Tx.Operations[0].Body.PaymentOp.Amount
		asset = envelope.V1.Tx.Operations[0].Body.PaymentOp.Asset.String()
	} else if createAccountOp != nil {
		amount = envelope.V1.Tx.Operations[0].Body.CreateAccountOp.StartingBalance
		asset = "native"
	}

	if asset != "native" {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "No XLM transaction",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	if input_tx.Memo != stellar_tx.Memo {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Memos do not match",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	delta_amounts := input_tx.AmountXdr - amount
	if delta_amounts > 10000000 {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "TX amount not valid",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	// All the checks have been successfull

	// Upgrading user to Premium
	upgrade_sql := `
		UPDATE users
		SET plan = 'premium'
		WHERE id = $1;`
	_, err = DbWebApp.Exec(upgrade_sql, user.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Update user to premium failed",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	// Adding a payment into table
	payment_sql := `
		INSERT INTO payments (userid, blockchain, currency, transactionid, amount, months, createdat, endat)  
		VALUES ($1, $2, $3, $4, $5, $6, current_timestamp, current_timestamp + interval '1 month' * $7);`
	_, err = DbWebApp.Exec(
		payment_sql,
		user.Id,
		input_tx.Blockchain,
		input_tx.Currency,
		input_tx.Id,
		amount,
		input_tx.Months,
		input_tx.Months)
	if err != nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       input_tx.Id,
			"custom_msg": "Insert new payment failed",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	log.WithFields(log.Fields{
		"session": session.Id,
		"user":    user.Id,
		"txid":    input_tx.Id,
	}).Info("Successfully validated XLM transaction")
	json.NewEncoder(w).Encode("OK")
}

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

type StellarTx struct {
	Memo          string    `json:"memo"`
	ID            string    `json:"id"`
	Successful    bool      `json:"successful"`
	CreatedAt     time.Time `json:"created_at"`
	SourceAccount string    `json:"source_account"`
	EnvelopeXdr   string    `json:"envelope_xdr"`
}

func (tx TxBuyPremium) ValidateStellarTransaction() (status string) {
	fmt.Println(Gray(8-1, "Starting ValidateTransaction..."))

	// Extract transaction data from horizon API
	stellar_tx_url := fmt.Sprintf(
		"https://horizon.stellar.org/transactions/%s",
		tx.Id,
	)

	res, err := http.Get(stellar_tx_url)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Failed fetching TX from Horizon API",
		}).Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Failed converting TX into struct",
		}).Error(err)
		return "KO"
	}

	var stellar_tx StellarTx
	json.Unmarshal(body, &stellar_tx)

	// Check status
	if !stellar_tx.Successful {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Unsucsessfull TX",
		}).Error(err)
		return "KO"
	}

	// Check envelope
	envelopeXDR := stellar_tx.EnvelopeXdr
	envelope := xdr.TransactionEnvelope{}
	err = xdr.SafeUnmarshalBase64(envelopeXDR, &envelope)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Corrupted EnvelopeXDR",
		}).Error(err)
		return "KO"
	}

	// Check type
	paymentOp := envelope.V1.Tx.Operations[0].Body.PaymentOp
	createAccountOp := envelope.V1.Tx.Operations[0].Body.CreateAccountOp
	if paymentOp == nil && createAccountOp == nil {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Not a tx payment type",
		}).Error(err)
		return "KO"
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

	// Check asset
	if asset != "native" {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "No XLM transaction",
		}).Error(err)
		return "KO"
	}

	// Check memo
	if tx.Memo != stellar_tx.Memo {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Memos do not match",
		}).Error(err)
		return "KO"
	}

	// Check amount
	tx_amount_xdr := xdr.Int64(tx.Amount)
	delta_amounts := tx_amount_xdr - amount
	if delta_amounts > 10000000 {
		log.WithFields(log.Fields{
			"sessionCode": tx.SessionCode,
			"userCode":    tx.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "TX amount not valid",
		}).Error(err)
		return "KO"
	}

	// All the checks positive
	log.WithFields(log.Fields{
		"sessionCode": tx.SessionCode,
		"userCode":    tx.UserCode,
		"txid":        tx.Id,
	}).Info("Successfully validated XLM transaction")
	return "OK"
}

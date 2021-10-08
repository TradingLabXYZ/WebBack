package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

type TxBuyPremium struct {
	SessionId  int
	Userid     int
	Id         string  `json:"Id"`
	Memo       string  `json:"Memo"`
	Months     int     `json:"Months"`
	Amount     float64 `json:"Amount"`
	Blockchain string  `json:"Blockchain"`
	Asset      string  `json:"Asset"`
}

/**
1 --> Parse data coming from frontend
2 --> Validate the transaction based on the blockchain selected
3 --> Change the status of the user to premium
4 --> Register the payment into the database
FAIL --> RETURN KO
SUCCESS --> RETURN OK
*/
func BuyPremiumMonths(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting BuyPremiumMonths..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	var tx TxBuyPremium
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tx)
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

	tx.Userid = user.Id
	tx.SessionId = session.Id

	var status string
	if tx.Blockchain == "Stellar" {
		status = tx.ValidateStellarTransaction()
	}

	if status == "KO" {
		json.NewEncoder(w).Encode("KO")
		return
	}

	user.UpdateUserStatus("premium")
	if err != nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       tx.Id,
			"custom_msg": "Update user to premium failed",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	err = tx.InsertPayment("basicToPremium")
	if err != nil {
		log.WithFields(log.Fields{
			"session":    session.Id,
			"user":       user.Id,
			"txid":       tx.Id,
			"custom_msg": "Insert new payment failed",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	log.WithFields(log.Fields{
		"session":    session.Id,
		"user":       user.Id,
		"txid":       tx.Id,
		"blockchain": tx.Blockchain,
		"asset":      tx.Asset,
	}).Info("Successfully validated transaction")
	json.NewEncoder(w).Encode("OK")
}

func (user User) UpdateUserStatus(new_status string) error {
	upgrade_sql := `
		UPDATE users
		SET plan = $1
		WHERE id = $2;`
	_, err := DbWebApp.Exec(upgrade_sql, new_status, user.Id)
	return err
}

func (tx TxBuyPremium) InsertPayment(reason string) error {
	payment_sql := `
		INSERT INTO payments (userid, type, blockchain, currency, transactionid, amount, months, createdat, endat)  
		VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp + interval '1 month' * $8);`
	_, err := DbWebApp.Exec(
		payment_sql,
		tx.Userid,
		reason,
		tx.Blockchain,
		tx.Currency,
		tx.Id,
		tx.Amount,
		tx.Months,
		tx.Months)
	return err
}

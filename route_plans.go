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
	SessionCode string
	UserCode    string
	Id          string  `json:"Id"`
	Memo        string  `json:"Memo"`
	Months      int     `json:"Months"`
	Amount      float64 `json:"Amount"`
	Blockchain  string  `json:"Blockchain"`
	Asset       string  `json:"Asset"`
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

	session, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if session.Code == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user, err := SelectUser("code", session.UserCode)
	if err != nil {
		log.Warn("User not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var tx TxBuyPremium
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&tx)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"userCode":    session.UserCode,
			"custom_msg":  "Failed decoding data from user",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	tx.UserCode = session.UserCode
	tx.SessionCode = session.Code

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
			"sessionCode": session.Code,
			"userCode":    session.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Update user to premium failed",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	err = tx.InsertPayment("basicToPremium")
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"userCode":    session.UserCode,
			"txid":        tx.Id,
			"custom_msg":  "Insert new payment failed",
		}).Error(err)
		json.NewEncoder(w).Encode("KO")
		return
	}

	log.WithFields(log.Fields{
		"sessionCode": session.Code,
		"userCode":    session.UserCode,
		"txid":        tx.Id,
		"blockchain":  tx.Blockchain,
		"asset":       tx.Asset,
	}).Info("Successfully validated transaction")
	json.NewEncoder(w).Encode("OK")
}

func (user User) UpdateUserStatus(new_status string) error {
	upgrade_sql := `
		UPDATE users
		SET plan = $1
		WHERE code = $2;`
	_, err := Db.Exec(upgrade_sql, new_status, user.Code)
	return err
}

func (tx TxBuyPremium) InsertPayment(reason string) error {
	payment_sql := `
		INSERT INTO payments (usercode, type, blockchain, currency, transactionid, amount, months, createdat, endat)  
		VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp + interval '1 month' * $8);`
	_, err := Db.Exec(
		payment_sql,
		tx.UserCode,
		reason,
		tx.Blockchain,
		tx.Asset,
		tx.Id,
		tx.Amount,
		tx.Months,
		tx.Months)
	return err
}

func GetUserPremiumData(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetUserPremiumData..."))

	session, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if session.Code == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	type Payment struct {
		CreatedAt     string
		Amount        float64
		Months        int
		Blockchain    string
		Currency      string
		TransactionId string
	}

	user_premium_data := struct {
		Payments      []Payment
		RemainingDays int
	}{}

	payments_sql := `
		SELECT
			TO_CHAR(createdat::DATE, 'YYYY-MM-DD') AS createdat,
			blockchain,
			currency,
			amount,
			months,
			transactionid
		FROM payments
		WHERE usercode = $1
		ORDER BY 1;`

	rows, err := Db.Query(payments_sql, session.UserCode)
	defer rows.Close()
	if err != nil {
		log.Error(err)
	}
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(
			&payment.CreatedAt,
			&payment.Blockchain,
			&payment.Currency,
			&payment.Amount,
			&payment.Months,
			&payment.TransactionId,
		); err != nil {
			log.Error(err)
		}
		user_premium_data.Payments = append(user_premium_data.Payments, payment)
	}

	remaining_days_sql := `
		SELECT
			EXTRACT(DAY FROM MAX(endat)::date - now()) AS remaining_days
		FROM payments
		WHERE usercode = $1;`

	err = Db.QueryRow(
		remaining_days_sql,
		session.UserCode).Scan(&user_premium_data.RemainingDays)
	if err != nil {
		log.Error(err)
	}

	json.NewEncoder(w).Encode(user_premium_data)
}

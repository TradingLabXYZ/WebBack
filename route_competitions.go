package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func InsertPrediction(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed inserting prediction, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "web" {
		log.Error("Failed inserting prediction, origin not web")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitionname := mux.Vars(r)["competition"]
	prediction := mux.Vars(r)["prediction"]
	s_prediction := map[string]string{"prediction": prediction}

	s_payload, err := json.Marshal(s_prediction)
	if err != nil {
		return
	}

	statement := `
		INSERT INTO submissions (
			updatedat, competitionname, userwallet, payload)
		VALUES (current_timestamp, $1, $2, $3)
		ON CONFLICT (userwallet) DO UPDATE
		SET updatedat = now(), payload = EXCLUDED.payload;`
	_, err = Db.Exec(
		statement,
		competitionname,
		session.UserWallet,
		s_payload)
	if err != nil {
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SelectPrediction(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed selecting prediction, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "web" {
		log.Error("Failed selecting prediction, origin not web")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitionname := mux.Vars(r)["competition"]

	var prediction string
	_ = Db.QueryRow(`
			SELECT payload#>>'{prediction}'
			FROM submissions
			WHERE competitionname = $1
			AND userwallet = $2;`,
		competitionname,
		session.UserWallet).Scan(&prediction)
	if err != nil {
		log.Warn(err)
		return
	}

	json.NewEncoder(w).Encode(prediction)
}

func DeletePrediction(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed deleting prediction, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session.Origin != "web" {
		log.Error("Failed deleting prediction, origin not web")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitionname := mux.Vars(r)["competition"]

	statement := `
		DELETE FROM submissions
		WHERE userwallet = $1
		AND competitionname = $2;`
	_, err = Db.Exec(
		statement,
		session.UserWallet,
		competitionname)
	if err != nil {
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetCountPartecipants(w http.ResponseWriter, r *http.Request) {
	competitionname := mux.Vars(r)["competition"]

	var count int
	err := Db.QueryRow(`
			SELECT COUNT(*)
			FROM submissions
			WHERE competitionname = $1;`,
		competitionname).Scan(&count)
	if err != nil {
		log.Warn(err)
		return
	}

	json.NewEncoder(w).Encode(count)
}

func GetPartecipants(w http.ResponseWriter, r *http.Request) {
	competitionname := mux.Vars(r)["competition"]

	type Prediction struct {
		CreatedAt      string
		Username       string
		Wallet         string
		ProfilePicture string
		Prediction     string
		BtcPrice       string
		DeltaPerc      float32
		AbsDeltaPrice  string
	}

	predictions := []Prediction{}

	prediction_query := `
		WITH
			last_btc_price AS (
				SELECT
					price
				FROM lastprices
				WHERE coinid = 1)
		SELECT
			u.updatedat,
			u.username,
			CASE WHEN u.profilepicture IS NULL THEN '' ELSE u.profilepicture END AS profilepicture,
			LEFT(s.userwallet, 3) || '...' || RIGHT(s.userwallet, 3) AS userwallet,
			TO_CHAR(ROUND((s.payload#>>'{prediction}')::NUMERIC, 2), '999,999,999') AS prediction,
			TO_CHAR(l.price, '999,999,999') || '$' AS btc_price,
			ROUND(((s.payload#>>'{prediction}')::NUMERIC / l.price - 1) * 100, 2) AS deltaprice,
			ABS((s.payload#>>'{prediction}')::NUMERIC / l.price - 1) AS absdeltaprice
		FROM submissions s
		LEFT JOIN users u ON(s.userwallet = u.wallet)
		LEFT JOIN last_btc_price l ON(1 = 1)
		WHERE competitionname = $1
		ORDER BY 8;`

	predictions_rows, err := Db.Query(
		prediction_query,
		competitionname)
	defer predictions_rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"competitionName": competitionname,
			"custom_msg":      "Failed running prediction_sql",
		}).Error(err)
	}
	for predictions_rows.Next() {
		prediction := Prediction{}
		if err = predictions_rows.Scan(
			&prediction.CreatedAt,
			&prediction.Username,
			&prediction.ProfilePicture,
			&prediction.Wallet,
			&prediction.Prediction,
			&prediction.BtcPrice,
			&prediction.DeltaPerc,
			&prediction.AbsDeltaPrice,
		); err != nil {
			log.WithFields(log.Fields{
				"competitionName": competitionname,
				"custom_msg":      "Failed parsing prediction_sql",
			}).Error(err)
		}
		predictions = append(predictions, prediction)
	}

	json.NewEncoder(w).Encode(predictions)
}

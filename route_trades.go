package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

func CheckUserPrivacy(next http.Handler) http.Handler {
	fmt.Println(Gray(8-1, "Starting CheckUserPrivacy..."))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]
		userA := UserByUsername(username)
		privacy := userA.Privacy
		if privacy == "all" {
			next.ServeHTTP(w, r)
			return
		} else {
			session := SelectSession(r)
			if session.Id == 0 {
				w.Write([]byte(`{"Status": "denied", "Reason": "login"}`))
				return
			}
			userB := UserByEmail(session.Email)
			if userA.Id == userB.Id {
				next.ServeHTTP(w, r)
				return
			}
			switch privacy {
			case "private":
				w.Write([]byte(`{"Status": "denied", "Reason": "private"}`))
				return
			case "followers":
				var isfollower bool
				_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE followto = $1
					AND followfrom = $2;`, userA.Id, userB.Id).Scan(
					&isfollower,
				)
				if isfollower {
					next.ServeHTTP(w, r)
				} else {
					w.Write([]byte(`{"Status": "denied", "Reason": "follow"}`))
					return
				}
			case "subscribers":
				var issubscriber bool
				_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE subscribeto = $1
					AND subscribefrom = $2;`, userA.Id, userB.Id).Scan(
					&issubscriber,
				)
				if issubscriber {
					next.ServeHTTP(w, r)
					return
				} else {
					w.Write([]byte(`{"Status": "denied", "Reason": "subscribe"}`))
					return
				}
			case "individuals":
				var isindividual bool
				_ = DbWebApp.QueryRow(`
					SELECT TRUE
					FROM individuals
					WHERE individualto = $1
					AND individualfrom = $2;`, userA.Id, userB.Id).Scan(
					&isindividual,
				)
				if isindividual {
					next.ServeHTTP(w, r)
					return
				} else {
					w.Write([]byte(`{"Status": "denied", "Reason": "individual"}`))
					return
				}
			}
		}
	})
}

func InsertTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting InsertTrade..."))

	session := SelectSession(r)
	user := UserByEmail(session.Email)

	trade := struct {
		Exchange     string `json:"Exchange"`
		FirstPairId  int    `json:"FirstPair"`
		SecondPairId int    `json:"SecondPair"`
		Subtrades    []struct {
			Timestamp string      `json:"Timestamp"`
			Type      string      `json:"Type"`
			Reason    string      `json:"Reason"`
			Quantity  json.Number `json:"Quantity"`
			AvgPrice  json.Number `json:"AvgPrice"`
			Total     json.Number `json:"Total"`
		} `json:"subtrades"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&trade)
	if err != nil {
		log.Error(err)
	}

	var trade_id string
	trade_sql := `
		INSERT INTO trades (id, userid, exchange, firstpair, secondpair, createdat, updatedat, isopen)
		VALUES (SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4, current_timestamp, current_timestamp, true)
		RETURNING id;`
	err = DbWebApp.QueryRow(
		trade_sql,
		user.Id,
		trade.Exchange,
		trade.FirstPairId,
		trade.SecondPairId,
	).Scan(&trade_id)
	if err != nil {
		log.Error(err)
	}

	for _, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (tradeid, tradetimestamp, type, reason, quantity, avgprice, total, createdat, updatedat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp)`
		DbWebApp.Exec(
			subtrade_sql,
			trade_id,
			subtrade.Timestamp,
			subtrade.Type,
			subtrade.Reason,
			subtrade.Quantity,
			subtrade.AvgPrice,
			subtrade.Total,
		)
	}
	if err != nil {
		log.Error(err)
	}

	json.NewEncoder(w).Encode("OK")
}

func UpdateTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateTrade..."))

	trade := struct {
		Id           string `json:"Id"`
		Exchange     string `json:"Exchange"`
		FirstPairId  int    `json:"FirstPairId"`
		SecondPairId int    `json:"SecondPairId"`
		Subtrades    []struct {
			Timestamp string      `json:"Timestamp"`
			Type      string      `json:"Type"`
			Reason    string      `json:"Reason"`
			Quantity  json.Number `json:"Quantity"`
			AvgPrice  json.Number `json:"AvgPrice"`
			Total     json.Number `json:"Total"`
		} `json:"Subtrades"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&trade)
	if err != nil {
		log.Error(err)
	}

	DbWebApp.Exec(`
		DELETE FROM subtrades
		WHERE tradeid = $1;
	`, trade.Id)

	for _, subtrade := range trade.Subtrades {
		subtrade_sql := `
		INSERT INTO subtrades (tradeid, tradetimestamp, type, reason, quantity, avgprice, total, createdat, updatedat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp)`
		DbWebApp.Exec(
			subtrade_sql,
			trade.Id,
			subtrade.Timestamp,
			subtrade.Type,
			subtrade.Reason,
			subtrade.Quantity,
			subtrade.AvgPrice,
			subtrade.Total,
		)
	}

	json.NewEncoder(w).Encode("OK")
}

func CloseTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting CloseTrade..."))

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.Exec(`
		UPDATE trades
		SET
			isopen = False,
			updatedat = current_timestamp
		WHERE id = $1;
		`, tradeid)

	json.NewEncoder(w).Encode("OK")
}

func OpenTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting OpenTrade..."))

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.Exec(`
		UPDATE trades
		SET
			isopen = True,
			updatedat = current_timestamp
		WHERE id = $1;
		`, tradeid)

	json.NewEncoder(w).Encode("OK")
}

func DeleteTrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting DeleteTrade..."))

	tradeid := mux.Vars(r)["tradeid"]

	DbWebApp.Exec(`
		DELETE FROM subtrades
		WHERE tradeid IN (
			SELECT id
			FROM trades
			WHERE id = $1);
		`, tradeid)

	DbWebApp.Exec(`
		UPDATE users
		SET updatedat = current_timestamp
		WHERE id = (
			SELECT
				userid
			FROM trades
			WHERE id = $1);
		`, tradeid)

	DbWebApp.Exec(`
		DELETE FROM trades
		WHERE id = $1;
		`, tradeid)

	json.NewEncoder(w).Encode("OK")
}

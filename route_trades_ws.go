package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsTrade struct {
	UserToSee User
	SessionId string
	Channel   chan TradesSnapshot
	Ws        *websocket.Conn
}

func StartTradesWs(w http.ResponseWriter, r *http.Request) {
	wallet := mux.Vars(r)["wallet"]
	session_id := mux.Vars(r)["sessionid"]

	session := Session{}
	user := User{}
	if session_id != "undefined" {
		session.Code = session_id
		session.Select()
		user, _ = SelectUser("wallet", session.UserWallet)
		/* if err != nil {
			log.WithFields(log.Fields{
				"urlPath": r.URL.Path,
			}).Warn("Failed starting ws, user has cookie but not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		} */
	}

	url_split := strings.Split(r.URL.Path, "/")

	if len(url_split) < 4 {
		log.WithFields(log.Fields{
			"urlPath": r.URL.Path,
		}).Warn("Failed starting ws, wrong url")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wallet = url_split[2]

	userToSee, err := SelectUser("wallet", wallet)
	if err != nil {
		log.WithFields(log.Fields{
			"urlPath":   r.URL.Path,
			"userToSee": userToSee,
		}).Warn("Failed starting ws, user not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ws, err := InstanciateTradeWs(w, r)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed instanciating TradeWs",
		}).Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c := make(chan TradesSnapshot)
	ws_trade := WsTrade{userToSee, session_id, c, ws}
	ws_trade_output := ws_trade.UserToSee.GetSnapshot()

	ws_trade_output.CheckPrivacy(user, userToSee)

	if ws_trade_output.PrivacyStatus.Status == "KO" {
		ws_trade_output.Trades = nil
		ws_trade.SendInitialSnapshot(ws_trade_output)
		return
	}
	ws_trade.SendInitialSnapshot(ws_trade_output)
	trades_wss[wallet] = append(trades_wss[wallet], ws_trade)
	go ws_trade.WaitToTerminate()
	go ws_trade.WaitToSendMessage()
}

func (snapshot *TradesSnapshot) CheckPrivacy(user User, userToSee User) {
	if userToSee.Privacy == "all" {
		snapshot.PrivacyStatus.Status = "OK"
		snapshot.PrivacyStatus.Reason = "userToSee ALL"
		return
	}

	if user.Wallet == "" {
		snapshot.PrivacyStatus.Status = "KO"
		snapshot.PrivacyStatus.Reason = "user is not logged in"
		return
	}

	if user.Wallet == userToSee.Wallet {
		snapshot.PrivacyStatus.Status = "OK"
		snapshot.PrivacyStatus.Reason = "user access its own profile"
		return
	}

	switch userToSee.Privacy {
	case "private":
		snapshot.PrivacyStatus.Status = "KO"
		snapshot.PrivacyStatus.Reason = "private"
		return
	case "followers":
		var isfollower bool
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE followfrom = $1
					AND followto = $2;`, user.Wallet, userToSee.Wallet).Scan(
			&isfollower,
		)
		if isfollower {
			snapshot.PrivacyStatus.Status = "OK"
			snapshot.PrivacyStatus.Reason = "user is follower"
			return
		} else {
			snapshot.PrivacyStatus.Status = "KO"
			snapshot.PrivacyStatus.Reason = "user is not follower"
			return
		}
	case "subscribers":
		var issubscriber bool
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE subscribefrom = $1
					AND subscribeto = $2;`, user.Wallet, userToSee.Wallet).Scan(
			&issubscriber,
		)
		if issubscriber {
			snapshot.PrivacyStatus.Status = "OK"
			snapshot.PrivacyStatus.Reason = "user is subscriber"
			return
		} else {
			snapshot.PrivacyStatus.Status = "KO"
			snapshot.PrivacyStatus.Reason = "user is not subscriber"
			return
		}
	default:
		log.WithFields(log.Fields{
			"userToSee": userToSee.Wallet,
			"user":      user.Wallet,
		}).Warn("Not possible to determine user's privacy")
		snapshot.PrivacyStatus.Status = "KO"
		snapshot.PrivacyStatus.Reason = "unknown reason"
		return
	}
}

func InstanciateTradeWs(w http.ResponseWriter, r *http.Request) (ws *websocket.Conn, err error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		for _, origin := range Origins {
			if origin == r.Header["Origin"][0] {
				return true
			}
		}
		return false
	}
	if upgrader.CheckOrigin(r) {
		ws, err = upgrader.Upgrade(w, r, nil)
	} else {
		err = errors.New("CheckOrigin not accepted")
	}
	return
}

func (ws_trade *WsTrade) SendInitialSnapshot(snapshot TradesSnapshot) {
	err := ws_trade.Ws.WriteJSON(snapshot)
	if err != nil {
		ws_trade.Ws.Close()
		log.WithFields(log.Fields{
			"wallet":     ws_trade.UserToSee.Wallet,
			"custom_msg": "Failed running sending initial snapshot",
		}).Error(err)
		return
	}
}

func (ws_trade *WsTrade) WaitToTerminate() {
	for {
		_, _, err := ws_trade.Ws.ReadMessage()
		if err != nil {
			for i, v := range trades_wss[ws_trade.UserToSee.Wallet] {
				if v.SessionId == ws_trade.SessionId {
					trades_wss[ws_trade.UserToSee.Wallet] = append(
						trades_wss[ws_trade.UserToSee.Wallet][:i],
						trades_wss[ws_trade.UserToSee.Wallet][i+1:]...)
				}
			}
			ws_trade.Ws.Close()
			return
		} else {
			ws_trade.Ws.Close()
			log.WithFields(log.Fields{
				"wallet":     ws_trade.UserToSee.Wallet,
				"custom_msg": "Failed terminating trade ws",
			}).Error(err)
			return
		}
	}
}

func (ws_trade *WsTrade) WaitToSendMessage() {
	for {
		s1 := <-ws_trade.Channel
		if s1.UserDetails.Username == ws_trade.UserToSee.Wallet {
			err := ws_trade.Ws.WriteJSON(s1)
			if err != nil {
				ws_trade.Ws.Close()
				log.WithFields(log.Fields{
					"wallet":     ws_trade.UserToSee.Wallet,
					"custom_msg": "Failed running sending snapshot",
				}).Error(err)
				return
			}
		}
	}
}

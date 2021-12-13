package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsTrade struct {
	Observer  User
	Observed  User
	SessionId string
	Channel   chan TradesSnapshot
	Ws        *websocket.Conn
}

func StartTradesWs(w http.ResponseWriter, r *http.Request) {
	url_split := strings.Split(r.URL.Path, "/")

	// UNDERSTAND WHICH USER PROFILE NEEDS TO BE DISPLAYED
	if len(url_split) < 4 {
		log.WithFields(log.Fields{
			"urlPath": r.URL.Path,
		}).Warn("Failed starting ws, wrong url")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	wallet := url_split[2]

	observed, err := SelectUser("wallet", wallet)
	if err != nil {
		log.WithFields(log.Fields{
			"urlPath":  r.URL.Path,
			"observed": observed,
		}).Warn("Failed starting ws, user not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// UNDERSTAND WHICH USER WANTS TO RECEIVE THE DATA
	session_id := url_split[3]
	session := Session{}
	observer := User{}
	session.Code = session_id
	err = session.Select()
	if err == nil {
		observer, err = SelectUser("wallet", session.UserWallet)
		if err != nil {
			log.WithFields(log.Fields{
				"urlPath": r.URL.Path,
			}).Error("Failed starting ws, user has cookie but not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		// USER HAS NOT BEEN FOUND
	}

	// INSTANCIATE WS
	ws, err := InstanciateTradeWs(w, r)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed instanciating TradeWs",
		}).Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c := make(chan TradesSnapshot)
	ws_trade := WsTrade{observer, observed, session_id, c, ws}
	snapshot := observed.GetSnapshot()

	snapshot.CheckPrivacy(observer, observed)

	if snapshot.PrivacyStatus.Status == "KO" {
		snapshot.Trades = nil
		ws_trade.SendInitialSnapshot(snapshot)
		return
	}

	ws_trade.SendInitialSnapshot(snapshot)
	trades_wss[wallet] = append(trades_wss[wallet], ws_trade)
	go ws_trade.WaitToTerminate()
	go ws_trade.WaitToSendMessage()
}

func (snapshot *TradesSnapshot) CheckPrivacy(observer User, observed User) {
	if observed.Privacy == "all" {
		snapshot.PrivacyStatus.Status = "OK"
		snapshot.PrivacyStatus.Reason = "observed ALL"
		return
	}

	if observer.Wallet == "" {
		snapshot.PrivacyStatus.Status = "KO"
		snapshot.PrivacyStatus.Reason = "user is not logged in"
		return
	}

	if observer.Wallet == observed.Wallet {
		snapshot.PrivacyStatus.Status = "OK"
		snapshot.PrivacyStatus.Reason = "user access its own profile"
		return
	}

	switch observed.Privacy {
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
					AND followto = $2;`, observer.Wallet, observed.Wallet).Scan(
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
					AND subscribeto = $2;`, observer.Wallet, observed.Wallet).Scan(
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
			"observed": observed.Wallet,
			"observer": observer.Wallet,
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
			"sessionId":  ws_trade.SessionId,
			"custom_msg": "Failed running sending initial snapshot",
		}).Error(err)
		return
	}
}

func (ws_trade *WsTrade) WaitToTerminate() {
	for {
		_, message, err := ws_trade.Ws.ReadMessage()
		if string(message) == "ping" {
			continue
		} else if err != nil {
			observers := []WsTrade{}
			for _, observer := range trades_wss[ws_trade.Observed.Wallet] {
				if observer.SessionId != ws_trade.SessionId {
					observers = append(observers, observer)
				}
			}
			trades_wss[ws_trade.Observed.Wallet] = observers

			if len(trades_wss[ws_trade.Observed.Wallet]) == 0 {
				delete(trades_wss, ws_trade.Observed.Wallet)
			}

			ws_trade.Ws.Close()
			return
		} else {
			ws_trade.Ws.Close()
			log.WithFields(log.Fields{
				"sessionId":  ws_trade.SessionId,
				"custom_msg": "Failed terminating trade ws",
			}).Error(err)
			return
		}
	}
}

func (ws_trade *WsTrade) WaitToSendMessage() {
	for {
		s1 := <-ws_trade.Channel
		err := ws_trade.Ws.WriteJSON(s1)
		if err != nil {
			ws_trade.Ws.Close()
			log.WithFields(log.Fields{
				"sessionId":  ws_trade.SessionId,
				"custom_msg": "Failed running sending snapshot",
			}).Error(err)
			return
		}
	}
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func StartTradesWs(w http.ResponseWriter, r *http.Request) {
	url_split := strings.Split(r.URL.Path, "/")

	// OBSERVED
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

	// OBSERVER
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
	}

	fmt.Println(session)

	if session.Origin != "web" {
		log.Error("Failed starting ws, origin not web")
		w.WriteHeader(http.StatusBadRequest)
		return
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

	user_connection := Connection{
		Observer: observer,
		Observed: observed,
	}
	user_connection.CheckConnection()
	user_connection.CheckPrivacy()

	snapshot := observed.GetSnapshot()
	snapshot.IsFollower = user_connection.IsFollower
	snapshot.IsSubscriber = user_connection.IsSubscriber
	snapshot.PrivacyStatus = user_connection.Privacy

	if observer.Wallet != observed.Wallet {
		observed.CheckVisibility(&snapshot)
	}

	c := make(chan TradesSnapshot)
	ws_trade := WsTrade{observer, observed, session_id, c, ws}

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

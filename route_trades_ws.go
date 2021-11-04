package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsTrade struct {
	UserToSee User
	RequestId string
	Channel   chan TradesSnapshot
	Ws        *websocket.Conn
}

type PrivacyStatus struct {
	Status string
	Reason string
}

func StartTradesWs(w http.ResponseWriter, r *http.Request) {
	url_split := strings.Split(r.URL.Path, "/")

	if len(url_split) < 4 {
		log.WithFields(log.Fields{
			"urlPath": r.URL.Path,
		}).Warn("Failed starting ws, wrong url")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username := url_split[2]
	request_id := url_split[3]

	userToSee, err := SelectUser("username", username)
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

	status := CheckPrivacy(r, userToSee)
	if status.Status != "OK" {
		err = ws.WriteJSON(status)
		if err != nil {
			log.WithFields(log.Fields{
				"urlPath":    r.URL.Path,
				"userToSee":  userToSee,
				"custom_msg": "Failed returning status",
			}).Error(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	c := make(chan TradesSnapshot)
	ws_trade := WsTrade{userToSee, request_id, c, ws}
	trades_wss[username] = append(trades_wss[username], ws_trade)

	go ws_trade.SendInitialSnapshot()
	go ws_trade.WaitToTerminate()
	go ws_trade.WaitToSendMessage()
}

func CheckPrivacy(request *http.Request, userToSee User) (status PrivacyStatus) {
	if userToSee.Privacy == "all" {
		status.Status = "OK"
		status.Reason = "userToSee ALL"
		return
	}

	session, err := GetSession(request, "cookie")
	if err != nil {
		status.Status = "KO"
		status.Reason = "cookie"
		return
	}

	user, err := SelectUser("code", session.UserCode)
	if err != nil {
		status.Status = "KO"
		status.Reason = "invalid user code"
		return
	}

	if user.Code == userToSee.Code {
		status.Status = "OK"
		status.Reason = "user access its own profile"
		return
	}

	switch userToSee.Privacy {
	case "private":
		status.Status = "KO"
		status.Reason = "private"
		return
	case "followers":
		var isfollower bool
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE followfrom = $1
					AND followto = $2;`, user.Code, userToSee.Code).Scan(
			&isfollower,
		)
		if isfollower {
			status.Status = "OK"
			status.Reason = "user is follower"
			return
		} else {
			status.Status = "KO"
			status.Reason = "user is not follower"
			return
		}
	case "subscribers":
		var issubscriber bool
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE subscribefrom = $1
					AND subscribeto = $2;`, user.Code, userToSee.Code).Scan(
			&issubscriber,
		)
		if issubscriber {
			status.Status = "OK"
			status.Reason = "user is subscriber"
			return
		} else {
			status.Status = "KO"
			status.Reason = "user is not subscriber"
			return
		}
	default:
		log.WithFields(log.Fields{
			"userToSee": userToSee.Code,
			"user":      user.Code,
		}).Warn("Not possible to determine user's privacy")
		status.Status = "KO"
		status.Reason = "unknown reason"
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

func (ws_trade *WsTrade) SendInitialSnapshot() {
	ws_trade_output := ws_trade.UserToSee.GetSnapshot()
	if len(ws_trade_output.Trades) > 0 {
		err := ws_trade.Ws.WriteJSON(ws_trade_output)
		if err != nil {
			ws_trade.Ws.Close()
			log.WithFields(log.Fields{
				"username":   ws_trade.UserToSee.UserName,
				"custom_msg": "Failed running sending initial snapshot",
			}).Error(err)
			return
		}
	}
}

func (ws_trade *WsTrade) WaitToTerminate() {
	for {
		_, _, err := ws_trade.Ws.ReadMessage()
		if err != nil {
			for i, v := range trades_wss[ws_trade.UserToSee.UserName] {
				if v.RequestId == ws_trade.RequestId {
					trades_wss[ws_trade.UserToSee.UserName] = append(
						trades_wss[ws_trade.UserToSee.UserName][:i],
						trades_wss[ws_trade.UserToSee.UserName][i+1:]...)
				}
			}
			ws_trade.Ws.Close()
			return
		} else {
			ws_trade.Ws.Close()
			log.WithFields(log.Fields{
				"username":   ws_trade.UserToSee.UserName,
				"custom_msg": "Failed terminating trade ws",
			}).Error(err)
			return
		}
	}
}

func (ws_trade *WsTrade) WaitToSendMessage() {
	for {
		s1 := <-ws_trade.Channel
		if s1.UserDetails.Username == ws_trade.UserToSee.UserName {
			err := ws_trade.Ws.WriteJSON(s1)
			if err != nil {
				ws_trade.Ws.Close()
				log.WithFields(log.Fields{
					"username":   ws_trade.UserToSee.UserName,
					"custom_msg": "Failed running sending snapshot",
				}).Error(err)
				return
			}
		}
	}
}

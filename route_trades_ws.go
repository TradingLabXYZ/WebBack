package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

func NomeDaDecidere() {
	// Check if user exist
	// Check privacy
	// Add user to websockets list
	// Check if origin is authoized

}

type WsTrade struct {
	UserToSee User
	RequestId string
	Channel   chan TradesSnapshot
	Ws        *websocket.Conn
}

func StartTradesWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetTrades..."))

	request_id := mux.Vars(r)["requestid"]
	username := mux.Vars(r)["username"]

	userToSee, err := SelectUser("username", username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Warn("User not found")
		return
	}

	status := CheckPrivacy(r, userToSee)
	if status != "OK" {
		w.Write([]byte(status))
		return
	}

	ws, err := InstanciateTradeWs(w, r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.WithFields(log.Fields{
			"custom_msg": "Failed instanciating TradeWs",
		}).Error(err)
		return
	}

	c := make(chan TradesSnapshot)
	ws_trade := WsTrade{userToSee, request_id, c, ws}
	trades_wss[username] = append(trades_wss[username], ws_trade)

	go ws_trade.SendInitialSnapshot()
	go ws_trade.WaitToTerminate()
	go ws_trade.WaitToSendMessage()
}

func CheckPrivacy(request *http.Request, userToSee User) (status string) {
	fmt.Println(Gray(8-1, "Starting CheckUserPrivacy..."))

	if userToSee.Privacy == "all" {
		return "OK"
	}

	session, err := GetSession(request, "cookie")
	if err != nil {
		return "KO"
	}

	user, err := SelectUser("email", session.Email)
	if err != nil {
		return "KO"
	}

	if user.Id == userToSee.Id {
		return "OK"
	}

	switch userToSee.Privacy {
	case "private":
		return `{"Status": "denied", "Reason": "private"}`
	case "followers":
		var isfollower bool
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE followto = $1
					AND followfrom = $2;`, user.Id, userToSee.Id).Scan(
			&isfollower,
		)
		if isfollower {
			return "OK"
		} else {
			return `{"Status": "denied", "Reason": "follow"}`
		}
	case "subscribers":
		var issubscriber bool
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE subscribeto = $1
					AND subscribefrom = $2;`, user.Id, userToSee.Id).Scan(
			&issubscriber,
		)
		if issubscriber {
			return "OK"
		} else {
			return `{"Status": "denied", "Reason": "subscribe"}`
		}
	default:
		return `{"Status": "denied", "Reason": "unknown"}`
	}
}

func InstanciateTradeWs(w http.ResponseWriter, r *http.Request) (ws *websocket.Conn, err error) {
	var upgrader = websocket.Upgrader{
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

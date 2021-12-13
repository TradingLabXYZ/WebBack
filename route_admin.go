package main

import (
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type OnlineUser struct {
	Wallet   string
	Observed []string
}

type OnlineUsers struct {
	Count int
	Users []OnlineUser
}

func SelectActivity(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	admin_token := os.Getenv("ADMIN_TOKEN")

	if token != admin_token {
		log.Warn("Attempted accessing admin with invalid token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// clean trade_wss
	unique_observers := make(map[string]bool)
	all := [][]string{}
	for observed := range trades_wss {
		for _, q := range trades_wss[observed] {
			var observer string
			if q.Observer.Wallet != "" {
				observer = q.Observer.Wallet
			} else {
				observer = q.SessionId
			}
			if !unique_observers[observer] {
				unique_observers[observer] = true
			}
			all = append(all, []string{observer, q.Observed.Wallet})
		}
	}

	// prepare output for html template
	online_users := OnlineUsers{}
	for observer := range unique_observers {
		online_user := OnlineUser{}
		online_user.Wallet = observer
		for _, pair := range all {
			if observer == pair[0] {
				online_user.Observed = append(online_user.Observed, pair[1])
			}
		}
		online_users.Users = append(online_users.Users, online_user)
	}

	online_users.Count = len(online_users.Users)

	tmpl := template.Must(template.ParseFiles("templates/admin_dashboard.html"))
	tmpl.Execute(w, online_users)
	http.ListenAndServe(":80", nil)
}

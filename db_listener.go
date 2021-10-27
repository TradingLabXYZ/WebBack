package main

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DbListener struct {
	Listener *pq.Listener
}

func InstanciateActivityMonitor() {
	listener := DbListener{}
	listener.Instanciate()
	listener.Listen()
}

func (l *DbListener) Instanciate() {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.WithFields(log.Fields{
				"custom_msg": "Failed instanciating database listener",
			}).Error(err.Error())
		}
	}
	l.Listener = pq.NewListener(DbUrl, 10*time.Second, time.Minute, reportProblem)
	err := l.Listener.Listen("activity_update")
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed listening database",
		}).Error(err.Error())
	}
}

func (db_listener *DbListener) Listen() {
	for {
		n := <-db_listener.Listener.Notify
		user_code := n.Extra
		fmt.Println("RECEIVED DB CALL", user_code)
		DistpachSnapshots(user_code)
	}
}

func DistpachSnapshots(user_code string) {
	user, _ := SelectUser("code", user_code)
	user_snapshot := user.GetSnapshot()
	for _, q := range trades_wss[user.UserName] {
		q.Channel <- user_snapshot
	}
}

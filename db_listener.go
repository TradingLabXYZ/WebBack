package main

import (
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
		user_wallet := n.Extra
		DistpachSnapshots(user_wallet)
	}
}

func DistpachSnapshots(user_wallet string) {
	observed, _ := SelectUser("wallet", user_wallet)
	snapshot := observed.GetSnapshot()
	for _, q := range trades_wss[observed.Wallet] {
		snapshot.CheckRelation(q.Observer, observed)
		snapshot.CheckPrivacy(q.Observer, observed)
		if snapshot.PrivacyStatus.Status == "KO" {
			snapshot.Trades = nil
		}
		q.Channel <- snapshot
	}
}

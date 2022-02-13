package main

import (
	"context"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func (l *DbListener) Instanciate() {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.WithFields(log.Fields{
				"custom_msg": "Failed instanciating database listener",
			}).Error(err.Error())
			discordNotifier.Send(
				context.Background(),
				"db_listener | Instanciate",
				"Failed instanciating database listener",
			)
		}
	}
	l.Listener = pq.NewListener(DbUrl, 10*time.Second, time.Minute, reportProblem)
	err := l.Listener.Listen("activity_update")
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed listening database",
		}).Error(err.Error())
		discordNotifier.Send(
			context.Background(),
			"db_listener | Instanciate",
			"Failed listening database",
		)
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
		user_connection := Connection{
			Observer: q.Observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		snapshot.IsFollower = user_connection.IsFollower
		snapshot.IsSubscriber = user_connection.IsSubscriber
		snapshot.PrivacyStatus = user_connection.Privacy
		if snapshot.PrivacyStatus.Status == "KO" {
			snapshot.Trades = nil
		}
		q.Channel <- snapshot
	}
}

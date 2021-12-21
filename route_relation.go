package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func UpdateFollower(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed updating follower, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	string_status := mux.Vars(r)["status"]
	status, err := strconv.ParseBool(string_status)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed updating follower, wrong status",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	wallet := mux.Vars(r)["wallet"]
	if status {
		Db.Exec(`
			DELETE FROM followers
			WHERE followfrom = $1 AND followto = $2;`,
			session.UserWallet, wallet)
	} else {
		Db.Exec(`
			INSERT INTO followers (followfrom, followto, createdat)
			VALUES ($1, $2, current_timestamp);`,
			session.UserWallet, wallet)
	}
	w.Write([]byte("OK"))
}

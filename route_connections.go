package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

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

func SelectConnections(w http.ResponseWriter, r *http.Request) {
	wallet := mux.Vars(r)["wallet"]
	observed, err := SelectUser("wallet", wallet)
	if err != nil {
		log.WithFields(log.Fields{
			"urlPath":  r.URL.Path,
			"observed": observed,
		}).Warn("Failed selecting relations, user not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, _ := GetSession(r, "header")
	observer, err := SelectUser("wallet", session.UserWallet)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed selecting relations, not available session",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user_connection := Connection{
		Observer: observer,
		Observed: observed,
	}
	user_connection.CheckConnection()
	user_connection.CheckPrivacy()
  
	var relations Relations

	relations.Privacy = user_connection.Privacy

	if relations.Privacy.Status != "OK" {
		json.NewEncoder(w).Encode(relations)
		return
	}

	followers_sql := func(wg *sync.WaitGroup) {
		defer wg.Done()
		query := `
			SELECT
				f.followfrom,
				u.profilepicture,
				COUNT(t.*)
			FROM followers f
			LEFT JOIN users u ON (f.followfrom = u.wallet)
			LEFT JOIN trades t ON(f.followfrom = t.userwallet)
			WHERE f.followto = $1
			GROUP BY 1, 2
			ORDER BY 3 DESC;`
		followers_rows, err := Db.Query(query, wallet)
		defer followers_rows.Close()
		if err != nil {
			log.Error(err)
		}
		for followers_rows.Next() {
			var follower Follower
			if err = followers_rows.Scan(
				&follower.Wallet,
				&follower.ProfilePicture,
				&follower.CountTrades,
			); err != nil {
				log.Error(err)
			}
			relations.Followers = append(relations.Followers, follower)
		}
	}
	following_sql := func(wg *sync.WaitGroup) {
		defer wg.Done()
		query := `
			SELECT
				f.followto,
				u.profilepicture,
				COUNT(t.*)
			FROM followers f
			LEFT JOIN users u ON (f.followto = u.wallet)
			LEFT JOIN trades t ON(f.followto = t.userwallet)
			WHERE f.followfrom = $1
			GROUP BY 1, 2
			ORDER BY 3 DESC;`
		following_rows, err := Db.Query(query, wallet)
		defer following_rows.Close()
		if err != nil {
			log.Error(err)
		}
		for following_rows.Next() {
			var following Following
			if err = following_rows.Scan(
				&following.Wallet,
				&following.ProfilePicture,
				&following.CountTrades,
			); err != nil {
				log.Error(err)
			}
			relations.Following = append(relations.Following, following)
		}
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go followers_sql(&wg)
	go following_sql(&wg)
	wg.Wait()

	json.NewEncoder(w).Encode(relations)
}

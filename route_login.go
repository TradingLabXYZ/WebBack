package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Login(w http.ResponseWriter, r *http.Request) {
	wallet := mux.Vars(r)["wallet"]
	timezone := mux.Vars(r)["timezone"]

	clean_timezone := strings.ReplaceAll(timezone, "_", "/")
	timeLoc, err := time.LoadLocation(clean_timezone)
	if err != nil {
		log.WithFields(log.Fields{
			"wallet":     wallet,
			"custom_msg": "Attemped accessing with invalid timezone",
		}).Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_ = timeLoc

	user_wallet := UserWallet{
		Wallet: wallet,
	}

	validate := validator.New()
	err = validate.Struct(user_wallet)
	if err != nil {
		log.WithFields(log.Fields{
			"wallet":     wallet,
			"custom_msg": "Attemped accessing with invald wallet",
		}).Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsWalletInSessions(wallet) {
		InsertUser(wallet)
		InsertVisibility(wallet)
	}

	user, err := SelectUser("wallet", wallet)

	session, err := user.InsertSession("web")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user_data := struct {
		SessionId      string
		Wallet         string
		Username       string
		Twitter        string
		Discord        string
		Github         string
		ProfilePicture string
		Privacy        string
		Followers      int
		Subscribers    int
		MonthlyFee     string
		Visibility     VisibilityStatus
	}{
		session.Code,
		user.Wallet,
		user.Username,
		user.Twitter,
		user.Discord,
		user.Github,
		user.ProfilePicture,
		user.Privacy,
		user.Followers,
		user.Subscribers,
		user.MonthlyFee,
		user.Visibility,
	}

	json.NewEncoder(w).Encode(user_data)
}

package main

import (
	"encoding/json"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Login(w http.ResponseWriter, r *http.Request) {
	wallet := mux.Vars(r)["wallet"]

	user_wallet := UserWallet{
		Wallet: wallet,
	}

	validate := validator.New()
	err := validate.Struct(user_wallet)
	if err != nil {
		log.WithFields(log.Fields{
			"wallet":     wallet,
			"custom_msg": "Attemped accessing with invald wallet",
		}).Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := SelectUser("wallet", user_wallet.Wallet)
	if user == (User{}) {
		InsertUser(wallet)
		InsertVisibility(wallet)
		user, err = SelectUser("wallet", user_wallet.Wallet)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

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

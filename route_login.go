package main

import (
	"encoding/json"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type UserWallet struct {
	Wallet string `validate="eth_addr"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	wallet := mux.Vars(r)["wallet"]

	user_wallet := UserWallet{
		Wallet: wallet,
	}

	validate := validator.New()
	err := validate.Struct(user_wallet)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		first_error := validationErrors[0].Tag()
		w.Write([]byte(first_error))
		return
	}

	user, err := SelectUser("wallet", user_wallet.Wallet)
	if user == (User{}) {
		InsertUser(wallet)
		user, err = SelectUser("wallet", user_wallet.Wallet)
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	session, err := user.InsertSession()
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
		Plan           string
	}{
		session.Code,
		user.Wallet,
		user.UserName,
		user.Twitter,
		user.Discord,
		user.Github,
		user.ProfilePicture,
		user.Privacy,
		user.Plan,
	}

	json.NewEncoder(w).Encode(user_data)
}

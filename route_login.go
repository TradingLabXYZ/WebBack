package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	if err == sql.ErrNoRows {
		fmt.Println("OK SONO QUI")
		InsertUser(wallet)
		user, err = SelectUser("wallet", user_wallet.Wallet)
		if err != nil {
			fmt.Println(err)
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
		Username       string
		Wallet         string
		ProfilePicture string
		Twitter        string
		Website        string
		Privacy        string
		Plan           string
	}{
		session.Code,
		user.UserName,
		user.Wallet,
		user.ProfilePicture,
		user.Twitter,
		user.Website,
		user.Privacy,
		user.Plan,
	}

	json.NewEncoder(w).Encode(user_data)
}

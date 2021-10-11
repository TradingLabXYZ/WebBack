package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Login..."))

	decoder := json.NewDecoder(r.Body)
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Error(err)
	}

	user := UserByEmail(body.Email)
	user.LoginPassword = body.Password
	if user.Id != 0 && user.Password == Encrypt(user.LoginPassword) {
		log.Info("Log in valid user " + strconv.Itoa(user.Id))
		session := user.CreateSession()
		user_data := struct {
			SessionId      string
			UserName       string
			Code           string
			ProfilePicture string
			Twitter        string
			Website        string
		}{
			session.Uuid,
			user.UserName,
			user.Code,
			user.ProfilePicture,
			user.Twitter,
			user.Website,
		}
		json.NewEncoder(w).Encode(user_data)
	} else {
		log.Info("Log in not valid user " + user.Email)
	}
}

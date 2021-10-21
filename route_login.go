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

	user := SelectUser("email", body.Email)
	user.LoginPassword = body.Password
	if user.Id != 0 && user.Password == Encrypt(user.LoginPassword) {
		log.Info("Log in valid user " + strconv.Itoa(user.Id))
		session := user.CreateSession()
		user_data := struct {
			SessionId      string
			Username       string
			Email          string
			Code           string
			ProfilePicture string
			Twitter        string
			Website        string
			Privacy        string
			Plan           string
		}{
			session.Uuid,
			user.UserName,
			user.Email,
			user.Code,
			user.ProfilePicture,
			user.Twitter,
			user.Website,
			user.Privacy,
			user.Plan,
		}
		json.NewEncoder(w).Encode(user_data)
	} else {
		log.Info("Log in not valid user " + strconv.Itoa(user.Id))
	}
}

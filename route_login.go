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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user := SelectUser("email", body.Email)
	if user.Id == 0 {
		log.Warn("Email not present")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user.LoginPassword = body.Password
	encrypted_password, err := Encrypt(user.LoginPassword)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if user.Password != encrypted_password {
		log.Warn("Passwords do not match")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	log.Info("Log in valid user " + strconv.Itoa(user.Id))
	session, err := user.CreateSession()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
}

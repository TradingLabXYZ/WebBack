package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,containsany=!?*()_$,containsany=1234567890"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	credentials := LoginCredentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(credentials)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		first_error := validationErrors[0].Tag()
		w.Write([]byte(first_error))
		return
	}

	user, err := SelectUser("email", credentials.Email)
	if err != nil {
		w.Write([]byte("User not found"))
		return
	}

	encrypted_password, err := Encrypt(credentials.Password)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if user.Password != encrypted_password {
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

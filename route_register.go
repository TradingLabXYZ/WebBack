package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/logrusorgru/aurora"
)

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Register..."))
	decoder := json.NewDecoder(r.Body)
	s := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&s)
	if err != nil {
		panic(err)
	}
	InsertUser(s.Email, s.Username, s.Password)
	json.NewEncoder(w).Encode("OK")
}

func InsertUser(email string, username string, password string) {
	fmt.Println(Gray(8-1, "Starting InsertUser..."))
	statement := `
		INSERT INTO users (email, username, password, createdat, updatedat)
		VALUES ($1, $2, $3, current_timestamp, current_timestamp);`
	_, err := DbWebApp.Exec(statement, email, username, Encrypt(password))
	if err != nil {
		panic(err)
	}
}

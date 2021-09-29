package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	permission := "all"
	default_profile_picture := os.Getenv("CDN_PATH") + "/profile_pictures/default_picture.png"
	statement := `
		INSERT INTO users (code, email, username, password, permission, profilepicture, createdat, updatedat)
		VALUES (SUBSTR(MD5(RANDOM()::TEXT), 0, 12), $1, $2, $3, $4, $5, current_timestamp, current_timestamp);`
	_, err := DbWebApp.Exec(statement, email, username, Encrypt(password), permission, default_profile_picture)
	if err != nil {
		panic(err)
	}
}

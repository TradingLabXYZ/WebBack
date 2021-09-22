package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/logrusorgru/aurora"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting Authenticate..."))

	decoder := json.NewDecoder(r.Body)
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}

	user := UserByEmail(body.Email)
	user.LoginPassword = body.Password
	if user.Id != 0 && user.Password == Encrypt(user.LoginPassword) {
		fmt.Println(Blue("Log in valid..."))
		session := user.CreateSession()
		user_data := struct {
			SessionId string
			UserName  string
		}{
			session.Uuid,
			user.UserName,
		}
		json.NewEncoder(w).Encode(user_data)
	} else {
		fmt.Println(Yellow("Log in not valid..."))
	}
}

func (user *User) CreateSession() (session Session) {
	fmt.Println(Gray(8-1, "Starting CreateSession..."))
	session_sql := `
		INSERT INTO sessions (uuid, email, userid, createdat)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, email, userid, createdat;`
	DbWebApp.QueryRow(
		session_sql,
		createUUID(),
		user.Email,
		user.Id,
		time.Now(),
	).Scan(
		&session.Id,
		&session.Uuid,
		&session.Email,
		&session.UserId,
		&session.CreatedAt,
	)
	return
}

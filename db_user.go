package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

type User struct {
	Id            int
	Email         string
	UserName      string
	LoginPassword string
	Password      string
	Permission    string
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func UserByEmail(email string) (user User) {
	fmt.Println(Gray(8-1, "Starting UserByEmail..."))
	_ = DbWebApp.QueryRow(`
					SELECT
						id,
						email,
						password,
						username,
						permission
					FROM users
					WHERE email = $1;`, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.UserName,
		&user.Permission,
	)

	return
}

func UserByUsername(username string) (user User) {
	fmt.Println(Gray(8-1, "Starting UserByUsername..."))

	rows, err := DbWebApp.Query(`
		SELECT
			id,
			email,
			permission
		FROM users
		WHERE username = $1;`, username)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		if err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Permission); err != nil {
			panic(err.Error())
		}
	}

	return
}

func SelectSession(r *http.Request) (session Session) {
	fmt.Println(Gray(8-1, "Starting SelectSession..."))

	auth := r.Header["Authorization"][0]
	session.Uuid = strings.Split(auth, "sessionId=")[1]

	err := DbWebApp.QueryRow(`
      SELECT
        id,
        uuid,
        email,
        userid,
        createdat
      FROM sessions
      WHERE uuid = $1;`, session.Uuid).Scan(
		&session.Id,
		&session.Uuid,
		&session.Email,
		&session.UserId,
		&session.CreatedAt,
	)
	if err != nil {
		fmt.Println(Gray(8-1, "No session found, user not logged in..."))
	}

	return
}

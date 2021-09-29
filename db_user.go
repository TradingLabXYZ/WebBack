package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

type User struct {
	Id             int
	Code           string
	Email          string
	UserName       string
	LoginPassword  string
	Password       string
	Permission     string
	ProfilePicture string
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
						code,
						email,
						password,
						username,
						permission,
						profilepicture
					FROM users
					WHERE email = $1;`, email).Scan(
		&user.Id,
		&user.Code,
		&user.Email,
		&user.Password,
		&user.UserName,
		&user.Permission,
		&user.ProfilePicture,
	)

	return
}

func UserByUsername(username string) (user User) {
	fmt.Println(Gray(8-1, "Starting UserByUsername..."))

	rows, err := DbWebApp.Query(`
		SELECT
			id,
			code,
			email,
			permission
		FROM users
		WHERE username = $1;`, username)
	defer rows.Close()
	if err != nil {
		log.Error(err)
	}
	for rows.Next() {
		if err := rows.Scan(
			&user.Id,
			&user.Code,
			&user.Email,
			&user.Permission); err != nil {
			log.Error(err)
		}
	}

	return
}

func SelectSession(r *http.Request) (session Session) {
	fmt.Println(Gray(8-1, "Starting SelectSession..."))
	var auth string
	if len(r.Header["Authorization"]) == 0 {
		log.Warn("User not authenticated")
		return
	} else {
		auth = r.Header["Authorization"][0]
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
			log.Info("No session found, user not logged in...")
		}
		return
	}
}

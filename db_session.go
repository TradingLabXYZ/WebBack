package main

import (
	"fmt"
	"net/http"
	"strings"

	"time"

	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
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
			log.Warn("No session found, user not logged in...")
			return
		}
		return
	}
}

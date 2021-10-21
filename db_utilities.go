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

type User struct {
	Id             int
	Code           string
	Email          string
	UserName       string
	LoginPassword  string
	Password       string
	Privacy        string
	Plan           string
	ProfilePicture string
	Twitter        string
	Website        string
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
	if len(r.Header["Authorization"]) > 0 {
		split_auth := strings.Split(r.Header["Authorization"][0], "sessionId=")
		if len(split_auth) >= 1 {
			auth = split_auth[1]
		}
	} else {
		for _, cookie := range r.Cookies() {
			if cookie.Name == "sessionId" {
				auth = cookie.Value
			}
		}
	}
	if len(auth) > 0 {
		session.Uuid = auth
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
	} else {
		log.Warn("No session found, user not logged in...")
		return
	}
}

func SelectUser(by string, value string) (user User) {
	fmt.Println(Gray(8-1, "Starting SelectUser..."))
	user_sql := fmt.Sprintf(`
		SELECT
			id,
			code,
			email,
			password,
			username,
			privacy,
			plan,
			profilepicture,
			CASE WHEN twitter IS NULL THEN '' ELSE twitter END AS twitter,
			CASE WHEN website IS NULL THEN '' ELSE website END AS website
		FROM users
		WHERE %s = $1;`, by)
	err := DbWebApp.QueryRow(user_sql, value).Scan(
		&user.Id,
		&user.Code,
		&user.Email,
		&user.Password,
		&user.UserName,
		&user.Privacy,
		&user.Plan,
		&user.ProfilePicture,
		&user.Twitter,
		&user.Website,
	)
	if err != nil {
		log.Warn("No user found...")
		return
	}
	return
}

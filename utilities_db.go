package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"time"
)

type Session struct {
	Code      string
	UserCode  string
	CreatedAt time.Time
}

type User struct {
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

func (user *User) CreateSession() (session Session, err error) {
	uuid, err := CreateUUID()
	if err != nil {
		return
	}

	session_sql := `
		INSERT INTO sessions (code, usercode, createdat)
		VALUES ($1, $2, $3)
		RETURNING code, usercode, createdat;`

	err = Db.QueryRow(
		session_sql,
		uuid,
		user.Code,
		time.Now()).Scan(
		&session.Code,
		&session.UserCode,
		&session.CreatedAt,
	)
	if err != nil {
		err = errors.New("Error inserting new session in db")
		return
	}
	return
}

func GetSession(r *http.Request, using string) (session Session, err error) {
	err = session.ExtractFromRequest(r, using)
	if err != nil {
		return
	}
	err = session.Select()
	return
}

func (session *Session) ExtractFromRequest(r *http.Request, using string) (err error) {
	switch using {
	case "header":
		err = session.ExtractFromHeader(r)
		if err != nil {
			return
		}
	case "cookie":
		err = session.ExtractFromCookie(r)
		if err != nil {
			return
		}
	}
	return
}

func (session *Session) ExtractFromHeader(r *http.Request) (err error) {
	if len(r.Header["Authorization"]) > 0 {
		split_auth := strings.Split(r.Header["Authorization"][0], "sessionId=")
		if len(split_auth) >= 1 {
			session.Code = split_auth[1]
			return
		} else {
			err = errors.New("Could not find sessionId in header")
			return
		}
	} else {
		err = errors.New("Could not find authorization in header")
		return
	}
}

func (session *Session) ExtractFromCookie(r *http.Request) (err error) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "sessionId" {
			session.Code = cookie.Value
		}
	}
	if session.Code == "" {
		err = errors.New("Empty sessionId in cookie")
	}
	return
}

func (session *Session) Select() (err error) {
	err = Db.QueryRow(`
			SELECT
				usercode
			FROM sessions
			WHERE code = $1;`, session.Code).Scan(
		&session.UserCode,
	)
	return
}

func SelectUser(by string, value string) (user User, err error) {
	user_sql := fmt.Sprintf(`
		SELECT
			code,
			email,
			password,
			username,
			privacy,
			plan,
			CASE WHEN profilepicture IS NULL THEN '' ELSE profilepicture END AS profilepicture,
			CASE WHEN twitter IS NULL THEN '' ELSE twitter END AS twitter,
			CASE WHEN website IS NULL THEN '' ELSE website END AS website
		FROM users
		WHERE %s = $1;`, by)
	err = Db.QueryRow(user_sql, value).Scan(
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
	return
}

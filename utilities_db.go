package main

import (
	"errors"
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

func (user *User) CreateSession() (session Session, err error) {
	fmt.Println(Gray(8-1, "Starting CreateSession..."))

	if user.Email == "" {
		err = errors.New("Empty email")
		return
	}

	uuid, err := CreateUUID()
	if err != nil {
		return
	}

	session_sql := `
		INSERT INTO sessions (uuid, email, userid, createdat)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, email, userid, createdat;`

	err = DbWebApp.QueryRow(
		session_sql,
		uuid,
		user.Email,
		user.Id,
		time.Now()).Scan(
		&session.Id,
		&session.Uuid,
		&session.Email,
		&session.UserId,
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
			session.Uuid = split_auth[1]
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
			session.Uuid = cookie.Value
		}
	}
	if session.Uuid == "" {
		err = errors.New("Empty sessionId in cookie")
	}
	return
}

func (session *Session) Select() (err error) {
	err = DbWebApp.QueryRow(`
			SELECT
				email,
				userid
			FROM sessions
			WHERE uuid = $1;`, session.Uuid).Scan(
		&session.Email,
		&session.UserId,
	)
	return
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
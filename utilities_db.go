package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (user *User) InsertSession() (session Session, err error) {
	uuid, err := CreateUUID()
	if err != nil {
		return
	}

	session_sql := `
		INSERT INTO sessions (code, userwallet, createdat)
		VALUES ($1, $2, $3)
		RETURNING code, userwallet, createdat;`

	err = Db.QueryRow(
		session_sql,
		uuid,
		user.Wallet,
		time.Now()).Scan(
		&session.Code,
		&session.UserWallet,
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
		if len(split_auth) > 1 {
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
				userwallet
			FROM sessions
			WHERE code = $1;`, session.Code).Scan(
		&session.UserWallet,
	)
	return
}

func InsertUser(wallet string) {
	default_profile_picture := os.Getenv("CDN_PATH") + "/profile_pictures/default_picture.png"
	statement := `
		INSERT INTO users (
			wallet, profilepicture, username, privacy,
			plan, createdat, updatedat)
		VALUES (
			$1, $2, '', 'all', 'basic',
			current_timestamp, current_timestamp);`
	_, err := Db.Exec(statement, wallet, default_profile_picture)
	if err != nil {
		log.Error(err)
		return
	}
}

func SelectUser(by string, value string) (user User, err error) {
	user_sql := fmt.Sprintf(`
		SELECT
			wallet,
			TO_CHAR(createdat, 'Month') || ' ' || TO_CHAR(createdat, 'YYYY') AS jointime,
			CASE WHEN username IS NULL THEN '' ELSE username END AS username,
			CASE WHEN twitter IS NULL THEN '' ELSE twitter END AS twitter,
			CASE WHEN discord IS NULL THEN '' ELSE discord END AS discord,
			CASE WHEN github IS NULL THEN '' ELSE github END AS github,
			privacy,
			plan,
			CASE WHEN profilepicture IS NULL THEN '' ELSE profilepicture END AS profilepicture,
			f.count_followers,
			fo.count_followings,
			s.count_subscribers,
			mf.monthly_fee
		FROM users
		LEFT JOIN (
			SELECT
				COUNT(*) AS count_followers
			FROM followers
			WHERE followto = $1) f ON(1=1)
		LEFT JOIN (
			SELECT
				COUNT(*) AS count_followings
			FROM followers
			WHERE followfrom = $1) fo ON(1=1)
		LEFT JOIN (
			SELECT
				COUNT(*) AS count_subscribers
			FROM subscribers
			WHERE subscribeto = $1) s ON(1=1)
		LEFT JOIN (
			SELECT DISTINCT ON(sender)
				sender,
				createdat,
				payload#>>'{Value}' AS monthly_fee
			FROM smartcontractevents
			WHERE name = 'ChangePlan'
			AND sender = $1
			ORDER BY 1, 2 DESC) mf ON(1=1)
		WHERE %s = $1;`, by)
	err = Db.QueryRow(user_sql, value).Scan(
		&user.Wallet,
		&user.JoinTime,
		&user.Username,
		&user.Twitter,
		&user.Discord,
		&user.Github,
		&user.Privacy,
		&user.Plan,
		&user.ProfilePicture,
		&user.Followers,
		&user.Followings,
		&user.Subscribers,
		&user.MonthlyFee,
	)

	return
}

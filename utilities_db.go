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

func (user *User) InsertSession(origin string) (session Session, err error) {
	uuid, err := CreateUUID()
	if err != nil {
		return
	}

	session_sql := `
		INSERT INTO sessions (code, userwallet, origin, createdat)
		VALUES ($1, $2, $3, $4)
		RETURNING code, userwallet, origin, createdat;`

	err = Db.QueryRow(
		session_sql,
		uuid,
		user.Wallet,
		origin,
		time.Now()).Scan(
		&session.Code,
		&session.UserWallet,
		&session.Origin,
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
		auth := strings.Split(r.Header["Authorization"][0], "sessionId=")
		if auth[1] != "" {
			session.Code = strings.Split(auth[1], ";")[0]
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
				userwallet,
				origin
			FROM sessions
			WHERE code = $1;`, session.Code).Scan(
		&session.UserWallet,
		&session.Origin,
	)
	return
}

func IsWalletInSessions(wallet string) (exists bool) {
	err := Db.QueryRow(`
			SELECT
				true
			FROM sessions
			WHERE userwallet = $1;`, wallet).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func InsertUser(wallet string) {
	default_profile_picture := os.Getenv("CDN_PATH") + "/profile_pictures/default_picture.png"
	statement := `
		INSERT INTO users (
			wallet, profilepicture, username, privacy,
			createdat, updatedat)
		VALUES (
			$1, $2, '', 'all',
			current_timestamp, current_timestamp);`
	_, err := Db.Exec(statement, wallet, default_profile_picture)
	if err != nil {
		log.Error(err)
		return
	}
}

func InsertVisibility(wallet string) {
	statement := `
		INSERT INTO visibilities (
			wallet,
			totalcounttrades,
			totalportfolio,
			totalreturn,
			totalroi,
			tradeqtyavailable,
			tradevalue,
			tradereturn,
			traderoi,
			subtradesall,
			subtradereasons,
			subtradequantity,
			subtradeavgprice,
			subtradetotal)
		VALUES (
			$1, TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE,
			TRUE, TRUE, TRUE, TRUE);`
	_, err := Db.Exec(statement, wallet)
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
				CASE WHEN profilepicture IS NULL THEN '' ELSE profilepicture END AS profilepicture,
				f.count_followers,
				fo.count_followings,
				s.count_subscribers,
				CASE WHEN mf.monthly_fee IS NULL THEN '0' ELSE mf.monthly_fee END AS monthly_fee
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
					createdat AS eventcreatedat,
					payload#>>'{Value}' AS monthly_fee
				FROM smartcontractevents
				WHERE name = 'ChangePlan'
				AND LOWER(sender) = LOWER($1)
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
		&user.ProfilePicture,
		&user.Followers,
		&user.Followings,
		&user.Subscribers,
		&user.MonthlyFee,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"by":         by,
			"value":      value,
			"custom_msg": "Failed selecting user",
		}).Error(err)
		return
	}

	visibility_sql := `
			SELECT
				totalcounttrades,
				totalportfolio,
				totalreturn,
				totalroi,
				tradeqtyavailable,
				tradevalue,
				tradereturn,
				traderoi,
				subtradesall,
				subtradereasons,
				subtradequantity,
				subtradeavgprice,
				subtradetotal
			FROM visibilities
			WHERE wallet = $1;`

	err = Db.QueryRow(
		visibility_sql,
		&user.Wallet).Scan(
		&user.Visibility.TotalCountTrades,
		&user.Visibility.TotalPortfolio,
		&user.Visibility.TotalReturn,
		&user.Visibility.TotalRoi,
		&user.Visibility.TradeQtyAvailable,
		&user.Visibility.TradeValue,
		&user.Visibility.TradeReturn,
		&user.Visibility.TradeRoi,
		&user.Visibility.SubtradesAll,
		&user.Visibility.SubtradeReasons,
		&user.Visibility.SubtradeQuantity,
		&user.Visibility.SubtradeAvgPrice,
		&user.Visibility.SubtradeTotal,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"by":         by,
			"value":      value,
			"custom_msg": "Failed selecting visibilities",
		}).Error(err)
		return
	}

	return
}

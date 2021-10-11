package main

import (
	"fmt"

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
	Privacy        string
	Plan           string
	ProfilePicture string
	Twitter        string
	Website        string
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
						privacy,
						plan,
						profilepicture,
						CASE WHEN twitter IS NULL THEN '' ELSE twitter END AS twitter,
						CASE WHEN website IS NULL THEN '' ELSE website END AS website
					FROM users
					WHERE email = $1;`, email).Scan(
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

	return
}

func UserByUsername(username string) (user User) {
	fmt.Println(Gray(8-1, "Starting UserByUsername..."))

	rows, err := DbWebApp.Query(`
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
			&user.Password,
			&user.UserName,
			&user.Privacy,
			&user.Plan,
			&user.ProfilePicture,
			&user.Twitter,
			&user.Website,
		); err != nil {
			log.Error(err)
		}
	}

	return
}

func UserByUsercode(usercode string) (user User) {
	fmt.Println(Gray(8-1, "Starting UserByUsercode..."))

	rows, err := DbWebApp.Query(`
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
		WHERE code = $1;`, usercode)
	defer rows.Close()
	if err != nil {
		log.Error(err)
	}
	for rows.Next() {
		if err := rows.Scan(
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
		); err != nil {
			log.Error(err)
		}
	}

	return
}

package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

func DbConnect() (DbWebApp *sqlx.DB) {

	env := os.Getenv("TL_APP_ENV")
	var DB_NAME string
	if env == "production" {
		DB_NAME = "webappconnectionpool"
	} else if env == "test" {
		DB_NAME = "testwebappconnectionpool"
	}

	WEBAPP_DATABASE_URL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("TL_DB_USER"),
		os.Getenv("TL_DB_PASS"),
		os.Getenv("TL_DB_HOST"),
		os.Getenv("TL_DB_PORT"),
		DB_NAME,
	)

	DbWebApp, err := sqlx.Connect("postgres", WEBAPP_DATABASE_URL)
	if err != nil {
		panic(err.Error())
	}

	if err = DbWebApp.Ping(); err != nil {
		fmt.Println(Bold(Red("Unsuccessfully connected to")), DB_NAME)
		return
	}

	fmt.Println(Bold(Green("Successfully connected to")), DB_NAME)
	return
}

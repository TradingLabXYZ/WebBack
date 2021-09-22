package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

func DbConnect() (DbWebApp *sqlx.DB) {

	WEBAPP_DATABASE_URL := "***REMOVED***"
	DbWebApp, err := sqlx.Connect("postgres", WEBAPP_DATABASE_URL)
	if err != nil {
		panic(err.Error())
	}

	if err = DbWebApp.Ping(); err != nil {
		DbWebApp.Close()
		fmt.Println(Bold(Red("Unsuccessfully connected to WebApp database")))
		return
	}

	fmt.Println(Bold(Green("Successfully connected to WebApp database")))

	return
}

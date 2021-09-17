package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

var DbWebApp *sqlx.DB

func DbConnect() {

	WEBAPP_DATABASE_URL := "***REMOVED***"
	dbWebApp, err := sqlx.Connect("postgres", WEBAPP_DATABASE_URL)
	if err != nil {
		panic(err.Error())
	}
	DbWebApp = dbWebApp

	if err = DbWebApp.Ping(); err != nil {
		DbWebApp.Close()
		fmt.Println(Bold(Red("Unsuccessfully connected to WebApp database")))
		return
	}

	fmt.Println(Bold(Green("Successfully connected to WebApp database")))
}

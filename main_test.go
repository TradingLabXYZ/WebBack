package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	Db = *setUpTestDb()
	log_file := SetUpTestLog()
	defer log_file.Close()
	destroyTestTables()
	createTestTables()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setUpTestDb() (db *sqlx.DB) {
	WEBAPP_DATABASE_URL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("TL_DB_USER"),
		os.Getenv("TL_DB_PASS"),
		os.Getenv("TL_DB_HOST"),
		"25060",
		"testwebapp",
	)
	db, _ = sqlx.Connect("postgres", WEBAPP_DATABASE_URL)
	return
}

func createTestTables() {
	query, err := ioutil.ReadFile("sql/test/create_tables.sql")
	if err != nil {
		panic(err)
	}
	if _, err := Db.Exec(string(query)); err != nil {
		panic(err)
	}
}

func destroyTestTables() {
	query, err := ioutil.ReadFile("sql/test/destroy_tables.sql")
	if err != nil {
		panic(err)
	}
	if _, err := Db.Exec(string(query)); err != nil {
		panic(err)
	}
}

func SetUpTestLog() (file *os.File) {
	file, err := os.OpenFile(
		"logs_test.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0666,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed setting up log file",
		}).Error(err)
		return
	}
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	return
}

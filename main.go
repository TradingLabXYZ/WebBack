package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	. "github.com/logrusorgru/aurora"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var Db sqlx.DB

var Origins = []string{
	"http://localhost:9000",
	"https://tradinglab.xyz",
	"https://www.tradinglab.xyz",
	"https://staging.tradinglab.xyz",
	"https://hoofcoffee-wet0e0.stormkit.dev",
	"https://hoofcoffee-wet0e0--staging.stormkit.dev",
}

func main() {

	r := SetupRoutes()
	c := SetUpCors()
	h := c.Handler(r)

	log_file := SetUpLog()
	defer log_file.Close()
	Db = *setUpDb()
	defer Db.Close()

	go InstanciateTradesDispatcher()

	fmt.Println(Bold(Green("Application is running on port 8080")))
	log.Fatal(http.ListenAndServe(":8080", h))
}

func SetupRoutes() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")

	router.HandleFunc("/user_settings", GetUserSettings).Methods("GET")
	router.HandleFunc("/user_settings", UpdateUserSettings).Methods("POST")
	router.HandleFunc("/update_password", UpdateUserPassword).Methods("POST")
	router.HandleFunc("/update_privacy", UpdateUserPrivacy).Methods("POST")
	router.HandleFunc("/insert_profile_picture", InsertProfilePicture).Methods("PUT")
	router.HandleFunc("/user_premium_data", GetUserPremiumData).Methods("GET")

	router.HandleFunc("/get_trades/{username}/{requestid}", GetTrades)
	router.HandleFunc("/insert_trade", InsertTrade).Methods("POST")
	router.HandleFunc("/close_trade/{tradeid}", CloseTrade).Methods("GET")
	router.HandleFunc("/open_trade/{tradeid}", OpenTrade).Methods("GET")
	router.HandleFunc("/delete_trade/{tradeid}", DeleteTrade).Methods("GET")
	router.HandleFunc("/update_trade", UpdateTrade).Methods("POST")

	router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")
	router.HandleFunc("/stellar_price", SelectStellarPrice).Methods("GET")
	router.HandleFunc("/transaction_credentials", SelectTransactionCredentials).Methods("GET")

	router.HandleFunc("/buy_months", BuyPremiumMonths).Methods("POST")

	return
}

func SetUpLog() (file *os.File) {
	file, err := os.OpenFile(
		"logs.log",
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

func setUpDb() (db *sqlx.DB) {
	env := os.Getenv("TL_APP_ENV")
	var DB_NAME string
	if env == "production" {
		DB_NAME = "webappconnectionpool"
	} else if env == "staging" {
		DB_NAME = "stagingwebappconnectionpool"
	}

	WEBAPP_DATABASE_URL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("TL_DB_USER"),
		os.Getenv("TL_DB_PASS"),
		os.Getenv("TL_DB_HOST"),
		os.Getenv("TL_DB_PORT"),
		DB_NAME,
	)

	db, err := sqlx.Connect("postgres", WEBAPP_DATABASE_URL)
	if err != nil {
		log.WithFields(log.Fields{
			"dbname":     WEBAPP_DATABASE_URL,
			"custom_msg": "Failed setting up database",
		}).Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.WithFields(log.Fields{
			"dbname":     WEBAPP_DATABASE_URL,
			"custom_msg": "Unsucessfully connected with db",
		}).Error(err)
		return
	}

	return
}

func SetUpCors() (c *cors.Cors) {
	return cors.New(cors.Options{
		AllowedOrigins:   Origins,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})
}

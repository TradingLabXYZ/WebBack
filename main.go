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

var DbWebApp sqlx.DB

var Origins = []string{
	"http://localhost:9000",
	"https://tradinglab.xyz",
	"https://www.tradinglab.xyz",
	"https://staging.tradinglab.xyz",
	"https://hoofcoffee-wet0e0.stormkit.dev",
	"https://hoofcoffee-wet0e0--staging.stormkit.dev",
}

func main() {
	log_file := setUpLog()
	defer log_file.Close()
	DbWebApp = *setUpDb()
	defer DbWebApp.Close()

	go InstanciateTradesDispatcher()

	r := setupRoutes()
	c := setUpCors()
	h := c.Handler(r)
	fmt.Println(Bold(Green("Application is running on port 8080")))
	log.Fatal(http.ListenAndServe(":8080", h))
}

func setUpLog() (file *os.File) {
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

func setUpCors() (c *cors.Cors) {
	return cors.New(cors.Options{
		AllowedOrigins:   Origins,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})
}

func setUpAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SelectSession(r)
		if session.Id == 0 {
			log.Warn("Attempted open url without sessionid")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func setupRoutes() (router *mux.Router) {
	router = mux.NewRouter()

	// Define subrouter middlewares
	auth_router := router.PathPrefix("/").Subrouter()
	auth_router.Use(setUpAuthMiddleware)

	trades_router := router.PathPrefix("/get_trades/{username}/{requestid}").Subrouter()
	trades_router.Use(CheckUserPrivacy)

	// API endpoints
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")

	auth_router.HandleFunc("/user_settings", GetUserSettings).Methods("GET")
	auth_router.HandleFunc("/user_settings", UpdateUserSettings).Methods("POST")
	auth_router.HandleFunc("/update_password", UpdateUserPassword).Methods("POST")
	auth_router.HandleFunc("/update_privacy", UpdateUserPrivacy).Methods("POST")
	auth_router.HandleFunc("/insert_profile_picture", InsertProfilePicture).Methods("PUT")
	auth_router.HandleFunc("/user_premium_data", GetUserPremiumData).Methods("GET")

	trades_router.HandleFunc("", GetTrades)
	auth_router.HandleFunc("/insert_trade", InsertTrade).Methods("POST")
	auth_router.HandleFunc("/close_trade/{tradeid}", CloseTrade).Methods("GET")
	auth_router.HandleFunc("/open_trade/{tradeid}", OpenTrade).Methods("GET")
	auth_router.HandleFunc("/delete_trade/{tradeid}", DeleteTrade).Methods("GET")
	auth_router.HandleFunc("/update_trade", UpdateTrade).Methods("POST")

	router.HandleFunc("/get_prices/{usercode}", GetPrices)
	auth_router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")
	auth_router.HandleFunc("/stellar_price", SelectStellarPrice).Methods("GET")
	auth_router.HandleFunc("/transaction_credentials", SelectTransactionCredentials).Methods("GET")

	auth_router.HandleFunc("/buy_months", BuyPremiumMonths).Methods("POST")

	return
}

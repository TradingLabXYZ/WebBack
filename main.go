package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kz/discordrus"
	. "github.com/logrusorgru/aurora"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var (
	DbUrl      string
	Db         sqlx.DB
	trades_wss = make(map[string][]WsTrade)
)

func main() {
	r := SetupRoutes()
	c := SetUpCors()
	h := c.Handler(r)

	log_file := SetUpLog()
	defer log_file.Close()
	Db = *SetUpDb()
	defer Db.Close()

	// go TrackContractEvents() temporary paused
	// go ManageUnsubscriptions() temporary paused
	go InstanciateActivityMonitor()

	fmt.Println(Bold(Green("Application running on port 8080")))
	log.Fatal(http.ListenAndServe(":8080", h))
}

func SetupRoutes() (router *mux.Router) {
	router = mux.NewRouter()

	// Web
	router.HandleFunc("/login/{wallet}", Login).Methods("GET")
	router.HandleFunc("/get_trades/{wallet}/{sessionid}", StartTradesWs)
	router.HandleFunc("/get_explore/{offset}", SelectExplore).Methods("GET")
	router.HandleFunc("/user_settings", UpdateUserSettings).Methods("POST")
	router.HandleFunc("/update_privacy", UpdateUserPrivacy).Methods("POST")
	router.HandleFunc("/update_visibility", UpdateUserVisibility).Methods("POST")
	router.HandleFunc("/insert_profile_picture", InsertProfilePicture).Methods("PUT")
	router.HandleFunc("/follow/{wallet}/{status}", UpdateFollower).Methods("GET")
	router.HandleFunc("/subscribe/{wallet}/{status}", UpdateSubscriber).Methods("GET")
	router.HandleFunc("/get_connections/{wallet}", SelectConnections).Methods("GET")
	router.HandleFunc("/subscription/{wallet}", SelectSubscriptionMonthlyPrice).Methods("GET")
	router.HandleFunc("/admin/{token}", SelectActivity).Methods("GET")
	router.HandleFunc("/generate_api_token", GenerateApiToken).Methods("GET")

	// Web & API
	router.HandleFunc("/insert_trade", CreateTrade).Methods("POST")
	router.HandleFunc("/delete_trade/{tradecode}", DeleteTrade).Methods("GET")
	router.HandleFunc("/insert_subtrade/{tradecode}", CreateSubtrade).Methods("GET")
	router.HandleFunc("/update_subtrade", UpdateSubtrade).Methods("POST")
	router.HandleFunc("/delete_subtrade/{subtradecode}", DeleteSubtrade).Methods("GET")
	router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")
	router.HandleFunc("/get_pair_ratio/{firstPairCoinId}/{secondPairCoinId}", SelectPairRatio).Methods("GET")

	// API
	router.HandleFunc("/list_trades", ListTrades).Methods("GET")
	router.HandleFunc("/list_subtrades/{tradecode}", ListSubtrades).Methods("GET")
	router.HandleFunc("/get_snapshot", GetSnapshot).Methods("GET")
	/* TODOS Only APIs
	   /get_trade
	   /get_subtrades
		 /get_results*/

	files := http.FileServer(http.Dir("templates/public"))
	s := http.StripPrefix("/static/", files)
	router.PathPrefix("/static/").Handler(s)

	return
}

var Origins = []string{
	"http://127.0.0.1",
	"http://localhost:9000",
	"https://tradinglab.xyz",
	"https://www.tradinglab.xyz",
	"https://staging.tradinglab.xyz",
}

func SetUpCors() (c *cors.Cors) {
	return cors.New(cors.Options{
		AllowedOrigins:   Origins,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})
}

func SetUpDb() (db *sqlx.DB) {
	env := os.Getenv("TL_APP_ENV")
	var DB_NAME string
	if env == "production" {
		DB_NAME = "webappconnectionpool"
	} else if env == "staging" {
		DB_NAME = "stagingwebappconnectionpool"
	} else if env == "test" {
		DB_NAME = "testwebapp"
	}

	DbUrl = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("TL_DB_USER"),
		os.Getenv("TL_DB_PASS"),
		os.Getenv("TL_DB_HOST"),
		os.Getenv("TL_DB_PORT"),
		DB_NAME,
	)

	db, err := sqlx.Connect("postgres", DbUrl)
	if err != nil {
		log.WithFields(log.Fields{
			"dbname":     DB_NAME,
			"custom_msg": "Failed setting up database",
		}).Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.WithFields(log.Fields{
			"dbname":     DB_NAME,
			"custom_msg": "Unsucessfully connected with db",
		}).Error(err)
		return
	}

	return
}

func SetUpLog() (file *os.File) {
	env := os.Getenv("TL_APP_ENV")
	var DISCORD_WEBHOOK_URL string
	if env == "production" {
		DISCORD_WEBHOOK_URL = os.Getenv("DISCORD_WEBHOOK_URL")
	} else if env == "staging" {
		DISCORD_WEBHOOK_URL = os.Getenv("DISCORD_STAGING_WEBHOOK_URL")
	} else if env == "test" {
		DISCORD_WEBHOOK_URL = os.Getenv("DISCORD_STAGING_WEBHOOK_URL")
	}
	file, err := os.OpenFile(
		"logs.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0o666,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Failed setting up log file",
		}).Error(err)
		return
	}
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.AddHook(discordrus.NewHook(
		DISCORD_WEBHOOK_URL,
		log.TraceLevel,
		&discordrus.Opts{
			DisableTimestamp:   false,
			EnableCustomColors: true,
			CustomLevelColors: &discordrus.LevelColors{
				Trace: 3092790,
				Debug: 10170623,
				Info:  3581519,
				Warn:  14327864,
				Error: 13631488,
				Panic: 13631488,
				Fatal: 13631488,
			},
			DisableInlineFields: false,
		},
	))
	return
}

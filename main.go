package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	. "github.com/logrusorgru/aurora"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/discord"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var (
	DbUrl           string
	Db              sqlx.DB
	trades_wss      = make(map[string][]WsTrade)
	discordNotifier *notify.Notify
)

func main() {
	r := SetupRoutes()
	c := SetUpCors()
	h := c.Handler(r)

	log_file := SetUpLog()
	defer log_file.Close()
	Db = *setUpDb()
	defer Db.Close()
	discordNotifier = SetUpDiscordNotifier()

	go TrackContractEvents()
	go InstanciateActivityMonitor()
	go ManageUnsubscriptions()

	fmt.Println(Bold(Green("Application running on port 8080")))
	log.Fatal(http.ListenAndServe(":8080", h))
}

func SetupRoutes() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/login/{wallet}", Login).Methods("GET")
	router.HandleFunc("/get_trades/{wallet}/{sessionid}", StartTradesWs)
	router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")
	router.HandleFunc("/get_pair_ratio/{firstPairCoinId}/{secondPairCoinId}", SelectPairRatio).Methods("GET")
	router.HandleFunc("/get_explore/{offset}", SelectExplore).Methods("GET")

	router.HandleFunc("/insert_trade", CreateTrade).Methods("POST")
	router.HandleFunc("/delete_trade/{tradecode}", DeleteTrade).Methods("GET")
	router.HandleFunc("/update_subtrade", UpdateSubtrade).Methods("POST")
	router.HandleFunc("/insert_subtrade/{tradecode}", CreateSubtrade).Methods("GET")
	router.HandleFunc("/delete_subtrade/{subtradecode}", DeleteSubtrade).Methods("GET")

	router.HandleFunc("/user_settings", UpdateUserSettings).Methods("POST")
	router.HandleFunc("/update_privacy", UpdateUserPrivacy).Methods("POST")
	router.HandleFunc("/insert_profile_picture", InsertProfilePicture).Methods("PUT")

	router.HandleFunc("/admin/{token}", SelectActivity).Methods("GET")

	router.HandleFunc("/follow/{wallet}/{status}", UpdateFollower).Methods("GET")
	router.HandleFunc("/subscribe/{wallet}/{status}", UpdateSubscriber).Methods("GET")
	router.HandleFunc("/get_connections/{wallet}", SelectConnections).Methods("GET")

	router.HandleFunc("/subscription/{wallet}", SelectSubscriptionMonthlyPrice).Methods("GET")

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

func setUpDb() (db *sqlx.DB) {
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
	return
}

func SetUpDiscordNotifier() *notify.Notify {
	discordService := discord.New()
	_ = discordService.AuthenticateWithBotToken(os.Getenv("DISCORD_BOT_ID"))
	discordService.AddReceivers(os.Getenv("DISCORD_CHANNEL_LOG"))
	notifier := notify.New()
	notifier.UseServices(discordService)
	return notifier
}

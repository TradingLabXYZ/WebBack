package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
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
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://99a5eb64ecb041abb66d2809bcd4e101@o1054584.ingest.sentry.io/6040036",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")

	r := SetupRoutes()
	c := SetUpCors()
	h := c.Handler(r)

	log_file := SetUpLog()
	defer log_file.Close()
	Db = *setUpDb()
	defer Db.Close()

	go InstanciateActivityMonitor()

	fmt.Println(Bold(Green("Application is running on port 8080")))
	log.Fatal(http.ListenAndServe(":8080", h))
}

func SetupRoutes() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/login/{wallet}", Login).Methods("GET")
	router.HandleFunc("/get_trades/{username}/{requestid}", StartTradesWs)
	router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")

	router.HandleFunc("/insert_trade", CreateTrade).Methods("POST")
	router.HandleFunc("/change_trade/{tradecode}/{tostatus}", ChangeTradeStatus).Methods("GET")
	router.HandleFunc("/delete_trade/{tradecode}", DeleteTrade).Methods("GET")
	router.HandleFunc("/update_subtrade", UpdateSubtrade).Methods("POST")
	router.HandleFunc("/insert_subtrade/{tradecode}", CreateSubtrade).Methods("GET")
	router.HandleFunc("/delete_subtrade/{subtradecode}", DeleteSubtrade).Methods("GET")

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

var Origins = []string{
	"http://127.0.0.1",
	"http://localhost:9000",
	"https://tradinglab.xyz",
	"https://www.tradinglab.xyz",
	"https://staging.tradinglab.xyz",
	"https://hoofcoffee-wet0e0.stormkit.dev",
	"https://hoofcoffee-wet0e0--staging.stormkit.dev",
}

func SetUpCors() (c *cors.Cors) {
	return cors.New(cors.Options{
		AllowedOrigins:   Origins,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})
}

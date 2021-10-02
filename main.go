package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var DbWebApp = DbConnect()

func main() {

	f, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	defer DbWebApp.Close()

	router := mux.NewRouter()

	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")

	router.HandleFunc("/user_settings", GetUserSettings).Methods("GET")
	router.HandleFunc("/user_settings", UpdateUserSettings).Methods("POST")
	router.HandleFunc("/update_password", UpdateUserPassword).Methods("POST")
	router.HandleFunc("/update_privacy", UpdateUserPrivacy).Methods("POST")

	selectTradesRouter := router.PathPrefix("/select_trades/{username}").Subrouter()
	selectTradesRouter.Use(CheckUserPrivacy)
	selectTradesRouter.HandleFunc("", SelectTrades)

	router.HandleFunc("/insert_trade", InsertTrade).Methods("POST")
	router.HandleFunc("/close_trade/{tradeid}", CloseTrade).Methods("GET")
	router.HandleFunc("/open_trade/{tradeid}", OpenTrade).Methods("GET")
	router.HandleFunc("/delete_trade/{tradeid}", DeleteTrade).Methods("GET")
	router.HandleFunc("/update_trade", UpdateTrade).Methods("POST")

	router.HandleFunc("/get_prices/{username}", GetPrices)
	router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")

	router.HandleFunc("/insert_profile_picture", InsertProfilePicture).Methods("PUT")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:9000",
			"https://fishsunset-fmujji.stormkit.dev",
			"http://www.tradinglab.xyz",
			"https://tradinglab.xyz",
			"https://www.tradinglab.xyz",
			"http://test.tradinglab.xyz",
			"https://test.tradinglab.xyz",
		},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Info("Application is running on port 8080..")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

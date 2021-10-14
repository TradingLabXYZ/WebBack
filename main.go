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

var Origins = []string{
	"http://localhost:9000",
	"https://tradinglab.xyz",
	"https://www.tradinglab.xyz",
	"https://staging.tradinglab.xyz",
	"https://hoofcoffee-wet0e0.stormkit.dev",
	"https://hoofcoffee-wet0e0--staging.stormkit.dev",
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SelectSession(r)
		if session.Id == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

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

	// Define subrouter middlewares
	auth_router := router.PathPrefix("/").Subrouter()
	auth_router.Use(AuthMiddleware)

	trades_router := router.PathPrefix("/select_trades/{username}").Subrouter()
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

	trades_router.HandleFunc("", SelectTrades)
	auth_router.HandleFunc("/insert_trade", InsertTrade).Methods("POST")
	auth_router.HandleFunc("/close_trade/{tradeid}", CloseTrade).Methods("GET")
	auth_router.HandleFunc("/open_trade/{tradeid}", OpenTrade).Methods("GET")
	auth_router.HandleFunc("/delete_trade/{tradeid}", DeleteTrade).Methods("GET")
	auth_router.HandleFunc("/update_trade", UpdateTrade).Methods("POST")

	router.HandleFunc("/get_trades/{username}", GetTrades)

	router.HandleFunc("/get_prices/{usercode}", GetPrices)
	auth_router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")
	auth_router.HandleFunc("/stellar_price", SelectStellarPrice).Methods("GET")
	auth_router.HandleFunc("/transaction_credentials", SelectTransactionCredentials).Methods("GET")

	auth_router.HandleFunc("/buy_months", BuyPremiumMonths).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   Origins,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Info("Application is running on port 8080..")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

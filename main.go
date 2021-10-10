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

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SelectSession(r)
		if session.Id == 0 {
			fmt.Println("BAD SESSION")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println("GOOD SESSION", r.URL)
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

	router.HandleFunc("/user_settings", GetUserSettings).Methods("GET")
	router.HandleFunc("/user_settings", UpdateUserSettings).Methods("POST")
	router.HandleFunc("/update_password", UpdateUserPassword).Methods("POST")
	router.HandleFunc("/update_privacy", UpdateUserPrivacy).Methods("POST")
	router.HandleFunc("/insert_profile_picture", InsertProfilePicture).Methods("PUT")
	router.HandleFunc("/user_premium_data", GetUserPremiumData).Methods("GET")

	trades_router.HandleFunc("", SelectTrades)

	router.HandleFunc("/insert_trade", InsertTrade).Methods("POST")

	auth_router.HandleFunc("/close_trade/{tradeid}", CloseTrade).Methods("GET")

	router.HandleFunc("/open_trade/{tradeid}", OpenTrade).Methods("GET")
	router.HandleFunc("/delete_trade/{tradeid}", DeleteTrade).Methods("GET")
	router.HandleFunc("/update_trade", UpdateTrade).Methods("POST")

	router.HandleFunc("/get_prices/{usercode}", GetPrices)
	router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")
	router.HandleFunc("/stellar_price", SelectStellarPrice).Methods("GET")
	router.HandleFunc("/transaction_credentials", SelectTransactionCredentials).Methods("GET")

	router.HandleFunc("/buy_months", BuyPremiumMonths).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:9000",
			"https://tradinglab.xyz",
			"https://www.tradinglab.xyz",
			"https://staging.tradinglab.xyz",
			"https://hoofcoffee-wet0e0.stormkit.dev",
			"https://hoofcoffee-wet0e0--staging.stormkit.dev",
		},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Info("Application is running on port 8080..")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

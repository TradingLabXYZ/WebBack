package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	DbConnect()
	defer DbWebApp.Close()

	router := mux.NewRouter()
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")

	router.HandleFunc("/select_trades/{isopen}", SelectTrades).Methods("GET")
	router.HandleFunc("/insert_trade", InsertTrade).Methods("POST")
	router.HandleFunc("/close_trade/{tradeid}", CloseTrade).Methods("GET")
	router.HandleFunc("/open_trade/{tradeid}", OpenTrade).Methods("GET")
	router.HandleFunc("/delete_trade/{tradeid}", DeleteTrade).Methods("GET")
	router.HandleFunc("/update_trade", UpdateTrade).Methods("POST")

	router.HandleFunc("/get_price/{firstpair}/{secondpair}", SelectPrice).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:9000",
			"https://rabbitflicker-y6ho6u.stormkit.dev",
			"http://www.tradinglab.xyz",
			"https://tradinglab.xyz",
			"https://www.tradinglab.xyz",
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	fmt.Println("Application is running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}

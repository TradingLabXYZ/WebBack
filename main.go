package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var DbWebApp = DbConnect()

func main() {
	defer DbWebApp.Close()

	router := mux.NewRouter()
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")

	router.HandleFunc("/select_trades/{username}/{isopen}", SelectTrades).Methods("GET")
	router.HandleFunc("/insert_trade", InsertTrade).Methods("POST")
	router.HandleFunc("/close_trade/{tradeid}", CloseTrade).Methods("GET")
	router.HandleFunc("/open_trade/{tradeid}", OpenTrade).Methods("GET")
	router.HandleFunc("/delete_trade/{tradeid}", DeleteTrade).Methods("GET")
	router.HandleFunc("/update_trade", UpdateTrade).Methods("POST")

	router.HandleFunc("/get_price/{firstpairid}/{secondpairid}", SelectPrice).Methods("GET")
	router.HandleFunc("/get_pairs", SelectPairs).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:9000",
			"https://fishsunset-fmujji.stormkit.dev",
			"https://fishsunset-fmujji-44079058204449.stormkit.dev",
			"http://www.tradinglab.xyz",
			"https://tradinglab.xyz",
			"https://www.tradinglab.xyz",
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	fmt.Println("Application is running on port 8080..")

	log.Fatal(http.ListenAndServe(":8080", handler))
}

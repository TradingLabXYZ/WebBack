package main

import (
	"net/http"
)

type ActivitySnapshot struct {
	ActiveUsers int
	User        []string
}

func SelectActivity(w http.ResponseWriter, r *http.Request) {
	/* activeUsers := len(trades_wss)
	// fmt.Println(activeUsers)
	for userWatcher, _ := range trades_wss {
		for userToSee, _ := range trades_wss {
			fmt.Println("userWatcher:", userWatcher)
			fmt.Println("\tuserToSee:", userToSee)
		}
	}
	json.NewEncoder(w).Encode("OK") */
}

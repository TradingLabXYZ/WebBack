package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func SelectActivity(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	admin_token := os.Getenv("ADMIN_TOKEN")
	if token != admin_token {
		log.Warn("Attempted accessing admin with invalid token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output := make(map[string][]string)
	for userToSee, _ := range trades_wss {
		for _, q := range trades_wss[userToSee] {
			output[userToSee] = append(output[userToSee], q.SessionId)
		}
	}
	json.NewEncoder(w).Encode(output)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestLogin(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			wallet, username, privacy,
			twitter, profilepicture,
			plan, createdat, updatedat)
		VALUES
			('0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B', 'usera', 'all', 
			'testtwitter', 'testprofilepicture',
			'basic', current_timestamp, current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test login invalid eth wallet"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		vars := map[string]string{
			"wallet": "ABC",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed login invalid eth wallet")
		}
	})
	t.Run(fmt.Sprintf("Test login new wallet address"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		vars := map[string]string{
			"wallet": "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		bytes_body, _ := ioutil.ReadAll(w.Body)
		var body map[string]string
		_ = json.Unmarshal([]byte(bytes_body), &body)
		if res.StatusCode != 200 {
			t.Fatal("failed login new wallet address, response")
		}
		if body["SessionId"] == "" {
			t.Fatal("failed login new wallet, sessionid")
		}
		if body["ProfilePicture"] != "https://tradinglab.fra1.digitaloceanspaces.com/profile_pictures/default_picture.png" {
			t.Fatal("failed login new wallet, profile_picture")
		}
		if body["Wallet"] != "0x71C7656EC7ab88b098defB751B7401B5f6d8976F" {
			t.Fatal("failed login new wallet, wallet")
		}
	})
	t.Run(fmt.Sprintf("Test login existing user"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		vars := map[string]string{
			"wallet": "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
		}
		req = mux.SetURLVars(req, vars)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		bytes_body, _ := ioutil.ReadAll(w.Body)
		var body map[string]string
		_ = json.Unmarshal([]byte(bytes_body), &body)
		if res.StatusCode != 200 {
			t.Fatal("failed login existing user, response")
		}
		if body["SessionId"] == "" {
			t.Fatal("failed login existing user, sessionid")
		}
		if body["ProfilePicture"] != "testprofilepicture" {
			t.Fatal("failed login existing user, profile_picture")
		}
		if body["Wallet"] != "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B" {
			t.Fatal("failed login existing user, wallet")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

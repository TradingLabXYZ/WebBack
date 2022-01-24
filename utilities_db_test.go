package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSession(t *testing.T) {
	// <setup code>
	Db.Exec(`
		INSERT INTO users (
			wallet,
			username,
			privacy,
			plan,
			createdat,
			updatedat)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A',
			'r',
			'all',
			'basic',
			current_timestamp,
			current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test not existing user code"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		_, err := user.InsertSession()
		if err == nil {
			t.Fatal("Failed test not existing user code")
		}
	})

	t.Run(fmt.Sprintf("Test successfully creation of session"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession()
		if session.UserWallet == "" {
			t.Fatal("Failed successfully create session")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestExtractFromHeader(t *testing.T) {
	// <test code>
	t.Run(fmt.Sprintf("Test extracting session from wrong header"), func(t *testing.T) {
		session := Session{}
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		err := session.ExtractFromHeader(req)
		if err.Error() != "Could not find sessionId in header" {
			t.Fatal("Failed extracting session from wrong header")
		}
	})
	t.Run(fmt.Sprintf("Test extracting session missing header"), func(t *testing.T) {
		session := Session{}
		req := httptest.NewRequest("GET", "/", nil)
		err := session.ExtractFromHeader(req)
		if err.Error() != "Could not find authorization in header" {
			t.Fatal("Failed extracting session missing header")
		}
	})
}

func TestExtractFromCookie(t *testing.T) {
	// <test code>
	t.Run(fmt.Sprintf("Test wrong cookie"), func(t *testing.T) {
		session := Session{}
		req := httptest.NewRequest("GET", "/", nil)
		var cookie *http.Cookie = new(http.Cookie)
		cookie.Name = "Test="
		cookie.Value = "Test"
		req.AddCookie(cookie)
		err := session.ExtractFromCookie(req)
		if err.Error() != "Empty sessionId in cookie" {
			t.Fatal("Failed wrong cookie")
		}
	})
	t.Run(fmt.Sprintf("Test correct cookie"), func(t *testing.T) {
		session := Session{}
		req := httptest.NewRequest("GET", "/", nil)
		var cookie *http.Cookie = new(http.Cookie)
		cookie.Name = "sessionId="
		cookie.Value = "correct"
		req.AddCookie(cookie)
		_ = session.ExtractFromCookie(req)
		if session.Code != "=correct" {
			t.Fatal("Failed correct cookie")
		}
	})
}

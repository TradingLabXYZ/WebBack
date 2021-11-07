package main

import (
	"fmt"
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

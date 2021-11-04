package main

import (
	"fmt"
	"testing"
)

func TestCreateSession(t *testing.T) {
	// <setup code>
	Db.Exec(`
		INSERT INTO users (
			code,
			email,
			username,
			password,
			privacy,
			plan,
			createdat,
			updatedat)
		VALUES (
			'HAHAHAH',
			'r@r.r',
			'r',
			'rrrr',
			'all',
			'basic',
			current_timestamp,
			current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test not existing user code"), func(t *testing.T) {
		user := User{Code: "ABABABAB"}
		_, err := user.CreateSession()
		if err == nil {
			t.Fatal("Failed test not existing user code")
		}
	})

	t.Run(fmt.Sprintf("Test successfully creation of session"), func(t *testing.T) {
		user := User{Code: "HAHAHAH"}
		session, _ := user.CreateSession()
		if session.UserCode == "" {
			t.Fatal("Failed successfully create session")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM coins WHERE 1 = 1;`)
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

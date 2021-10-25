package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {

	// Test empty email
	user := User{
		Id: 99999999,
	}
	_, err := user.CreateSession()
	if err.Error() != "Empty email" {
		t.Error(err.Error())
	}

	// Test not valid id
	user = User{
		Id: 0,
	}
	_, err = user.CreateSession()
	if err == nil {
		t.Error(err)
	}

	// Test creation of session
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
			'XXXXX',
			'r@r.r',
			'r',
			'rrrr',
			'all',
			'basic',
			current_timestamp,
			current_timestamp);`)

	user = User{
		Id:    1,
		Email: "r@r.r",
	}
	session, err := user.CreateSession()
	assert.Greater(t, session.Id, 0)
}

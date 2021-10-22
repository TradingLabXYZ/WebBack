package main

import (
	"testing"
)

type SSUser struct {
	Id             int
	Code           string
	Email          string
	UserName       string
	LoginPassword  string
	Password       string
	Privacy        string
	Plan           string
	ProfilePicture string
	Twitter        string
	Website        string
}

func TestCreateSession(t *testing.T) {
	DbWebApp = *setUpDb()
	defer DbWebApp.Close()

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
	user = User{
		Id:    1,
		Email: "test@test.com",
	}
	session, err := user.CreateSession()
	if session.Id == 0 {
		t.Error(err)
	}
}

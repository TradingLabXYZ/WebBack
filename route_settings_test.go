package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func TestInsertProfilePicture(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			code, email, username, password, privacy,
			plan, createdat, updatedat)
		VALUES (
			'JFJFJF', 'jsjsjs@r.r', 'jsjsjsj', 'testpassword',
			'all', 'basic', current_timestamp, current_timestamp);`)

	user := User{Code: "JFJFJF"}
	session, _ := user.CreateSession()

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/insert_profile_picture", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		InsertProfilePicture(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test insert profile picture, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test process input file"), func(t *testing.T) {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		go func() {
			defer writer.Close()
			part, err := writer.CreateFormFile("THISNOTEXSIST", "someimg.png")
			if err != nil {
				t.Error(err)
			}
			upLeft := image.Point{0, 0}
			lowRight := image.Point{100, 100}
			img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}
		}()
		req := httptest.NewRequest("PUT", "/insert_profile_picture", pr)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		InsertProfilePicture(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test insert profile picture, wrong input file")
		}
	})
	// <tear-down code>
}

func TestGetUserSettings(t *testing.T) {
	// <setup code>
	// <test code>
	// <tear-down code>
}

func TestUpdateUserSettings(t *testing.T) {
	// <setup code>
	// <test code>
	// <tear-down code>
}

func TestUpdateUserPrivacy(t *testing.T) {
	// <setup code>
	// <test code>
	// <tear-down code>
}

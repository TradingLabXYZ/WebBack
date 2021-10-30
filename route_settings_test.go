package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

func TestInsertProfilePicture(t *testing.T) {
	// <setup code>
	do_key := os.Getenv("DO_KEY")
	do_secret := os.Getenv("DO_SECRET")
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(do_key, do_secret, ""),
		Endpoint:    aws.String("https://fra1.digitaloceanspaces.com"),
		Region:      aws.String("fra1"),
	}
	sess := aws_session.New(s3Config)
	svc := s3.New(sess)

	Db.Exec(
		`INSERT INTO users (
			code, email, username, password, privacy,
			plan, createdat, updatedat)
		VALUES (
			'testusertest', 'jsjsjs@r.r', 'jsjsjsj', 'testpassword',
			'all', 'basic', current_timestamp, current_timestamp);`)

	user := User{Code: "testusertest"}
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

	t.Run(fmt.Sprintf("Test file name wrong file form"), func(t *testing.T) {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		go func() {
			defer w.Close()
			_, _ = w.CreateFormFile("", "someimg.png")
		}()
		req := httptest.NewRequest("PUT", "/insert_profile_picture", &b)
		req.Header.Add("Content-Type", w.FormDataContentType())
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		x := httptest.NewRecorder()
		InsertProfilePicture(x, req)
		res := x.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test insert profile picture, invalid file form")
		}
	})

	t.Run(fmt.Sprintf("Test file name wrong extension"), func(t *testing.T) {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		go func() {
			defer w.Close()
			_, _ = w.CreateFormFile("file", "someimg.pdf")
		}()
		req := httptest.NewRequest("PUT", "/insert_profile_picture", &b)
		req.Header.Add("Content-Type", w.FormDataContentType())
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		x := httptest.NewRecorder()
		InsertProfilePicture(x, req)
		res := x.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test insert profile picture, wrong extension")
		}
	})

	t.Run(fmt.Sprintf("Test file name too many dots"), func(t *testing.T) {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		go func() {
			defer w.Close()
			_, _ = w.CreateFormFile("file", "someimg.pdf.png")
		}()
		req := httptest.NewRequest("PUT", "/insert_profile_picture", &b)
		req.Header.Add("Content-Type", w.FormDataContentType())
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		x := httptest.NewRecorder()
		InsertProfilePicture(x, req)
		res := x.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test insert profile picture, too many dots")
		}
	})

	t.Run(fmt.Sprintf("Test successfully upload file"), func(t *testing.T) {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		go func() {
			defer w.Close()
			_, _ = w.CreateFormFile("file", "someimg.png")
		}()
		req := httptest.NewRequest("PUT", "/insert_profile_picture", &b)
		req.Header.Add("Content-Type", w.FormDataContentType())
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		x := httptest.NewRecorder()
		InsertProfilePicture(x, req)
		resp, _ := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String("tradinglab"),
			Prefix: aws.String("profile_pictures/" + session.UserCode),
		})
		count_files := len(resp.Contents)
		if count_files != 1 {
			t.Fatal("Failed successfully upload file")
		}
	})

	t.Run(fmt.Sprintf("Test delete previous file"), func(t *testing.T) {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		go func() {
			defer w.Close()
			_, _ = w.CreateFormFile("file", "someimg.jpg")
		}()
		req := httptest.NewRequest("PUT", "/insert_profile_picture", &b)
		req.Header.Add("Content-Type", w.FormDataContentType())
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		x := httptest.NewRecorder()
		InsertProfilePicture(x, req)
		resp, _ := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String("tradinglab"),
			Prefix: aws.String("profile_pictures/" + session.UserCode),
		})
		count_files := len(resp.Contents)
		if count_files > 1 {
			t.Fatal("Failed deleting previous profile picture")
		}
	})

	t.Run(fmt.Sprintf("Test upload file path in database"), func(t *testing.T) {
		var file_path string
		_ = Db.QueryRow(`
			SELECT
				profilepicture
			FROM users
			WHERE code = $1`,
			session.UserCode).Scan(&file_path)
		fmt.Println(file_path)
		if file_path == "" {
			t.Fatal("Failed uploading file_path into db")
		}
	})

	// <tear-down code>
	resp, _ := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("tradinglab"),
		Prefix: aws.String("profile_pictures/" + session.UserCode),
	})
	for _, item := range resp.Contents {
		input := &s3.DeleteObjectInput{
			Bucket: aws.String("tradinglab"),
			Key:    aws.String(*item.Key),
		}
		_, _ = svc.DeleteObject(input)
	}
	Db.Exec(`DELETE FROM users WHERE 0 = 1;`)
}

func TestGetUserSettings(t *testing.T) {

	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			code, email, username, password, privacy,
			plan, twitter, website, createdat, updatedat)
		VALUES (
			'testusertest', 'jsjsjs@r.r', 'jsjsjsj', 'testpassword',
			'all', 'basic', 'thisistwitter', 'thisiswebsite',
			current_timestamp, current_timestamp);`)

	user := User{Code: "testusertest"}
	session, _ := user.CreateSession()

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user_settings", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		GetUserSettings(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test get users settings, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test invalid user"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user_settings", nil)
		req.Header.Set("Authorization", "Bearer sessionId=thisisaninvalid")
		w := httptest.NewRecorder()
		GetUserSettings(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test get users settings, invalid user")
		}
	})

	t.Run(fmt.Sprintf("Test successfully get user settings"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user_settings", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		GetUserSettings(w, req)
		res := w.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		if !strings.Contains(string(body), "thisistwitter") {
			t.Fatal("Failed successfully get user settings")
		}
		if !strings.Contains(string(body), "thisiswebsite") {
			t.Fatal("Failed successfully get user settings")
		}
		if !strings.Contains(string(body), "jsjsjs@r.r") {
			t.Fatal("Failed successfully get user settings")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestUpdateUserSettings(t *testing.T) {

	// <setup code>
	Db.Exec(
		`INSERT INTO users (
			code, email, username, password, privacy,
			plan, twitter, website, createdat, updatedat)
		VALUES (
			'testusertest', 'jsjsjs@r.r', 'jsjsjsj', 'testpassword',
			'all', 'basic', 'thisistwitter', 'thisiswebsite',
			current_timestamp, current_timestamp);`)

	user := User{Code: "testusertest"}
	session, _ := user.CreateSession()

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("POST", "/user_settings", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		UpdateUserSettings(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test update users settings, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test invalid user"), func(t *testing.T) {
		params := []byte(`{
			"Email": "new_email",
			"Twitter": "new_twitter",
			"Website": "new_website",
		}`)
		req := httptest.NewRequest("POST", "/user_settings", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId=thisisaninvalid")
		w := httptest.NewRecorder()
		UpdateUserSettings(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test update users settings, invalid user")
		}
	})

	t.Run(fmt.Sprintf("Test already present email"), func(t *testing.T) {
		Db.Exec(
			`INSERT INTO users (
					code, email, username, password, privacy,
					plan, twitter, website, createdat, updatedat)
				VALUES (
					'test2', 'emailalreadytaken', 'jsjsjsj2', 'testpassword',
					'all', 'basic', 'thisistwitter', 'thisiswebsite',
					current_timestamp, current_timestamp);`)

		params := []byte(`{
			"Email": "emailalreadytaken",
			"Twitter": "new_twitter",
			"Website": "new_website"
		}`)
		req := httptest.NewRequest("POST", "/user_settings", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateUserSettings(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test update users settings, email already presents")
		}
	})

	t.Run(fmt.Sprintf("Test already present website"), func(t *testing.T) {
		Db.Exec(
			`INSERT INTO users (
			code, email, username, password, privacy,
			plan, twitter, website, createdat, updatedat)
		VALUES (
			'test3', 'jdjadew', 'dqdjwq', 'testpassword',
			'all', 'basic', 'thisistwitter', 'websitealreadytaken',
			current_timestamp, current_timestamp);`)
		params := []byte(`{
			"Email": "ajsdhkad",
			"Twitter": "new_twitter",
			"Website": "websitealreadytaken"
		}`)
		req := httptest.NewRequest("POST", "/user_settings", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateUserSettings(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test update users settings, website already presents")
		}
	})

	t.Run(fmt.Sprintf("Test already present twitter"), func(t *testing.T) {
		Db.Exec(
			`INSERT INTO users (
			code, email, username, password, privacy,
			plan, twitter, website, createdat, updatedat)
		VALUES (
			'test4', 'wpskdhj', 'q2jdj', 'testpassword',
			'all', 'basic', 'twitteralreadytaken', 'thisiswebsite',
			current_timestamp, current_timestamp);`)
		params := []byte(`{
			"Email": "emailrandom",
			"Twitter": "twitteralreadytaken",
			"Website": "werbsiterandom"
		}`)
		req := httptest.NewRequest("POST", "/user_settings", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateUserSettings(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test update users settings, twitter already presents")
		}
	})

	t.Run(fmt.Sprintf("Test successfully update user settings"), func(t *testing.T) {
		params := []byte(`{
			"Email": "emailresult",
			"Twitter": "twitterresult",
			"Website": "websiteresult"
		}`)
		req := httptest.NewRequest("POST", "/user_settings", bytes.NewBuffer(params))
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		UpdateUserSettings(w, req)
		var email_result, twitter_result, website_result string
		_ = Db.QueryRow(`SELECT email, twitter, website FROM users WHERE code = $1;`, session.UserCode).Scan(&email_result, &twitter_result, &website_result)
		if email_result != "emailresult" {
			t.Fatal("Failed updating user settings, email")
		}
		if website_result != "websiteresult" {
			t.Fatal("Failed updating user settings, website")
		}
		if twitter_result != "twitterresult" {
			t.Fatal("Failed updating user settings, twitter")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestUpdateUserPassword(t *testing.T) {
	// <setup code>
	// <test code>
	// <tear-down code>
}

func TestUpdateUserPrivacy(t *testing.T) {
	// <setup code>
	// <test code>
	// <tear-down code>
}

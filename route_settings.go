package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/logrusorgru/aurora"
)

func InsertProfilePicture(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting InsertProfilePicture..."))

	session := SelectSession(r)
	if session.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user := UserByEmail(session.Email)

	// PROCESS INPUT FILE
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Warn("Error Retrieving the File")
	}
	defer file.Close()
	file_extension := strings.Split(handler.Header["Content-Type"][0], "/")[1]
	filename := user.Code + "." + file_extension
	filepath := "profile_pictures/" + filename

	// CONNECT AWS S3
	do_key := os.Getenv("DO_KEY")
	do_secret := os.Getenv("DO_SECRET")
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(do_key, do_secret, ""),
		Endpoint:    aws.String("https://fra1.digitaloceanspaces.com"),
		Region:      aws.String("fra1"),
	}
	sess := aws_session.New(s3Config)

	// DELETE OLD
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("tradinglab"),
		Prefix: aws.String("profile_pictures/" + user.Code),
	})
	if err != nil {
		log.Error(err)
	}
	for _, item := range resp.Contents {
		input := &s3.DeleteObjectInput{
			Bucket: aws.String("tradinglab"),
			Key:    aws.String(*item.Key),
		}
		_, err = svc.DeleteObject(input)
		if err != nil {
			log.Error(err)
		}
	}

	// UPLOAD NEW
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("tradinglab"),
		Key:    aws.String(filepath),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Error(err)
	}

	// SAVE PICTURE IN DB
	cdn_path := os.Getenv("CDN_PATH")
	file_cdn_path := cdn_path + "/" + filepath
	statement := `
		UPDATE users
		SET profilepicture = $1
		WHERE id = $2;`
	_, err = DbWebApp.Exec(statement, file_cdn_path, user.Id)
	if err != nil {
		log.Error(err)
	}

	w.Write([]byte(file_cdn_path))
}

func GetUserSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting GetUserSettings..."))

	session := SelectSession(r)
	if session.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := UserByEmail(session.Email)

	settings := struct {
		Email   string `json:"Email"`
		Twitter string `json:"Twitter"`
		Website string `json:"Website"`
		Privacy string `json:"Privacy"`
		Plan    string `json:"Plan"`
	}{
		user.Email,
		user.Twitter,
		user.Website,
		user.Privacy,
		user.Plan,
	}

	json.NewEncoder(w).Encode(settings)
}

func UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateUserSettings..."))

	/** TODOs
	- Check if Twitter URL is already taken
	- Check if Website is already taken
	- Check if email is already taken
	*/

	session := SelectSession(r)
	if session.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := UserByEmail(session.Email)

	settings := struct {
		Email   string `json:"Email"`
		Twitter string `json:"Twitter"`
		Website string `json:"Website"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&settings)
	if err != nil {
		log.Error(err)
	}

	statement := `
		UPDATE users
		SET email = $1,
		twitter = $2,
		website = $3
		WHERE id = $4;`
	_, err = DbWebApp.Exec(
		statement,
		settings.Email,
		settings.Twitter,
		settings.Website,
		user.Id)
	if err != nil {
		log.Error(err)
	}

	w.Write([]byte("OK"))
}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateUserPassword..."))

	session := SelectSession(r)
	if session.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := UserByEmail(session.Email)

	passwords := struct {
		OldPassword       string `json:"OldPassword"`
		NewPassword       string `json:"NewPassword"`
		RepeatNewPassword string `json:"RepeatNewPassword"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&passwords)
	if err != nil {
		log.Error(err)
	}

	if passwords.NewPassword != passwords.RepeatNewPassword {
		w.Write([]byte("KO"))
		return
	}

	if user.Password != Encrypt(passwords.OldPassword) {
		w.Write([]byte("KO"))
		return
	}

	statement := `
		UPDATE users
		SET password = $1
		WHERE id = $2;`
	_, err = DbWebApp.Exec(
		statement,
		Encrypt(passwords.NewPassword),
		user.Id)
	if err != nil {
		log.Error(err)
	}

	w.Write([]byte("OK"))
}

func UpdateUserPrivacy(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting UpdateUserPrivacy..."))

	session := SelectSession(r)
	if session.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := UserByEmail(session.Email)

	privacy := struct {
		Privacy string `json:"Privacy"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&privacy)
	if err != nil {
		log.Error(err)
	}

	statement := `
		UPDATE users
		SET privacy = $1
		WHERE id = $2;`
	_, err = DbWebApp.Exec(
		statement,
		privacy.Privacy,
		user.Id)
	if err != nil {
		log.Error(err)
	}

	w.Write([]byte("OK"))
}

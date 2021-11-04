package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

func InsertProfilePicture(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed inserting profile picture, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, handler, err := r.FormFile("file")
	if handler == nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed inserting profile picture, wrong file form",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file_name := handler.Filename
	if !strings.HasSuffix(file_name, "jpg") && !strings.HasSuffix(file_name, "png") {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"fileName":    file_name,
			"customMsg":   "Failed inserting profile picture, wrong file",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	file_extensions := strings.Split(file_name, ".")
	if len(file_extensions) != 2 {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"fileName":    file_name,
			"customMsg":   "Failed inserting profile picture, invalid file name",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file_path := "profile_pictures/" + session.UserCode + "." + file_extensions[1]

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
		Prefix: aws.String("profile_pictures/" + session.UserCode),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed inserting profile picture, error listing user pictures",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, item := range resp.Contents {
		input := &s3.DeleteObjectInput{
			Bucket: aws.String("tradinglab"),
			Key:    aws.String(*item.Key),
		}
		_, err = svc.DeleteObject(input)
		if err != nil {
			log.WithFields(log.Fields{
				"sessionCode": session.Code,
				"customMsg":   "Failed inserting profile picture, error deleting old",
			}).Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// UPLOAD NEW
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("tradinglab"),
		Key:    aws.String(file_path),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed inserting profile picture, error uploading new",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// SAVE PICTURE IN DB
	cdn_path := os.Getenv("CDN_PATH")
	file_cdn_path := cdn_path + "/" + file_path
	statement := `
		UPDATE users
		SET profilepicture = $1
		WHERE code = $2;`
	_, err = Db.Exec(statement, file_cdn_path, session.UserCode)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed inserting profile picture, error saving picture in db",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(file_cdn_path))
}

func GetUserSettings(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed getting settings, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := SelectUser("code", session.UserCode)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed getting settings, missing user",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
	/** TODOs
	- Check if Twitter URL is already taken
	- Check if Website is already taken
	- Check if email is already taken
	*/

	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed getting settings, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	settings := struct {
		Email   string `json:"Email"`
		Twitter string `json:"Twitter"`
		Website string `json:"Website"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&settings)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed getting settings, wrong payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var is_email_already_taken bool
	err = Db.QueryRow(`
		SELECT
			TRUE
		FROM users
		WHERE code != $1
		AND email = $2;`,
		session.UserCode,
		settings.Email).Scan(&is_email_already_taken)
	if is_email_already_taken {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"email":       settings.Email,
			"customMsg":   "Failed getting settings, email taken",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var is_twitter_already_taken bool
	err = Db.QueryRow(`
		SELECT
			TRUE
		FROM users
		WHERE code != $1
		AND twitter = $2;`,
		session.UserCode,
		settings.Twitter).Scan(&is_twitter_already_taken)
	if is_twitter_already_taken {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"twitter":     settings.Twitter,
			"customMsg":   "Failed getting settings, twitter taken",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var is_website_already_taken bool
	err = Db.QueryRow(`
		SELECT
			TRUE
		FROM users
		WHERE code != $1
		AND website = $2;`,
		session.UserCode,
		settings.Website).Scan(&is_website_already_taken)
	if is_website_already_taken {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"website":     settings.Website,
			"customMsg":   "Failed getting settings, website taken",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statement := `
		UPDATE users
		SET
			email = $1,
			twitter = $2,
			website = $3
		WHERE code = $4;`
	_, err = Db.Exec(
		statement,
		settings.Email,
		settings.Twitter,
		settings.Website,
		session.UserCode)
	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed changing password, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := SelectUser("code", session.UserCode)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed changing password, missing user",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	passwords := struct {
		OldPassword       string `json:"OldPassword"`
		NewPassword       string `json:"NewPassword"`
		RepeatNewPassword string `json:"RepeatNewPassword"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&passwords)
	if err != nil ||
		passwords.NewPassword == "" ||
		passwords.OldPassword == "" ||
		passwords.RepeatNewPassword == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed changing password, wrong payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if passwords.NewPassword != passwords.RepeatNewPassword {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encrypted_old_password, err := Encrypt(passwords.OldPassword)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed changing password, failed encrypting old password",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Password != encrypted_old_password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encrypted_new_password, err := Encrypt(passwords.NewPassword)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed changing password, failed encrypting new password",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statement := `
		UPDATE users
		SET password = $1
		WHERE code = $2;`
	_, err = Db.Exec(
		statement,
		encrypted_new_password,
		user.Code)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed changing password, failed sql",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateUserPrivacy(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed updating privacy, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	privacy := struct {
		Privacy string `json:"Privacy"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&privacy)
	if err != nil || privacy.Privacy == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed updating privacy, empty payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if privacy.Privacy != "all" &&
		privacy.Privacy != "private" &&
		privacy.Privacy != "subscribers" &&
		privacy.Privacy != "followers" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"payload":     privacy.Privacy,
			"customMsg":   "Failed updating privacy, wrong payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statement := `
		UPDATE users
		SET privacy = $1
		WHERE code = $2;`
	_, err = Db.Exec(
		statement,
		privacy.Privacy,
		session.UserCode)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed updating privacy, failed sql",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

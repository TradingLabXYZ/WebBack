package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed inserting profile picture, wrong file",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	content_types := handler.Header["Content-Type"]
	if content_types[0] == "" {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed inserting profile picture, bad content type",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	content_type := content_types[0]
	temp_file_extension := strings.Split(content_type, "/")
	if len(temp_file_extension) == 0 {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed inserting profile picture, bad file extension",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file_extension := temp_file_extension[1]

	filename := session.UserCode + "." + file_extension
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
		Key:    aws.String(filepath),
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
	file_cdn_path := cdn_path + "/" + filepath
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
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user, err := SelectUser("code", session.UserCode)
	if err != nil {
		log.Warn("User not found")
		w.WriteHeader(http.StatusNotFound)
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
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
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
		log.Error(err)
	}

	statement := `
		UPDATE users
		SET email = $1,
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

	w.Write([]byte("OK"))
}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {

	session, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user, err := SelectUser("code", session.UserCode)
	if err != nil {
		log.Warn("User not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	passwords := struct {
		OldPassword       string `json:"OldPassword"`
		NewPassword       string `json:"NewPassword"`
		RepeatNewPassword string `json:"RepeatNewPassword"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&passwords)
	if err != nil {
		log.Error(err)
	}

	if passwords.NewPassword != passwords.RepeatNewPassword {
		w.Write([]byte("KO"))
		return
	}

	encrypted_old_password, err := Encrypt(passwords.OldPassword)
	if err != nil {
		log.Error(err)
		w.Write([]byte("KO"))
		return
	}

	if user.Password != encrypted_old_password {
		w.Write([]byte("KO"))
		return
	}

	encrypted_new_password, err := Encrypt(passwords.NewPassword)
	if err != nil {
		log.Error(err)
		w.Write([]byte("KO"))
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
		log.Error(err)
	}

	w.Write([]byte("OK"))
}

func UpdateUserPrivacy(w http.ResponseWriter, r *http.Request) {

	session, err := GetSession(r, "header")
	if err != nil {
		log.Warn("User not log in")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	privacy := struct {
		Privacy string `json:"Privacy"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&privacy)
	if err != nil {
		log.Error(err)
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
		log.Error(err)
	}

	w.Write([]byte("OK"))
}

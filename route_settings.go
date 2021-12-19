package main

import (
	"encoding/json"
	"fmt"
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

	file_path := "profile_pictures/" + session.UserWallet + "." + file_extensions[1]

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
		Prefix: aws.String("profile_pictures/" + session.UserWallet),
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
		WHERE wallet = $2;`
	_, err = Db.Exec(statement, file_cdn_path, session.UserWallet)
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

func UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r, "header")
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed getting settings, wrong header",
		}).Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	settings := struct {
		Username string `json:"Username"`
		Twitter  string `json:"Twitter"`
		Discord  string `json:"Discord"`
		Github   string `json:"Github"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&settings)
	if err != nil {
		fmt.Println(err)
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed getting settings, wrong payload",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if settings.Username != "" {
		var is_username_already_taken bool
		err = Db.QueryRow(`
		SELECT
			TRUE
		FROM users
		WHERE wallet != $1
		AND username = $2;`,
			session.UserWallet,
			settings.Username).Scan(&is_username_already_taken)
		if is_username_already_taken {
			log.WithFields(log.Fields{
				"sessionCode": session.Code,
				"username":    settings.Username,
				"customMsg":   "Failed getting settings, username taken",
			}).Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if settings.Twitter != "" {
		var is_twitter_already_taken bool
		err = Db.QueryRow(`
		SELECT
			TRUE
		FROM users
		WHERE wallet != $1
		AND twitter = $2;`,
			session.UserWallet,
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
	}

	if settings.Discord != "" {
		var is_discord_already_taken bool
		err = Db.QueryRow(`
		SELECT
			TRUE
		FROM users
		WHERE wallet != $1
		AND discord = $2;`,
			session.UserWallet,
			settings.Discord).Scan(&is_discord_already_taken)
		if is_discord_already_taken {
			log.WithFields(log.Fields{
				"sessionCode": session.Code,
				"twitter":     settings.Discord,
				"customMsg":   "Failed getting settings, discord taken",
			}).Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if settings.Github != "" {
		var is_github_already_taken bool
		err = Db.QueryRow(`
		SELECT
			TRUE
		FROM users
		WHERE wallet != $1
		AND github = $2;`,
			session.UserWallet,
			settings.Github).Scan(&is_github_already_taken)
		if is_github_already_taken {
			log.WithFields(log.Fields{
				"sessionCode": session.Code,
				"github":      settings.Github,
				"customMsg":   "Failed getting settings, github taken",
			}).Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	statement := `
		UPDATE users
		SET
			username = $1,
			twitter = $2,
			discord = $3,
			github = $4
		WHERE wallet = $5;`
	_, err = Db.Exec(
		statement,
		settings.Username,
		settings.Twitter,
		settings.Discord,
		settings.Github,
		session.UserWallet)
	if err != nil {
		log.Error(err)
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
		WHERE wallet = $2;`
	_, err = Db.Exec(
		statement,
		privacy.Privacy,
		session.UserWallet)
	if err != nil {
		log.WithFields(log.Fields{
			"sessionCode": session.Code,
			"customMsg":   "Failed updating privacy, failed sql",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

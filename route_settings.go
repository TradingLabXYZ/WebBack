package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/logrusorgru/aurora"
)

func InsertProfilePicture(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Gray(8-1, "Starting InsertProfilePicture..."))

	user_session := SelectSession(r)
	user := UserByEmail(user_session.Email)

	// PROCESS INPUT FILE
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
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
	sess := session.New(s3Config)

	// DELETE OLD
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("tradinglab"),
		Prefix: aws.String("profile_pictures/" + user.Code),
	})
	if err != nil {
		panic(err)
	}
	for _, item := range resp.Contents {
		input := &s3.DeleteObjectInput{
			Bucket: aws.String("tradinglab"),
			Key:    aws.String(*item.Key),
		}
		_, err = svc.DeleteObject(input)
		if err != nil {
			panic(err)
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
		panic(err)
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
		panic(err)
	}

	w.Write([]byte(file_cdn_path))
}

package main

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	mathrand "math/rand"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func Encrypt(plaintext string) (cryptext string, err error) {
	if plaintext == "" {
		err = errors.New("Empty string")
		return
	}
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

func CreateUUID() (uuid string, err error) {
	u := new([16]byte)
	_, err = rand.Read(u[:])
	if err != nil {
		log.WithFields(log.Fields{
			"custom_msg": "Error during UUID creation",
		}).Error(err)
		err = errors.New("Error creating UUID")
		return
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mathrand.Intn(len(letterBytes))]
	}
	return string(b)
}

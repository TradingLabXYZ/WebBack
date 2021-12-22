package main

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	mathrand "math/rand"
	"sync"
	"time"

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
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

type Relation struct {
	Observer     User
	Observed     User
	Privacy      PrivacyStatus
	IsFollower   bool
	IsSubscriber bool
}

func (relation *Relation) CheckRelation() {
	follow_sql := func(wg *sync.WaitGroup) {
		defer wg.Done()
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE followfrom = $1
					AND followto = $2;`, relation.Observer.Wallet, relation.Observed.Wallet).Scan(
			&relation.IsFollower,
		)
	}
	subscribe_sql := func(wg *sync.WaitGroup) {
		defer wg.Done()
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE subscribefrom = $1
					AND subscribeto = $2;`, relation.Observer.Wallet, relation.Observed.Wallet).Scan(
			&relation.IsSubscriber,
		)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go follow_sql(&wg)
	go subscribe_sql(&wg)
	wg.Wait()
}

func (relation *Relation) CheckPrivacy() {
	if relation.Observer.Privacy == "all" {
		relation.Privacy.Status = "OK"
		relation.Privacy.Reason = "observed ALL"
		return
	}

	if relation.Observer.Wallet == "" {
		relation.Privacy.Status = "KO"
		relation.Privacy.Reason = "user is not logged in"
		return
	}

	if relation.Observer.Wallet == relation.Observed.Wallet {
		relation.Privacy.Status = "OK"
		relation.Privacy.Reason = "user access its own profile"
		return
	}

	switch relation.Observed.Privacy {
	case "private":
		relation.Privacy.Status = "KO"
		relation.Privacy.Reason = "private"
		return
	case "followers":
		if relation.IsFollower {
			relation.Privacy.Status = "OK"
			relation.Privacy.Reason = "user is follower"
			return
		} else {
			relation.Privacy.Status = "KO"
			relation.Privacy.Reason = "user is not follower"
			return
		}
	case "subscribers":
		if relation.IsSubscriber {
			relation.Privacy.Status = "OK"
			relation.Privacy.Reason = "user is subscriber"
			return
		} else {
			relation.Privacy.Status = "KO"
			relation.Privacy.Reason = "user is not subscriber"
			return
		}
	default:
		log.WithFields(log.Fields{
			"observed": relation.Observer.Wallet,
			"observer": relation.Observer.Wallet,
		}).Warn("Not possible to determine user's privacy")
		relation.Privacy.Status = "KO"
		relation.Privacy.Reason = "unknown reason"
		return
	}
}

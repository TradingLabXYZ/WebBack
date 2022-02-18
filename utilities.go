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

func (connection *Connection) CheckConnection() {
	follow_sql := func(wg *sync.WaitGroup) {
		defer wg.Done()
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM followers
					WHERE followfrom = $1
					AND followto = $2;`, connection.Observer.Wallet, connection.Observed.Wallet).Scan(
			&connection.IsFollower,
		)
	}
	subscribe_sql := func(wg *sync.WaitGroup) {
		defer wg.Done()
		_ = Db.QueryRow(`
					SELECT TRUE
					FROM subscribers
					WHERE subscribefrom = $1
					AND subscribeto = $2;`, connection.Observer.Wallet, connection.Observed.Wallet).Scan(
			&connection.IsSubscriber,
		)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go follow_sql(&wg)
	go subscribe_sql(&wg)
	wg.Wait()
}

func (connection *Connection) CheckPrivacy() {
	if connection.Observed.Privacy == "all" {
		connection.Privacy.Status = "OK"
		connection.Privacy.Reason = "observed ALL"
		return
	}

	if connection.Observer.Wallet == "" {
		connection.Privacy.Status = "KO"
		connection.Privacy.Reason = "user is not logged in"
		connection.Privacy.Message = "You need to login to visualise these infos!"
		return
	}

	if connection.Observer.Wallet == connection.Observed.Wallet {
		connection.Privacy.Status = "OK"
		connection.Privacy.Reason = "user access its own profile"
		return
	}

	switch connection.Observed.Privacy {
	case "private":
		connection.Privacy.Status = "KO"
		connection.Privacy.Reason = "private"
		connection.Privacy.Message = "This user prefers to keep things private!"
		return
	case "followers":
		if connection.IsFollower {
			connection.Privacy.Status = "OK"
			connection.Privacy.Reason = "user is follower"
			return
		} else {
			connection.Privacy.Status = "KO"
			connection.Privacy.Reason = "user is not follower"
			connection.Privacy.Message = "This user shares infos only with followers!"
			return
		}
	case "subscribers":
		if connection.IsSubscriber {
			connection.Privacy.Status = "OK"
			connection.Privacy.Reason = "user is subscriber"
			return
		} else {
			connection.Privacy.Status = "KO"
			connection.Privacy.Reason = "user is not subscriber"
			connection.Privacy.Message = "This user shares infos only with subscribers!"
			return
		}
	default:
		log.WithFields(log.Fields{
			"observed": connection.Observer.Wallet,
			"observer": connection.Observer.Wallet,
		}).Warn("Not possible to determine user's privacy")
		connection.Privacy.Status = "KO"
		connection.Privacy.Reason = "unknown reason"
		return
	}
}

func (observed *User) CheckVisibilities(snapshot TradesSnapshot) {
	var visibilities VisibilityStatus
	visibility_sql := `
		SELECT
			totalcounttrades,
			totalportfolio,
			totalreturn,
			totalroi,
			tradeqtyavailable,
			tradevalue,
			tradereturn,
			traderoi,
			subtradesall,
			subtradereasons,
			subtradequantity,
			subtradeavgprice,
			subtradetotal)
		FROM visibilities
		WHERE wallet = $1;`

	err := Db.QueryRow(
		visibility_sql,
		observed.Wallet).Scan(&visibilities)
	if err != nil {
		log.WithFields(log.Fields{
			"wallet":    observed.Wallet,
			"customMsg": "Failed extracting visibilities",
		}).Error(err)
		return
	}

	if !visibilities.TotalCountTrades {
		snapshot.CountTrades = 0
	}
	if !visibilities.TotalPortfolio {
		snapshot.TotalPortfolioUsd = "0"
	}
	if !visibilities.TotalReturn {
		snapshot.TotalReturnBtc = "0"
		snapshot.TotalReturnUsd = "0"
	}
	if !visibilities.TotalRoi {
		snapshot.Roi = 0
	}
	if !visibilities.TradeQtyAvailable {
		for _, trade := range snapshot.Trades {
			trade.QtyAvailable = "0"
		}
	}
	if !visibilities.TradeValue {
		for _, trade := range snapshot.Trades {
			trade.TotalValueUsd = 0
			trade.TotalValueUsdS = "0"
		}
	}
	if !visibilities.TradeReturn {
		for _, trade := range snapshot.Trades {
			trade.TotalReturn = 0
			trade.TotalReturnUsd = 0
			trade.TotalReturnBtc = 0
			trade.TotalReturnS = "0"
		}
	}
	if !visibilities.TradeRoi {
		for _, trade := range snapshot.Trades {
			trade.Roi = 0
		}
	}
	if !visibilities.SubtradesAll {
		for _, trade := range snapshot.Trades {
			trade.Subtrades = []Subtrade{}
		}
	}
	if !visibilities.SubtradeReasons {
		for _, trade := range snapshot.Trades {
			for _, subtrade := range trade.Subtrades {
				subtrade.Reason = ""
			}
		}
	}
	if !visibilities.SubtradeQuantity {
		for _, trade := range snapshot.Trades {
			for _, subtrade := range trade.Subtrades {
				subtrade.Quantity = 0
			}
		}
	}
	if !visibilities.SubtradeAvgPrice {
		for _, trade := range snapshot.Trades {
			for _, subtrade := range trade.Subtrades {
				subtrade.AvgPrice = 0
			}
		}
	}
	if !visibilities.SubtradeTotal {
		for _, trade := range snapshot.Trades {
			for _, subtrade := range trade.Subtrades {
				subtrade.Total = 0
			}
		}
	}
}

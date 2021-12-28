/** HOW TO
1 - Run Node
			truffle run moonbeam start
2 - Modify contract
3 - Deploy contract
			truffle migrate --network dev --reset
4 - Copy contract to FrontEnd
			cp build/contracts/Store.json $HOME/Code/TradingLab/WebFront/src/functions
5 - Create ABI
			truffle run abigen Store
6 - Update the contract in this file as well as the event params
7 - Run this file, and you should see logs when interacting with contract
*/

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func TrackContractTransaction() {
	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}
	go TrackSubscription(*client)
}

func TrackSubscription(client ethclient.Client) {
	subscriptionContractAddress := common.HexToAddress("0x42e2EE7Ba8975c473157634Ac2AF4098190fc741")
	subscriptionQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subscriptionContractAddress},
	}
	subscriptionLogs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), subscriptionQuery, subscriptionLogs)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed instanciating context contract ChangePlan",
		}).Error(err)
		return
	}
	subscriptionPath, _ := filepath.Abs("contracts/Subscription.abi")
	subscriptionFile, err := ioutil.ReadFile(subscriptionPath)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed reading ChangePlan abi file",
		}).Error(err)
		return
	}
	subscriptionAbi, err := abi.JSON(strings.NewReader(string(subscriptionFile)))
	if err != nil {
		fmt.Println("Invalid abi:", err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.WithFields(log.Fields{
				"customMsg": "Failed receiving vLog data",
			}).Warn(err)
		case vLog := <-subscriptionLogs:
			event := struct {
				Sender common.Address
				Value  *big.Int
			}{}
			err := subscriptionAbi.UnpackIntoInterface(&event, "ChangePlan", vLog.Data)
			if err != nil {
				log.WithFields(log.Fields{
					"vLog":      string(vLog.Data),
					"customMsg": "Failed unpacking vLog data",
				}).Warn(err)
			}
			value := event.Value.String()
			tx := vLog.TxHash.String()
			address := vLog.Address
			_, err = Db.Exec(`
				INSERT INTO changeplans (
					createdat,
					transaction,
					sender,
					value,
					contract)
				VALUES(current_timestamp, $1, $2, $3, $4);`,
				tx, event.Sender.String(), value, address.String())
			if err != nil {
				log.WithFields(log.Fields{
					"transaction": tx,
					"address":     address,
					"sender":      event.Sender.String(),
					"customMsg":   "Failed inserting ChangePlan into db",
				}).Warn(err)
			}
		}
	}
}

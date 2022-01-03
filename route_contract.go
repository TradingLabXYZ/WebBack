package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
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
	events_json, err := os.Open("contracts/subscription_events.json")
	defer events_json.Close()
	if err != nil {
		// MANAGE ERROR!!!
	}
	events_json_byte, err := ioutil.ReadAll(events_json)
	var events_data map[string]interface{}
	json.Unmarshal([]byte(events_json_byte), &events_data)
	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}
	go TrackSubscriptionContract(*client, events_data)
}

func TrackSubscriptionContract(client ethclient.Client, events_data map[string]interface{}) {
	subscriptionContractAddress := common.HexToAddress("0xfE5D3c52F7ee9aa32a69b96Bfbb088Ba0bCd8EfC")
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
			"customMsg": "Failed reading Subscription abi file",
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
			event_hex := vLog.Topics[0].String()
			event_map := events_data[event_hex].(map[string]interface{})
			event_name := ""
			for k, v := range event_map {
				if k == "name" {
					event_name = v.(string)
				}
			}

			fmt.Println(event_hex, event_name)

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

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

type SmartContract struct {
	Contract string `json:"contract"`
	Event    []struct {
		Signature string `json:"signature"`
		Name      string `json:"name"`
	} `json:"event"`
}

func TrackContractTransaction() {
	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}
	go TrackSubscriptionContract(*client)
}

func TrackSubscriptionContract(client ethclient.Client) {
	events_json, err := os.Open("contracts/subscription_info.json")
	events_json_byte, err := ioutil.ReadAll(events_json)
	var subscription_contract SmartContract
	json.Unmarshal([]byte(events_json_byte), &subscription_contract)
	defer events_json.Close()
	if err != nil {
		// MANAGE ERROR!!!
	}

	subscriptionContractAddress := common.HexToAddress(subscription_contract.Contract)
	subscriptionQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subscriptionContractAddress},
	}
	subscriptionLogs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(
		context.Background(),
		subscriptionQuery,
		subscriptionLogs,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed instanciating context contract ChangePlan",
		}).Error(err)
		return
	}
	subscriptionPath, _ := filepath.Abs("contracts/SubscriptionModel.abi")
	subscriptionFile, err := ioutil.ReadFile(subscriptionPath)
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed reading Subscription abi file",
		}).Error(err)
		return
	}
	subscriptionAbi, err := abi.JSON(
		strings.NewReader(string(subscriptionFile)),
	)
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
			fmt.Println("RECEIVED EVENT")
			event_signature := vLog.Topics[0].String()
			event_name := ""
			for _, v := range subscription_contract.Event {
				if v.Signature == event_signature {
					event_name = v.Name
				}
			}
			event_sender := ""
			event_payload := ""
			fmt.Println("EVENT NAME", event_name)
			switch {
			case event_name == "ChangePlan":
				event := struct {
					Sender common.Address
					Value  *big.Int
				}{}
				err := subscriptionAbi.UnpackIntoInterface(
					&event,
					"ChangePlan",
					vLog.Data,
				)
				if err != nil {
					log.WithFields(log.Fields{
						"event":     event_name,
						"customMsg": "Failed unpacking vLog data",
					}).Warn(err)
				}
				s_event, err := json.Marshal(event)
				if err != nil {
					fmt.Println(err)
					return
				}
				event_sender = event.Sender.Hex()
				event_payload = string(s_event)
			case event_name == "Subscribe":
				event := struct {
					Sender common.Address
					To     common.Address
					Weeks  *big.Int
					Amount *big.Int
				}{}
				err := subscriptionAbi.UnpackIntoInterface(
					&event,
					"Subscribe",
					vLog.Data,
				)
				if err != nil {
					log.WithFields(log.Fields{
						"event":     event_name,
						"customMsg": "Failed unpacking vLog data",
					}).Warn(err)
				}
				s_event, err := json.Marshal(event)
				if err != nil {
					fmt.Println(err)
					return
				}
				event_sender = event.Sender.Hex()
				event_payload = string(s_event)
			}

			tx := vLog.TxHash.String()
			contract_address := vLog.Address.String()
			_, err = Db.Exec(`
			INSERT INTO smartcontractevents (
				createdat,
				transaction,
				contract,
				name,
				signature,
				sender,
				payload)
			VALUES(current_timestamp, $1, $2, $3, $4, $5, $6);`,
				tx,
				contract_address,
				event_name,
				event_signature,
				event_sender,
				event_payload)
			if err != nil {
				log.WithFields(log.Fields{
					"transaction":     tx,
					"contractAddress": contract_address,
					"eventName":       event_name,
					"eventSignature":  event_signature,
					"eventSender":     event_sender,
					"eventPayload":    event_payload,
					"customMsg":       "Failed inserting smart contract event into db",
				}).Warn(err)
			}
		}
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func TrackContractEvents() {
	client, err := ethclient.Dial(os.Getenv("MOONBEAM_ENDPOINT"))
	if err != nil {
		time.Sleep(2 * time.Second)
		fmt.Println("Created a new dial after a fatal connection")
		go TrackContractEvents()
	}
	go TrackSubscriptionContract(client)
	go KeepRpcClientAlive(client)
}

func KeepRpcClientAlive(client *ethclient.Client) {
	for {
		block_number, err := client.BlockNumber(context.Background())
		if err != nil {
			client.Close()
			fmt.Println("Created a new dial connection after failed block extraction")
			go TrackContractEvents()
		}
		fmt.Println("Connection Rpc active:", block_number)
		time.Sleep(2 * time.Second)
	}
}

func TrackSubscriptionContract(client *ethclient.Client) {
	contract_address := os.Getenv("CONTRACT_SUBSCRIPTION")
	subscriptionAbi, err := abi.JSON(
		strings.NewReader(string(SubscriptionModelABI)),
	)
	if err != nil {
		log.Warn("Invalid abi:", err)
	}

	subscriptionContractAddress := common.HexToAddress(contract_address)
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

	for {
		select {
		case err := <-sub.Err():
			log.WithFields(log.Fields{
				"customMsg": "Failed receiving vLog data",
			}).Warn(err)
			break
		case vLog := <-subscriptionLogs:
			event_signature := vLog.Topics[0].String()
			event_hash := common.HexToHash(event_signature)
			event_details, err := subscriptionAbi.EventByID(event_hash)
			event := struct {
				Sender common.Address
				Value  *big.Int
			}{}
			err = subscriptionAbi.UnpackIntoInterface(
				&event,
				event_details.Name,
				vLog.Data,
			)
			if err != nil {
				log.WithFields(log.Fields{
					"event":     event_details.Name,
					"customMsg": "Failed unpacking vLog data",
				}).Warn(err)
			}

			event_sender := ""
			event_payload := ""
			switch {
			case event_details.Name == "ChangePlan":
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
						"event":     event_details.Name,
						"customMsg": "Failed unpacking vLog data",
					}).Warn(err)
				}
				s_event, err := json.Marshal(event)
				if err != nil {
					return
				}
				event_sender = event.Sender.Hex()
				event_payload = string(s_event)
			case event_details.Name == "Subscribe":
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
						"event":     event_details.Name,
						"customMsg": "Failed unpacking vLog data",
					}).Warn(err)
				}
				s_event, err := json.Marshal(event)
				if err != nil {
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
				event_details.Name,
				event_signature,
				event_sender,
				event_payload)
			if err != nil {
				log.WithFields(log.Fields{
					"transaction":     tx,
					"contractAddress": contract_address,
					"eventName":       event_details.Name,
					"eventSignature":  event_signature,
					"eventSender":     event_sender,
					"eventPayload":    event_payload,
					"customMsg":       "Failed inserting smart contract event into db",
				}).Warn(err)
			}
		}
	}
}

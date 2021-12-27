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
	"log"
	"math/big"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TrackContractTransaction() {
	go TrackStore()
	go TrackSubscription()

	// STORE CONTRACTS
}

func TrackStore() {
	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}

	storeContractAddress := common.HexToAddress("0xF8cef78E923919054037a1D03662bBD884fF4edf")
	storeQuery := ethereum.FilterQuery{
		Addresses: []common.Address{storeContractAddress},
	}
	storeLogs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), storeQuery, storeLogs)
	if err != nil {
		log.Fatal(err)
	}
	storePath, _ := filepath.Abs("Store.abi")
	storeFile, err := ioutil.ReadFile(storePath)
	if err != nil {
		fmt.Println("Failed to read file:", err)
	}
	storeAbi, err := abi.JSON(strings.NewReader(string(storeFile)))
	if err != nil {
		fmt.Println("Invalid abi:", err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-storeLogs:
			event := struct {
				Key   string
				Value string
			}{}
			err := storeAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Store Contract", string(event.Key))
			fmt.Println("Store Contract", string(event.Value))
		}
	}
}

func TrackSubscription() {
	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}
	subscriptionContractAddress := common.HexToAddress("0x42e2EE7Ba8975c473157634Ac2AF4098190fc741")
	subscriptionQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subscriptionContractAddress},
	}
	subscriptionLogs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), subscriptionQuery, subscriptionLogs)
	if err != nil {
		log.Fatal(err)
	}
	subscriptionPath, _ := filepath.Abs("Subscription.abi")
	subscriptionFile, err := ioutil.ReadFile(subscriptionPath)
	if err != nil {
		fmt.Println("Failed to read file:", err)
	}
	subscriptionAbi, err := abi.JSON(strings.NewReader(string(subscriptionFile)))
	if err != nil {
		fmt.Println("Invalid abi:", err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-subscriptionLogs:
			fmt.Println(vLog)
			event := struct {
				Sender common.Address
				Value  *big.Int
			}{}
			err := subscriptionAbi.UnpackIntoInterface(&event, "ChangePlan", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Subscription Contract", event.Sender)
			fmt.Println("Subscription Contract", event.Value)
		}
	}
}

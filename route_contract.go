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
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TrackContractTransaction() {
	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0xBF1fdCb2A1CAf1CA5662222417f0351043EEc19A")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	// FROM ME
	path, _ := filepath.Abs("/home/dolphin/Code/TradingLab/Contracts/abigenBindings/abi/Store.abi")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Failed to read file:", err)
	}
	edabi, err := abi.JSON(strings.NewReader(string(file)))
	if err != nil {
		fmt.Println("Invalid abi:", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			event := struct {
				Key   string
				Value string
			}{}
			err := edabi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(event.Key))   // foo
			fmt.Println(string(event.Value)) // barn("event", event)
		}
	}
}

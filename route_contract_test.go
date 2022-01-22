package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

func TestTT(t *testing.T) {
	// <setup code>
	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}
	/* events_json, err := os.Open("contracts/subscription_info.json")
	events_json_byte, err := ioutil.ReadAll(events_json)
	var subscription_contract SmartContract
	json.Unmarshal([]byte(events_json_byte), &subscription_contract)
	defer events_json.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"customMsg": "Failed loading subscription ABI",
		}).Error(err)
		return
	} */
	subscriptionContractAddress := common.HexToAddress("0x42e2EE7Ba8975c473157634Ac2AF4098190fc741")

	// instance, err := NewSubscriptionModel(subscriptionContractAddress, client)
	instance, err := NewPlansStorage(subscriptionContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("9d53abc69f2b6cb3ce693956433d3d64992a2d042323eb3249c8531d300e2413")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	_ = nonce

	session := &PlansStorageSession{
		Contract: instance,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592,
		},
	}

	ciao := common.HexToAddress("0xeF36DD9C9615447474aAE3D152e15188359D8e98")
	a := big.NewInt(10)
	session.AddPlan(ciao, a)

	// <test code>
	t.Run(fmt.Sprintf("Test BLOCKCHAIN"), func(t *testing.T) {
	})

	// <tear-down code>
}

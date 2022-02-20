package main

import (
	"testing"
)

func TestTrackContractEvents(t *testing.T) {
	// <setup code>
	/* baltathar_address := "0x3Cd0A705a2DC65e5b1E1205896BaA2be8A07c6e0"
	baltathar_private_key := "8075991ce870b93a8870eca0c0f91913d12f47948ca0fd25b49c6fa7cdbeee8b"
	go TrackContractEvents()

	client, err := ethclient.Dial("ws://127.0.0.1:9944")
	if err != nil {
		log.Fatal(err)
	}
	contract_address := os.Getenv("CONTRACT_SUBSCRIPTION")
	subscriptionContractAddress := common.HexToAddress(contract_address)

	instance, err := NewSubscriptionModel(subscriptionContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	testPrivateKey, err := crypto.HexToECDSA(baltathar_private_key)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := testPrivateKey.Public()
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

	auth := bind.NewKeyedTransactor(testPrivateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	_ = nonce

	session := &SubscriptionModelSession{
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

	// <test code>
	t.Run(fmt.Sprintf("Test smart contract ChangePlan"), func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		random_value := rand.Intn(10000-1) + 1
		plan_value := big.NewInt(int64(random_value))
		session.ChangePlan(plan_value)
		time.Sleep(3 * time.Second)
		var insert_plan_value int
		_ = Db.QueryRow(`
			SELECT DISTINCT
				payload#>>'{Value}' AS monthly_fee
			FROM smartcontractevents
			WHERE name = 'ChangePlan';`).Scan(&insert_plan_value)
		if random_value != insert_plan_value {
			t.Fatal("Failed smart contract ChangePlan")
		}
	})

	t.Run(fmt.Sprintf("Test smart contract Subscribe"), func(t *testing.T) {
		subscribe_to_address := common.HexToAddress(baltathar_address)
		weeks_value := 10
		big_weeks_value := big.NewInt(int64(weeks_value))
		session.Subscribe(subscribe_to_address, big_weeks_value)
		time.Sleep(3 * time.Second)
		var insert_address string
		var insert_weeks int
		_ = Db.QueryRow(`
			SELECT DISTINCT
				payload#>>'{To}',
				payload#>>'{Weeks}'
			FROM smartcontractevents
			WHERE name = 'Subscribe';`).Scan(&insert_address, &insert_weeks)
		insert_to_address_modified := common.HexToAddress(insert_address)
		if subscribe_to_address != insert_to_address_modified {
			t.Fatal("Failed smart contract Subscribe address")
		}
		if weeks_value != insert_weeks {
			t.Fatal("Failed smart contract Subscribe weeks")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM smartcontractevents WHERE 1 = 1;`) */
}

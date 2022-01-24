# Purpose
This repository contains the backend of TradingLab.

It is written in Golang using Postgresql as database.

The code is hosted on DigitalOceans.

There are two main branches: production and staging.

The code is automatically re-compiled at any changes, so it is possibile to simply refresh the webpage to run the new code.

## Architecture

The router is Gorilla Mux, using http and websocket.

To allow a dynamic experience to the users, the trading section creates a websocket which is called every time there is a change in the page, returning fresh processed data. In this way if multiple users are watching the same page, they will be all updated almost in real-time, almost at the same time.

In order to know when a user interacts with the platform, a specific Postgresql function is triggered, requiring to the server to activate a specific websocket.

Each new websocket is stored in the variable `trades_wss`:
```golang
type TradesSnapshot struct {
	UserDetails    UserDetails
	Trades         []Trade
	CountTrades    int
	TotalReturnUsd float64
	TotalReturnBtc float64
	Roi            float64
}

type WsTrade struct {
	UserToSee User
	RequestId string
	Channel   chan TradesSnapshot
	Ws        *websocket.Conn
}

trades_wss = make(map[string][]WsTrade)
```
In this way, every time `user_a` wants to see the profile of `user_b`:

1. send initial snapshot
2. instanciate websocket
3. add to `trades_wss`: `user_b` as key and `user_a` as value
4. if `user_b` makes a change, get `trades_wss[user_b]`, obtaining `user_a` 
5. if `user_a` closes the page the websocket is deleted from `trades_wss` 

# Run

Set environmental variables:
```bash
export TL_APP_ENV=
export TL_DB_USER=
export TL_DB_PASS=
export TL_DB_HOST=
export TL_DB_PORT=
export DO_KEY=
export DO_SECRET=
export CDN_PATH=
export ADMIN_TOKEN=
```

Build and run the program:
```bash
modd
```

# Test
```bash
go test -v -cover -parallel 1
```

# Migrate
Use `Makefile` to migrate up or down the database

# Run with smart contracts

1 - Run Node (truffle run moonbeam start)

2 - Modify contract

3 - Deploy contract (truffle migrate --network dev --reset)

4 - Copy contract to FrontEnd (cp build/contracts/Store.json $HOME/Code/TradingLab/WebFront/src/functions)

5 - Create ABI (truffle run abigen Store)

6 - Update the contract in this file as well as the event params

7 - Run this file, and you should see logs when interacting with contract

package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"
)

type NewSubtrade struct {
	CreatedAt  string      `json:"CreatedAt"`
	Type       string      `json:"Type"`
	Reason     string      `json:"Reason"`
	Quantity   json.Number `json:"Quantity"`
	AvgPrice   json.Number `json:"AvgPrice"`
	Total      json.Number `json:"Total"`
	UserWallet string
}

type NewTrade struct {
	Exchange     string        `json:"Exchange"`
	FirstPairId  int           `json:"FirstPair"`
	SecondPairId int           `json:"SecondPair"`
	Subtrades    []NewSubtrade `json:"Subtrades"`
	UserWallet   string
	Code         string
}

type WsTrade struct {
	Observer  User
	Observed  User
	SessionId string
	Channel   chan TradesSnapshot
	Ws        *websocket.Conn
}

type Session struct {
	Code       string
	UserWallet string
	CreatedAt  time.Time
}

type User struct {
	Wallet         string
	JoinTime       string
	Username       string
	Twitter        string
	Discord        string
	Github         string
	Privacy        string
	Plan           string
	ProfilePicture string
	Followers      int
	Followings     int
	Subscribers    int
	MonthlyFee     string
}

type Connection struct {
	Observer     User
	Observed     User
	Privacy      PrivacyStatus
	IsFollower   bool
	IsSubscriber bool
}

type OnlineUser struct {
	Wallet   string
	Observed []string
}

type OnlineUsers struct {
	Count int
	Users []OnlineUser
}

type DbListener struct {
	Listener *pq.Listener
}

func InstanciateActivityMonitor() {
	listener := DbListener{}
	listener.Instanciate()
	listener.Listen()
}

type UserWallet struct {
	Wallet string `validate:"eth_addr"`
}

type PairInfo struct {
	Symbol string
	Name   string
	Slug   string
}

type Follower struct {
	ProfilePicture string
	CountTrades    int
	Wallet         string
}

type Following struct {
	ProfilePicture string
	CountTrades    int
	Wallet         string
}

type Relations struct {
	Followers []Follower
	Following []Following
	Privacy   PrivacyStatus
}

type UserDetails struct {
	Username       string
	Twitter        string
	Github         string
	Discord        string
	Followers      int
	Followings     int
	Subscribers    int
	ProfilePicture string
	JoinTime       string
}

type Subtrade struct {
	Code      string
	TradeCode string
	CreatedAt string
	Type      string
	Reason    string
	Quantity  float64
	AvgPrice  float64
	Total     float64
}

type Trade struct {
	Code              string
	Username          string
	Userwallet        string
	Exchange          string
	FirstPairId       int
	SecondPairId      int
	FirstPairName     string
	SecondPairName    string
	FirstPairSymbol   string
	SecondPairSymbol  string
	FirstPairPrice    float64
	SecondPairPrice   float64
	FirstPairUrlIcon  string
	SecondPairUrlIcon string
	CurrentPrice      string
	QtyBuys           float64
	QtySells          float64
	QtyAvailable      string
	TotalBuys         float64
	TotalBuysBtc      float64
	TotalBuysUsd      float64
	TotalSells        float64
	TotalSellsBtc     float64
	TotalSellsUsd     float64
	ActualReturn      float64
	FutureReturn      float64
	FutureReturnBtc   float64
	FutureReturnUsd   float64
	TotalReturn       float64
	TotalReturnS      string
	TotalReturnBtc    float64
	TotalReturnUsd    float64
	TotalValueUsd     float64
	TotalValueUsdS    string
	Roi               float64
	BtcPrice          float64
	Subtrades         []Subtrade
}

type PrivacyStatus struct {
	Status  string
	Reason  string
	Message string
}

type TradesSnapshot struct {
	UserDetails       UserDetails
	PrivacyStatus     PrivacyStatus
	IsFollower        bool
	IsSubscriber      bool
	Trades            []Trade
	CountTrades       int
	TotalReturnUsd    string
	TotalReturnBtc    string
	Roi               float64
	TotalPortfolioUsd string
}

type SmartContract struct {
	Contract string `json:"contract"`
	Event    []struct {
		Signature string `json:"signature"`
		Name      string `json:"name"`
	} `json:"event"`
}

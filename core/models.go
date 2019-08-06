package core

import (
	"time"
)

const (
	OrdLimit        = "Limit"
	OrdMarket       = "Market"
	OrdStop         = "Stop"
	OrdStopLimit    = "StopLimit"
	LimitIfTouched  = "LimitIfTouched"
	MarketIfTouched = "MarketIfTouched"
)

const (
	SideBuy  = "Buy"
	SideSell = "Sell"
)

type Order struct {
	Symbol    string
	OrderID   string
	ClOrdID   string
	Price     float64
	Side      string
	Qty       float64
	OrdType   string
	OrdStatus string
}

type Trade struct {
	Symbol    string
	TID       string
	Side      string
	Price     float64
	Amount    float64
	Timestamp time.Time
}

type Quote struct {
	Symbol    string
	BidSize   float64
	BidPrice  float64
	AskSize   float64
	AskPrice  float64
	Timestamp time.Time
}

type Account struct {
}

type Tick struct {
	Symbol    string
	Last      float64
	Buy       float64
	Sell      float64
	High      float64
	Low       float64
	Vol       float64
	Timestamp time.Time
}

type Position struct {
}

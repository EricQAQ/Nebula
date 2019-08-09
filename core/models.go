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

	SideBuy  = "Buy"
	SideSell = "Sell"
)

type Order struct {
	Symbol       string
	OrderID      string
	ClOrdID      string
	OrdType      string
	Price        float64
	AvgPrice     float64
	Amount       float32
	FilledAmount float32
	Side         string
	OrdStatus    string
	Timestamp    time.Time
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

type Tick struct {
	Symbol    string
	Open      float64
	Close     float64
	High      float64 // 最高价
	Low       float64 // 最低价
	Vol       float64 // 量能
	Timestamp time.Time
}

type Position struct {
	Account        float64
	Symbol         string
	Currency       string
	LeverRate      float64 // 杠杆率
	ForceLiquPrice float64 //预估爆仓价

	SellAmount       float64 // 空单量
	SellAvailable    float64 // 可用空单量
	SellPriceAvg     float64 // 空单开仓均价
	SellPriceCost    float64 // 空单持仓金额
	SellProfitReal   float64 // 空单浮盈
	OpenOrderSellQty float64 // 委托空单平仓数量

	BuyAmount       float64 // 多单量
	BuyAvailable    float64 // 可用多单量
	BuyPriceAvg     float64 // 多单开仓均价
	BuyPriceCost    float64 // 多餐持仓金额
	BuyProfitReal   float64 // 多单浮盈
	OpenOrderBuyQty float64 // 委托多单平仓数量
}

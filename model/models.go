package model

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

type DepthRecord struct {
	Symbol string
	ID     float64
	Side   string
	Price  float64
	Amount float64
}

type Depth struct {
	// order: High -> Low
	Sell []*DepthRecord
	// order: low -> high
	Buy []*DepthRecord
}

type DepthRecordList []*DepthRecord

func (s DepthRecordList) Len() int {
	return len(s)
}

func (s DepthRecordList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DepthRecordList) Less(i, j int) bool {
	return s[i].Price < s[j].Price
}

type Order struct {
	Symbol       string
	OrderID      string
	ClOrdID      string
	OrdType      string
	Price        float64
	AvgPrice     float64
	Amount       float64
	FilledAmount float64
	Side         string
	OrdStatus    string
	Timestamp    time.Time
}

type Trade struct {
	Symbol    string    `csv:"symbol"`
	TID       string    `csv:"tid"`
	Side      string    `csv:"side"`
	Price     float64   `csv:"price"`
	Amount    float64   `csv:"amount"`
	Timestamp time.Time `csv:"timestamp"`
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

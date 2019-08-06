package bitmex

import (
	"github.com/EricQAQ/Traed/config"
	"github.com/EricQAQ/Traed/core"
)

const (
	dataLength = 2048
)

type Bitmex struct {
	APIKey    string
	APISecret string
	Proxy     string

	auth      *BitmexAuth
	subscribe *BitmexSubscribe

	tickKeys     map[string][]string
	tradeKeys    map[string][]string
	quoteKeys    map[string][]string
	orderKeys    map[string][]string
	positionKeys map[string][]string

	tickData     map[string][]*core.Tick
	tradeData    map[string][]*core.Trade
	quoteData    map[string][]*core.Quote
	orderData    map[string][]*core.Order
	positionData map[string][]*core.Position
}

func CreateBitmex(
	exchangeConfig *config.ExchangeConfig,
	httpConfig *config.HttpConfig) *Bitmex {
	bm := new(Bitmex)
	bm.APIKey = exchangeConfig.APIKey
	bm.APISecret = exchangeConfig.APISecret
	bm.Proxy = httpConfig.Proxy
	bm.auth = NewBitmexAuth(bm.APIKey, bm.APISecret, 24)
	bm.subscribe = NewBitmexSubscribe(exchangeConfig.Symbols, exchangeConfig.Topic...)

	bm.tickKeys = make(map[string][]string)
	bm.tradeKeys = make(map[string][]string)
	bm.quoteKeys = make(map[string][]string)
	bm.orderKeys = make(map[string][]string)
	bm.positionKeys = make(map[string][]string)

	bm.tickData = make(map[string][]*core.Tick)
	bm.tradeData = make(map[string][]*core.Trade)
	bm.quoteData = make(map[string][]*core.Quote)
	bm.orderData = make(map[string][]*core.Order)
	bm.positionData = make(map[string][]*core.Position)
	for _, symbol := range exchangeConfig.Symbols {
		bm.tickKeys[symbol] = make([]string, 0, 16)
		bm.tradeKeys[symbol] = make([]string, 0, 16)
		bm.quoteKeys[symbol] = make([]string, 0, 16)
		bm.orderKeys[symbol] = make([]string, 0, 16)
		bm.positionKeys[symbol] = make([]string, 0, 16)

		bm.tickData[symbol] = make([]*core.Tick, 0, dataLength)
		bm.tradeData[symbol] = make([]*core.Trade, 0, dataLength)
		bm.quoteData[symbol] = make([]*core.Quote, 0, dataLength)
		bm.orderData[symbol] = make([]*core.Order, 0, dataLength)
		bm.positionData[symbol] = make([]*core.Position, 0, dataLength)
	}
	return bm
}

func (bm *Bitmex) GetWsAuthHandler() core.WsAuthSubscribeHandler {
	return bm.auth
}

func (bm *Bitmex) GetWsSubscribeHandler() core.WsAuthSubscribeHandler {
	return bm.subscribe
}

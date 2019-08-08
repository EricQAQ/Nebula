package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/config"
	"github.com/EricQAQ/Traed/core"
	"github.com/orcaman/concurrent-map"
)

const (
	dataLength = 2048
)

type Bitmex struct {
	APIKey    string
	APISecret string
	Proxy     string
	timeout   time.Duration
	BaseUrl   string

	auth      *BitmexAuth
	subscribe *BitmexSubscribe

	tickKeys     map[string][]string
	tradeKeys    map[string][]string
	quoteKeys    map[string][]string
	orderKeys    map[string][]string
	positionKeys map[string][]string

	tickData     cmap.ConcurrentMap
	tradeData    cmap.ConcurrentMap
	quoteData    cmap.ConcurrentMap
	orderData    cmap.ConcurrentMap
	positionData cmap.ConcurrentMap
}

func CreateBitmex(
	exchangeConfig *config.ExchangeConfig,
	httpConfig *config.HttpConfig) *Bitmex {
	bm := new(Bitmex)
	bm.APIKey = exchangeConfig.APIKey
	bm.APISecret = exchangeConfig.APISecret
	bm.Proxy = httpConfig.Proxy
	bm.timeout = time.Duration(httpConfig.Timeout) * time.Millisecond
	bm.BaseUrl = exchangeConfig.HttpUrl + routeUrl
	bm.auth = NewBitmexAuth(bm.APIKey, bm.APISecret, 24)
	bm.subscribe = NewBitmexSubscribe(exchangeConfig.Symbols, exchangeConfig.Topic...)

	bm.tickKeys = make(map[string][]string)
	bm.tradeKeys = make(map[string][]string)
	bm.quoteKeys = make(map[string][]string)
	bm.orderKeys = make(map[string][]string)
	bm.positionKeys = make(map[string][]string)

	bm.tickData = cmap.New()
	bm.tradeData = cmap.New()
	bm.quoteData = cmap.New()
	bm.orderData = cmap.New()
	bm.positionData = cmap.New()
	for _, symbol := range exchangeConfig.Symbols {
		bm.tickKeys[symbol] = wsInstrumentKeys
		bm.tradeKeys[symbol] = wsTradeKeys
		bm.quoteKeys[symbol] = wsQuoteKeys
		bm.orderKeys[symbol] = wsOrderKeys
		bm.positionKeys[symbol] = wsPositionKeys

		bm.tickData.Set(symbol, make([]*core.Tick, 0, dataLength))
		bm.tradeData.Set(symbol, make([]*core.Trade, 0, dataLength))
		bm.quoteData.Set(symbol, make([]*core.Quote, 0, dataLength))
		bm.orderData.Set(symbol, make([]*core.Order, 0, dataLength))
		bm.positionData.Set(symbol, make([]*core.Position, 0, dataLength))
	}
	return bm
}

func (bm *Bitmex) GetWsAuthHandler() core.WsAuthSubscribeHandler {
	return bm.auth
}

func (bm *Bitmex) GetWsSubscribeHandler() core.WsAuthSubscribeHandler {
	return bm.subscribe
}

package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/config"
	"github.com/EricQAQ/Traed/core"
	"github.com/EricQAQ/Traed/core/kline"
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

	tickData     *kline.Ticker
	tradeData    *trade
	quoteData    *quote
	orderData    *order
	positionData *position
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

	bm.tickData = kline.NewTicker(exchangeConfig.Symbols)
	bm.tradeData = newTrade(exchangeConfig.Symbols)
	bm.quoteData = newQuote(exchangeConfig.Symbols)
	bm.orderData = newOrder(exchangeConfig.Symbols)
	bm.positionData = newPosition(exchangeConfig.Symbols)
	return bm
}

func (bm *Bitmex) GetWsAuthHandler() core.WsAuthSubscribeHandler {
	return bm.auth
}

func (bm *Bitmex) GetWsSubscribeHandler() core.WsAuthSubscribeHandler {
	return bm.subscribe
}

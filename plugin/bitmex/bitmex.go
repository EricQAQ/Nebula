package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/config"
	"github.com/EricQAQ/Traed/core"
	"github.com/EricQAQ/Traed/kline"
	"github.com/EricQAQ/Traed/model"
)

const (
	dataLength = 2048
)

type Bitmex struct {
	APIKey    string
	APISecret string
	proxy     string
	timeout   time.Duration
	retryCount int
	retryInterval time.Duration
	BaseUrl   string

	auth      *BitmexAuth
	subscribe *BitmexSubscribe

	tickData     *kline.Ticker
	tickCh       chan model.Tick
	tradeData    *trade
	quoteData    *quote
	orderData    *order
	depthData    *depth
	positionData *position
}

func CreateBitmex(
	exchangeConfig *config.ExchangeConfig,
	httpConfig *config.HttpConfig) *Bitmex {
	bm := new(Bitmex)
	bm.APIKey = exchangeConfig.APIKey
	bm.APISecret = exchangeConfig.APISecret
	bm.proxy = httpConfig.Proxy
	bm.retryCount = httpConfig.RetryCount
	bm.retryInterval = time.Duration(httpConfig.RetryInterval) * time.Millisecond
	bm.timeout = time.Duration(httpConfig.Timeout) * time.Millisecond
	bm.BaseUrl = exchangeConfig.HttpUrl + routeUrl
	bm.auth = NewBitmexAuth(bm.APIKey, bm.APISecret, 24)
	bm.subscribe = NewBitmexSubscribe(exchangeConfig.Symbols, exchangeConfig.Topic...)

	bm.tickData = kline.NewTicker(exchangeConfig.Symbols)
	bm.tradeData = newTrade(exchangeConfig.Symbols)
	bm.quoteData = newQuote(exchangeConfig.Symbols)
	bm.orderData = newOrder(exchangeConfig.Symbols)
	bm.positionData = newPosition(exchangeConfig.Symbols)
	bm.depthData = newDepth(exchangeConfig.Symbols)

	bm.tickCh = make(chan model.Tick, 1024)
	return bm
}

func (bm *Bitmex) GetWsAuthHandler() core.WsAuthSubscribeHandler {
	return bm.auth
}

func (bm *Bitmex) GetWsSubscribeHandler() core.WsAuthSubscribeHandler {
	return bm.subscribe
}

func (bm *Bitmex) GetExchangeName() string {
	return "bitmex"
}

func (bm *Bitmex) GetTickChan() chan model.Tick {
	return bm.tickCh
}

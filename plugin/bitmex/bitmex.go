package bitmex

import (
	"github.com/EricQAQ/Traed/config"
	"github.com/EricQAQ/Traed/core"
)

type Bitmex struct {
	APIKey          string
	APISecret       string
	Proxy           string

	auth *BitmexAuth
	subscribe *BitmexSubscribe
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
	return bm
}

func (bm *Bitmex) GetWsAuthHandler() core.WsAuthSubscribeHandler {
	return bm.auth
}

func (bm *Bitmex) GetWsSubscribeHandler() core.WsAuthSubscribeHandler {
	return bm.subscribe
}

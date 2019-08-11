package core

import (
	"time"

	"github.com/EricQAQ/Traed/model"
	"github.com/EricQAQ/Traed/kline"
)

const (
	WelcomeMsg = iota
	AuthMsg
	SubscribeMsg
	ErrorMsg
	Message
)

type ParsedData struct {
	Type int
	Data map[string]interface{}
}

type WsAuthSubscribeHandler interface {
	GetOperateArgs() []string
	Serialize() ([]byte, error)
}

type OrderAPI interface {
	LimitBuy(symbol string, price float64, quantity float32) (*model.Order, error)
	LimitSell(symbol string, price float64, quantity float32) (*model.Order, error)
	MarketBuy(symbol string, quantity float32) (*model.Order, error)
	MarketSell(symbol string, quantity float32) (*model.Order, error)
	LimitStopBuy(symbol string, stopPx float64, price float64, quantity float32) (*model.Order, error)
	LimitStopSell(symbol string, stopPx float64, price float64, quantity float32) (*model.Order, error)
	MarketStopBuy(symbol string, stopPx float64, quantity float32) (*model.Order, error)
	MarketStopSell(symbol string, stopPx float64, quantity float32) (*model.Order, error)
	LimitIfTouchedBuy(symbol string, stopPx float64, price float64, quantity float32) (*model.Order, error)
	LimitIfTouchedSell(symbol string, stopPx float64, price float64, quantity float32) (*model.Order, error)
}

type ExportAPI interface {
	GetTick(symbol string) (*model.Tick, bool)
	GetQuote(symbol string) (*model.Quote, bool)
	GetTrade(symbol string) (*model.Trade, bool)
	GetOrders(symbol string) ([]*model.Order, bool)
	GetPosition(symbol string) ([]*model.Position, bool)
}

type ExchangeAPI interface {
	GetExchangeName() string
	GetWsAuthHandler() WsAuthSubscribeHandler
	GetWsSubscribeHandler() WsAuthSubscribeHandler
	Parse(data []byte) (*ParsedData, error)
	HandleMessage(ParsedData)
	GetTickChan() chan model.Tick
	GetHistoryKline(symbol, period string, start, end time.Time) ([]*kline.Kline, error)

	OrderAPI
	ExportAPI
}

package bitmex

import (
	"sync/atomic"
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

type trade struct {
	tradeKeys map[string][]string
	tradeData cmap.ConcurrentMap
	isUpdate  int32
}

func newTrade(symbols []string) *trade {
	t := new(trade)
	t.tradeKeys = make(map[string][]string)
	t.tradeData = cmap.New()
	t.isUpdate = 0
	for _, symbol := range symbols {
		t.tradeKeys[symbol] = wsTradeKeys
		t.tradeData.Set(symbol, make([]*core.Trade, 0, dataLength))
	}
	return t
}

func (td *trade) getTradeList(symbol string) []*core.Trade {
	data, _ := td.tradeData.Get(symbol)
	tradeList := data.([]*core.Trade)
	return tradeList
}

func (td *trade) insertTrade(symbol string, trade *core.Trade) {
	tradeList := td.getTradeList(symbol)
	length := len(tradeList)
	if length >= dataLength {
		tradeList = tradeList[length-dataLength:]
	}
	tradeList = append(tradeList, trade)
	td.tradeData.Set(symbol, tradeList)
	atomic.StoreInt32(&td.isUpdate, 1)
}

func (td *trade) makeTrade(data map[string]interface{}) *core.Trade {
	trade := new(core.Trade)
	trade.Symbol = data["symbol"].(string)

	if side, ok := data["side"]; ok {
		trade.Side = side.(string)
	}
	if size, ok := data["size"]; ok {
		trade.Amount = size.(float64)
	}
	if price, ok := data["price"]; ok {
		trade.Price = price.(float64)
	}
	if trdMatchID, ok := data["trdMatchID"]; ok {
		trade.TID = trdMatchID.(string)
	}
	if ts, ok := data["timestamp"].(string); ok {
		trade.Timestamp, _ = time.Parse(time.RFC3339, ts)
	}
	return trade
}

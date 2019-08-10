package bitmex

import (
	"sync/atomic"
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/model"
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
		t.tradeData.Set(symbol, make([]*model.Trade, 0, dataLength))
	}
	return t
}

func (td *trade) getTradeList(symbol string) []*model.Trade {
	data, _ := td.tradeData.Get(symbol)
	tradeList := data.([]*model.Trade)
	return tradeList
}

func (td *trade) insertTrade(symbol string, trade *model.Trade) {
	tradeList := td.getTradeList(symbol)
	length := len(tradeList)
	if length >= dataLength {
		tradeList = tradeList[length-dataLength:]
	}
	tradeList = append(tradeList, trade)
	td.tradeData.Set(symbol, tradeList)
	atomic.StoreInt32(&td.isUpdate, 1)
}

func (td *trade) makeTrade(data map[string]interface{}) *model.Trade {
	trade := new(model.Trade)
	trade.Symbol = data["symbol"].(string)
	trade.Side = data["side"].(string)
	trade.Amount = data["size"].(float64)
	trade.Price = data["price"].(float64)
	trade.TID = data["trdMatchID"].(string)

	loc, _ := time.LoadLocation("Asia/Chongqing")
	ts, _ := time.Parse(time.RFC3339, data["timestamp"].(string))
	trade.Timestamp = ts.In(loc)
	return trade
}

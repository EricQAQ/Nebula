package bitmex

import (
	"time"
	"sync/atomic"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

type quote struct {
	quoteKeys map[string][]string
	quoteData cmap.ConcurrentMap
	isUpdate  int32
}

func newQuote(symbols []string) *quote {
	q := new(quote)
	q.quoteKeys = make(map[string][]string)
	q.quoteData = cmap.New()
	q.isUpdate = 0
	for _, symbol := range symbols {
		q.quoteKeys[symbol] = wsQuoteKeys
		q.quoteData.Set(symbol, make([]*core.Quote, 0, dataLength))
	}
	return q
}

func (qu *quote) getQuoteList(symbol string) []*core.Quote {
	data, _ := qu.quoteData.Get(symbol)
	quoteList := data.([]*core.Quote)
	return quoteList
}

func (qu *quote) insertQuote(symbol string, quote *core.Quote) {
	quoteList := qu.getQuoteList(symbol)
	length := len(quoteList)
	if length >= dataLength {
		quoteList = quoteList[length-dataLength:]
	}
	quoteList = append(quoteList, quote)
	qu.quoteData.Set(symbol, quoteList)
	atomic.StoreInt32(&qu.isUpdate, 1)
}

func (qu *quote) makeQuote(data map[string]interface{}) *core.Quote {
	quote := new(core.Quote)
	quote.Symbol = data["symbol"].(string)
	quote.BidSize = data["bidSize"].(float64)
	quote.BidPrice = data["bidPrice"].(float64)
	quote.AskPrice = data["askPrice"].(float64)
	quote.AskSize = data["askSize"].(float64)

	loc, _ := time.LoadLocation("Asia/Chongqing")
	ts, _ := time.Parse(time.RFC3339, data["timestamp"].(string))
	quote.Timestamp = ts.In(loc)
	return quote
}

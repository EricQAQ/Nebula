package bitmex

import (
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

type quote struct {
	quoteKeys     map[string][]string
	quoteData     cmap.ConcurrentMap
}

func newQuote(symbols []string) *quote {
	t := new(quote)
	t.quoteKeys = make(map[string][]string)
	t.quoteData = cmap.New()
	for _, symbol := range symbols {
		t.quoteKeys[symbol] = wsQuoteKeys
		t.quoteData.Set(symbol, make([]*core.Quote, 0, dataLength))
	}
	return t
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
}

func (qu *quote) makeQuote(data map[string]interface{}) *core.Quote {
	quote := new(core.Quote)
	quote.Symbol = data["symbol"].(string)

	if bidSize, ok := data["bidSize"]; ok {
		quote.BidSize = bidSize.(float64)
	}
	if bidPrice, ok := data["bidPrice"]; ok {
		quote.BidPrice = bidPrice.(float64)
	}
	if askPrice, ok := data["askPrice"]; ok {
		quote.AskPrice = askPrice.(float64)
	}
	if askSize, ok := data["askSize"]; ok {
		quote.AskSize = askSize.(float64)
	}
	if ts, ok := data["timestamp"].(string); ok {
		quote.Timestamp, _ = time.Parse(time.RFC3339, ts)
	}
	return quote
}

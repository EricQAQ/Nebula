package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) makeQuote(data map[string]interface{}) *core.Quote {
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

func (bm *Bitmex) insertQuote(symbol string, quote *core.Quote) {
	data, _ := bm.quoteData.Get(symbol)
	quoteList := data.([]*core.Quote)
	length := len(quoteList)
	if length >= dataLength {
		quoteList = quoteList[length-dataLength:]
	}
	quoteList = append(quoteList, quote)
	bm.quoteData.Set(symbol, quoteList)
}

package bitmex

import (
	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) makeQuote(data []map[string]interface{}) []*core.Quote {
	resp := make([]*core.Quote, 0, len(data))

	for _, item := range data {
		quote := new(core.Quote)
		quote.Symbol = item["symbol"].(string)

		if bidSize, ok := item["bidSize"]; ok {
			quote.BidSize = bidSize.(float64)
		}
		if bidPrice, ok := item["bidPrice"]; ok {
			quote.BidPrice = bidPrice.(float64)
		}
		if askPrice, ok := item["askPrice"]; ok {
			quote.AskPrice = askPrice.(float64)
		}
		if askSize, ok := item["askSize"]; ok {
			quote.AskSize = askSize.(float64)
		}
		if ts, ok := item["timestamp"].(string); ok {
			quote.Timestamp = time.Parse(time.RFC3339, ts)
		}
		resp = append(resp, quote)
	}
	return resp
}

func (bm *Bitmex) insertQuoteList(symbol string, quoteList []*core.Quote) {
	updateLength := len(quoteList)
	length = len(bm.quoteData[symbol])
	if length+updateLength >= dataLength {
		bm.quoteData[symbol] = bm.quoteData[symbol][length+updateLength-dataLength:]
	}
	bm.quoteData[symbol] = append(bm.quoteData[symbol], quoteList...)
}

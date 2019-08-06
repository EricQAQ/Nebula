package bitmex

import (
	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) makeTrade(data []map[string]interface{}) []*core.Trade {
	resp := make([]*core.Trade, 0, len(data))

	for _, item := range data {
		trade := new(core.Trade)
		trade.Symbol = item["symbol"].(string)

		if side, ok := item["side"]; ok {
			trade.Side = side.(string)
		}
		if size, ok := item["size"]; ok {
			trade.Amount = size.(float64)
		}
		if price, ok := item["price"]; ok {
			trade.Price = price.(float64)
		}
		if trdMatchID, ok := item["trdMatchID"]; ok {
			trade.TID = trdMatchID.(float64)
		}
		if ts, ok := item["timestamp"].(string); ok {
			trade.Timestamp = time.Parse(time.RFC3339, ts)
		}
		resp = append(resp, trade)
	}
	return resp
}

func (bm *Bitmex) insertTradeList(symbol string, tradeList []*core.Trade) {
	updateLength := len(tradeList)
	length = len(bm.tradeData[symbol])
	if length+updateLength >= dataLength {
		bm.tradeData[symbol] = bm.tradeData[symbol][length+updateLength-dataLength:]
	}
	bm.tradeData[symbol] = append(bm.tradeData[symbol], tradeList...)
}

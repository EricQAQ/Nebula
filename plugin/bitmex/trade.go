package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) makeTrade(data map[string]interface{}) *core.Trade {
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

func (bm *Bitmex) insertTrade(symbol string, trade *core.Trade) {
	length := len(bm.tradeData[symbol])
	if length >= dataLength {
		bm.tradeData[symbol] = bm.tradeData[symbol][length-dataLength:]
	}
	bm.tradeData[symbol] = append(bm.tradeData[symbol], trade)
}

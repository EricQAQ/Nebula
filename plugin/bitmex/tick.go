package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) makeInstrument(data []map[string]interface{}) []*core.Tick {
	resp := make([]*core.Tick, 0, len(data))
	for _, item := range data {
		tick := new(core.Tick)
		tick.Symbol = item["symbol"].(string)

		if last, ok := item["lastPrice"]; ok {
			tick.Last = last.(float64)
		}
		if high, ok := item["highPrice"]; ok {
			tick.High = high.(float64)
		}
		if low, ok := item["lowPrice"]; ok {
			tick.Low = low.(float64)
		}
		if buy, ok := item["bidPrice"]; ok {
			tick.Buy = buy.(float64)
		}
		if sell, ok := item["askPrice"]; ok {
			tick.Sell = sell.(float64)
		}
		if vol, ok := item["homeNotional24h"]; ok {
			tick.Vol = vol.(float64)
		}
		if ts, ok := item["timestamp"].(string); ok {
			tick.Timestamp, _ = time.Parse(time.RFC3339, ts)
		}
		resp = append(resp, tick)
	}

	return resp
}

func (bm *Bitmex) insertTickList(symbol string, tickList []*core.Tick) {
	updateLength := len(tickList)
	length := len(bm.tickData[symbol])
	if length+updateLength >= dataLength {
		bm.tickData[symbol] = bm.tickData[symbol][length+updateLength-dataLength:]
	}
	bm.tickData[symbol] = append(bm.tickData[symbol], tickList...)
}

func (bm *Bitmex) updateTick(tick *core.Tick, data map[string]interface{}) {
	for name, value := range data {
		if name == "lastPrice" {
			tick.Last = value.(float64)
		} else if name == "highPrice" {
			tick.High = value.(float64)
		} else if name == "lowPrice" {
			tick.Low = value.(float64)
		} else if name == "bidPrice" {
			tick.Buy = value.(float64)
		} else if name == "askPrice" {
			tick.Sell = value.(float64)
		} else if name == "homeNotional24h" {
			tick.Vol = value.(float64)
		} else if name == "timestamp" {
			tick.Timestamp, _ = time.Parse(time.RFC3339, value.(string))
		}
	}
}

func (bm *Bitmex) findTickItemByKeys(
	symbol string, updateData map[string]interface{}) (int, *core.Tick) {
	for index, val := range bm.tickData[symbol] {
		if val.Symbol == updateData["symbol"].(string) {
			return index, val
		}
	}
	return 0, nil
}

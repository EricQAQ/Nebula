package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) makeInstrument(data map[string]interface{}) *core.Tick {
	tick := new(core.Tick)
	tick.Symbol = data["symbol"].(string)

	if last, ok := data["lastPrice"]; ok {
		tick.Last = last.(float64)
	}
	if high, ok := data["highPrice"]; ok {
		tick.High = high.(float64)
	}
	if low, ok := data["lowPrice"]; ok {
		tick.Low = low.(float64)
	}
	if buy, ok := data["bidPrice"]; ok {
		tick.Buy = buy.(float64)
	}
	if sell, ok := data["askPrice"]; ok {
		tick.Sell = sell.(float64)
	}
	if vol, ok := data["homeNotional24h"]; ok {
		tick.Vol = vol.(float64)
	}
	if ts, ok := data["timestamp"].(string); ok {
		tick.Timestamp, _ = time.Parse(time.RFC3339, ts)
	}
	return tick
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

func (bm *Bitmex) insertTick(symbol string, tick *core.Tick) {
	data, _ := bm.tickData.Get(symbol)
	tickList := data.([]*core.Tick)
	length := len(tickList)
	if length >= dataLength {
		tickList = tickList[length-dataLength:]
	}
	tickList = append(tickList, tick)
	bm.tickData.Set(symbol, tickList)
}

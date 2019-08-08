package bitmex

import (
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

type ticker struct {
	tickKeys     map[string][]string
	tickData     cmap.ConcurrentMap
}

func newTicker(symbols []string) *ticker {
	t := new(ticker)
	t.tickKeys = make(map[string][]string)
	t.tickData = cmap.New()
	for _, symbol := range symbols {
		t.tickKeys[symbol] = wsInstrumentKeys
		t.tickData.Set(symbol, make([]*core.Tick, 0, dataLength))
	}
	return t
}

func (t *ticker) getTickList(symbol string) []*core.Tick {
	data, _ := t.tickData.Get(symbol)
	tickList := data.([]*core.Tick)
	return tickList
}

func (t *ticker) insertTick(symbol string, tick *core.Tick) {
	tickList := t.getTickList(symbol)
	length := len(tickList)
	if length >= dataLength {
		tickList = tickList[length-dataLength:]
	}
	tickList = append(tickList, tick)
	t.tickData.Set(symbol, tickList)
}

func (t *ticker) updateTick(symbol string, data map[string]interface{}) {
	tickList := t.getTickList(symbol)
	length := len(tickList)
	if length <= 0 {
		return
	}

	for name, value := range data {
		if name == "lastPrice" {
			tickList[length-1].Last = value.(float64)
		} else if name == "highPrice" {
			tickList[length-1].High = value.(float64)
		} else if name == "lowPrice" {
			tickList[length-1].Low = value.(float64)
		} else if name == "bidPrice" {
			tickList[length-1].Buy = value.(float64)
		} else if name == "askPrice" {
			tickList[length-1].Sell = value.(float64)
		} else if name == "homeNotional24h" {
			tickList[length-1].Vol = value.(float64)
		} else if name == "timestamp" {
			tickList[length-1].Timestamp, _ = time.Parse(time.RFC3339, value.(string))
		}
	}
}

func (t *ticker) deleteLastTick(symbol string) {
	tickList := t.getTickList(symbol)
	length := len(tickList)
	if length <= 0 {
		return
	}
	tickList = tickList[:length-2]
	t.tickData.Set(symbol, tickList)
}

func (t *ticker) makeInstrument(data map[string]interface{}) *core.Tick {
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

package kline

import (
	"sync/atomic"
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

const (
	tickInterval = 5 // 5ms
	tickLength = 4096
)

type Ticker struct {
	tickData cmap.ConcurrentMap
	isUpdate int32
	interval time.Duration
}

func NewTicker(symbols []string) *Ticker {
	t := new(Ticker)
	t.tickData = cmap.New()
	t.isUpdate = 0
	t.interval = time.Duration(tickInterval) * time.Millisecond
	for _, symbol := range symbols {
		t.tickData.Set(symbol, make([]*core.Tick, 0, tickLength))
	}
	return t
}

func (t *Ticker) GetTickerList(symbol string) []*core.Tick {
	data, _ := t.tickData.Get(symbol)
	tickList := data.([]*core.Tick)
	return tickList
}

func (t *Ticker) makeTick(trade *core.Trade) *core.Tick {
	tick := new(core.Tick)
	tick.Symbol = trade.Symbol
	tick.Open = trade.Price
	tick.Close = trade.Price
	tick.High = trade.Price
	tick.Low = trade.Price
	tick.Vol = trade.Amount
	tick.Timestamp = trade.Timestamp
	return tick
}

func (t *Ticker) appendTick(symbol string, tick *core.Tick) {
	tickList := t.GetTickerList(symbol)
	if len(tickList) >= tickLength {
		tickList = tickList[1:]
	}
	tickList = append(tickList, tick)
	t.tickData.Set(symbol, tickList)
}

func (t *Ticker) updateTick(tick, newTick *core.Tick) {
	tick.Close = newTick.Close
	tick.High = maxFloat(tick.High, newTick.High)
	tick.Low = minFloat(tick.Low, newTick.Low)
	tick.Vol += newTick.Vol
	tick.Timestamp = newTick.Timestamp
}

func (t *Ticker) UpdateTicker(symbol string, trade *core.Trade) {
	tick := t.makeTick(trade)
	tickList := t.GetTickerList(symbol)

	length := len(tickList)
	if length == 0 {
		t.appendTick(symbol, tick)
		return
	}
	lastTicker := tickList[length-1]
	if tick.Timestamp.Sub(lastTicker.Timestamp) > t.interval {
		t.appendTick(symbol, tick)
	} else {
		t.updateTick(lastTicker, tick)
	}
	t.SetUpdateFlag()
}

func (t *Ticker) SetUpdateFlag() {
	atomic.StoreInt32(&t.isUpdate, 1)
}

func (t *Ticker) IsUpdate() bool {
	return atomic.CompareAndSwapInt32(&t.isUpdate, 1, 0)
}

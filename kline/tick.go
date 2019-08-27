package kline

import (
	"sync/atomic"
	"time"

	"github.com/EricQAQ/Nebula/model"

	"github.com/orcaman/concurrent-map"
)

const (
	tickInterval = 500 // 500ms
	tickLength   = 4096
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
		t.tickData.Set(symbol, make([]*model.Tick, 0, tickLength))
	}
	return t
}

func (t *Ticker) GetTickerList(symbol string) []*model.Tick {
	data, _ := t.tickData.Get(symbol)
	tickList := data.([]*model.Tick)
	return tickList
}

func (t *Ticker) makeTick(trade *model.Trade) *model.Tick {
	tick := new(model.Tick)
	tick.Symbol = trade.Symbol
	tick.Open = trade.Price
	tick.Close = trade.Price
	tick.High = trade.Price
	tick.Low = trade.Price
	tick.Vol = trade.Amount
	tick.Timestamp = trade.Timestamp.Truncate(tickInterval * time.Millisecond)
	return tick
}

func (t *Ticker) appendTick(symbol string, tick *model.Tick, tCh chan model.Tick) {
	tickList := t.GetTickerList(symbol)
	if len(tickList) >= tickLength {
		tickList = tickList[1:]
	}
	tickList = append(tickList, tick)
	t.tickData.Set(symbol, tickList)
	tCh <- *tick
	t.SetUpdateFlag()
}

func (t *Ticker) updateTick(tick, newTick *model.Tick, tCh chan model.Tick) {
	tick.Close = newTick.Close
	tick.High = maxFloat(tick.High, newTick.High)
	tick.Low = minFloat(tick.Low, newTick.Low)
	tick.Vol += newTick.Vol
	tCh <- *newTick
	t.SetUpdateFlag()
}

func (t *Ticker) UpdateTicker(
	symbol string, trade *model.Trade, tCh chan model.Tick) {
	tick := t.makeTick(trade)
	tickList := t.GetTickerList(symbol)

	length := len(tickList)
	if length == 0 {
		t.appendTick(symbol, tick, tCh)
		return
	}
	lastTicker := tickList[length-1]
	if tick.Timestamp.Sub(lastTicker.Timestamp) > t.interval {
		t.appendTick(symbol, tick, tCh)
	} else {
		t.updateTick(lastTicker, tick, tCh)
	}
}

func (t *Ticker) SetUpdateFlag() {
	atomic.StoreInt32(&t.isUpdate, 1)
}

func (t *Ticker) IsUpdate() bool {
	return atomic.CompareAndSwapInt32(&t.isUpdate, 1, 0)
}

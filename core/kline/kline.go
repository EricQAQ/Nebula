package kline

import (
	"strconv"
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

const (
	klineLength = 2048
)

type Kline struct {
	Symbol    string
	Start     float64
	End       float64
	High      float64
	Low       float64
	Vol       float64
	Timestamp time.Time
}

type KlineManager struct {
	symbol    string
	intervals []string
	klineMap  cmap.ConcurrentMap
}

func newKlineManager(symbol string, intervals []int) *KlineManager {
	km := new(KlineManager)
	km.symbol = symbol
	km.intervals = make([]string, 0, len(intervals))
	km.klineMap = cmap.New()

	for _, val := range intervals {
		str := strconv.Itoa(val)
		km.intervals = append(km.intervals, str)
		km.klineMap.Set(str, make([]*Kline, 0, klineLength))
	}
	return km
}

func (km *KlineManager) getKline(interval int) []*Kline {
	str := strconv.Itoa(interval)
	data, _ := km.klineMap.Get(str)
	klines := data.([]*Kline)
	return klines
}

func (kv *KlineManager) newKline(tick *core.Tick) *Kline {
	kline := new(Kline)
	kline.Symbol = tick.Symbol
	kline.Start = tick.Last
	kline.End = tick.Last
	kline.High = tick.High
	kline.Low = tick.Low
	kline.Vol = tick.Vol
	kline.Timestamp = tick.Timestamp
	return kline
}

func (kv *KlineManager) updateKline(kline *Kline, tick *core.Tick) {
	kline.High = maxFloat(kline.High, tick.High)
	kline.Low = minFloat(kline.Low, tick.Low)
	kline.End = tick.Last
	kline.Timestamp = tick.Timestamp
}

func (km *KlineManager) appendKlines(interval int, klines []*Kline, tick *core.Tick) {
	kline := km.newKline(tick)
	if len(klines) >= klineLength {
		klines = klines[1:]
	}
	klines = append(klines, kline)

	str := strconv.Itoa(interval)
	km.klineMap.Set(str, klines)
}

func (km *KlineManager) updateKlines(interval int, tick *core.Tick) {
	klines := km.getKline(interval)
	length := len(klines)
	if length == 0 {
		km.appendKlines(interval, klines, tick)
	}
	iv := time.Duration(interval) * time.Second
	lastKline := klines[length-1]
	if tick.Timestamp.Sub(lastKline.Timestamp) > iv {
		km.appendKlines(interval, klines, tick)
	} else {
		km.updateKline(lastKline, tick)
	}
}

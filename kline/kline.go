package kline

import (
	"strconv"
	"time"

	"github.com/EricQAQ/Traed/model"

	"github.com/orcaman/concurrent-map"
)

const (
	klineLength = 2048
)

type Kline struct {
	Symbol    string
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Vol       float64
	Timestamp time.Time
}

func NewKline(data map[string]interface{}) *Kline {
	kline := new(Kline)
	kline.Symbol = data["symbol"].(string)
	kline.Open = data["open"].(float64)
	kline.Close = data["close"].(float64)
	kline.High = data["high"].(float64)
	kline.Low = data["low"].(float64)
	kline.Vol = data["volume"].(float64)

	loc, _ := time.LoadLocation("Asia/Chongqing")
	ts, _ := time.Parse(time.RFC3339, data["timestamp"].(string))
	kline.Timestamp = ts.In(loc)
	return kline
}

type KlineManager struct {
	symbol    string
	intervals []string
	updateMap cmap.ConcurrentMap
	klineMap  cmap.ConcurrentMap
}

func NewKlineManager(symbol string, intervals []int) *KlineManager {
	km := new(KlineManager)
	km.symbol = symbol
	km.intervals = make([]string, 0, len(intervals))
	km.updateMap = cmap.New()
	km.klineMap = cmap.New()

	for _, val := range intervals {
		str := strconv.Itoa(val)
		km.intervals = append(km.intervals, str)
		km.klineMap.Set(str, make([]*Kline, 0, klineLength))
		km.updateMap.Set(str, false)
	}
	return km
}

func (km *KlineManager) GetUpdate(interval int) bool {
	str := strconv.Itoa(interval)
	data, _ := km.updateMap.Get(str)
	km.updateMap.Set(str, false)
	return data.(bool)
}

func (km *KlineManager) GetKline(interval int) []*Kline {
	str := strconv.Itoa(interval)
	data, _ := km.klineMap.Get(str)
	klines := data.([]*Kline)
	return klines
}

func (km *KlineManager) SetKline(interval int, klines []*Kline) {
	str := strconv.Itoa(interval)
	km.klineMap.Set(str, klines)
}

func (km *KlineManager) newKline(interval int, tick *model.Tick) *Kline {
	kline := new(Kline)
	kline.Symbol = tick.Symbol
	kline.Open = tick.Open
	kline.Close = tick.Close
	kline.High = tick.High
	kline.Low = tick.Low
	kline.Vol = tick.Vol
	kline.Timestamp = tick.Timestamp.Truncate(
		time.Duration(interval) * time.Second)
	return kline
}

func (km *KlineManager) updateKline(interval int, kline *Kline, tick *model.Tick) {
	kline.Close = tick.Close
	kline.High = maxFloat(kline.High, tick.High)
	kline.Low = minFloat(kline.Low, tick.Low)
	kline.Vol += tick.Vol
	str := strconv.Itoa(interval)
	km.updateMap.Set(str, true)
}

func (km *KlineManager) appendKlines(interval int, klines []*Kline, tick *model.Tick) {
	kline := km.newKline(interval, tick)
	if len(klines) >= klineLength {
		klines = klines[1:]
	}
	klines = append(klines, kline)

	str := strconv.Itoa(interval)
	km.klineMap.Set(str, klines)
	km.updateMap.Set(str, true)
}

func (km *KlineManager) UpdateKlines(interval int, tick *model.Tick) {
	klines := km.GetKline(interval)
	length := len(klines)
	if length == 0 {
		km.appendKlines(interval, klines, tick)
		return
	}
	iv := time.Duration(interval) * time.Second
	lastKline := klines[length-1]
	if tick.Timestamp.Sub(lastKline.Timestamp) > iv {
		km.appendKlines(interval, klines, tick)
	} else {
		km.updateKline(interval, lastKline, tick)
	}
}

package core

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
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

type klineManager struct {
	symbol    string
	intervals []string
	updateMap cmap.ConcurrentMap
	klineMap  cmap.ConcurrentMap
}

func newKlineManager(symbol string, intervals []int) *klineManager {
	km := new(klineManager)
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

func (km *klineManager) getUpdate(interval int) bool {
	str := strconv.Itoa(interval)
	data, _ := km.updateMap.Get(str)
	km.updateMap.Set(str, false)
	return data.(bool)
}

func (km *klineManager) getKline(interval int) []*Kline {
	str := strconv.Itoa(interval)
	data, _ := km.klineMap.Get(str)
	klines := data.([]*Kline)
	return klines
}

func (km *klineManager) newKline(tick *Tick) *Kline {
	kline := new(Kline)
	kline.Symbol = tick.Symbol
	kline.Open = tick.Open
	kline.Close = tick.Close
	kline.High = tick.High
	kline.Low = tick.Low
	kline.Vol = tick.Vol
	kline.Timestamp = tick.Timestamp
	return kline
}

func (km *klineManager) updateKline(interval int, kline *Kline, tick *Tick) {
	kline.Close = tick.Close
	kline.High = maxFloat(kline.High, tick.High)
	kline.Low = minFloat(kline.Low, tick.Low)
	kline.Vol += tick.Vol
	str := strconv.Itoa(interval)
	km.updateMap.Set(str, true)
}

func (km *klineManager) appendKlines(interval int, klines []*Kline, tick *Tick) {
	kline := km.newKline(tick)
	if len(klines) >= klineLength {
		klines = klines[1:]
	}
	klines = append(klines, kline)

	str := strconv.Itoa(interval)
	km.klineMap.Set(str, klines)
	km.updateMap.Set(str, true)
}

func (km *klineManager) updateKlines(interval int, tick *Tick) {
	klines := km.getKline(interval)
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

type SymbolsKlineManager struct {
	ctx       context.Context
	wg        sync.WaitGroup
	exchange  ExchangeAPI
	symbols   []string
	intervals []int
	manager   cmap.ConcurrentMap
}

func NewSymbolsKlineManager(
	ctx context.Context, ex ExchangeAPI,
	symbols []string, intervals []int) *SymbolsKlineManager {
	skm := new(SymbolsKlineManager)
	skm.ctx = ctx
	skm.exchange = ex
	skm.symbols = symbols

	skm.intervals = make([]int, 0, len(intervals))
	skm.intervals = append(skm.intervals, intervals...)

	skm.manager = cmap.New()
	for _, symbol := range skm.symbols {
		skm.manager.Set(symbol, newKlineManager(symbol, intervals))
	}
	return skm
}

func (skm *SymbolsKlineManager) getKlineManager(symbol string) *klineManager {
	data, _ := skm.manager.Get(symbol)
	return data.(*klineManager)
}

func (skm *SymbolsKlineManager) update(tick *Tick) {
	km := skm.getKlineManager(tick.Symbol)
	for _, interval := range skm.intervals {
		km.updateKlines(interval, tick)
	}
}

func (skm *SymbolsKlineManager) GetKline(symbol string, interval int) (*Kline, bool) {
	mng := skm.getKlineManager(symbol)
	klineList := mng.getKline(interval)
	length := len(klineList)
	if length == 0 {
		return nil, false
	}
	return klineList[length-1], mng.getUpdate(interval)
}

func (skm *SymbolsKlineManager) startKlineManager() {
	tickCh := skm.exchange.GetTickChan()
	for {
		select {
		case <-skm.ctx.Done():
			return
		case tick := <-tickCh:
			skm.update(&tick)
		}
	}
}

func (skm *SymbolsKlineManager) Start() {
	skm.wg.Add(1)
	go func() {
		defer skm.wg.Done()
		skm.startKlineManager()
	}()
	log.Infof(
		"[Kline Generator(%s)] Start K Line Generator.",
		skm.exchange.GetExchangeName())
}

func (skm *SymbolsKlineManager) Stop() {
	skm.wg.Wait()
	log.Infof(
		"[Kline Generator(%s)] Stop K Line Generator.",
		skm.exchange.GetExchangeName())
}

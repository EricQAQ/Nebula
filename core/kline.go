package core

import (
	"context"
	"sync"

	"github.com/EricQAQ/Traed/model"
	"github.com/EricQAQ/Traed/kline"

	log "github.com/sirupsen/logrus"
	"github.com/orcaman/concurrent-map"
)

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
		skm.manager.Set(symbol, kline.NewKlineManager(symbol, intervals))
	}
	return skm
}

func (skm *SymbolsKlineManager) getKlineManager(symbol string) *kline.KlineManager {
	data, _ := skm.manager.Get(symbol)
	return data.(*kline.KlineManager)
}

func (skm *SymbolsKlineManager) update(tick *model.Tick) {
	km := skm.getKlineManager(tick.Symbol)
	for _, interval := range skm.intervals {
		km.UpdateKlines(interval, tick)
	}
}

func (skm *SymbolsKlineManager) GetKline(symbol string, interval int) (*kline.Kline, bool) {
	mng := skm.getKlineManager(symbol)
	klineList := mng.GetKline(interval)
	length := len(klineList)
	if length == 0 {
		return nil, false
	}
	return klineList[length-1], mng.GetUpdate(interval)
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

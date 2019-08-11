package core

import (
	"context"
	"sync"
	"time"

	"github.com/EricQAQ/Traed/kline"
	"github.com/EricQAQ/Traed/model"
	"github.com/EricQAQ/Traed/storage"

	"github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
)

type SymbolsKlineManager struct {
	ctx       context.Context
	wg        sync.WaitGroup
	exchange  ExchangeAPI
	store     storage.StorageAPI
	symbols   []string
	intervals []int
	manager   cmap.ConcurrentMap
}

func NewSymbolsKlineManager(
	ctx context.Context, ex ExchangeAPI, store storage.StorageAPI,
	symbols []string, intervals []int) *SymbolsKlineManager {
	skm := new(SymbolsKlineManager)
	skm.ctx = ctx
	skm.exchange = ex
	skm.store = store
	skm.symbols = symbols

	skm.intervals = make([]int, 0, len(intervals))
	skm.intervals = append(skm.intervals, intervals...)

	skm.manager = cmap.New()
	for _, symbol := range skm.symbols {
		skm.manager.Set(symbol, kline.NewKlineManager(symbol, intervals))
		skm.loadHistoryData(symbol)
	}
	return skm
}

func (skm *SymbolsKlineManager) loadHistoryData(symbol string) error {
	log.Infof("[Traed Kline(%s)] Start load history k-line", skm.exchange.GetExchangeName())
	endTime := time.Now().Truncate(time.Minute)
	startTime := endTime.AddDate(0, 0, -1)
	klist, err := skm.store.GetKlines(
		skm.exchange.GetExchangeName(), symbol, startTime, endTime)
	if err != nil {
		return LoadHistoryErr.FastGen(err.Error())
	}
	km := skm.getKlineManager(symbol)
	km.SetKline(60, klist)

	for _, i := range skm.intervals {
		km.SetKline(i, kline.AggregateKlines(60, i, klist))
	}
	log.Infof("[Traed Kline(%s)] Finish load history k-line", skm.exchange.GetExchangeName())
	return nil
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

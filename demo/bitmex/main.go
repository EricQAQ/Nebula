package main

import (
	"os"
	"flag"
	"time"

	"github.com/EricQAQ/Traed/core"
	"github.com/EricQAQ/Traed/plugin/bitmex"

	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		configPath = flag.String("config", "", "config file path")

		app *core.TraedApp
		err error
	)
	flag.Parse()

	app = core.NewTraedApp(*configPath)
	bxCfg := app.Cfg.ExchangeMap["bitmex"]
	bm := bitmex.CreateBitmex(bxCfg, app.Cfg.Http)
	app.SetExchange("bitmex", bm)

	if err = app.Start(); err != nil {
		os.Exit(0)
	}
	go func() {
		for {
			tick, isUpdate := bm.GetTick("XBTUSD")
			if tick == nil {
				continue
			}
			if isUpdate {
				log.Infof(`
Receive tick data:
symbol: %s, open: %f, close: %f, high: %f, low: %f, vol: %f, time: %s`,
					tick.Symbol, tick.Open, tick.Close, tick.High,
					tick.Low, tick.Vol, tick.Timestamp)
			}
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()

	go func() {
		for {
			positions, isUpdate := bm.GetPosition("XBTUSD")
			if isUpdate {
				for _, pos := range positions {
					log.Infof(`
Position: leverage:%f,
sell_amount:%f, sell_avaiable:%f, sell_price_avg:%f, sell_profit_real:%f
buy_amount:%f, buy_avaiable:%f, buy_price_avg:%f, buy_profit_real:%f`,
						pos.LeverRate, pos.SellAmount, pos.SellAvailable, pos.SellPriceAvg, pos.SellProfitReal,
						pos.BuyAmount, pos.BuyAvailable, pos.BuyPriceAvg, pos.BuyProfitReal)
				}
			}
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()

	go func() {
		for {
			for _, interval := range app.Cfg.KlineInterval {
				kline, isUpdate := app.GetKline("bitmex", "XBTUSD", interval)
				if isUpdate {
					log.Infof("K line %d, symbol:%s, open:%f, close:%f, high:%f, low:%f, vol:%f, time:%s", interval, kline.Symbol, kline.Open, kline.Close, kline.High, kline.Low, kline.Vol, kline.Timestamp)
				}
			}
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()
	app.Stop()
	os.Exit(0)
}

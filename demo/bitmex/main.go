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
		var ts time.Time
		for {
			tick := bm.GetTick("XBTUSD")
			if tick == nil {
				continue
			}
			if tick.Timestamp.After(ts) {
				ts = tick.Timestamp
				log.Infof(
					"Receive tick data: symbol: %s, last: %f, buy: %f, sell: %f, high: %f, low: %f, vol: %f, time: %s", tick.Symbol, tick.Last, tick.Buy, tick.Sell, tick.High, tick.Low, tick.Vol, tick.Timestamp)
			}
		}
	}()
	app.Stop()
	os.Exit(0)
}

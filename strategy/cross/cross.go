package main

import (
	"os"
	"time"

	"github.com/EricQAQ/Nebula/core"
	"github.com/EricQAQ/Nebula/finIndex"
	"github.com/EricQAQ/Nebula/plugin/bitmex"

	log "github.com/sirupsen/logrus"
)

func createApp() *core.NebulaApp {
	app := core.NewNebulaApp("bin/nebula.toml")
	bxCfg := app.Cfg.ExchangeMap["bitmex"]
	bm := bitmex.CreateBitmex(bxCfg, app.Cfg.Http)
	app.SetExchange("bitmex", bm)
	return app
}

func strategy(app *core.NebulaApp) {
	klines, _ := app.GetKline("bitmex", "XBTUSD", 60)
	ma := fin.NewMA(klines, 99)
	ema := fin.NewEMA(klines, 17)
	ema.Calculation()
	ma.Calculation()
	bm := app.Exchange["bitmex"]
	for {
		klines, isUpdate := app.GetKline("bitmex", "XBTUSD", 60)
		if isUpdate {
			ma.InsertKline(klines[len(klines)-1])
			ema.InsertKline(klines[len(klines)-1])
			maLen := len(ma.Points)
			emaLen := len(ema.Points)
			log.Debugf(
				"(ema, ma): (%f, %f), (%f, %f)\n",
				ema.Points[emaLen-2].Value, ma.Points[maLen-2].Value,
				ema.Points[emaLen-1].Value, ma.Points[maLen-1].Value)

			if ema.Points[emaLen-2].Value < ma.Points[maLen-2].Value &&
				ema.Points[emaLen-1].Value > ma.Points[maLen-1].Value {
				depth, _ := bm.GetDepth("XBTUSD")
				_, err := bm.LimitBuy("XBTUSD", depth.Buy[len(depth.Buy)-1].Price, 10)
				if err != nil {
					log.Warnf("Buy failed: %s", err.Error())
				}
			} else if ema.Points[emaLen-2].Value > ma.Points[maLen-2].Value &&
				ema.Points[emaLen-1].Value < ma.Points[maLen-1].Value {
				depth, _ := bm.GetDepth("XBTUSD")
				_, err := bm.LimitSell("XBTUSD", depth.Sell[0].Price, 10)
				if err != nil {
					log.Warnf("Sell failed: %s", err.Error())
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	app := createApp()
	if err := app.Start(); err != nil {
		os.Exit(0)
	}
	go strategy(app)
	app.Stop()
	os.Exit(0)
}

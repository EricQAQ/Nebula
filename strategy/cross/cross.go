package main

import (
	"os"
	"fmt"
	"time"

	"github.com/EricQAQ/Nebula/core"
	"github.com/EricQAQ/Nebula/finIndex"
	"github.com/EricQAQ/Nebula/plugin/bitmex"
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
			fmt.Printf(
				"(ema, ma): (%f, %f), (%f, %f)\n",
				ema.Points[emaLen-2].Value, ma.Points[maLen-2].Value,
				ema.Points[emaLen-1].Value, ma.Points[maLen-1].Value)

			if ema.Points[emaLen-2].Value < ma.Points[maLen-2].Value &&
				ema.Points[emaLen-1].Value > ma.Points[maLen-1].Value {
				fmt.Println("Buy!!")
				depth, _ := bm.GetDepth("XBTUSD")
				_, err := bm.LimitBuy("XBTUSD", depth.Buy[len(depth.Buy)-1].Price, 10)
				fmt.Println("Buy!! %v", err)
			} else if ema.Points[emaLen-2].Value > ma.Points[maLen-2].Value &&
				ema.Points[emaLen-1].Value < ma.Points[maLen-1].Value {
				fmt.Println("Sell!!")
				depth, _ := bm.GetDepth("XBTUSD")
				_, err := bm.LimitSell("XBTUSD", depth.Sell[0].Price, 10)
				fmt.Println("Sell!! %v", err)
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

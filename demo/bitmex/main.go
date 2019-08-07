package main

import (
	"os"
	"flag"

	"github.com/EricQAQ/Traed/core"
	"github.com/EricQAQ/Traed/plugin/bitmex"
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
	app.Stop()
	os.Exit(0)
}

package core

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/EricQAQ/Traed/config"
	"github.com/EricQAQ/Traed/kline"
	"github.com/EricQAQ/Traed/logger"
	"github.com/EricQAQ/Traed/storage"
	"github.com/EricQAQ/Traed/storage/csv"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

const (
	VERSION = "v0.1.0"
)

var shutdownCtx, cancel = context.WithCancel(context.Background())

func printInfo() {
	level := log.GetLevel()
	log.SetLevel(log.InfoLevel)
	PrintLogo()
	PrintInfo()
	log.SetLevel(level)
}

func LoadConfig(configPath string) *config.TraedConfig {
	cfg := config.GetTraedConfig()
	if configPath != "" {
		if err := cfg.LoadFromToml(configPath); err != nil {
			log.Fatalf(errors.ErrorStack(err))
		}
	}
	return cfg
}

func SetupLog(cfg *config.TraedConfig) {
	err := logger.CreateLoggerFromConfig(cfg)
	if err != nil {
		log.Fatalf(errors.ErrorStack(err))
	}
}

type TraedApp struct {
	Cfg      *config.TraedConfig
	Exchange map[string]ExchangeAPI
	wsMap    map[string]*WsClient
	klineMng map[string]*SymbolsKlineManager
	store    storage.StorageAPI
}

func NewTraedApp(cfgPath string) *TraedApp {
	cfg := LoadConfig(cfgPath)
	SetupLog(cfg)

	app := new(TraedApp)
	app.Cfg = cfg
	app.Exchange = make(map[string]ExchangeAPI)
	app.wsMap = make(map[string]*WsClient)
	app.klineMng = make(map[string]*SymbolsKlineManager)
	app.setStorage()
	return app
}

func (app *TraedApp) setStorage() {
	switch app.Cfg.Storage.StorageType {
	case "csv":
		app.store = csv.NewCsvStorage(app.Cfg.Storage.Csv.DataDir)
	}
}

func (app *TraedApp) SetExchange(exchangeName string, exchange ExchangeAPI) error {
	_, ok := app.Cfg.ExchangeMap[exchangeName]
	if !ok {
		return ExchangeNotExistErr.FastGen(exchangeName)
	}
	app.Exchange[exchangeName] = exchange
	return nil
}

func (app *TraedApp) CreateWsClient() error {
	for name, exchange := range app.Exchange {
		exCfg := app.Cfg.ExchangeMap[name]
		worker := NewWorker(shutdownCtx, name, exchange)
		ws, err := NewWsClient(
			shutdownCtx, name, exCfg.Address,
			exchange, app.Cfg.Websocket, worker)
		if err != nil {
			return err
		}
		app.wsMap[name] = ws
	}
	return nil
}

func (app *TraedApp) setupSingalHandler() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Kill, os.Interrupt,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		log.Infof("Got signal [%s], prepare to exit.", sig)
		app.Shutdown()
	}()
}

func (app *TraedApp) GetKline(exchange, symbol string, interval int) (*kline.Kline, bool) {
	return app.klineMng[exchange].GetKline(symbol, interval)
}

func (app *TraedApp) Start() error {
	app.setupSingalHandler()
	printInfo()
	for name, _ := range app.Exchange {
		if err := app.CreateWsClient(); err != nil {
			return CreateWsErr.FastGen(name, err.Error())
		}
	}
	for name, exCfg := range app.Cfg.ExchangeMap {
		exchange := app.Exchange[name]
		app.klineMng[name] = NewSymbolsKlineManager(
			shutdownCtx, exchange, app.store,
			exCfg.Symbols, app.Cfg.KlineInterval)
	}
	for _, ws := range app.wsMap {
		ws.StartClient(app)
	}
	for _, mng := range app.klineMng {
		mng.Start()
	}
	log.Infof("[April] Start April App.")
	return nil
}

func (app *TraedApp) Stop() {
	for _, ws := range app.wsMap {
		ws.StopClient()
	}
	for _, mng := range app.klineMng {
		mng.Stop()
	}
	log.Infof("[April] Stop April App.")
}

func (app *TraedApp) Shutdown() {
	cancel()
}

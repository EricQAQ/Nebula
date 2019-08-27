package core

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/EricQAQ/Nebula/config"
	"github.com/EricQAQ/Nebula/kline"
	"github.com/EricQAQ/Nebula/logger"
	"github.com/EricQAQ/Nebula/storage"
	"github.com/EricQAQ/Nebula/storage/csv"

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

func LoadConfig(configPath string) *config.NebulaConfig {
	cfg := config.GetNebulaConfig()
	if configPath != "" {
		if err := cfg.LoadFromToml(configPath); err != nil {
			log.Fatalf(errors.ErrorStack(err))
		}
	}
	return cfg
}

func SetupLog(cfg *config.NebulaConfig) {
	err := logger.CreateLoggerFromConfig(cfg)
	if err != nil {
		log.Fatalf(errors.ErrorStack(err))
	}
}

type NebulaApp struct {
	Cfg      *config.NebulaConfig
	Exchange map[string]ExchangeAPI
	wsMap    map[string]*WsClient
	klineMng map[string]*SymbolsKlineManager
	store    storage.StorageAPI
}

func NewNebulaApp(cfgPath string) *NebulaApp {
	cfg := LoadConfig(cfgPath)
	SetupLog(cfg)

	app := new(NebulaApp)
	app.Cfg = cfg
	app.Exchange = make(map[string]ExchangeAPI)
	app.wsMap = make(map[string]*WsClient)
	app.klineMng = make(map[string]*SymbolsKlineManager)
	app.setStorage()
	return app
}

func (app *NebulaApp) setStorage() {
	switch app.Cfg.Storage.StorageType {
	case "csv":
		app.store = csv.NewCsvStorage(app.Cfg.Storage.Csv.DataDir)
	}
}

func (app *NebulaApp) SetExchange(exchangeName string, exchange ExchangeAPI) error {
	_, ok := app.Cfg.ExchangeMap[exchangeName]
	if !ok {
		return ExchangeNotExistErr.FastGen(exchangeName)
	}
	app.Exchange[exchangeName] = exchange
	return nil
}

func (app *NebulaApp) CreateWsClient() error {
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

func (app *NebulaApp) setupSingalHandler() {
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

func (app *NebulaApp) GetKline(exchange, symbol string, interval int) ([]*kline.Kline, bool) {
	return app.klineMng[exchange].GetKline(symbol, interval)
}

func (app *NebulaApp) Start() error {
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
	log.Infof("[Nebula] Start Nebula App.")
	return nil
}

func (app *NebulaApp) Stop() {
	for _, ws := range app.wsMap {
		ws.StopClient()
	}
	for _, mng := range app.klineMng {
		mng.Stop()
	}
	log.Infof("[Nebula] Stop Nebula App.")
}

func (app *NebulaApp) Shutdown() {
	cancel()
}

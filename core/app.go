package core

import (
	"context"

	"github.com/EricQAQ/Traed/config"
	"github.com/EricQAQ/Traed/logger"

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
}

func NewTraedApp(cfgPath string) *TraedApp {
	cfg := LoadConfig(cfgPath)
	SetupLog(cfg)

	app := new(TraedApp)
	app.Cfg = cfg
	app.Exchange = make(map[string]ExchangeAPI)
	app.wsMap = make(map[string]*WsClient)
	return app
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
		worker := NewWorker(shutdownCtx, name, exchange.GetCallbackHandler())
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

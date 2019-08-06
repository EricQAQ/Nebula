package core

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Worker struct {
	exchange         string
	ctx              context.Context
	workerCh         chan ParsedData
	msgHandler ExchangeAPI
}

func NewWorker(ctx context.Context, exchange string,
	msgHandler ExchangeAPI) *Worker {
	w := new(Worker)
	w.ctx = ctx
	w.exchange = exchange
	w.workerCh = make(chan ParsedData, 1024)
	w.msgHandler = msgHandler
	return w
}

func (w *Worker) StartWorker(app *TraedApp) {
	log.Infof("[Traed Worker(%s)] start.", w.exchange)
	for {
		select {
		case <-w.ctx.Done():
			return
		case data := <-w.workerCh:
			if data.Type != Message || data.Type != ErrorMsg {
				continue
			}
			w.msgHandler(data)
		}
	}
}

func (w *Worker) StopWorker() {
	close(w.workerCh)
	log.Infof("[Traed Worker(%s)] stop.", w.exchange)
}


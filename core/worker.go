package core

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Worker struct {
	exchange         string
	ctx              context.Context
	workerCh         chan ParsedData
	callbackHandlers CallbackHandler
}

func NewWorker(ctx context.Context, exchange string,
	callbackHandlers CallbackHandler) *Worker {
	w := new(Worker)
	w.ctx = ctx
	w.exchange = exchange
	w.workerCh = make(chan ParsedData, 1024)
	w.callbackHandlers = callbackHandlers
	return w
}

func (w *Worker) StartWorker(app *TraedApp) {
	log.Infof("[Traed Worker(%s)] start.", w.exchange)
	for {
		select {
		case <-w.ctx.Done():
			return
		case data := <-w.workerCh:
			switch data.Dtype {
			case WelcomeType:
				w.callbackHandlers.HandleWelcomeMsg(data)
			case SubscribeType:
				w.callbackHandlers.HandleSubscribeMsg(data)
			case AuthType:
				w.callbackHandlers.HandleAuthMsg(data)
			case StatusType:
				w.callbackHandlers.HandleStatusMsg(data)
			case RetType:
				w.callbackHandlers.HandleRetMsg(app, data)
			}
		}
	}
}

func (w *Worker) StopWorker() {
	close(w.workerCh)
	log.Infof("[Traed Worker(%s)] stop.", w.exchange)
}


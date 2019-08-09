package core

import (
	"context"
	"sync"
	"time"

	"github.com/EricQAQ/Traed/config"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsClient struct {
	wg           sync.WaitGroup
	ctx          context.Context
	exchangeName string
	wsAddress    string
	exchange     ExchangeAPI
	WsConfig     *config.WebsocketConfig

	ReadWait          time.Duration
	WriteWait         time.Duration
	HeartbeatDuration time.Duration

	Conn   *websocket.Conn
	tick   *ticker
	worker *Worker
}

func NewWsClient(
	ctx context.Context, exchangeName string,
	wsAddress string, exchange ExchangeAPI,
	wsConfig *config.WebsocketConfig,
	worker *Worker) (*WsClient, error) {
	ws := new(WsClient)
	ws.ctx = ctx
	ws.exchangeName = exchangeName
	ws.exchange = exchange
	ws.wsAddress = wsAddress
	ws.WsConfig = wsConfig
	ws.worker = worker

	ws.ReadWait = time.Duration(ws.WsConfig.ReadWait) * time.Millisecond
	ws.WriteWait = time.Duration(ws.WsConfig.WriteWait) * time.Millisecond
	ws.HeartbeatDuration = time.Duration(ws.WsConfig.HeartbeatDuration) * time.Second

	if err := ws.connect(false); err != nil {
		return nil, err
	}
	ws.tick = newTicker(ws)
	return ws, nil
}

func (ws *WsClient) connect(reconnect bool) error {
	if reconnect {
		ws.Conn.Close()
	}
	var err error
	ws.Conn, _, err = websocket.DefaultDialer.Dial(ws.wsAddress, nil)
	if err != nil {
		return err
	}
	if err = ws.Auth(); err != nil {
		return err
	}
	return ws.Subscribe()
}

func (ws *WsClient) retryReconnect() error {
	var err error
	retryCount := ws.WsConfig.RetryCount
	for retryCount > 0 {
		if err = ws.connect(true); err != nil {
			retryCount--
			log.Errorf("[Traed WsClient(%s)] failed to reconnect: %s", ws.exchangeName, err.Error())
			time.Sleep(time.Second * 1)
		} else {
			return nil
		}
	}
	return RetryMaxErr
}

func (ws *WsClient) WriteTextMsg(msg []byte) error {
	if err := ws.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return err
	}
	return nil
}

func (ws *WsClient) WriteBinaryMsg(msg []byte) error {
	if err := ws.Conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		return err
	}
	return nil
}

func (ws *WsClient) ReadMsg() error {
	for {
		select {
		case <-ws.ctx.Done():
			return nil
		default:
			ws.Conn.SetReadDeadline(time.Now().Add(ws.ReadWait))
			_, msg, err := ws.Conn.ReadMessage()
			if err != nil {
				log.Warnf("[Traed WsClient(%s)] receive message error: %s, reconnect.", ws.exchangeName, err.Error())
				if err = ws.retryReconnect(); err != nil {
					log.Errorf("[Traed WsClient(%s)] reconnect failed: %s", ws.exchangeName, err.Error())
					return err
				}
				continue
			}
			data, err := ws.exchange.Parse(msg)
			if err != nil {
				log.Errorf("[Traed WsClient(%s)] parse data error: %s", ws.exchangeName, err.Error())
				continue
			}
			if data.Data == nil {
				continue
			}
			ws.worker.workerCh <- *data
		}
	}
}

func (ws *WsClient) StartClient(app *TraedApp) {
	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()
		ws.worker.StartWorker(app)
	}()

	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()
		ws.tick.startTicker()
	}()

	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()
		ws.ReadMsg()
	}()

	log.Infof("[Traed WsClient(%s)] start.", ws.exchangeName)
}

func (ws *WsClient) StopClient() {
	ws.wg.Wait()
	ws.Conn.Close()
	ws.worker.StopWorker()
	log.Infof("[Traed WsClient(%s)] stop.", ws.exchangeName)
}

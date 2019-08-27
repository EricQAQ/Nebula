package core

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ticker struct {
	tc       *time.Ticker
	wsClient *WsClient
}

func newTicker(wsClient *WsClient) *ticker {
	tc := time.NewTicker(wsClient.HeartbeatDuration)
	t := &ticker{
		tc:       tc,
		wsClient: wsClient,
	}
	return t
}

func (ticker *ticker) heartBeat() error {
	ticker.wsClient.Conn.SetWriteDeadline(
		time.Now().Add(ticker.wsClient.WriteWait))
	err := ticker.wsClient.Conn.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		log.Errorf("[Nebula Ticker] send hearbeat failed. %s", err.Error())
		return err
	}
	return nil
}

func (ticker *ticker) startTicker() {
	defer ticker.tc.Stop()
	for {
		select {
		case <-ticker.tc.C:
			if err := ticker.heartBeat(); err != nil {
				if err = ticker.wsClient.connect(true); err != nil {
					break
				}
			}
		case <-ticker.wsClient.ctx.Done():
			return
		}
	}
}


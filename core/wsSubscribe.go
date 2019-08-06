package core

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func (ws *WsClient) Auth() error {
	log.Infof("[Traed WsClient(%s)] Authenticate.", ws.exchangeName)
	au := ws.exchange.GetWsAuthHandler()
	data, err := au.Serialize()
	if err != nil {
		return err
	}
	return ws.WriteTextMsg(data)
}

func (ws *WsClient) Subscribe() error {
	sub := ws.exchange.GetWsSubscribeHandler()
	log.Infof(
		"[Traed WsClient(%s)] Subscribe: %s",
		ws.exchangeName,
		strings.Join(sub.GetOperateArgs(), ","))
	data, err := sub.Serialize()
	if err != nil {
		return err
	}
	if err = ws.WriteTextMsg(data); err != nil {
		return err
	}
	return nil
}


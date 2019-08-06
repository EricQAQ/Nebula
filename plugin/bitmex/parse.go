package bitmex

import (
	"encoding/json"

	"github.com/EricQAQ/Traed/core"

	log "github.com/sirupsen/logrus"
)

func (bm *Bitmex) isAuthMsg(rv map[string]interface{}) bool {
	if _, ok := rv["success"]; !ok {
		return false
	}
	req, ok := rv["request"]
	if !ok {
		return false
	}
	reqData := req.(map[string]interface{})
	if reqData["op"].(string) != "authKeyExpires" {
		return false
	}
	return true
}

func (bm *Bitmex) Parse(data []byte) (*core.ParsedData, error) {
	pd := new(core.ParsedData)
	rv := make(map[string]interface{})
	log.Debugf("receive msg: %s", string(data))
	if err := json.Unmarshal(data, &rv); err != nil {
		return nil, err
	}
	pd.Data = rv

	if _, ok := rv["info"]; ok {
		pd.Type = core.WelcomeMsg
		return pd, nil
	}

	if val, ok := rv["subscribe"]; ok {
		pd.Type = core.SubscribeMsg
		if rv["success"].(bool) {
			log.Infof("[Traed Bitmex] Subscribe success: %s", val.Subscribe)
		} else {
			log.Infof("[April Bitmex] Subscribe failed: %s", val.Subscribe)
		}
		return pd, nil
	}

	if isAuthMsg(rv) {
		pd.Type = core.AuthMsg
		if rv["success"].(bool) {
			log.Infof("[April Bitmex] Auth success.")
		} else {
			log.Infof("[April Bitmex] Auth failed.")
		}
		return pd, nil
	}

	if _, ok := rv["status"]; ok {
		pd.Type = core.ErrorMsg
		return pd, nil
	}

	pd.Type = core.Message
	return pd, nil
}

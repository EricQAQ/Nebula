package bitmex

import (
	"time"
	"fmt"
)

type BitmexAuth struct {
	op string
	apiKey string
	apiSecret string
	expire int
}

func NewBitmexAuth(apiKey, apiSecret string, expire int) *BitmexAuth {
	ba := new(BitmexAuth)
	ba.op = "authKeyExpires"
	ba.expire = expire
	ba.apiKey = apiKey
	ba.apiSecret = apiSecret
	return ba
}

func (ba *BitmexAuth) GetOperateArgs() []string {
	argList := make([]string, 0, 3)
	return argList
}

func (ba *BitmexAuth) Serialize() ([]byte, error) {
	expTs := time.Now().Add(time.Duration(ba.expire) * time.Hour).Unix()
	signature, err := GenerateSign(ba.apiSecret, "GET", "/realtime", expTs, nil)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(
		"{\"op\":\"%s\", \"args\":[\"%s\",%d,\"%x\"]}",
		ba.op, ba.apiKey, expTs, signature,
	)), nil
}

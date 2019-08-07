package bitmex

import (
	"strconv"
	"fmt"
	"time"
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

type makeResp func(map[string]interface{}) interface{}

func (bm *Bitmex) makeAuthHeader(
	req *gorequest.SuperAgent, method, url string, expireTs int64,
	data map[string]interface{}) error {
	sig, err := GenerateSign(
		bm.APISecret, method, url, expireTs, data)
	if err != nil {
		return err
	}
	req.Set("api-key", bm.APIKey)
	req.Set("api-expires", strconv.FormatInt(expireTs, 10))
	req.Set("api-signature", fmt.Sprintf("%x", sig))
	return nil
}

func (bm *Bitmex) doAuthRequest(
	method, url string, data map[string]interface{},
	expire int64, makeRespFn makeResp) (interface{}, error) {
	body, err := bm.sendRequest(method, url, data, true, expire)
	if err != nil {
		return nil, err
	}
	bodyMap := make(map[string]interface{})
	if err = json.Unmarshal(body, &bodyMap); err != nil {
		return nil, err
	}
	return makeRespFn(bodyMap), nil
}

func (bm *Bitmex) sendRequest(
	method, url string, data map[string]interface{},
	needAuth bool, expire int64) ([]byte, error) {
	uri := bm.BaseUrl + url
	request := gorequest.New().CustomMethod(method, uri)
	if len(bm.Proxy) > 0 {
		request = request.Proxy(bm.Proxy)
	}
	if bm.timeout > 0 {
		request = request.Timeout(bm.timeout)
	}
	if needAuth {
		expTs := time.Now().Add(time.Duration(expire) * time.Hour).Unix()
		if err := bm.makeAuthHeader(request, method, url, expTs, data); err != nil {
			return nil, err
		}
	}
	var respBody []byte
	var errs []error
	_, respBody, errs = request.Send(data).EndBytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return respBody, nil
}

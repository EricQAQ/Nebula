package bitmex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
)

type makeResp func(map[string]interface{}) interface{}
type makeRespList func([]map[string]interface{}) interface{}

func (bm *Bitmex) makeAuthHeader(
	req *gorequest.SuperAgent, method, url string, expireTs int64,
	data map[string]interface{}) error {
	sig, err := GenerateSign(
		bm.APISecret, method, routeUrl+url, expireTs, data)
	if err != nil {
		return err
	}
	req.Set("api-key", bm.APIKey)
	req.Set("api-expires", strconv.FormatInt(expireTs, 10))
	req.Set("api-signature", fmt.Sprintf("%x", sig))
	return nil
}

func (bm *Bitmex) doRequestGetList(
	method, url string, data map[string]interface{},
	makeRespFn makeRespList) (interface{}, error) {
	body, err := bm.sendRequest(method, url, data, false, 0)
	if err != nil {
		return nil, err
	}
	bodyMap := make([]map[string]interface{}, 0, 1024)
	if err = json.Unmarshal(body, &bodyMap); err != nil {
		return nil, err
	}
	return makeRespFn(bodyMap), nil
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
	request := gorequest.New().
		CustomMethod(method, uri).
		Retry(bm.retryCount, bm.retryInterval,
			http.StatusBadRequest,
			http.StatusInternalServerError,
			http.StatusGatewayTimeout)
	if len(bm.proxy) > 0 {
		request = request.Proxy(bm.proxy)
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
	resp, respBody, errs := request.Send(data).EndBytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode == http.StatusNotFound{
		return nil, ResponseErr.FastGen(string(respBody))
	}
	return respBody, nil
}

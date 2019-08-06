package bitmex

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/url"
	"encoding/json"
)

// Generates an API signature.
// A signature is HMAC_SHA256(secret, verb + path + nonce + data), hex encoded.
// Verb must be uppercased, url is relative, nonce must be an increasing 64-bit integer
// and the data, if present, must be JSON without whitespace between keys.
//
// For example, in psuedocode (and in real code below):
//
// verb=POST
// url=/api/v1/order
// nonce=1416993995705
// data={"symbol":"XBTZ14","quantity":1,"price":395.01}
// signature = HEX(HMAC_SHA256(secret, 'POST/api/v1/order1416993995705{"symbol":"XBTZ14","quantity":1,"price":395.01}'))
func GenerateSign(apiSecret, verb, urlpath string,
	nonce int64, data map[string]interface{}) ([]byte, error) {
	urlObj, err := url.Parse(urlpath)
	if err != nil {
		return nil, err
	}
	path := urlObj.Path
	if len(urlObj.RawQuery) > 0 {
		path = fmt.Sprintf("%s?%s", path, urlObj.RawQuery)
	}

	var val []byte
	if data != nil {
		val, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	message := []byte(
		fmt.Sprintf("%s%s%d%s", verb, path, nonce, string(val)))
	signature := hmac.New(sha256.New, []byte(apiSecret))
	_, err = signature.Write(message)
	if err != nil {
		return nil, err
	}
	return signature.Sum(nil), nil
}


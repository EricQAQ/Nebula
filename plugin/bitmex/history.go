package bitmex

import (
	"time"

	"github.com/EricQAQ/Nebula/kline"
)

func (bm *Bitmex) makeHistoryKlineList(data []map[string]interface{}) interface{} {
	tl := make([]*kline.Kline, 0, len(data))
	for _, ret := range data {
		tl = append(tl, kline.NewKline(ret))
	}
	return tl
}

func (bm *Bitmex) GetHistoryKline(
	symbol, binSize string, start, end time.Time) ([]*kline.Kline, error) {
	data := make(map[string]interface{})
	data["symbol"] = symbol
	data["binSize"] = binSize
	data["startTime"] = start
	data["endTime"] = end

	tl, err := bm.doRequestGetList(
		"GET", "/trade/bucketed", data, bm.makeHistoryKlineList)
	if err != nil {
		return nil, err
	}
	return tl.([]*kline.Kline), nil
}

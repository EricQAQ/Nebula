package bitmex

import (
	"sync/atomic"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

type position struct {
	positionKeys map[string][]string
	positionData cmap.ConcurrentMap
	isUpdate     int32
}

func newPosition(symbols []string) *position {
	pos := new(position)
	pos.positionKeys = make(map[string][]string)
	pos.positionData = cmap.New()
	pos.isUpdate = 0
	for _, symbol := range symbols {
		pos.positionKeys[symbol] = wsPositionKeys
		pos.positionData.Set(symbol, make([]*core.Position, 0, dataLength))
	}
	return pos
}

func (pos *position) getPositionList(symbol string) []*core.Position {
	data, _ := pos.positionData.Get(symbol)
	posList := data.([]*core.Position)
	return posList
}

func (pos *position) makePosition(data map[string]interface{}) *core.Position {
	p := new(core.Position)
	p.Symbol = data["symbol"].(string)
	p.Account = data["account"].(float64)
	p.Currency = data["currency"].(string)
	p.LeverRate = data["leverage"].(float64)
	p.ForceLiquPrice = data["liquidationPrice"].(float64)
	p.OpenOrderBuyQty = data["openOrderBuyQty"].(float64)
	p.OpenOrderSellQty = data["openOrderSellQty"].(float64)

	currQry := data["currentQty"].(float64)
	if currQry > 0 {
		// hold long position
		p.BuyAmount = currQry
		p.BuyPriceCost = data["avgCostPrice"].(float64)
		p.BuyPriceAvg = data["avgEntryPrice"].(float64)
		p.BuyProfitReal = data["unrealisedPnlPcnt"].(float64)
		p.BuyAvailable = p.BuyAmount - p.OpenOrderBuyQty
	} else {
		p.SellAmount = -currQry
		p.SellPriceCost = data["avgCostPrice"].(float64)
		p.SellPriceAvg = data["avgEntryPrice"].(float64)
		p.SellProfitReal = data["unrealisedPnlPcnt"].(float64)
		p.SellAvailable = p.SellAmount - p.OpenOrderSellQty
	}
	return p
}

func (pos *position) insertPosition(symbol string, position *core.Position) {
	posList := pos.getPositionList(symbol)
	length := len(posList)
	if length >= dataLength {
		posList = posList[length-dataLength:]
	}
	posList = append(posList, position)
	pos.positionData.Set(symbol, posList)
	atomic.StoreInt32(&pos.isUpdate, 1)
}

func (pos *position) findPositionItemByKeys(
	symbol string, updateData map[string]interface{}) (int, *core.Position) {
	posList := pos.getPositionList(symbol)
	for index, val := range posList {
		if val.Account == updateData["account"].(float64) &&
			val.Symbol == updateData["symbol"].(string) &&
			val.Currency == updateData["currency"].(string) {
			return index, val
		}
	}
	return 0, nil
}

func (pos *position) updatePosition(symbol string, data map[string]interface{}) {
	_, position := pos.findPositionItemByKeys(symbol, data)
	if position == nil {
		return
	}

	if leverage, ok := data["leverage"]; ok {
		position.LeverRate = leverage.(float64)
	}
	if liquPrice, ok := data["liquidationPrice"]; ok {
		position.ForceLiquPrice = liquPrice.(float64)
	}
	if oob, ok := data["openOrderBuyQty"]; ok {
		position.OpenOrderBuyQty = oob.(float64)
	}
	if oos, ok := data["openOrderSellQty"]; ok {
		position.OpenOrderSellQty = oos.(float64)
	}
	if currentQty, ok := data["currentQty"]; ok {
		qty := currentQty.(float64)
		if qty > 0 {
			position.SellAmount = 0
			position.SellPriceCost = 0
			position.SellPriceAvg = 0
			position.SellProfitReal = 0
			position.SellAvailable = 0
			position.BuyAmount = qty
			if acg, ok := data["avgCostPrice"]; ok {
				position.BuyPriceCost = acg.(float64)
			}
			if aep, ok := data["avgEntryPrice"]; ok {
				position.BuyPriceAvg = aep.(float64)
			}
			if urp, ok := data["unrealisedPnlPcnt"]; ok {
				position.BuyProfitReal = urp.(float64)
			}
			position.BuyAvailable = position.BuyAmount - position.OpenOrderBuyQty
		} else {
			position.BuyAmount = 0
			position.BuyPriceCost = 0
			position.BuyPriceAvg = 0
			position.BuyProfitReal = 0
			position.BuyAvailable = 0
			position.SellAmount = -qty
			if acg, ok := data["avgCostPrice"]; ok {
				position.SellPriceCost = acg.(float64)
			}
			if aep, ok := data["avgEntryPrice"]; ok {
				position.SellPriceAvg = aep.(float64)
			}
			if urp, ok := data["unrealisedPnlPcnt"]; ok {
				position.SellProfitReal = urp.(float64)
			}
			position.SellAvailable = position.SellAmount - position.OpenOrderSellQty
		}
	}
	atomic.StoreInt32(&pos.isUpdate, 1)
}

func (pos *position) deletePosition(symbol string, data map[string]interface{}) {
	index, position := pos.findPositionItemByKeys(symbol, data)
	if position == nil {
		return
	}
	posList := pos.getPositionList(symbol)
	posList = append(posList[:index], posList[index+1:]...)
	pos.positionData.Set(symbol, posList)
	atomic.StoreInt32(&pos.isUpdate, 1)
}

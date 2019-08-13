package bitmex

import (
	"sync/atomic"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/model"
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
		pos.positionData.Set(symbol, make([]*model.Position, 0, dataLength))
	}
	return pos
}

func (pos *position) getPositionList(symbol string) []*model.Position {
	data, _ := pos.positionData.Get(symbol)
	posList := data.([]*model.Position)
	return posList
}

func (pos *position) makePosition(data map[string]interface{}) *model.Position {
	p := new(model.Position)
	p.Symbol = data["symbol"].(string)
	p.Account = data["account"].(float64)
	p.Currency = data["currency"].(string)
	if value, ok := data["leverage"]; ok && value != nil {
		p.LeverRate = value.(float64)
	}
	if value, ok := data["liquidationPrice"]; ok && value != nil {
		p.ForceLiquPrice = value.(float64)
	}
	if value, ok := data["openOrderBuyQty"]; ok && value != nil {
		p.OpenOrderBuyQty = value.(float64)
	}
	if value, ok := data["openOrderSellQty"]; ok && value != nil {
		p.OpenOrderSellQty = value.(float64)
	}

	currQry := data["currentQty"].(float64)
	if currQry > 0 {
		// hold long position
		p.BuyAmount = currQry
		if acg, ok := data["avgCostPrice"]; ok && acg != nil {
			p.BuyPriceCost = acg.(float64)
		}
		if aep, ok := data["avgEntryPrice"]; ok && aep != nil {
			p.BuyPriceAvg = aep.(float64)
		}
		if urp, ok := data["unrealisedPnlPcnt"]; ok && urp != nil {
			p.BuyProfitReal = urp.(float64)
		}
		p.BuyAvailable = p.BuyAmount - p.OpenOrderBuyQty
	} else {
		p.SellAmount = -currQry
		if acg, ok := data["avgCostPrice"]; ok && acg != nil {
			p.SellPriceCost = acg.(float64)
		}
		if aep, ok := data["avgEntryPrice"]; ok && aep != nil {
			p.SellPriceAvg = aep.(float64)
		}
		if urp, ok := data["unrealisedPnlPcnt"]; ok && urp != nil {
			p.SellProfitReal = urp.(float64)
		}
		p.SellAvailable = p.SellAmount - p.OpenOrderSellQty
	}
	return p
}

func (pos *position) insertPosition(symbol string, position *model.Position) {
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
	symbol string, updateData map[string]interface{}) (int, *model.Position) {
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

	if leverage, ok := data["leverage"]; ok && leverage != nil {
		position.LeverRate = leverage.(float64)
	}
	if liquPrice, ok := data["liquidationPrice"]; ok && liquPrice != nil {
		position.ForceLiquPrice = liquPrice.(float64)
	}
	if oob, ok := data["openOrderBuyQty"]; ok && oob != nil {
		position.OpenOrderBuyQty = oob.(float64)
	}
	if oos, ok := data["openOrderSellQty"]; ok && oos != nil {
		position.OpenOrderSellQty = oos.(float64)
	}
	if currentQty, ok := data["currentQty"]; ok && currentQty != nil {
		qty := currentQty.(float64)
		if qty > 0 {
			position.SellAmount = 0
			position.SellPriceCost = 0
			position.SellPriceAvg = 0
			position.SellProfitReal = 0
			position.SellAvailable = 0
			position.BuyAmount = qty
			if acg, ok := data["avgCostPrice"]; ok && acg != nil {
				position.BuyPriceCost = acg.(float64)
			}
			if aep, ok := data["avgEntryPrice"]; ok && aep != nil {
				position.BuyPriceAvg = aep.(float64)
			}
			if urp, ok := data["unrealisedPnlPcnt"]; ok && urp != nil {
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
			if acg, ok := data["avgCostPrice"]; ok && acg != nil {
				position.SellPriceCost = acg.(float64)
			}
			if aep, ok := data["avgEntryPrice"]; ok && aep != nil {
				position.SellPriceAvg = aep.(float64)
			}
			if urp, ok := data["unrealisedPnlPcnt"]; ok && urp != nil {
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

package bitmex

import (
	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) makePosition(data []map[string]interface{}) []*core.Position {
	resp := make([]*core.Position, 0, len(data))

	for _, item := range data {
		pos := new(core.Position)
		pos.Symbol = item["symbol"].(string)
		pos.Account = item["account"].(float32)
		pos.Currency = item["currency"].(string)
		pos.LeverRate = item["leverage"].(float64)
		pos.ForceLiquPrice = item["liquidationPrice"].(float64)
		pos.OpenOrderBuyQty = item["openOrderBuyQty"].(float64)
		pos.OpenOrderSellQty = item["openOrderSellQty"].(float64)

		currQry := item["currentQty"].(float64)
		if currQry > 0 {
			// hold long position
			pos.BuyAmount = currQry
			pos.BuyPriceCost = item["avgCostPrice"].(float64)
			pos.BuyPriceAvg = item["avgEntryPrice"].(float64)
			pos.BuyProfitReal = item["unrealisedPnlPcnt"].(float64)
			pos.BuyAvailable = pos.BuyAmount - pos.OpenOrderBuyQty
		} else {
			pos.SellAmount = -currQry
			pos.SellPriceCost = item["avgCostPrice"].(float64)
			pos.SellPriceAvg = item["avgEntryPrice"].(float64)
			pos.SellProfitReal = item["unrealisedPnlPcnt"].(float64)
			pos.SellAvailable = pos.SellAmount - pos.OpenOrderSellQty
		}
		resp = append(resp, pos)
	}
	return resp
}

func (bm *Bitmex) insertPositionList(symbol string, PositionList []*core.Position) {
	updateLength := len(PositionList)
	length = len(bm.PositionData[symbol])
	if length+updateLength >= dataLength {
		bm.PositionData[symbol] = bm.PositionData[symbol][length+updateLength-dataLength:]
	}
	bm.PositionData[symbol] = append(bm.PositionData[symbol], PositionList...)
}

func (bm *Bitmex) findPositionItemByKeys(
	symbol string, updateData map[string]interface{}) (int, *core.Position) {
	for index, val := range bm.positionData[symbol] {
		if val.Account == updateData["account"].(float32) &&
			val.Symbol == updateData["symbol"].(string) &&
			val.Currency == updateData["currency"].(string) {
			return index, val
		}
	}
	return 0, nil
}

func (bm *Bitmex) updatePosition(pos *core.Position, data map[string]interface{}) {
	if leverage, ok := data["leverage"]; ok {
		pos.Leverage = leverage.(float64)
	}
	if liquPrice, ok := data["liquidationPrice"]; ok {
		pos.ForceLiquPrice = liquPrice.(float64)
	}
	if oob, ok := data["openOrderBuyQty"]; ok {
		pos.OpenOrderBuyQty = oob.(float64)
	}
	if oos, ok := data["openOrderSellQty"]; ok {
		pos.OpenOrderSellQty = oos.(float64)
	}
	if currentQty, ok := data["currentQty"]; ok {
		qty := currentQty.(float64)
		if qty > 0 {
			pos.SellAmount = 0
			pos.SellPriceCost = 0
			pos.SellPriceAvg = 0
			pos.SellProfitReal = 0
			pos.SellAvailable = 0
			pos.BuyAmount = qty
			if acg, ok := data["avgCostPrice"]; ok {
				pos.BuyPriceCost = acg.(float64)
			}
			if aep, ok := data["avgEntryPrice"]; ok {
				pos.BuyPriceAvg = aep.(float64)
			}
			if urp, ok := data["unrealisedPnlPcnt"]; ok {
				pos.BuyProfitReal = urp.(float64)
			}
			pos.BuyAvailable = pos.BuyAmount - pos.OpenOrderBuyQty
		} else {
			pos.BuyAmount = 0
			pos.BuyPriceCost = 0
			pos.BuyPriceAvg = 0
			pos.BuyProfitReal = 0
			pos.BuyAvailable = 0
			pos.SellAmount = -qty
			if acg, ok := data["avgCostPrice"]; ok {
				pos.SellPriceCost = acg.(float64)
			}
			if aep, ok := data["avgEntryPrice"]; ok {
				pos.SellPriceAvg = aep.(float64)
			}
			if urp, ok := data["unrealisedPnlPcnt"]; ok {
				pos.SellProfitReal = urp.(float64)
			}
			pos.SellAvailable = pos.SellAmount - pos.OpenOrderSellQty
		}
	}
}

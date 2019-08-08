package bitmex

import (
	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) GetTick(symbol string) *core.Tick {
	tickList := bm.tickData.getTickList(symbol)
	length := len(tickList)
	if length == 0 {
		return nil
	}
	return tickList[length-1]
}

func (bm *Bitmex) GetQuote(symbol string) *core.Quote {
	quoteList := bm.quoteData.getQuoteList(symbol)
	length := len(quoteList)
	if length == 0 {
		return nil
	}
	return quoteList[length-1]
}

func (bm *Bitmex) GetTrade(symbol string) *core.Trade {
	tradeList := bm.tradeData.getTradeList(symbol)
	length := len(tradeList)
	if length == 0 {
		return nil
	}
	return tradeList[length-1]
}

func (bm *Bitmex) GetOrders(symbol string) []*core.Order {
	orderList := bm.orderData.getOrderList(symbol)
	resp := make([]*core.Order, 0, len(orderList))
	// filter orders
	for _, order := range orderList {
		if !bm.orderData.needDeleteOrder(order) {
			resp = append(resp, order)
		}
	}
	return resp
}

func (bm *Bitmex) GetPosition(symbol string) []*core.Position {
	positionList := bm.positionData.getPositionList(symbol)
	return positionList
}

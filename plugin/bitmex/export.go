package bitmex

import (
	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) GetTick(symbol string) (*core.Tick, error) {
	data, ok := bm.tickData.Get(symbol)
	if !ok {
		return nil, SymbolErr
	}
	tickList := data.([]*core.Tick)
	length := len(tickList)
	if length == 0 {
		return nil, nil
	}
	return tickList[length-1], nil
}

func (bm *Bitmex) GetQuote(symbol string) (*core.Quote, error) {
	data, ok := bm.quoteData.Get(symbol)
	if !ok {
		return nil, SymbolErr
	}
	quoteList := data.([]*core.Quote)
	length := len(quoteList)
	if length == 0 {
		return nil, nil
	}
	return quoteList[length-1], nil
}

func (bm *Bitmex) GetTrade(symbol string) (*core.Trade, error) {
	data, ok := bm.tradeData.Get(symbol)
	if !ok {
		return nil, SymbolErr
	}
	tradeList := data.([]*core.Trade)
	length := len(tradeList)
	if length == 0 {
		return nil, nil
	}
	return tradeList[length-1], nil
}

func (bm *Bitmex) GetOrders(symbol string) ([]*core.Order, error) {
	data, ok := bm.orderData.Get(symbol)
	if !ok {
		return nil, SymbolErr
	}
	orderList := data.([]*core.Order)
	resp := make([]*core.Order, 0, len(orderList))
	// filter orders
	for _, order := range orderList {
		if !bm.needDeleteOrder(order) {
			resp = append(resp, order)
		}
	}
	return resp, nil
}

func (bm *Bitmex) GetPosition(symbol string) ([]*core.Position, error) {
	data, ok := bm.positionData.Get(symbol)
	if !ok {
		return nil, SymbolErr
	}
	positionList := data.([]*core.Position)
	return positionList, nil
}

package bitmex

import (
	"github.com/EricQAQ/Traed/core"
)

func (bm *Bitmex) GetTick(symbol string) (*core.Tick, error) {
	tickDataList, ok := bm.tickData[symbol]
	if !ok {
		return nil, SymbolErr
	}
	length := len(tickDataList)
	if length == 0 {
		return nil, nil
	}
	return tickDataList[length-1], nil
}

func (bm *Bitmex) GetQuote(symbol string) (*core.Quote, error) {
	quoteList, ok := bm.quoteData[symbol]
	if !ok {
		return nil, SymbolErr
	}
	length := len(quoteList)
	if length == 0 {
		return nil, nil
	}
	return quoteList[length-1], nil
}

func (bm *Bitmex) GetTrade(symbol string) (*core.Trade, error) {
	tradeList, ok := bm.tradeData[symbol]
	if !ok {
		return nil, SymbolErr
	}
	length := len(tradeList)
	if length == 0 {
		return nil, nil
	}
	return tradeList[length-1], nil
}

func (bm *Bitmex) GetOrders(symbol string) ([]*core.Order, error) {
	orderList, ok := bm.orderData[symbol]
	if !ok {
		return nil, SymbolErr
	}
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
	positionList, ok := bm.positionData[symbol]
	if !ok {
		return nil, SymbolErr
	}
	return positionList, nil
}

package bitmex

import (
	"sync/atomic"

	"github.com/EricQAQ/Traed/model"
)

func (bm *Bitmex) GetTick(symbol string) (*model.Tick, bool) {
	tickList := bm.tickData.GetTickerList(symbol)
	length := len(tickList)
	if length == 0 {
		return nil, false
	}
	isUpdate := bm.tickData.IsUpdate()
	return tickList[length-1], isUpdate
}

func (bm *Bitmex) GetQuote(symbol string) (*model.Quote, bool) {
	quoteList := bm.quoteData.getQuoteList(symbol)
	length := len(quoteList)
	if length == 0 {
		return nil, false
	}
	isUpdate := atomic.CompareAndSwapInt32(&bm.quoteData.isUpdate, 1, 0)
	return quoteList[length-1], isUpdate
}

func (bm *Bitmex) GetTrade(symbol string) (*model.Trade, bool) {
	tradeList := bm.tradeData.getTradeList(symbol)
	length := len(tradeList)
	if length == 0 {
		return nil, false
	}
	isUpdate := atomic.CompareAndSwapInt32(&bm.tradeData.isUpdate, 1, 0)
	return tradeList[length-1], isUpdate
}

func (bm *Bitmex) GetOrders(symbol string) ([]*model.Order, bool) {
	orderList := bm.orderData.getOrderList(symbol)
	resp := make([]*model.Order, 0, len(orderList))
	// filter orders
	for _, order := range orderList {
		if !bm.orderData.needDeleteOrder(order) {
			resp = append(resp, order)
		}
	}
	isUpdate := atomic.CompareAndSwapInt32(&bm.orderData.isUpdate, 1, 0)
	return resp, isUpdate
}

func (bm *Bitmex) GetPosition(symbol string) ([]*model.Position, bool) {
	positionList := bm.positionData.getPositionList(symbol)
	isUpdate := atomic.CompareAndSwapInt32(&bm.positionData.isUpdate, 1, 0)
	return positionList, isUpdate
}

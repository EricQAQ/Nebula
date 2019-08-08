package bitmex

import (
	"sync/atomic"
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Traed/core"
)

type order struct {
	orderKeys map[string][]string
	orderData cmap.ConcurrentMap
	isUpdate  int32
}

func newOrder(symbols []string) *order {
	ord := new(order)
	ord.orderKeys = make(map[string][]string)
	ord.orderData = cmap.New()
	ord.isUpdate = 0
	for _, symbol := range symbols {
		ord.orderKeys[symbol] = wsOrderKeys
		ord.orderData.Set(symbol, make([]*core.Order, 0, dataLength))
	}
	return ord
}

func (ord *order) getOrderList(symbol string) []*core.Order {
	data, _ := ord.orderData.Get(symbol)
	orderList := data.([]*core.Order)
	return orderList
}

func (ord *order) makeOrder(data map[string]interface{}) *core.Order {
	order := new(core.Order)
	order.Symbol = data["symbol"].(string)
	order.OrdStatus = data["ordStatus"].(string)
	order.Timestamp, _ = time.Parse(time.RFC3339, data["timestamp"].(string))
	order.Price = data["price"].(float64)
	order.Amount = data["orderQty"].(float32)
	order.OrdType = data["ordType"].(string)
	order.Side = data["side"].(string)
	order.OrderID = data["orderID"].(string)
	order.ClOrdID = data["clOrdID"].(string)
	order.AvgPrice = data["avgPx"].(float64)
	if filledAmount, ok := data["cumQty"].(float32); ok {
		order.FilledAmount = filledAmount
	} else {
		order.FilledAmount = 0.0
	}
	return order
}

func (ord *order) insertOrder(symbol string, order *core.Order) {
	orderList := ord.getOrderList(symbol)
	length := len(orderList)
	if length >= dataLength {
		orderList = orderList[length-dataLength:]
	}
	orderList = append(orderList, order)
	ord.orderData.Set(symbol, orderList)
	atomic.StoreInt32(&ord.isUpdate, 1)
}

func (ord *order) findOrderItemByKeys(
	symbol string, updateData map[string]interface{}) (int, *core.Order) {
	orderList := ord.getOrderList(symbol)
	for index, val := range orderList {
		if val.OrderID == updateData["orderID"].(string) {
			return index, val
		}
	}
	return 0, nil
}

func (ord *order) updateOrder(symbol string, data map[string]interface{}) {
	index, order := ord.findOrderItemByKeys(symbol, data)
	if order == nil {
		return
	}
	for name, value := range data {
		if name == "ordStatus" {
			order.OrdStatus = value.(string)
		} else if name == "price" {
			order.Price = value.(float64)
		} else if name == "orderQty" {
			order.Amount = value.(float32)
		} else if name == "ordType" {
			order.OrdType = value.(string)
		} else if name == "side" {
			order.Side = value.(string)
		} else if name == "avgPx" {
			order.AvgPrice = value.(float64)
		} else if name == "filledAmount" {
			order.FilledAmount = value.(float32)
		} else if name == "timestamp" {
			order.Timestamp, _ = time.Parse(time.RFC3339, value.(string))
		}
	}
	// Remove cancelled / filled orders
	ord.cleanOrder(index, order)
	atomic.StoreInt32(&ord.isUpdate, 1)
}

func (ord *order) deleteOrder(symbol string, data map[string]interface{}) {
	index, order := ord.findOrderItemByKeys(symbol, data)
	if order == nil {
		return
	}
	ordList := ord.getOrderList(symbol)
	ordList = append(ordList[:index], ordList[index+1:]...)
	ord.orderData.Set(symbol, ordList)
	atomic.StoreInt32(&ord.isUpdate, 1)
}

func (ord *order) needDeleteOrder(order *core.Order) bool {
	if orderStatusMap[order.OrdStatus] != "open" {
		return true
	}
	return false
}

func (ord *order) cleanOrder(index int, order *core.Order) {
	if ord.needDeleteOrder(order) {
		orderList := ord.getOrderList(order.Symbol)
		orderList = append(orderList[:index], orderList[index+1:]...)
		ord.orderData.Set(order.Symbol, orderList)
	}
}

func (bm *Bitmex) LimitBuy(
	symbol string, price float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideBuy, core.OrdLimit,
		price, 0, quantity)
}

func (bm *Bitmex) LimitSell(
	symbol string, price float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.OrdLimit,
		price, 0, quantity)
}

func (bm *Bitmex) MarketBuy(
	symbol string, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideBuy, core.OrdMarket,
		0, 0, quantity)
}

func (bm *Bitmex) MarketSell(
	symbol string, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.OrdMarket,
		0, 0, quantity)
}

func (bm *Bitmex) LimitStopBuy(
	symbol string, stopPx float64,
	price float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.OrdStopLimit,
		price, stopPx, quantity)
}

func (bm *Bitmex) LimitStopSell(
	symbol string, stopPx float64,
	price float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideBuy, core.OrdStopLimit,
		price, stopPx, quantity)
}

func (bm *Bitmex) MarketStopBuy(
	symbol string, stopPx float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.OrdStop,
		0, stopPx, quantity)
}

func (bm *Bitmex) MarketStopSell(
	symbol string, stopPx float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideBuy, core.OrdStop,
		0, stopPx, quantity)
}

func (bm *Bitmex) LimitIfTouchedBuy(
	symbol string, stopPx float64, price float64,
	quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.LimitIfTouched,
		price, stopPx, quantity)
}

func (bm *Bitmex) LimitIfTouchedSell(
	symbol string, stopPx float64, price float64,
	quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideBuy, core.LimitIfTouched,
		price, stopPx, quantity)
}

func (bm *Bitmex) MarketIfTouchedBuy(
	symbol string, stopPx float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.MarketIfTouched,
		0, stopPx, quantity)
}

func (bm *Bitmex) MarketIfTouchedSell(
	symbol string, stopPx float64, quantity float32) (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideBuy, core.MarketIfTouched,
		0, stopPx, quantity)
}

func (bm *Bitmex) create_order(
	symbol, side, ordType string,
	price, stopPx float64, quantity float32) (*core.Order, error) {
	data := make(map[string]interface{})
	data["symbol"] = symbol
	data["side"] = side
	data["ordType"] = ordType
	data["orderQty"] = quantity

	if ordType == core.OrdLimit || ordType == core.OrdStopLimit ||
		ordType == core.LimitIfTouched {
		data["price"] = price
	}

	if ordType != core.OrdLimit && ordType != core.OrdMarket {
		data["stopPx"] = stopPx
	}

	resp, err := bm.doAuthRequest("POST", "/order", data, 1, bm.makeSingleOrder)
	if err != nil {
		return nil, err
	}
	return resp.(*core.Order), nil
}

// func (bm *Bitmex) edit_order() {
// }

// func (bm *Bitmex) delete_order() {
// }

func (bm *Bitmex) makeSingleOrder(data map[string]interface{}) interface{} {
	return bm.orderData.makeOrder(data)
}

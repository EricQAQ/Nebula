package bitmex

import (
	"sync/atomic"
	"time"

	"github.com/orcaman/concurrent-map"

	"github.com/EricQAQ/Nebula/model"
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
		ord.orderData.Set(symbol, make([]*model.Order, 0, dataLength))
	}
	return ord
}

func (ord *order) getOrderList(symbol string) []*model.Order {
	data, _ := ord.orderData.Get(symbol)
	orderList := data.([]*model.Order)
	return orderList
}

func (ord *order) makeOrder(data map[string]interface{}) *model.Order {
	order := new(model.Order)
	order.Symbol = data["symbol"].(string)
	order.OrdStatus = data["ordStatus"].(string)
	if value, ok := data["price"]; ok && value != nil {
		order.Price = value.(float64)
	}
	if value, ok := data["orderQty"]; ok && value != nil {
		order.Amount = value.(float64)
	}
	if value, ok := data["orderType"]; ok && value != nil {
		order.OrdType = value.(string)
	}
	if value, ok := data["side"]; ok && value != nil {
		order.Side = value.(string)
	}
	if value, ok := data["orderID"]; ok && value != nil {
		order.OrderID = value.(string)
	}
	if value, ok := data["clOrdID"]; ok && value != nil {
		order.ClOrdID = value.(string)
	}
	if value, ok := data["avgPx"]; ok && value != nil {
		order.AvgPrice = value.(float64)
	}
	if filledAmount, ok := data["cumQty"].(float64); ok {
		order.FilledAmount = filledAmount
	} else {
		order.FilledAmount = 0.0
	}

	loc, _ := time.LoadLocation("Asia/Chongqing")
	ts, _ := time.Parse(time.RFC3339, data["timestamp"].(string))
	order.Timestamp = ts.In(loc)
	return order
}

func (ord *order) insertOrder(symbol string, order *model.Order) {
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
	symbol string, updateData map[string]interface{}) (int, *model.Order) {
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
			order.Amount = value.(float64)
		} else if name == "ordType" {
			order.OrdType = value.(string)
		} else if name == "side" {
			order.Side = value.(string)
		} else if name == "avgPx" {
			order.AvgPrice = value.(float64)
		} else if name == "filledAmount" {
			order.FilledAmount = value.(float64)
		} else if name == "timestamp" {
			loc, _ := time.LoadLocation("Asia/Chongqing")
			ts, _ := time.Parse(time.RFC3339, value.(string))
			order.Timestamp = ts.In(loc)
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

func (ord *order) needDeleteOrder(order *model.Order) bool {
	if orderStatusMap[order.OrdStatus] != "open" {
		return true
	}
	return false
}

func (ord *order) cleanOrder(index int, order *model.Order) {
	if ord.needDeleteOrder(order) {
		orderList := ord.getOrderList(order.Symbol)
		orderList = append(orderList[:index], orderList[index+1:]...)
		ord.orderData.Set(order.Symbol, orderList)
	}
}

func (bm *Bitmex) LimitBuy(
	symbol string, price float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideBuy, model.OrdLimit,
		price, 0, quantity)
}

func (bm *Bitmex) LimitSell(
	symbol string, price float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideSell, model.OrdLimit,
		price, 0, quantity)
}

func (bm *Bitmex) MarketBuy(
	symbol string, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideBuy, model.OrdMarket,
		0, 0, quantity)
}

func (bm *Bitmex) MarketSell(
	symbol string, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideSell, model.OrdMarket,
		0, 0, quantity)
}

func (bm *Bitmex) LimitStopBuy(
	symbol string, stopPx float64,
	price float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideSell, model.OrdStopLimit,
		price, stopPx, quantity)
}

func (bm *Bitmex) LimitStopSell(
	symbol string, stopPx float64,
	price float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideBuy, model.OrdStopLimit,
		price, stopPx, quantity)
}

func (bm *Bitmex) MarketStopBuy(
	symbol string, stopPx float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideSell, model.OrdStop,
		0, stopPx, quantity)
}

func (bm *Bitmex) MarketStopSell(
	symbol string, stopPx float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideBuy, model.OrdStop,
		0, stopPx, quantity)
}

func (bm *Bitmex) LimitIfTouchedBuy(
	symbol string, stopPx float64, price float64,
	quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideSell, model.LimitIfTouched,
		price, stopPx, quantity)
}

func (bm *Bitmex) LimitIfTouchedSell(
	symbol string, stopPx float64, price float64,
	quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideBuy, model.LimitIfTouched,
		price, stopPx, quantity)
}

func (bm *Bitmex) MarketIfTouchedBuy(
	symbol string, stopPx float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideSell, model.MarketIfTouched,
		0, stopPx, quantity)
}

func (bm *Bitmex) MarketIfTouchedSell(
	symbol string, stopPx float64, quantity float32) (*model.Order, error) {
	return bm.create_order(
		symbol, model.SideBuy, model.MarketIfTouched,
		0, stopPx, quantity)
}

func (bm *Bitmex) create_order(
	symbol, side, ordType string,
	price, stopPx float64, quantity float32) (*model.Order, error) {
	data := make(map[string]interface{})
	data["symbol"] = symbol
	data["side"] = side
	data["ordType"] = ordType
	data["orderQty"] = quantity

	if ordType == model.OrdLimit || ordType == model.OrdStopLimit ||
		ordType == model.LimitIfTouched {
		data["price"] = price
	}

	if ordType != model.OrdLimit && ordType != model.OrdMarket {
		data["stopPx"] = stopPx
	}

	resp, err := bm.doAuthRequest("POST", "/order", data, 1, bm.makeSingleOrder)
	if err != nil {
		return nil, err
	}
	return resp.(*model.Order), nil
}

// func (bm *Bitmex) edit_order() {
// }

// func (bm *Bitmex) delete_order() {
// }

func (bm *Bitmex) makeSingleOrder(data map[string]interface{}) interface{} {
	return bm.orderData.makeOrder(data)
}

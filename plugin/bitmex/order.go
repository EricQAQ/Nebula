package bitmex

import (
	"time"

	"github.com/EricQAQ/Traed/core"
)

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

func (bm *Bitmex) LimitStop(
	symbol string, stopPx float64,
	price float64, quantity float32) (*core.Order, error) {
}

func (bm *Bitmex) MarketStop(
	symbol string, stopPx float64, quantity float32) (*core.Order, error) {
}

func (bm *Bitmex) LimitIfTouched(
	symbol string, stopPx float64, price float64,
	quantity float32)  (*core.Order, error) {
}

func (bm *Bitmex) MarketIfTouched(
	symbol string, stopPx float64, quantity float32)  (*core.Order, error) {
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

	resp, err := bm.doAuthRequest("POST", "/order", data, 1, makeOrder)
	if err != nil {
		return nil, err
	}
	return resp.(*core.Order), nil
}

// func (bm *Bitmex) edit_order() {
// }

// func (bm *Bitmex) delete_order() {
// }

func (bm *Bitmex) makeOrder(data []map[string]interface{}) []*core.Order {
	resp := make([]*core.Order, 0, len(data))
	for _, item := range data {
		order := new(core.Order)
		order.Symbol = item["symbol"].(string)
		order.OrdStatus = item["ordStatus"].(string)
		order.Timestamp = time.Parse(time.RFC3339, item["timestamp"].(string))
		order.Price = item["price"].(float64)
		order.Amount = item["orderQty"].(float32)
		order.OrdType = item["ordType"].(string)
		order.Side = item["side"].(string)
		order.OrderID = item["orderID"].(string)
		order.ClOrderID = item["clOrdID"].(string)
		order.AvgPrice = item["avgPx"].(float64)
		if filledAmount, ok := item["cumQty"].(float32); ok {
			order.FilledAmount = filledAmount
		} else {
			order.FilledAmount = 0.0
		}
		resp = append(resp, order)
	}
	return resp
}

func (bm *Bitmex) insertOrderList(symbol string, orderList []*core.Order) {
	updateLength := len(orderList)
	length = len(bm.orderData[symbol])
	if length+updateLength >= dataLength {
		bm.orderData[symbol] = bm.orderData[symbol][length+updateLength-dataLength:]
	}
	bm.orderData[symbol] = append(bm.orderData[symbol], orderList...)
}

func (bm *Bitmex) findOrderItemByKeys(
	symbol string, updateData map[string]interface{}) (int, *core.Order) {
	for index, val := range bm.orderData[symbol] {
		if val.OrderID == updateData["orderID"].(string) {
			return index, val
		}
	}
	return 0, nil
}

func (bm *Bitmex) updateOrder(order *core.Order, data map[string]interface{}) {
	for name, value := range data {
		if name == "ordStatus" {
			order.OrderStatus = value.(string)
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
			order.Timestamp = time.Parse(time.RFC3339, value.(string))
		}
	}
}

func (bm *Bitmex) cleanOrder(index int, order *core.Order) {
	if orderStatusMap[order.OrderStatus] != "open" {
		bm.orderData = append(bm.orderData[:index], bm.orderData[index+1:]...)
	}
}

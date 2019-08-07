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
	quantity float32)  (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.LimitIfTouched,
		price, stopPx, quantity)
}

func (bm *Bitmex) LimitIfTouchedSell(
	symbol string, stopPx float64, price float64,
	quantity float32)  (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideBuy, core.LimitIfTouched,
		price, stopPx, quantity)
}

func (bm *Bitmex) MarketIfTouchedBuy(
	symbol string, stopPx float64, quantity float32)  (*core.Order, error) {
	return bm.create_order(
		symbol, core.SideSell, core.MarketIfTouched,
		0, stopPx, quantity)
}

func (bm *Bitmex) MarketIfTouchedSell(
	symbol string, stopPx float64, quantity float32)  (*core.Order, error) {
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

func (bm *Bitmex) makeOrder(data []map[string]interface{}) []*core.Order {
	resp := make([]*core.Order, 0, len(data))
	for _, item := range data {
		resp = append(resp, bm.makeSingleOrder(item).(*core.Order))
	}
	return resp
}

func (bm *Bitmex) insertOrderList(symbol string, orderList []*core.Order) {
	updateLength := len(orderList)
	length := len(bm.orderData[symbol])
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
}

func (bm *Bitmex) cleanOrder(index int, order *core.Order) {
	if orderStatusMap[order.OrdStatus] != "open" {
		bm.orderData[order.Symbol] = append(
			bm.orderData[order.Symbol][:index],
			bm.orderData[order.Symbol][index+1:]...)
	}
}

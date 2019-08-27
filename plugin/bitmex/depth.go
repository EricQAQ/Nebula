package bitmex

import (
	"sort"

	"github.com/EricQAQ/Nebula/model"
	"github.com/orcaman/concurrent-map"
)

type depth struct {
	depthKeys map[string][]string
	depthData cmap.ConcurrentMap
	updateMap cmap.ConcurrentMap
}

func newDepth(symbols []string) *depth {
	ob := new(depth)
	ob.depthKeys = make(map[string][]string)
	ob.depthData = cmap.New()
	ob.updateMap = cmap.New()
	for _, symbol := range symbols {
		depth := new(model.Depth)
		depth.Sell = make([]*model.DepthRecord, 0, 25)
		depth.Buy = make([]*model.DepthRecord, 0, 25)
		ob.depthKeys[symbol] = wsOrderBookKeys
		ob.depthData.Set(symbol, depth)
		ob.updateMap.Set(symbol, false)
	}
	return ob
}

func (ob *depth) isUpdate(symbol string) bool {
	flag, _ := ob.updateMap.Get(symbol)
	ob.updateMap.Set(symbol, false)
	return flag.(bool)
}

func (ob *depth) getDepth(symbol string) *model.Depth {
	data, _ := ob.depthData.Get(symbol)
	depth := data.(*model.Depth)
	return depth
}

func (ob *depth) insertDepthRecord(symbol string, dr *model.DepthRecord) {
	depth := ob.getDepth(symbol)
	if dr.Side == model.SideBuy {
		length := len(depth.Buy)
		if length >= 25 {
			depth.Buy = depth.Buy[1:]
		}
		depth.Buy = append(depth.Buy, dr)
		sort.Sort(model.DepthRecordList(depth.Buy))
		ob.depthData.Set(symbol, depth)
	} else {
		length := len(depth.Sell)
		if length >= 25 {
			depth.Sell = depth.Sell[1:]
		}
		depth.Sell = append(depth.Sell, dr)
		sort.Sort(model.DepthRecordList(depth.Sell))
		ob.depthData.Set(symbol, depth)
	}
	ob.updateMap.Set(symbol, true)
}

func (ob *depth) searchDepthRecord(
	depth *model.Depth, data map[string]interface{}) (int, *model.DepthRecord) {
	side := data["side"].(string)
	id := data["id"].(float64)
	symbol := data["symbol"].(string)

	var dlist []*model.DepthRecord
	if side == model.SideBuy {
		dlist = depth.Buy
	} else {
		dlist = depth.Sell
	}

	for index, d := range dlist {
		if d.ID == id && symbol == d.Symbol {
			return index, d
		}
	}
	return 0, nil
}

func (ob *depth) updateDepthRecord(data map[string]interface{}) {
	symbol := data["symbol"].(string)
	depth := ob.getDepth(symbol)
	_, dr := ob.searchDepthRecord(depth, data)
	if dr == nil {
		return
	}
	for name, value := range data {
		if name == "price" {
			dr.Price = value.(float64)
		} else if name == "size" {
			dr.Amount = value.(float64)
		}
	}
	ob.updateMap.Set(symbol, true)
}

func (ob *depth) deleteDepthRecord(data map[string]interface{}) {
	symbol := data["symbol"].(string)
	side := data["side"].(string)
	depth := ob.getDepth(symbol)
	index, dr := ob.searchDepthRecord(depth, data)
	if dr == nil {
		return
	}
	if side == model.SideBuy {
		depth.Buy = append(depth.Buy[:index], depth.Buy[index+1:]...)
	} else {
		depth.Sell = append(depth.Sell[:index], depth.Sell[index+1:]...)
	}
	ob.depthData.Set(symbol, depth)
	ob.updateMap.Set(symbol, true)
}

func (ob *depth) makeDepthRecord(data map[string]interface{}) *model.DepthRecord {
	dr := new(model.DepthRecord)
	dr.ID = data["id"].(float64)
	dr.Symbol = data["symbol"].(string)
	dr.Side = data["side"].(string)
	dr.Price = data["price"].(float64)
	dr.Amount = data["size"].(float64)
	return dr
}

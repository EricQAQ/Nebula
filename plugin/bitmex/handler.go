package bitmex

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/EricQAQ/Traed/core"
)

// There are four possible actions from the WS:
// 'partial' - full table image
// 'insert'  - new row
// 'update'  - update row
// 'delete'  - delete row

func (bm *Bitmex) callByTableName(table, action string, data map[string]interface{}) {
	fnName, exist := reflect.TypeOf(bm).MethodByName("Handle" + strings.Title(table))
	if !exist {
		log.Errorf("Unsupport topic: %s", table)
		return
	}
	args := make([]reflect.Value, 3)
	args[0] = reflect.ValueOf(bm)
	args[1] = reflect.ValueOf(action)
	args[2] = reflect.ValueOf(data)
	fnName.Func.Call(args)
}

func (bm *Bitmex) HandleMessage(data core.ParsedData) {
	if data.Type == core.ErrorMsg {
		log.Infof("Received Error Message: %s", data.Data)
		return
	}
	if data.Type != core.Message {
		return
	}
	table := data.Data["table"].(string)
	action := data.Data["action"].(string)
	bm.callByTableName(table, action, data.Data)
}

func (bm *Bitmex) HandleInstrument(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)

		if action == actionPartial || action == actionInsert {
			tick := bm.makeInstrument(ret)
			bm.insertTick(symbol, tick)
		} else if action == actionUpdate {
			length := len(bm.tickData[symbol])
			if length <= 0 {
				return
			}
			bm.updateTick(bm.tickData[symbol][length-1], ret)
		} else if action == actionDelete {
			length := len(bm.tickData[symbol])
			if length <= 0 {
				return
			}
			bm.tickData[symbol] = bm.tickData[symbol][:length-2]
		}
	}
}

func (bm *Bitmex) HandleTrade(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			trade := bm.makeTrade(ret)
			bm.insertTrade(symbol, trade)
		}
	}
}

func (bm *Bitmex) HandleQuote(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			quote := bm.makeQuote(ret)
			bm.insertQuote(symbol, quote)
		}
	}
}

func (bm *Bitmex) HandleOrder(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			order := bm.makeOrder(ret)
			bm.insertOrder(symbol, order)
		} else if action == actionUpdate {
			index, order := bm.findOrderItemByKeys(symbol, ret)
			if order == nil {
				return
			}
			bm.updateOrder(order, ret)
			// Remove cancelled / filled orders
			bm.cleanOrder(index, order)
		} else if action == actionDelete {
			index, order := bm.findOrderItemByKeys(symbol, ret)
			if order == nil {
				return
			}
			bm.orderData[symbol] = append(
				bm.orderData[symbol][:index], bm.orderData[symbol][index+1:]...)
		}
	}
}

func (bm *Bitmex) HandlePosition(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			pos := bm.makePosition(ret)
			bm.insertPosition(symbol, pos)
		} else if action == actionUpdate {
			_, pos := bm.findPositionItemByKeys(symbol, ret)
			if pos == nil {
				return
			}
			bm.updatePosition(pos, ret)
		} else if action == actionDelete {
			index, pos := bm.findPositionItemByKeys(symbol, ret)
			if pos == nil {
				return
			}
			bm.positionData[symbol] = append(
				bm.positionData[symbol][:index], bm.positionData[symbol][index+1:]...)
		}
	}
}

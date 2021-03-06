package bitmex

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/EricQAQ/Nebula/core"
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

func (bm *Bitmex) HandleTrade(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			trade := makeTrade(ret)
			bm.tradeData.insertTrade(symbol, trade)
			bm.tickData.UpdateTicker(symbol, trade, bm.tickCh)
		}
	}
}

func (bm *Bitmex) HandleQuote(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			quote := bm.quoteData.makeQuote(ret)
			bm.quoteData.insertQuote(symbol, quote)
		}
	}
}

func (bm *Bitmex) HandleOrder(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			order := bm.orderData.makeOrder(ret)
			bm.orderData.insertOrder(symbol, order)
		} else if action == actionUpdate {
			bm.orderData.updateOrder(symbol, ret)
		} else if action == actionDelete {
			bm.orderData.deleteOrder(symbol, ret)
		}
	}
}

func (bm *Bitmex) HandlePosition(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			if ret["currentQty"].(float64) == 0 {
				continue
			}
			pos := bm.positionData.makePosition(ret)
			bm.positionData.insertPosition(symbol, pos)
		} else if action == actionUpdate {
			bm.positionData.updatePosition(symbol, ret)
		} else if action == actionDelete {
			bm.positionData.deletePosition(symbol, ret)
		}
	}
}

func (bm *Bitmex) HandleOrderBookL2_25(action string, data map[string]interface{}) {
	retList := data["data"].([]interface{})
	for _, rv := range retList {
		ret := rv.(map[string]interface{})
		symbol := ret["symbol"].(string)
		if action == actionPartial || action == actionInsert {
			dr := bm.depthData.makeDepthRecord(ret)
			bm.depthData.insertDepthRecord(symbol, dr)
		} else if action == actionUpdate {
			bm.depthData.updateDepthRecord(ret)
		} else if action == actionDelete {
			bm.depthData.deleteDepthRecord(ret)
		}
	}
}

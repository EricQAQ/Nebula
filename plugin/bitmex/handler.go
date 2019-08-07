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
	value := reflect.ValueOf(bm).Elem()
	fnName := value.MethodByName("handle" + strings.Title(table))
	if fnName.IsNil() {
		log.Errorf("Unsupport topic: %s", table)
		return
	}
	args := []reflect.Value{
		reflect.ValueOf(action), reflect.ValueOf(data),
	}
	fnName.Call(args)
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

func (bm *Bitmex) handleInstrument(action string, data map[string]interface{}) {
	symbol := data["symbol"].(string)
	d := data["data"].([]map[string]interface{})

	if action == actionPartial || action == actionInsert {
		tickList := bm.makeInstrument(d)
		bm.insertTickList(symbol, tickList)
	} else if action == actionUpdate {
		for _, item := range d {
			_, tick := bm.findTickItemByKeys(symbol, item)
			if tick == nil {
				return
			}
			bm.updateTick(tick, item)
		}
	} else if action == actionDelete {
		for _, item := range d {
			index, _ := bm.findTickItemByKeys(symbol, item)
			bm.tickData[symbol] = append(
				bm.tickData[symbol][:index], bm.tickData[symbol][index+1:]...)
		}
	}
}

func (bm *Bitmex) handleTrade(action string, data map[string]interface{}) {
	symbol := data["symbol"].(string)
	d := data["data"].([]map[string]interface{})

	if action == actionPartial || action == actionInsert {
		tradeList := bm.makeTrade(d)
		bm.insertTradeList(symbol, tradeList)
	}
}

func (bm *Bitmex) handleQuote(action string, data map[string]interface{}) {
	symbol := data["symbol"].(string)
	d := data["data"].([]map[string]interface{})

	if action == actionPartial || action == actionInsert {
		quoteList := bm.makeQuote(d)
		bm.insertQuoteList(symbol, quoteList)
	}
}

func (bm *Bitmex) handleOrder(action string, data map[string]interface{}) {
	symbol := data["symbol"].(string)
	d := data["data"].([]map[string]interface{})

	if action == actionPartial || action == actionInsert {
		orderList := bm.makeOrder(d)
		bm.insertOrderList(symbol, orderList)
	} else if action == actionUpdate {
		for _, item := range d {
			index, order := bm.findOrderItemByKeys(symbol, item)
			if order == nil {
				return
			}
			bm.updateOrder(order, item)
			// Remove cancelled / filled orders
			bm.cleanOrder(index, order)
		}
	} else if action == actionDelete {
		for _, item := range d {
			index, order := bm.findOrderItemByKeys(symbol, item)
			if order == nil {
				return
			}
			bm.orderData[symbol] = append(
				bm.orderData[symbol][:index], bm.orderData[symbol][index+1:]...)
		}
	}
}

func (bm *Bitmex) handlePosition(action string, data map[string]interface{}) {
	symbol := data["symbol"].(string)
	d := data["data"].([]map[string]interface{})

	if action == actionPartial || action == actionInsert {
		posList := bm.makePosition(d)
		bm.insertPositionList(symbol, posList)
	} else if action == actionUpdate {
		for _, item := range d {
			_, pos := bm.findPositionItemByKeys(symbol, item)
			if pos == nil {
				return
			}
			bm.updatePosition(pos, item)
		}
	} else if action == actionDelete {
		for _, item := range d {
			index, _ := bm.findPositionItemByKeys(symbol, item)
			bm.positionData[symbol] = append(
				bm.positionData[symbol][:index], bm.positionData[symbol][index+1:]...)
		}
	}
}

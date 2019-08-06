package bitmex

import (
	"reflect"
	"strings"
	"time"

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
		reflect.Value(action), reflect.Value(data),
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

	length := len(bm.tickData[symbol])
	if action == actionPartial  {
		tickList := bm.makeInstrument(d)
		bm.insertTickList(symbol, tickList)
		bm.tickKeys[symbol] = append(
			bm.tickKeys[symbol], data["keys"].([]string)...)
		return
	} else if action == actionInsert {
		tickList := bm.makeInstrument(d)
		bm.insertTickList(symbol, tickList)
		return
	} else if action == updatePartial {
		for _, item := range d {
			_, tick := bm.findTickItemByKeys(symbol, d)
			if tick == nil {
				return
			}
			bm.updateTick(tick)
		}
		return
	} else if action == deletePartial {
		for _, item := range d {
			index, tick := bm.findTickItemByKeys(symbol, d)
			bm.tickData[symbol] = append(
				bm.tickData[symbol][:index], bm.tickData[symbol][index+1:]...)
		}
		return
	}
}

func (bm *Bitmex) handleTrade(action string, data []map[string]interface{}) {
}

func (bm *Bitmex) handleQuote(action string, data []map[string]interface{}) {
}

func (bm *Bitmex) handleOrder(action string, data []map[string]interface{}) {
}

func (bm *Bitmex) handlePosition(action string, data []map[string]interface{}) {
}

package bitmex

const (
	actionPartial = "partial"
	actionInsert =  "insert"
	actionUpdate =  "update"
	actionDelete =  "delete"
)

var (
	wsInstrumentKeys = []string{"symbol"}
	wsTradeKeys = []string{}
	wsQuoteKeys = []string{}
)

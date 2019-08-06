package bitmex

const (
	actionPartial = "partial"
	actionInsert =  "insert"
	actionUpdate =  "update"
	actionDelete =  "delete"

	routeUrl = "/api/v1"
)

var (
	wsInstrumentKeys = []string{"symbol"}
	wsTradeKeys = []string{}
	wsQuoteKeys = []string{}
)

package bitmex

const (
	actionPartial = "partial"
	actionInsert  = "insert"
	actionUpdate  = "update"
	actionDelete  = "delete"

	routeUrl = "/api/v1"
)

var (
	wsInstrumentKeys = []string{"symbol"}
	wsTradeKeys      = []string{}
	wsQuoteKeys      = []string{}
	wsOrderKeys      = []string{"orderID"}
	wsPositionKeys   = []string{"account", "symbol", "currency"}

	orderStatusMap = map[string]string{
		"New":             "open",
		"PartiallyFilled": "open",
		"Filled":          "closed",
		"DoneForDay":      "open",
		"Canceled":        "canceled",
		"PendingCancel":   "open",
		"PendingNew":      "open",
		"Rejected":        "rejected",
		"Expired":         "expired",
		"Stopped":         "open",
		"Untriggered":     "open",
		"Triggered":       "open",
	}
)

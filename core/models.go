package core

const (
	OrdLimit        = "Limit"
	OrdMarket       = "Market"
	OrdStop         = "Stop"
	OrdStopLimit    = "StopLimit"
	LimitIfTouched  = "LimitIfTouched"
	MarketIfTouched = "MarketIfTouched"
)

const (
	SideBuy  = "Buy"
	SideSell = "Sell"
)

type Order struct {
	Symbol    string
	OrderID   string
	ClOrdID   string
	Price     float64
	Side      string
	Qty       float64
	OrdType   string
	OrdStatus string
}

type Trade struct {
	Symbol    string
	Side      string
	Size      int64
	Price     float64
	Timestamp time.Time
}

type Account struct {
}

type Ticker struct {
}

type Position struct {
}

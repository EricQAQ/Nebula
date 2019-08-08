package core

const (
	WelcomeMsg = iota
	AuthMsg
	SubscribeMsg
	ErrorMsg
	Message
)

type ParsedData struct {
	Type int
	Data map[string]interface{}
}

type WsAuthSubscribeHandler interface {
	GetOperateArgs() []string
	Serialize() ([]byte, error)
}

type OrderAPI interface {
	LimitBuy(symbol string, price float64, quantity float32) (*Order, error)
	LimitSell(symbol string, price float64, quantity float32) (*Order, error)
	MarketBuy(symbol string, quantity float32) (*Order, error)
	MarketSell(symbol string, quantity float32) (*Order, error)
	LimitStopBuy(symbol string, stopPx float64, price float64, quantity float32) (*Order, error)
	LimitStopSell(symbol string, stopPx float64, price float64, quantity float32) (*Order, error)
	MarketStopBuy(symbol string, stopPx float64, quantity float32) (*Order, error)
	MarketStopSell(symbol string, stopPx float64, quantity float32) (*Order, error)
	LimitIfTouchedBuy(symbol string, stopPx float64, price float64, quantity float32) (*Order, error)
	LimitIfTouchedSell(symbol string, stopPx float64, price float64, quantity float32) (*Order, error)
}

type ExportAPI interface {
	GetTick(symbol string) *Tick
	GetQuote(symbol string) *Quote
	GetTrade(symbol string) *Trade
	GetOrders(symbol string) []*Order
	GetPosition(symbol string) []*Position
}

type ExchangeAPI interface {
	GetWsAuthHandler() WsAuthSubscribeHandler
	GetWsSubscribeHandler() WsAuthSubscribeHandler
	Parse(data []byte) (*ParsedData, error)
	HandleMessage(ParsedData)

	OrderAPI
	ExportAPI
}

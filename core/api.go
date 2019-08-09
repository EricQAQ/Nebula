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
	GetTick(symbol string) (*Tick, bool)
	GetQuote(symbol string) (*Quote, bool)
	GetTrade(symbol string) (*Trade, bool)
	GetOrders(symbol string) ([]*Order, bool)
	GetPosition(symbol string) ([]*Position, bool)
}

type ExchangeAPI interface {
	GetExchangeName() string
	GetWsAuthHandler() WsAuthSubscribeHandler
	GetWsSubscribeHandler() WsAuthSubscribeHandler
	Parse(data []byte) (*ParsedData, error)
	HandleMessage(ParsedData)
	GetTickChan() chan Tick

	OrderAPI
	ExportAPI
}

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
	LimitBuy(symbol string, price float64, quantity float32) (*core.Order, error)
	LimitSell(symbol string, price float64, quantity float32) (*core.Order, error)
	MarketBuy(symbol string, quantity float32) (*core.Order, error)
	MarketSell(symbol string, quantity float32) (*core.Order, error)
	LimitStopBuy(symbol string, stopPx float64, price float64, quantity float32) (*core.Order, error)
	LimitStopSell(symbol string, stopPx float64, price float64, quantity float32) (*core.Order, error)
	MarketStopBuy(symbol string, stopPx float64, quantity float32) (*core.Order, error)
	MarketStopSell(symbol string, stopPx float64, quantity float32) (*core.Order, error)
	LimitIfTouchedBuy(symbol string, stopPx float64, price float64, quantity float32)  (*core.Order, error)
	LimitIfTouchedSell(symbol string, stopPx float64, price float64, quantity float32)  (*core.Order, error)
}

type ExchangeAPI interface {
	GetWsSubscribeHandler() WsAuthSubscribeHandler
	Parse(data []byte) (*ParsedData, error)
	HandleMessage(core.ParsedData)

	OrderAPI
}

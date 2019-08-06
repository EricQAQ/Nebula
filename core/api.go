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

type ExchangeAPI interface {
	GetWsSubscribeHandler() WsAuthSubscribeHandler
	Parse(data []byte) (*ParsedData, error)
	HandleMessage(core.ParsedData)
}

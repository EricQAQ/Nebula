package core

type WsAuthSubscribeHandler interface {
	GetOperateArgs() []string
	Serialize() ([]byte, error)
}

type Parser interface {
}

type Handler interface {
}

type ExchangeAPI interface {
	GetWsSubscribeHandler() WsAuthSubscribeHandler
	GetParser() Parser
	GetHandler() Handler
}

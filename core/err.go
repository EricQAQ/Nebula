package core

import (
	"github.com/EricQAQ/Nebula/err"
)

const (
	RetryMaxErrCode         = 2001
	ExchangeNotExistErrCode = 2002
	CreateWsErrCode         = 2003
	LoadHistoryErrCode = 2004
)

var (
	RetryMaxErr = err.CreateNebulaError(
		RetryMaxErrCode, "Reach the maximum number of retries", nil)
	ExchangeNotExistErr = err.CreateNebulaError(
		ExchangeNotExistErrCode, "Callback not exist: %s", nil)
	CreateWsErr = err.CreateNebulaError(
		CreateWsErrCode, "Create websocket client failed: exchange: %s, %s", nil)
	LoadHistoryErr = err.CreateNebulaError(
		LoadHistoryErrCode, "Load history klines failed: %s", nil)
)

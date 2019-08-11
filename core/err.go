package core

import (
	"github.com/EricQAQ/Traed/err"
)

const (
	RetryMaxErrCode         = 2001
	ExchangeNotExistErrCode = 2002
	CreateWsErrCode         = 2003
	LoadHistoryErrCode = 2004
)

var (
	RetryMaxErr = err.CreateTraedError(
		RetryMaxErrCode, "Reach the maximum number of retries", nil)
	ExchangeNotExistErr = err.CreateTraedError(
		ExchangeNotExistErrCode, "Callback not exist: %s", nil)
	CreateWsErr = err.CreateTraedError(
		CreateWsErrCode, "Create websocket client failed: exchange: %s, %s", nil)
	LoadHistoryErr = err.CreateTraedError(
		LoadHistoryErrCode, "Load history klines failed: %s", nil)
)

package bitmex

import (
	"github.com/EricQAQ/Traed/err"
)

const (
	SymbolErrCode = 3001
	ResponseErrCode = 3002
)

var (
	SymbolErr = err.CreateTraedError(
		SymbolErrCode, "Symbol not subscribe in Bitmex.", nil)
	ResponseErr = err.CreateTraedError(
		ResponseErrCode, "%s", nil)
)


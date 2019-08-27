package bitmex

import (
	"github.com/EricQAQ/Nebula/err"
)

const (
	SymbolErrCode = 3001
	ResponseErrCode = 3002
)

var (
	SymbolErr = err.CreateNebulaError(
		SymbolErrCode, "Symbol not subscribe in Bitmex.", nil)
	ResponseErr = err.CreateNebulaError(
		ResponseErrCode, "%s", nil)
)


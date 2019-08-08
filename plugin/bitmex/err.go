package bitmex

import (
	"github.com/EricQAQ/Traed/err"
)

const (
	SymbolErrCode = 3001
)

var SymbolErr = err.CreateTraedError(
	SymbolErrCode, "Symbol not subscribe in Bitmex.", nil)

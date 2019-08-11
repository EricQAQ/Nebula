package storage

import (
	"time"

	"github.com/EricQAQ/Traed/kline"
)

type StorageAPI interface {
	GetDataDir() string
	SetKlines(exchange, symbol string, klines []*kline.Kline) error
	GetKlines(exchange, symbol string, start, end time.Time) ([]*kline.Kline, error)
}

package csv

import (
	"fmt"
	"path"
	"time"
	"os"

	"github.com/EricQAQ/Nebula/kline"
	"github.com/gocarina/gocsv"
)

const (
	timeFormat = "2006-01-02-15"
	// exchangeName_symbol_startTime_endTime
	csvFilename = "%s_%s_%s_%s.csv"
)

type CsvStorage struct {
	dataDir string
}

func NewCsvStorage(dir string) *CsvStorage {
	cs := new(CsvStorage)
	cs.dataDir = dir
	return cs
}

func (cs *CsvStorage) GetDataDir() string {
	return cs.dataDir
}

func (cs *CsvStorage) SetKlines(
	exchange, symbol string, klines []*kline.Kline) error {
	startTs := klines[0].Timestamp
	endTs := klines[len(klines)-1].Timestamp
	dir := path.Join(cs.dataDir, exchange, symbol)
	os.MkdirAll(dir, os.ModePerm)
	fileName := fmt.Sprintf(
		csvFilename, exchange, symbol,
		startTs.Format(timeFormat), endTs.Format(timeFormat))

	file, err := os.OpenFile(
		path.Join(dir, fileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	err = gocsv.MarshalFile(&klines, file)
	return nil
}

func (cs *CsvStorage) GetKlines(
	exchange, symbol string, start, end time.Time) ([]*kline.Kline, error) {
	dir := path.Join(cs.dataDir, exchange, symbol)
	fileName := fmt.Sprintf(
		path.Join(dir, csvFilename), exchange, symbol,
		start.Format(timeFormat), end.Format(timeFormat))

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	klines := []*kline.Kline{}
	if err := gocsv.UnmarshalFile(file, &klines); err != nil {
		return nil, err
	}
	return klines, nil
}

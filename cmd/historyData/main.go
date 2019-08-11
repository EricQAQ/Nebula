package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/EricQAQ/Traed/core"
	"github.com/EricQAQ/Traed/kline"
	"github.com/EricQAQ/Traed/plugin/bitmex"
	"github.com/EricQAQ/Traed/storage/csv"
	"github.com/EricQAQ/Traed/storage"
)

const (
	dateFormat = "2006-01-02"
)

func main() {
	var (
		exchange  = flag.String("exchange", "bitmex", "exchange name")
		symbol    = flag.String("symbols", "", "the symbols, use `,` to split")
		period    = flag.String("period", "1m", "1m, 5m, 1h, 1d")
		startDate = flag.String("start", "", "the start date, `2006-01-02`")
		endDate   = flag.String("end", "", "the end date, `2006-01-02`")
		config    = flag.String("config", "", "config path")
		dataPath  = flag.String("data-path", "", "storage data path")
	)
	flag.Parse()

	cfg := core.LoadConfig(*config)
	var ex core.ExchangeAPI
	if *exchange == "bitmex" {
		ex = bitmex.CreateBitmex(cfg.ExchangeMap[*exchange], cfg.Http)
	}
	st := csv.NewCsvStorage(*dataPath)

	start, _ := time.Parse(dateFormat, *startDate)
	end, _ := time.Parse(dateFormat, *endDate)

	if *period == "1h" {
		getOneHourPeriodData(
			*exchange, *symbol, *period, start, end, ex, st)
	} else if *period == "1m" {
		getOneMinutePeriodData(
			*exchange, *symbol, *period, start, end, ex, st)
	}
}

func getOneMinutePeriodData(
	exName, symbol, period string, start, end time.Time,
	exchange core.ExchangeAPI, st storage.StorageAPI) {
	end = end.Add(24*time.Hour)
	day := start
	klist := make([]*kline.Kline, 0, 60 * 24)
	for {
		tempEnd := start.Add(time.Hour)
		tl, err := exchange.GetHistoryKline(symbol, period, start, tempEnd)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		klist = append(klist, tl[:len(tl)-1]...)

		start = tempEnd
		if start.Sub(day) == 24 * time.Hour {
			if err = st.SetKlines(exName, symbol, klist); err != nil {
				fmt.Println(err.Error())
				return
			}
			day = start
			klist = klist[:1]
		}

		if start.After(end) || start.Equal(end) {
			return
		}
	}
}

func getOneHourPeriodData(
	exName, symbol, period string, start, end time.Time,
	exchange core.ExchangeAPI, st storage.StorageAPI) {
	for {
		tempEnd := start.Add(24 * time.Hour)
		tl, err := exchange.GetHistoryKline(symbol, period, start, tempEnd)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if err = st.SetKlines(exName, symbol, tl); err != nil {
			fmt.Println(err.Error())
			return
		}

		start = tempEnd
		if start.After(end) || start.Equal(end) {
			return
		}
	}
}

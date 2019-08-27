package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/EricQAQ/Nebula/core"
	"github.com/EricQAQ/Nebula/kline"
	"github.com/EricQAQ/Nebula/plugin/bitmex"
	"github.com/EricQAQ/Nebula/storage/csv"
	"github.com/EricQAQ/Nebula/storage"
)

const (
	dateFormat = "2006-01-02"
)

func main() {
	var (
		exchange  = flag.String("exchange", "bitmex", "exchange name")
		symbol    = flag.String("symbols", "", "the symbols, use `,` to split")
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

	getOneMinutePeriodData(
		*exchange, *symbol, "1m", start, end, ex, st)
}

func getOneMinutePeriodData(
	exName, symbol, period string, start, end time.Time,
	exchange core.ExchangeAPI, st storage.StorageAPI) {
	end = end.Add(24*time.Hour)
	day := start
	klist := make([]*kline.Kline, 0, 2048)
	for {
		tempEnd := start.Add(time.Hour)
		tl, err := exchange.GetHistoryKline(symbol, period, start, tempEnd)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Receive %d k-lines.\n", len(tl))
		klist = append(klist, tl[:len(tl)-1]...)

		start = tempEnd
		if start.Sub(day) == 24 * time.Hour {
			if err = st.SetKlines(exName, symbol, klist); err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("Flush k-line data: %s\n", day.String())
			day = start
			klist = klist[:0]
		}

		if start.After(end) || start.Equal(end) {
			return
		}

		if len(tl) == 0 || len(klist) % 60 != 0 {
			if err = st.SetKlines(exName, symbol, klist); err != nil {
				fmt.Println(err.Error())
			}
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

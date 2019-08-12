package fin

import (
	"github.com/EricQAQ/Traed/kline"
)

type MA struct {
	Period int
	Points []Point
	kline  []*kline.Kline
}

// NewMA new Func
func NewMA(value []*kline.Kline, period int) *MA {
	m := &MA{kline: value, Period: period}
	return m
}

func (e *MA) InsertKline(k *kline.Kline) {
	e.kline = append(e.kline)
	e.Calculate(len(e.kline)-1)
}

func (e *MA) Calculate(index int) {
	p := Point{}
	p.Time = e.kline[index].Timestamp
	if index < e.Period-1 {
		e.Points = append(e.Points, p)
		return
	}
	var sum float64
	for j := 0; j < e.Period; j++ {
		sum += e.kline[index-j].Close
	}
	p.Value = +(sum / float64(e.Period))
	e.Points = append(e.Points, p)
}

// Calculation Func
func (e *MA) Calculation() {
	for i := 0; i < len(e.kline); i++ {
		e.Calculate(i)
	}
}

package fin

import (
	"github.com/EricQAQ/Traed/kline"
)

type MA struct {
	Period  int
	Points []Point
	kline   []*kline.Kline
}

// NewMA new Func
func NewMA(value []*kline.Kline, period int) *MA {
	m := &MA{kline: value, Period: period}
	return m
}

// Calculation Func
func (e *MA) Calculation() *MA {
	for i := 0; i < len(e.kline); i++ {
		p := Point{}
		p.Time = e.kline[i].Timestamp
		if i < e.Period-1 {
			e.Points = append(e.Points, p)
			continue
		}
		var sum float64
		for j := 0; j < e.Period; j++ {

			sum += e.kline[i-j].Close
		}
		p.Value = +(sum / float64(e.Period))
		e.Points = append(e.Points, p)
	}
	return e
}

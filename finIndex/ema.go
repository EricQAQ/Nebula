package fin

import (
	"time"

	"github.com/EricQAQ/Traed/kline"
)

type EMA struct {
	Period int
	Points []Point
	kline  []*kline.Kline
}

// NewEMA new Func
func NewEMA(list []*kline.Kline, period int) *EMA {
	m := &EMA{kline: list, Period: period}
	return m
}

// Calculation Func
func (e *EMA) Calculation() *EMA {
	for _, v := range e.kline {
		e.Add(v.Timestamp, v.Close)
	}
	return e
}

// Add adds a new Value to Ema
func (e *EMA) Add(timestamp time.Time, value float64) {
	p := Point{}
	p.Time = timestamp

	//平滑指数，一般取作2/(N+1)
	alpha := 2.0 / (float64(e.Period) + 1.0)

	emaTminusOne := value
	if len(e.Points) > 0 {
		emaTminusOne = e.Points[len(e.Points)-1].Value
	}

	emaT := alpha*value + (1-alpha)*emaTminusOne
	p.Value = emaT
	e.Points = append(e.Points, p)
}

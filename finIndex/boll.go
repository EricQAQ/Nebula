package fin

import (
	"time"
	"math"

	"github.com/EricQAQ/Traed/kline"
)

// 日BOLL指标的计算公式
//		中轨线=N日的移动平均线
//		上轨线=中轨线+两倍的标准差
//		下轨线=中轨线－两倍的标准差

// 日BOLL指标的计算过程
//		1）计算MA MA=N日内的收盘价之和÷N
//		2）计算标准差MD MD=平方根N日的（C－MA）的两次方之和除以N
//		3）计算MB、UP、DN线
//				MB=（N－1）日的MA
//				UP=MB+2×MD
//				DN=MB－2×MD

type BOLL struct {
	PeriodN int     //计算周期
	PeriodK float64 //带宽
	Points  []BollPoint
	kline   []*kline.Kline
}

type BollPoint struct {
	UP   float64
	MID  float64
	Low  float64
	Time time.Time
}

func NewBOLL(list []*kline.Kline, periodN int, PeriodK float64) *BOLL {
	return &BOLL{
		PeriodN: periodN,
		PeriodK: PeriodK,
		kline: list,
	}
}

//sma 计算移动平均线
func (e *BOLL) sma(lines []*kline.Kline) float64 {
	s := len(lines)
	var sum float64
	for i := 0; i < s; i++ {
		sum += float64(lines[i].Close)
	}
	return sum / float64(s)
}

// dma MD=平方根N日的（C－MA）的两次方之和除以N
func (e *BOLL) dma(lines []*kline.Kline, ma float64) float64 {
	s := len(lines)
	var sum float64
	for i := 0; i < s; i++ {
		sum += (lines[i].Close - ma) * (lines[i].Close - ma)
	}
	return math.Sqrt(sum / float64(e.PeriodN))
}

func (e *BOLL) Calculation() *BOLL {
	l := len(e.kline)

	e.Points = make([]BollPoint, l)
	if l < e.PeriodN {
		for i := 0; i < len(e.kline); i++ {
			e.Points[i].Time = e.kline[i].Timestamp
		}
		return e
	}
	for i := l - 1; i > e.PeriodN-1; i-- {

		ps := e.kline[(i - e.PeriodN + 1) : i+1]
		e.Points[i].MID = e.sma(ps)

		//MD=平方根N日的（C－MA）的两次方之和除以N
		md := e.dma(ps, e.Points[i].MID)
		e.Points[i].UP = e.Points[i].MID + e.PeriodK*md
		e.Points[i].Low = e.Points[i].MID - e.PeriodK*md
		e.Points[i].Time = e.kline[i].Timestamp
	}
	return e
}

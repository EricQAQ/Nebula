package fin

import (
	"time"

	"github.com/EricQAQ/Nebula/kline"
)

// 1、计算移动平均值（EMA）
//		12日EMA的算式为 EMA（12）= 前一日EMA（12）×11/13+今日收盘价×2/13
//		26日EMA的算式为 EMA（26）= 前一日EMA（26）×25/27+今日收盘价×2/27
// 2、计算离差值（DIF）
//		DIF=今日EMA（12）－今日EMA（26）
// 3、计算DIF的9日EMA
//		根据离差值计算其9日的EMA，即离差平均值，是所求的MACD值。为了不与指标原名相混淆，此值又名 DEA或DEM。
//		今日DEA（MACD）=前一日DEA×8/10+今日DIF×2/10计算出的DIF和DEA的数值均为正值或负值。
//		用（DIF-DEA）×2即为MACD柱状图。
type MACD struct {
	PeriodShort  int //默认12
	PeriodSignal int //信号长度默认9
	PeriodLong   int //默认26
	Points       []MacdPoint
	kline        []*kline.Kline
}

type MacdPoint struct {
	Time time.Time
	DIF  float64
	DEA  float64
	MACD float64
}

// NewMACD new Func
func NewMACD(
	list []*kline.Kline,
	shortPeriod, longPeriod, signalPeriod int) *MACD {
	m := &MACD{
		PeriodShort:  shortPeriod,
		PeriodSignal: signalPeriod,
		PeriodLong:   longPeriod,
		kline:        list,
	}
	return m
}

func (e *MACD) Calculation() *MACD {
	emaShort := NewEMA(e.kline, e.PeriodShort)
	emaShort.Calculation()
	emaLong := NewEMA(e.kline, e.PeriodLong)
	emaLong.Calculation()
	//计算DIF
	for i := 0; i < len(e.kline); i++ {
		dif := emaShort.Points[i].Value - emaLong.Points[i].Value
		e.Points = append(
			e.Points, MacdPoint{DIF: dif, Time: emaShort.Points[i].Time})
	}

	//临时变量，用于计算DEA
	var difTempKline []*kline.Kline
	for _, v := range e.Points {
		difTempKline = append(
			difTempKline, &kline.Kline{Timestamp: v.Time, Close: v.DIF})
	}
	deaEMA := NewEMA(difTempKline, e.PeriodSignal)
	deaEMA.Calculation()

	//将DEA并入point，同时计算MACD
	for i := 0; i < len(e.Points); i++ {
		e.Points[i].DEA = deaEMA.Points[i].Value
		e.Points[i].MACD = (e.Points[i].DIF - e.Points[i].DEA) * 2
	}
	return e
}

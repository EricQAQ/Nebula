package kline

type Aggregator struct {
	sourcePeriod int
	targetPeriod int
}

func NewAggregator(source, target int) *Aggregator {
	ag := new(Aggregator)
	ag.sourcePeriod = source
	ag.targetPeriod = target
	return ag
}

func mergeKline(klines ...*Kline) *Kline {
	kline := new(Kline)
	kline.Symbol = klines[0].Symbol
	kline.Open = klines[0].Open
	kline.Close = klines[0].Close
	kline.High = klines[0].High
	kline.Low = klines[0].Low
	kline.Vol = klines[0].Vol
	kline.Timestamp = klines[0].Timestamp

	for _, k := range klines[1:] {
		kline.Close = k.Close
		kline.High = maxFloat(kline.High, k.High)
		kline.Low = minFloat(kline.Low, k.Low)
		kline.Vol += k.Vol
	}
	return kline
}

func AggregateKlines(
	sourceInterval, targetInterval int, sourceKline []*Kline) []*Kline {
	mergeCount := targetInterval / sourceInterval
	klist := make([]*Kline, 0, len(sourceKline) / mergeCount)

	for i := 0; i < len(sourceKline); i+=mergeCount {
		kline := mergeKline(sourceKline[i:i+mergeCount]...)
		klist = append(klist, kline)
	}
	return klist
}

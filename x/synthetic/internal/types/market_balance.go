package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"math"
	"strconv"
)

type MarketBalance struct {
	MarketDenom string    `json:"denom" yaml:"denom"`
	LongVolume  sdk.Dec   `json:"long_volume" yaml:"long_volume"`
	ShortVolume sdk.Dec   `json:"short_volume" yaml:"short_volume"`
	Imbalance   Imbalance `json:"imbalance" yaml:"imbalance"`
	FeePercent  float64   `json:"fee_percent"  yaml:"fee_percent"`
}

func (m *MarketBalance) IncreaseLongVolume(amount sdk.Dec) {
	m.LongVolume = m.LongVolume.Add(amount)
	m.recalculate()
}

func (m *MarketBalance) IncreaseShortVolume(amount sdk.Dec) {
	m.ShortVolume = m.ShortVolume.Add(amount)
	m.recalculate()
}

func (m *MarketBalance) GetImbalance() Imbalance {
	return m.Imbalance
}

func (m *MarketBalance) recalculate() {
	m.calculateImbalance()
}

func (m *MarketBalance) calculateImbalance() {
	if m.LongVolume.GT(m.ShortVolume) {
		ratio := m.getRatio(m.LongVolume, m.ShortVolume)
		m.Imbalance = Imbalance{
			matcheng.Bid,
			ratio,
		}
	}

	if m.ShortVolume.GT(m.LongVolume) {
		ratio := m.getRatio(m.ShortVolume, m.LongVolume)
		m.Imbalance = Imbalance{
			matcheng.Ask,
			ratio,
		}
	}
}

func (m MarketBalance) getRatio(num, den sdk.Dec) float64 {
	ratio := num.Quo(den)
	floatStr := ratio.String()
	flt, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		panic(err)
	}
	return flt
}

// we assume that percent cannot be more that math.MaxFloat64 (4503599627370496)
func (m *MarketBalance) CalculateFeePercent(marketImbalance float64) float64 {
	//exp(-1.61222872) * exp(3.22587251 * x)
	e1 := math.Exp(-1.61222872)
	e2 := math.Exp(3.22587251 * marketImbalance)
	return e1 * e2
}

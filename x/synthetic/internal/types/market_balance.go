package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"math"
	"math/big"
	"strconv"
	"time"
)

type MarketBalance struct {
	MarketDenom string          `json:"denom" yaml:"denom"`
	LongVolume  sdk.Int         `json:"long_volume" yaml:"long_volume"`
	ShortVolume sdk.Int         `json:"short_volume" yaml:"short_volume"`
	Imbalance   Imbalance       `json:"imbalance" yaml:"imbalance"`
	Snapshots   VolumeSnapshots `json:"snapshots"`
	Timer       IntervalTimer   `json:"timer" yaml:"timer"`
	Fee         float64         `json:"fee_percent"  yaml:"fee_percent"`
}

// if you prefer to ignore timer settings, just pass zero as an interval value
func EmptyMarketBalance(denom string, interval time.Duration) MarketBalance {
	return MarketBalance{
		denom,
		sdk.NewInt(0),
		sdk.NewInt(0),
		Imbalance{},
		NewVolumeSnapshots(1, []sdk.Int{sdk.NewInt(1)}),
		TimerFromInterval(interval),
		0,
	}
}

// creates snapshot if it a deadline has passed
func (m *MarketBalance) CheckForDeadline() {
	if m.Timer.IntervalIsZero() {
		return
	}

	if !m.Timer.IsExpired(time.Now()) {
		return
	}

	m.SnapshotAndFlash()
	m.Timer.Reset()
}

func (m *MarketBalance) IncreaseLongVolume(amount sdk.Int) {
	m.CheckForDeadline()

	m.LongVolume = m.LongVolume.Add(amount)
	m.recalculate()
}

func (m *MarketBalance) IncreaseShortVolume(amount sdk.Int) {
	m.CheckForDeadline()

	m.ShortVolume = m.ShortVolume.Add(amount)
	m.recalculate()
}

func (m *MarketBalance) Snapshot() VolumeSnapshot {
	return NewVolumeSnapshot(m.LongVolume, m.ShortVolume)
}

func (m *MarketBalance) SaveSnapshot(v VolumeSnapshot) {
	m.Snapshots.AddSnapshot(v)
}

func (m *MarketBalance) SnapshotAndFlash() {
	snapshot := m.Snapshot()
	m.SaveSnapshot(snapshot)
	m.FlashVolumes()
}

func (m *MarketBalance) FlashVolumes() {
	m.LongVolume = sdk.ZeroInt()
	m.ShortVolume = sdk.ZeroInt()
	m.Imbalance.Ratio = 0
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

func (m MarketBalance) getRatio(num1, den1 sdk.Int) float64 {
	num := sdk.NewDecFromBigInt(num1.BigInt())
	den := sdk.NewDecFromBigInt(den1.BigInt())

	if num.Equal(sdk.ZeroDec()) || den.Equal(sdk.ZeroDec()) {
		return 0
	}

	ratio := num.Quo(den).Sub(sdk.OneDec())
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
	if marketImbalance == 0 {
		return 0
	}

	e1 := math.Exp(-1.61222872)
	e2 := math.Exp(3.22587251 * marketImbalance)
	expProduct := e1 * e2
	return expProduct
}

func (m *MarketBalance) AddFee(amount sdk.Int) sdk.Int {
	//var scale int64 = 1000
	if m.Imbalance.Ratio == 0 {
		return amount
	}

	feePercent := m.CalculateFeePercent(m.Imbalance.Ratio)
	if feePercent == 0 { // in fact it is not possible to happen since m.Imbalance.Ratio has already been checked
		panic(amount)
	}

	num, denum := feePercentToNomDenom(feePercent)

	amt := amount.Mul(num).Quo(denum)

	return amt
}

func (m *MarketBalance) GetFeeForAmount(amount sdk.Int) sdk.Int {
	return sdk.ZeroInt()
}

// TODO: find a better name
func feePercentToNomDenom(fee float64) (sdk.Int, sdk.Int) {
	num := int64((100 + fee) * 1000)
	den := int64(100 * 1000)
	return sdk.NewInt(num), sdk.NewInt(den)
}

//
func bigintPow(a, b sdk.Int) sdk.Int {
	var c = big.NewInt(0)
	aBigint := a.BigInt()
	bBigint := b.BigInt()
	c.Exp(aBigint, bBigint, nil)
	return sdk.NewIntFromBigInt(c)
}

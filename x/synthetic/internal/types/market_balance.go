package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"math"
	"math/big"
	"strconv"
	"time"
)

// TODO: It is a good idea to make fields private and define a custom marshaller
//  it is not assumed that user has a permission to change a state of an object otherwise than via setters
type MarketBalance struct {
	MarketDenom     string          `json:"denom" yaml:"denom"`
	LongVolume      sdk.Int         `json:"long_volume" yaml:"long_volume"`
	ShortVolume     sdk.Int         `json:"short_volume" yaml:"short_volume"`
	Imbalance       Imbalance       `json:"imbalance" yaml:"imbalance"`
	VolumeSnapshots VolumeSnapshots `json:"VolumeSnapshots"`
	Timer           IntervalTimer   `json:"timer" yaml:"timer"` // unsafe logic! not tested well
	BlockLimit      int             `json:"block_limit"  yaml:"block_limit"`
	BlocksPassed    int             `json:"block_passed"  yaml:"block_passed"`
}

// if you prefer to ignore timer settings, just pass zero as an interval value
func EmptyMarketBalance(denom string, blockLimit int) MarketBalance {
	return MarketBalance{
		denom,
		sdk.NewInt(0),
		sdk.NewInt(0),
		Imbalance{},
		NewVolumeSnapshots(1, []sdk.Int{sdk.NewInt(1)}),
		TimerFromInterval(time.Duration(0)),
		blockLimit,
		0,
	}
} // if you prefer to ignore timer settings, just pass zero as an interval value
func NewMarketBalance(denom string, snapshots VolumeSnapshots, blockLimit int, interval time.Duration) MarketBalance {
	return MarketBalance{
		denom,
		sdk.NewInt(0),
		sdk.NewInt(0),
		Imbalance{},
		snapshots,
		TimerFromInterval(interval),
		blockLimit,
		0,
	}
}

// call this function on the end of the block or when a new block starts
func (m *MarketBalance) HandleBlockEvent(header abci.Header) {
	if m.BlockLimit != 0 {
		m.BlocksPassed++
		if m.BlocksPassed == m.BlockLimit {
			m.BlocksPassed = 0
			m.SnapshotAndFlash()
		}
	}

	if m.Timer.Interval != time.Duration(0) {
		if m.Timer.IsExpired(header.Time) {
			m.Timer.Reset()
			m.SnapshotAndFlash()
		}
	}
}

// Todo: Unsafe! Needs more testing to use
func (m *MarketBalance) Schedule() {
	f := func() {
		m.SnapshotAndFlash()
		m.Timer.Reset()
	}

	m.Timer.Schedule(f)
}

// creates snapshot if it a deadline has passed.
func (m *MarketBalance) CheckForDeadline() {
	if m.Timer.IsScheduling {
		return
	}

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
	m.LongVolume = m.LongVolume.Add(amount)
	m.Recalculate()
}

func (m *MarketBalance) IncreaseShortVolume(amount sdk.Int) {
	m.ShortVolume = m.ShortVolume.Add(amount)
	m.Recalculate()
}

func (m *MarketBalance) Snapshot() VolumeSnapshot {
	return NewVolumeSnapshot(m.LongVolume, m.ShortVolume)
}

func (m *MarketBalance) SaveSnapshot(v VolumeSnapshot) {
	m.VolumeSnapshots.AddSnapshot(v)
}

func (m *MarketBalance) SnapshotAndFlash() {
	snapshot := m.Snapshot()
	m.SaveSnapshot(snapshot)
	m.FlashVolumes()
}

func (m *MarketBalance) FlashVolumes() {
	m.LongVolume = sdk.ZeroInt()
	m.ShortVolume = sdk.ZeroInt()
	m.Imbalance.Ratio = sdk.ZeroDec()
}

func (m *MarketBalance) GetImbalance() Imbalance {
	return m.Imbalance
}

func (m *MarketBalance) Recalculate() {
	m.calculateImbalance()
}

func (m *MarketBalance) calculateImbalance() {
	latestSnapshot := m.VolumeSnapshots.AppendSnapshot(m.Snapshot())
	volumes := latestSnapshot.GetWeightedVolumes()

	if volumes.LongVolume.GT(volumes.ShortVolume) {
		ratio := m.getRatio(volumes.LongVolume, volumes.ShortVolume)
		m.Imbalance = Imbalance{
			matcheng.Bid,
			ratio,
		}
	}

	if volumes.ShortVolume.GT(volumes.LongVolume) {
		ratio := m.getRatio(volumes.ShortVolume, m.LongVolume)
		m.Imbalance = Imbalance{
			matcheng.Ask,
			ratio,
		}
	}
}

// num1/den1 - 1
func (m MarketBalance) getRatio(num1, den1 sdk.Int) sdk.Dec {
	num := sdk.NewDecFromBigInt(num1.BigInt())
	den := sdk.NewDecFromBigInt(den1.BigInt())

	if num.Equal(sdk.ZeroDec()) || den.Equal(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}

	ratio := num.Quo(den).Sub(sdk.OneDec())
	return ratio
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
	if m.Imbalance.Ratio.Equal(sdk.ZeroDec()) {
		return amount
	}

	floatStr := m.Imbalance.Ratio.String()
	flt, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		panic(err)
	}

	feePercent := m.CalculateFeePercent(flt)
	if feePercent == 0 { // in fact it is not possible to happen since m.Imbalance.Ratio has already been checked
		panic(amount)
	}

	num, denum := fractionForPercentAddition(feePercent)

	amt := amount.Mul(num).Quo(denum)

	return amt
}

func (m *MarketBalance) GetFeeForAmount(amount sdk.Int) sdk.Int {
	if m.Imbalance.Ratio.Equal(sdk.ZeroDec()) {
		return amount
	}

	floatStr := m.Imbalance.Ratio.String()
	flt, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		panic(err)
	}

	feePercent := m.CalculateFeePercent(flt)
	if feePercent == 0 { // in fact it is not possible to happen since m.Imbalance.Ratio has already been checked
		panic(amount)
	}

	num, denum := fractionForPercent(feePercent)

	amt := amount.Mul(num).Quo(denum)
	return amt
}

func (m *MarketBalance) GetFeeForDirection(amount sdk.Int, direction matcheng.Direction) sdk.Int {
	if direction == matcheng.Bid {
		if m.ShortVolume.LT(m.LongVolume) {
			return m.GetFeeForAmount(amount)
		}
		return sdk.ZeroInt()
	} else {
		if m.LongVolume.LT(m.ShortVolume) {
			return m.GetFeeForAmount(amount)
		}
		return sdk.ZeroInt()
	}
}

// TODO: find a better name
// returns a fraction p/q for a given percent
// it is assumed that 100% = 1.0
// so x * p/q = x + (x * percent)
func fractionForPercentAddition(percent float64) (sdk.Int, sdk.Int) {
	num := int64(100 + percent)
	den := int64(100)
	return sdk.NewInt(num), sdk.NewInt(den)
}

// returns a fraction p/q for a given percent
// so x * p/q = x * percent
func fractionForPercent(percent float64) (sdk.Int, sdk.Int) {
	num := int64(percent)
	den := int64(100)
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

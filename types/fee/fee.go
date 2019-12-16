package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Fee struct {
	Numerator   sdk.Int `json:"numerator" yaml:"numerator"`
	Denominator sdk.Int `json:"denominator" yaml:"denominator"`
	MinimumFee  sdk.Int `json:"minimum_fee"  yaml:"minimum_fee"`
}

func NewFee(num, den, minimumFee sdk.Int) Fee {
	if minimumFee.LT(sdk.ZeroInt()) {
		panic(MsgIncorrectMinimumFee)
	}

	if num.LTE(den) {
		panic(MsgNumeratorLTEDenominator)
	}

	return Fee{
		num,
		den,
		minimumFee,
	}
}

func validateNewFee(num, den, minimumFee sdk.Int) {
	if minimumFee.LT(sdk.ZeroInt()) {
		panic(MsgIncorrectMinimumFee)
	}

	if num.LTE(den) {
		panic(MsgNumeratorLTEDenominator)
	}
}

func NewDefaultFee() Fee {
	return Fee{
		Numerator:   sdk.NewInt(1004),
		Denominator: sdk.NewInt(1000),
		MinimumFee:  sdk.NewInt(1),
	}
}

// returns amount with a fee added.
// amount cannot be negative.
// for a specific cases when a fee is to small to be added to an int a MinimalFee variable is added.
// for example if a fee is 0.003 the formula would be : x * 1003 / 1000. So if x is less than 334 a fee would not be added (333 * 1003 / 1000 = int(333.999) = 333)
func (f Fee) AddToAmount(amount sdk.Int) sdk.Int {
	if amount.LTE(sdk.ZeroInt()) {
		panic(MsgIncorrectBaseAmountForFee)
	}

	amountWithFee := amount.Mul(f.Numerator).Quo(f.Denominator)
	if amountWithFee.Sub(amount).LT(f.MinimumFee) {
		return amount.Add(f.MinimumFee)
	}

	return amountWithFee
}

func (f Fee) SubFromAmount(amount sdk.Int) sdk.Int {
	if amount.LTE(sdk.ZeroInt()) {
		panic(MsgIncorrectBaseAmountForFee)
	}

	opposedNumerator := f.Denominator.Mul(sdk.NewInt(2)).Sub(f.Numerator)

	amountWithoutFee := amount.Mul(opposedNumerator).Quo(f.Denominator)
	if amountWithoutFee.LTE(sdk.ZeroInt()) {
		panic(MsgAmountSubFeeTooSmall)
	}

	return amountWithoutFee
}

// need to implement?
func (f Fee) FromRatio(ratio sdk.Dec) {
	ratio.IsInt64()
}

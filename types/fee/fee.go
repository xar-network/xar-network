package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Fee structure
// To add a fee to a certain amount the following operation is performed: amount * Denominator / Numerator.
// For example 335 * (1003/1000) = 335 * 1.003 = 336
//
// An opposite operation thus obvious: amount * Numerator / Denominator.
// For example 335 * (1000/1003) = 335 * 0.997 = 333.9 = 333
type Fee struct {
	Numerator            sdk.Int `json:"numerator" yaml:"numerator"`
	Denominator          sdk.Int `json:"denominator" yaml:"denominator"`
	MinimumAdditionalFee sdk.Int `json:"minimum_additional_fee"  yaml:"minimum_additional_fee"`
	MinimumSubFee        sdk.Int `json:"minimum_sub_fee"  yaml:"minimum_sub_fee"`
}

func NewFee(num, den, minAddFee, minSubFee sdk.Int) Fee {
	validateNewFee(num, den, minAddFee, minSubFee)

	return Fee{
		num,
		den,
		minAddFee,
		minSubFee,
	}
}

func validateNewFee(num, den, minimumFee, minSubFee sdk.Int) {
	if minimumFee.LT(sdk.ZeroInt()) {
		panic(MsgIncorrectMinimumFee)
	}
}

func NewDefaultFee() Fee {
	return Fee{
		Numerator:            sdk.NewInt(1004),
		Denominator:          sdk.NewInt(1000),
		MinimumAdditionalFee: sdk.NewInt(1),
		MinimumSubFee:        sdk.NewInt(1),
	}
}

// Returns amount with a fee added.
// Amount cannot be negative.
// For a specific cases when a fee is to small to be added to an int a MinimumAdditionalFee variable is added.
// For example if a fee is 0.003 the formula would be : x * 1003 / 1000. So if x is less than 334 a fee would not be added (333 * 1003 / 1000 = int(333.999) = 333).
// To evade such a case we will add MinimumAdditionalFee to 333.
// If you prefer to keep things as they are and want to ignore the cases such as described above - just set MinimumAdditionalFee to Zero (0)
//
// The other thing you can do with MinimumAdditionalFee is creating a constant fee (when it is not matter  how large amount is - a fee is always the same).
// To do it - just create a fee with zeroInt set as Numerator and Denominator and set MinimumAdditionalFee to a constant value.
func (f Fee) AddToAmount(amount sdk.Int) sdk.Int {
	if amount.LTE(sdk.ZeroInt()) {
		panic(MsgIncorrectBaseAmountForFee)
	}

	newAmt := amount.Mul(f.Numerator).Quo(f.Denominator)
	if newAmt.Sub(amount).LT(f.MinimumAdditionalFee) {
		return amount.Add(f.MinimumAdditionalFee)
	}

	return newAmt
}

func (f Fee) AddToCoin(coin sdk.Coin) sdk.Coin {
	coin.Amount = f.AddToAmount(coin.Amount)
	return coin
}

func (f Fee) SubFromAmount(amount sdk.Int) sdk.Int {
	if amount.LTE(sdk.ZeroInt()) {
		panic(MsgIncorrectBaseAmountForFee)
	}

	newAmt := amount.Mul(f.Denominator).Quo(f.Numerator)
	if amount.Sub(newAmt).LT(f.MinimumAdditionalFee) {
		return amount.Add(f.MinimumAdditionalFee)
	}

	return newAmt
}

func (f Fee) SubFromCoin(coin sdk.Coin) sdk.Coin {
	coin.Amount = f.SubFromAmount(coin.Amount)
	return coin
}

func (f Fee) GetAmountFee(amount sdk.Int) sdk.Int {
	newAmt := amount.Mul(f.Numerator).Quo(f.Denominator)
	return newAmt.Sub(amount)
}

func FromPercentString(feeRatioStr string) (f Fee) {
	feeRatio, err := sdk.NewDecFromStr(feeRatioStr)
	if err != nil {
		panic(err)
	}

	return FromPercent(feeRatio)
}

// builds fee numerator and denominator from percentToAdd. the formula is:
// percentToAdd = 1 + (numerator / denominator)
// An incoming variable represents a share of an original amount that would be added as fee
func FromPercent(percentToAdd sdk.Dec) (f Fee) {
	if percentToAdd.Equal(sdk.ZeroDec()) {
		f.Numerator = sdk.ZeroInt()
		f.Denominator = sdk.ZeroInt()
		return
	}
	ratio := percentToAdd.Add(sdk.OneDec()) // for example turns 0.001 (percent) to 1.001 (ratio)

	tenDec := sdk.NewDec(10)
	tenInt := sdk.NewInt(10)

	denominator := sdk.OneInt()

	trailingPart := ratio.Sub(ratio.TruncateDec())
	for !trailingPart.Equal(sdk.ZeroDec()) && HasTrail(trailingPart) {
		trailingPart = trailingPart.Mul(tenDec)
		denominator = denominator.Mul(tenInt)
	} //denominator

	num := denominator.Mul(ratio.TruncateInt()).Add(trailingPart.TruncateInt())

	f.Numerator = num
	f.Denominator = denominator
	return
}

// creates a constant Fee for operations that require a constant price to be payed
func MakeConstant(minimumAdditionalFee sdk.Int, minimumSubFee sdk.Int) Fee {
	return Fee{
		Numerator:            sdk.NewInt(0),
		Denominator:          sdk.NewInt(0),
		MinimumAdditionalFee: minimumAdditionalFee,
		MinimumSubFee:        minimumSubFee,
	}
}

func HasTrail(dec sdk.Dec) bool {
	return !dec.TruncateDec().Equal(dec)
}

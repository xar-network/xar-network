package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestFees(t *testing.T) {
	testBasicFuncs(t)
	testFeeAddition(t)
}

type feeTestCase struct {
	FloatStr    string
	Numerator   sdk.Int
	Denominator sdk.Int
}

// need more testcases?
func TestFeeFromRate(t *testing.T) {
	//num := sdk.NewInt(1003)
	//denom := sdk.NewInt(1000)
	//minAmtZero := sdk.NewInt(0)
	//
	//fee := NewFee(num, denom, minAmtZero, minAmtZero)
	testCases := []feeTestCase{
		{"0", sdk.NewInt(0), sdk.NewInt(0)},
		{"0.123", sdk.NewInt(1123), sdk.NewInt(1000)},
		{"0.103", sdk.NewInt(1103), sdk.NewInt(1000)},
		{"0.023", sdk.NewInt(1023), sdk.NewInt(1000)},
		{"0.003", sdk.NewInt(1003), sdk.NewInt(1000)},
		{"0.004", sdk.NewInt(1004), sdk.NewInt(1000)},
		{"0.005", sdk.NewInt(1005), sdk.NewInt(1000)},
		{"0.03", sdk.NewInt(103), sdk.NewInt(100)},
		{"0.13", sdk.NewInt(113), sdk.NewInt(100)},
		{"0.3", sdk.NewInt(13), sdk.NewInt(10)},
		{"1.123", sdk.NewInt(2123), sdk.NewInt(1000)},
		{"1.103", sdk.NewInt(2103), sdk.NewInt(1000)},
		{"1.023", sdk.NewInt(2023), sdk.NewInt(1000)},
		{"1.003", sdk.NewInt(2003), sdk.NewInt(1000)},
		{"1.004", sdk.NewInt(2004), sdk.NewInt(1000)},
		{"1.005", sdk.NewInt(2005), sdk.NewInt(1000)},
		{"1.03", sdk.NewInt(203), sdk.NewInt(100)},
		{"1.13", sdk.NewInt(213), sdk.NewInt(100)},
		{"1.3", sdk.NewInt(23), sdk.NewInt(10)},
		{"2.003", sdk.NewInt(3003), sdk.NewInt(1000)},
		{"2.03", sdk.NewInt(303), sdk.NewInt(100)},
		{"2.3", sdk.NewInt(33), sdk.NewInt(10)},
		{"3.003", sdk.NewInt(4003), sdk.NewInt(1000)},
		{"3.03", sdk.NewInt(403), sdk.NewInt(100)},
		{"3.3", sdk.NewInt(43), sdk.NewInt(10)},
	}

	for _, expected := range testCases {
		f := FromPercentString(expected.FloatStr)

		log.Println("Numerator", f.Numerator, expected.Numerator)
		log.Println("Denominator", f.Denominator, expected.Denominator)

		require.Equal(t, f.Numerator, expected.Numerator)
		require.Equal(t, f.Denominator, expected.Denominator)
	}
}

func testFeeAddition(t *testing.T) {
	testMinAmountAddition(t)
	testMinAmountRemoval(t)
}

func testBasicFuncs(t *testing.T) {
	amt := sdk.NewInt(10000)

	num := sdk.NewInt(1003)
	den := sdk.NewInt(1000)
	minAmtZero := sdk.NewInt(0)

	fee := NewFee(num, den, minAmtZero, minAmtZero)
	assumedAddFeeAmtVal := sdk.NewInt(10030)
	assumedSubFeeAmtVal := sdk.NewInt(9970)

	amtWithFee := fee.AddToAmount(amt)
	amtWithoutFee := fee.SubFromAmount(amt)
	log.Println(amtWithFee)
	log.Println(amtWithoutFee)

	require.True(t, assumedAddFeeAmtVal.Equal(amtWithFee))
	require.True(t, assumedSubFeeAmtVal.Equal(amtWithoutFee))
}

func testMinAmountAddition(t *testing.T) {
	amtWithoutFee := sdk.NewInt(333)

	num := sdk.NewInt(1003)
	den := sdk.NewInt(1000)
	minAmtZero := sdk.NewInt(0)
	minAmtNonZero := sdk.NewInt(1)

	fee := NewFee(num, den, minAmtZero, minAmtZero)
	amtWithFee := fee.AddToAmount(amtWithoutFee)
	require.True(t, amtWithFee.Equal(amtWithoutFee))

	fee = NewFee(num, den, minAmtNonZero, minAmtNonZero)
	amtWithFee = fee.AddToAmount(amtWithoutFee)
	require.False(t, amtWithFee.Equal(amtWithoutFee))

	require.True(t, amtWithFee.Sub(amtWithoutFee).Equal(minAmtNonZero))
}

func testMinAmountRemoval(t *testing.T) {
	amt := sdk.NewInt(333)

	num := sdk.NewInt(1003)
	expectedAntiDenum := sdk.NewInt(997)
	den := sdk.NewInt(1000)
	minAmtZero := sdk.NewInt(0)

	fee := NewFee(num, den, minAmtZero, minAmtZero)
	amtSubFee := fee.SubFromAmount(amt)
	amtSubFee2 := amt.Mul(expectedAntiDenum).Quo(den)

	require.False(t, amtSubFee.Equal(amt))
	require.True(t, amtSubFee.Equal(amtSubFee2))
}

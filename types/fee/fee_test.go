package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestFees(t *testing.T) {
	testFeeAddition(t)
}

func testFeeAddition(t *testing.T) {
	testMinAmountAddition(t)
	testMinAmountRemoval(t)
}

func testMinAmountAddition(t *testing.T) {
	amtWithoutFee := sdk.NewInt(333)

	nom := sdk.NewInt(1003)
	denom := sdk.NewInt(1000)
	minAmtZero := sdk.NewInt(0)
	minAmtNonZero := sdk.NewInt(1)

	fee := NewFee(nom, denom, minAmtZero)
	amtWithFee := fee.AddToAmount(amtWithoutFee)
	require.True(t, amtWithFee.Equal(amtWithoutFee))

	fee = NewFee(nom, denom, minAmtNonZero)
	amtWithFee = fee.AddToAmount(amtWithoutFee)
	require.False(t, amtWithFee.Equal(amtWithoutFee))

	require.True(t, amtWithFee.Sub(amtWithoutFee).Equal(minAmtNonZero))
}

func testMinAmountRemoval(t *testing.T) {
	amtWithoutFee := sdk.NewInt(333)

	nom := sdk.NewInt(1003)
	denom := sdk.NewInt(1000)
	minAmtZero := sdk.NewInt(0)
	minAmtNonZero := sdk.NewInt(1)

	fee := NewFee(nom, denom, minAmtZero)
	amtSubFee := fee.SubFromAmount(amtWithoutFee)
	log.Println(amtSubFee)
	require.True(t, amtSubFee.Equal(amtWithoutFee))

	fee = NewFee(nom, denom, minAmtNonZero)
	amtSubFee = fee.SubFromAmount(amtWithoutFee)
	require.False(t, amtSubFee.Equal(amtWithoutFee))

	require.True(t, amtSubFee.Sub(amtWithoutFee).Equal(minAmtNonZero))
}

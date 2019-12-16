package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFees(t *testing.T) {
	testFeeAddition(t)
}

func testFeeAddition(t *testing.T) {
	testMinAmountAddition(t)
}

func testMinAmountAddition(t *testing.T) {
	amtWithoutFee := sdk.NewInt(334)

	nom := sdk.NewInt(1003)
	denom := sdk.NewInt(1000)
	minAmt := sdk.NewInt(1)

	fee := NewFee(nom, denom, minAmt)

	amtWithFee := fee.MustAddToAmount(amtWithoutFee)

	require.False(t, amtWithFee.Equal(amtWithoutFee))
}

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNewParams1(t *testing.T) {
	is := NewInterestState("caps", "0.04", "kredits", "0.0")
	assert.NoError(t, ValidateInterestState(is))
	assert.Len(t, is.InterestAssets, 2)

	assert.Equal(t, sdk.NewDecWithPrec(4, 2), is.InterestAssets[0].Interest)
	assert.Equal(t, sdk.NewDec(0), is.InterestAssets[1].Interest)
}

func TestValidation(t *testing.T) {
	inflationStates := [...]InterestState{
		NewInterestState("caps", "-0.04"),
		NewInterestState("caps", "0.04", "CAPS", "0.10"),
	}

	for _, is := range inflationStates {
		err := ValidateInterestState(is)
		assert.Error(t, err)
	}
}

func TestFindAndChangeAssetByDenom(t *testing.T) {
	is := NewInterestState("caps", "0.04", "kredits", "0.0")

	kroner := is.FindByDenom("kroner")
	assert.Nil(t, kroner)

	caps := is.FindByDenom("caps")
	assert.NotNil(t, caps)
	caps.Interest, _ = sdk.NewDecFromStr("0.25")

	assert.Equal(t, sdk.NewDecWithPrec(25, 2), is.FindByDenom("caps").Interest)

}

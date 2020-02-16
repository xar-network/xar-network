package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestExponential_Truncate(t *testing.T) {
	var input uint64 = 20

	multiplied := NewExp(sdk.NewUint(input).Mul(ExpScale()))
	result := multiplied.Truncate()

	require.Equal(t, result, sdk.NewUint(input))
}

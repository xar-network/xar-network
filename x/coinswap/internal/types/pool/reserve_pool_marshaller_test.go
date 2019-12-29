package pool

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshaller(t *testing.T) {
	nDenom := "asd"
	nAmt := types.NewInt(1230)
	nc := types.NewCoin(nDenom, nAmt)

	nnDenom := "bsd"
	nnAmt := types.NewInt(1231)
	nnc := types.NewCoin(nnDenom, nnAmt)

	lDenom := "csd"
	lAmt := types.NewInt(1232)
	lc := types.NewCoin(lDenom, lAmt)

	rp := ReservePool{nc, nnc, lc}

	res, err := json.Marshal(rp)
	require.Nil(t, err)

	rp2 := ReservePool{}
	err = json.Unmarshal(res, &rp2)

	require.Nil(t, err)
	require.True(t, rp2.nativeCoins.IsEqual(nc))
	require.True(t, rp2.nonNativeCoins.IsEqual(nnc))
	require.True(t, rp2.liquidityCoins.IsEqual(lc))
}
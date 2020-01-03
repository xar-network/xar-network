/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package keeper

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/coinswap/internal/types"
)

const (
	moduleName = "swap1uftm2ubtc"
)

// test that the module account gets created with an initial
// balance of zero coins.
func TestCreateReservePool(t *testing.T) {
	testDenom := "asd"
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	moduleAcc := keeper.sk.GetModuleAccount(ctx, types.ModuleName)
	require.NotNil(t, moduleAcc)

	pool, found := keeper.GetReservePool(ctx, testDenom)
	require.False(t, found)

	keeper.CreateReservePool(ctx, testDenom)

	pool, found = keeper.GetReservePool(ctx, testDenom)
	require.True(t, found)
	require.NotNil(t, pool)
	poolName := pool.GetName()
	poolExpectedName := keeper.MustGetPoolName(testDenom)
	require.Equal(t, poolName, poolExpectedName)
}

func ttt(ctx *sdk.Context, k *Keeper) {
	k.CreateReservePool(*ctx, "swap1asd2stake")
}

// test that the params can be properly set and retrieved
func TestParams(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		params types.Params
	}{
		{types.DefaultParams()},
		{types.NewParams("pineapple", types.NewFeeParam(sdk.NewInt(5), sdk.NewInt(10)))},
	}

	for _, tc := range cases {
		keeper.SetParams(ctx, tc.params)

		feeParam := keeper.GetFeeParam(ctx)
		require.Equal(t, tc.params.Fee, feeParam)

		nativeDenom := keeper.GetNativeDenom(ctx)
		require.Equal(t, tc.params.NativeDenom, nativeDenom)
	}
}

// test that non existent reserve pool returns false and
// that balance is updated.
func TestGetReservePool(t *testing.T) {
	amt := sdk.NewInt(100)
	ctx, keeper, _ := createTestInput(t, amt, 1)
	log.Println("testinput is done")
	testDenom := "asd"

	reservePool, found := keeper.GetReservePool(ctx, testDenom)
	require.False(t, found)
	log.Println("reservePool 1", reservePool)

	keeper.CreateReservePool(ctx, testDenom)
	reservePool, found = keeper.GetReservePool(ctx, testDenom)
	require.True(t, found)
	log.Println("reservePool 2", reservePool)
}

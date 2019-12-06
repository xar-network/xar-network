/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/coinswap/internal/types"
)

// test that the params can be properly set and retrieved
func TestParams(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		params types.Params
	}{
		{types.DefaultParams()},
		{types.NewParams(sdk.NewRat(5, 10))},
	}

	for _, tc := range cases {
		keeper.SetParams(ctx, tc.params)

		feeParam := keeper.GetParams(ctx)
		require.Equal(t, tc.params.Fee, feeParam.Fee)
	}
}

func TestKeeper_UpdateLiquidity(t *testing.T) {
	total, _ := sdk.NewIntFromString("10000000000000000000")
	ctx, keeper, accs := createTestInput(t, total, 1)
	sender := accs[0].GetAddress()
	denom1 := "btc-min"
	denom2 := sdk.IrisAtto
	uniId, _ := types.GetUniId(denom1, denom2)
	poolAddr := getReservePoolAddr(uniId)

	btcAmt, _ := sdk.NewIntFromString("1")
	depositCoin := sdk.NewCoin("btc-min", btcAmt)

	ftmAmt, _ := sdk.NewIntFromString("10000000000000000000")
	minReward := sdk.NewInt(1)
	deadline := time.Now().Add(1 * time.Minute)
	msg := types.NewMsgAddLiquidity(depositCoin, ftmAmt, minReward, deadline.Unix(), sender)
	_, err := keeper.HandleAddLiquidity(ctx, msg)
	//assert
	require.Nil(t, err)
	reservePoolBalances := keeper.ak.GetAccount(ctx, poolAddr).GetCoins()
	require.Equal(t, "1btc-min,10000000000000000000uftm-atto,10000000000000000000uni:btc-min", reservePoolBalances.String())
	senderBlances := keeper.ak.GetAccount(ctx, sender).GetCoins()
	require.Equal(t, "9999999999999999999btc-min,10000000000000000000uni:btc-min", senderBlances.String())

	withdraw, _ := sdk.NewIntFromString("10000000000000000000")
	msgRemove := types.NewMsgRemoveLiquidity(sdk.NewInt(1), sdk.NewCoin("uni:btc-min", withdraw),
		sdk.NewInt(1), ctx.BlockHeader().Time.Unix(),
		sender)

	_, err = keeper.HandleRemoveLiquidity(ctx, msgRemove)
	require.Nil(t, err)

	poolAccout := keeper.ak.GetAccount(ctx, poolAddr)
	acc := keeper.ak.GetAccount(ctx, sender)
	require.Equal(t, "", poolAccout.GetCoins().String())
	require.Equal(t, "10000000000000000000btc-min,10000000000000000000uftm-atto", acc.GetCoins().String())
}

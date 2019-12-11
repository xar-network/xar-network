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
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

func TestNewQuerier(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(100), 2)

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(keeper)

	// query with incorrect path
	res, err := querier(ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// query for non existent reserve pool should return an error
	req.Path = fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidity)
	req.Data = keeper.cdc.MustMarshalJSON("btc")
	res, err = querier(ctx, []string{"liquidity"}, req)
	require.Error(t, err)
	require.Nil(t, res)

	// query for fee params
	fee := types.DefaultParams().Fee
	req.Path = fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParameters, types.ParamFee)
	req.Data = []byte{}
	res, err = querier(ctx, []string{types.QueryParameters, types.ParamFee}, req)
	if e := keeper.cdc.UnmarshalJSON(res, &fee); e != nil {
		t.Error(e.Error())
	}

	log.Println("res", string(res))
	require.Nil(t, err)
	require.Equal(t, fee, types.DefaultParams().Fee)

	// query for native denom param
	var nativeDenom string
	req.Path = fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParameters, types.ParamNativeDenom)
	res, err = querier(ctx, []string{types.QueryParameters, types.ParamNativeDenom}, req)
	keeper.cdc.UnmarshalJSON(res, &nativeDenom)
	require.Nil(t, err)
	require.Equal(t, nativeDenom, types.DefaultParams().NativeDenom)
	//log.Println("res", string(res))
	//require.Equal(t, nativeDenom, types.DefaultParams().NativeDenom)
}

// TODO: Add tests for valid liquidity queries

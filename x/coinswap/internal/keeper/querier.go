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

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/coinswap/internal/types"
)

// NewQuerier creates a querier for coinswap REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryLiquidity:
			return queryLiquidity(ctx, req, keeper)

		case types.QueryParameters:
			return queryParameters(ctx, path[1:], req, keeper)

		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
		}
	}
}

// queryLiquidity returns the total liquidity available for the provided denomination
// upon success or an error if the query fails.
func queryLiquidity(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var denom types.QueryLiquidityParams
	err := k.cdc.UnmarshalJSON(req.Data, &denom)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	nativeDenom := k.GetNativeDenom(ctx)
	moduleName, err := k.GetPoolName(nativeDenom, denom.NonNativeDenom)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not retrieve module name", err.Error()))
	}

	reservePool, found := k.GetReservePool(ctx, moduleName)
	if !found {
		return nil, types.ErrReservePoolNotFound(types.DefaultCodespace, moduleName)
	}

	// previous return data: reservePool.AmountOf(denom.NonNativeDenom)
	// I am not sure whether it is correct or not
	// as a user, when I query a reserve pool, I expect to use it for the purposes of making new Swap or Liquidity request
	// so I need to get info about an amount of NonNativeDenom AND I also seek for an info about NativeDenom (in order to understand ratio).
	// Thus I think that returning reserve pool is more suitable
	bz, err := k.cdc.MarshalJSONIndent(reservePool, "", " ")
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// queryParameters returns coinswap module parameter queried for upon success
// or an error if the query fails
func queryParameters(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	switch path[0] {
	case types.ParamFee:
		bz, err := k.cdc.MarshalJSONIndent(k.GetFeeParam(ctx), "", " ")
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamNativeDenom:
		bz, err := k.cdc.MarshalJSONIndent(k.GetNativeDenom(ctx), "", " ")
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	default:
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
	}
}

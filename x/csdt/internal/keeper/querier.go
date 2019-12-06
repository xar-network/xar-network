/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGetCsdts:
			return queryGetCsdts(ctx, req, keeper)
		case types.QueryGetParams:
			return queryGetParams(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown csdt query endpoint")
		}
	}
}

// queryGetCsdts fetches CSDTs, optionally filtering by any of the query params (in QueryCsdtsParams).
// While CSDTs do not have an ID, this method can be used to get one CSDT by specifying the collateral and owner.
func queryGetCsdts(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// Decode request
	var requestParams types.QueryCsdtsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &requestParams)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// Get CSDTs
	var csdts types.CSDTs
	if len(requestParams.Owner) != 0 {
		if len(requestParams.CollateralDenom) != 0 {
			// owner and collateral specified - get a single CSDT
			csdt, found := keeper.GetCSDT(ctx, requestParams.Owner, requestParams.CollateralDenom)
			if !found {
				csdt = types.CSDT{
					Owner:            requestParams.Owner,
					CollateralDenom:  requestParams.CollateralDenom,
					CollateralAmount: sdk.NewCoins(sdk.NewCoin(requestParams.CollateralDenom, sdk.ZeroInt())),
					Debt:             sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.ZeroInt())),
				}
			}
			csdts = types.CSDTs{csdt}
		} else {
			// owner, but no collateral specified - get all CSDTs for one address
			return nil, sdk.ErrInternal("getting all CSDTs belonging to one owner not implemented")
		}
	} else {
		// owner not specified -- get all CSDTs or all CSDTs of one collateral type, optionally filtered by price
		var errSdk sdk.Error // := doesn't work here
		csdts, errSdk = keeper.GetCSDTs(ctx, requestParams.CollateralDenom, requestParams.UnderCollateralizedAt)
		if errSdk != nil {
			return nil, errSdk
		}

	}

	// Encode results
	bz, err := codec.MarshalJSONIndent(keeper.cdc, csdts)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// queryGetParams fetches the csdt module parameters
// TODO does this need to exist? Can you use cliCtx.QueryStore instead?
func queryGetParams(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// Get params
	params := keeper.GetParams(ctx)

	// Encode results
	bz, err := codec.MarshalJSONIndent(keeper.cdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

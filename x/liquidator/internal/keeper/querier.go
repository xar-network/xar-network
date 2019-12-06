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
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGetOutstandingDebt:
			return queryGetOutstandingDebt(ctx, path[1:], req, keeper)
		// case QueryGetSurplus:
		// 	return queryGetSurplus()
		default:
			return nil, sdk.ErrUnknownRequest("unknown liquidator query endpoint")
		}
	}
}

func queryGetOutstandingDebt(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// Calculate the remaining seized debt after settling with the liquidator's stable coins.
	stableCoins := keeper.bankKeeper.GetCoins(
		ctx,
		keeper.sk.GetModuleAddress(types.ModuleName),
	).AmountOf(keeper.csdtKeeper.GetStableDenom())
	seizedDebt := keeper.GetSeizedDebt(ctx)
	settleAmount := sdk.MinInt(seizedDebt.Total, stableCoins)
	seizedDebt, err := seizedDebt.Settle(settleAmount)
	if err != nil {
		return nil, err // this shouldn't error in this context
	}

	// Get the available debt after settling
	oustandingDebt := seizedDebt.Available()

	// Encode and return
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, oustandingDebt)
	if err2 != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

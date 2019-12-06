/*

Copyright 2019 All in Bits, Inc
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

package market

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/market/types"
)

const (
	QueryList = "list"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryList:
			return queryList(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown market query endpoint")
		}
	}
}

func queryList(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res := types.ListQueryResult{
		Markets: make([]types.NamedMarket, 0),
	}

	var retErr sdk.Error
	keeper.Iterator(ctx, func(mkt types.Market) bool {
		name, err := keeper.Pair(ctx, mkt.ID)
		if err != nil {
			retErr = err
			return false
		}

		res.Markets = append(res.Markets, types.NamedMarket{
			ID:              mkt.ID.String(),
			BaseAssetDenom:  mkt.BaseAssetDenom,
			QuoteAssetDenom: mkt.QuoteAssetDenom,
			Name:            name,
		})
		return true
	})

	if retErr != nil {
		return nil, retErr
	}

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		panic(err)
	}
	return b, nil
}

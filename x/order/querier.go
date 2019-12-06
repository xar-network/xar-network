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

package order

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/order/types"
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
			return nil, sdk.ErrUnknownRequest("unknown order query endpoint")
		}
	}
}

func queryList(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res := types.ListQueryResult{
		Orders: make([]types.Order, 0),
	}

	keeper.Iterator(ctx, func(order types.Order) bool {
		res.Orders = append(res.Orders, order)
		return true
	})

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		panic("could not marshal result")
	}
	return b, nil
}

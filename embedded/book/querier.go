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

package book

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/xar-network/xar-network/embedded/order"
	"github.com/xar-network/xar-network/types/errs"
	"github.com/xar-network/xar-network/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryGet = "get"
)

func NewQuerier(keeper order.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGet(path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown spread query endpoint")
		}
	}
}

func queryGet(path []string, keeper order.Keeper) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, errs.ErrInvalidArgument("must specify a market ID")
	}

	mktId := store.NewEntityIDFromString(path[0])
	res := keeper.OpenOrdersByMarket(mktId)
	b, err := codec.MarshalJSONIndent(codec.New(), res)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result")
	}
	return b, nil
}

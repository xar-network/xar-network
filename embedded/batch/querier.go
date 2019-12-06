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

package batch

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/xar-network/xar-network/types/errs"
	"github.com/xar-network/xar-network/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryLatest = "latest"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryLatest:
			return queryLatest(path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown batch query endpoint")
		}
	}
}

func queryLatest(path []string, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, errs.ErrInvalidArgument("must specify a market ID")
	}

	marketID := store.NewEntityIDFromString(path[0])
	res, sdkErr := keeper.LatestByMarket(marketID)
	if sdkErr != nil {
		if sdkErr.Code() == errs.CodeNotFound {
			return nil, nil
		}

		return nil, sdkErr
	}

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, errs.ErrMarshalFailure("failed to marshal batch")
	}
	return b, nil
}

/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
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

	"github.com/xar-network/xar-network/x/record/internal/types"
)

//NewQuerier Instance
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryRecord:
			return QueryRecord(ctx, path[1], keeper)
		case types.QueryRecords:
			return QueryRecords(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown record query endpoint")
		}
	}
}

func QueryRecord(ctx sdk.Context, hash string, keeper Keeper) ([]byte, sdk.Error) {
	record := keeper.GetRecord(ctx, hash)
	if record == nil {
		return nil, types.ErrUnknownRecord(hash)
	}

	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), record)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func QueryRecords(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.RecordQueryParams
	err := keeper.Getcdc().UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	records := keeper.List(ctx, params)
	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), records)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

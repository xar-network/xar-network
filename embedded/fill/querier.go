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

package fill

import (
	"math"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/xar-network/xar-network/types/errs"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryGet = "get"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGet(ctx, keeper, req.Data)
		default:
			return nil, sdk.ErrUnknownRequest("unknown fill query endpoint")
		}
	}
}

func queryGet(ctx sdk.Context, keeper Keeper, reqB []byte) ([]byte, sdk.Error) {
	var req QueryRequest
	err := keeper.cdc.UnmarshalBinaryBare(reqB, &req)
	if err != nil {
		return nil, errs.ErrUnmarshalFailure("failed to unmarshal fill query request")
	}

	var start int64
	var end int64

	if req.StartBlock == 0 && req.EndBlock == 0 {
		end = ctx.BlockHeight()
		start = int64(math.Max(float64(end-50), 0))
	} else if req.StartBlock != 0 && req.EndBlock != 0 {
		start = req.StartBlock
		end = req.EndBlock
	} else {
		return nil, errs.ErrInvalidArgument("start and end must either both be defined or neither defined")
	}
	if start > end {
		return nil, errs.ErrInvalidArgument("start must not exceed end")
	}

	res := QueryResult{
		Fills: make([]Fill, 0),
	}
	keeper.IterOverBlockNumbers(start, end, func(fill Fill) bool {
		if !req.Owner.Empty() && !req.Owner.Equals(fill.Owner) {
			return true
		}

		res.Fills = append(res.Fills, fill)
		return true
	})

	b, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdk.ErrInternal("could not marshal result")
	}
	return b, nil
}

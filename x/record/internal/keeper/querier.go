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

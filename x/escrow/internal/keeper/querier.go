package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box/errors"
	"github.com/hashgard/hashgard/x/box/queriers"
	"github.com/hashgard/hashgard/x/box/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/hashgard/hashgard/x/box/keeper"
)

//New Querier Instance
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return queriers.QueryParams(ctx, keeper)
		case types.QueryBox:
			return queriers.QueryBox(ctx, path[1], keeper)
		case types.QuerySearch:
			return queriers.QueryName(ctx, path[1], path[2], keeper)
		case types.QueryList:
			return queriers.QueryList(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown box query endpoint")
		}
	}
}

func QueryParams(ctx sdk.Context, keeper keeper.Keeper) ([]byte, sdk.Error) {
	params := keeper.GetParams(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func QueryBox(ctx sdk.Context, id string, keeper keeper.Keeper) ([]byte, sdk.Error) {
	box := keeper.GetBox(ctx, id)
	if box == nil {
		return nil, errors.ErrUnknownBox(id)
	}

	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), box)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func QueryName(ctx sdk.Context, boxType string, name string, keeper keeper.Keeper) ([]byte, sdk.Error) {
	box := keeper.SearchBox(ctx, boxType, name)
	if box == nil {
		return nil, errors.ErrUnknownBox(name)
	}

	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), box)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

//func QueryDepositAmountFromDepositBox(ctx sdk.Context, id string, accAddress string, keeper keeper.Keeper) ([]byte, sdk.Error) {
//	address, err := sdk.AccAddressFromBech32(accAddress)
//	if err != nil {
//		return nil, sdk.ErrInvalidAddress(accAddress)
//	}
//	amount := keeper.GetDepositByAddress(ctx, id, address)
//
//	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), amount)
//	if err != nil {
//		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
//	}
//	return bz, nil
//}
func QueryList(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
	var params params.BoxQueryParams
	err := keeper.Getcdc().UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	boxs := keeper.List(ctx, params)
	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), boxs)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

//func QueryDepositList(ctx sdk.Context, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, sdk.Error) {
//	var params params.BoxQueryDepositListParams
//	err := keeper.Getcdc().UnmarshalJSON(req.Data, &params)
//	if err != nil {
//		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
//	}
//
//	boxs := keeper.QueryDepositListFromDepositBox(ctx, params.Id, params.Owner)
//	bz, err := codec.MarshalJSONIndent(keeper.Getcdc(), boxs)
//	if err != nil {
//		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
//	}
//	return bz, nil
//}

func GetQueryBoxPath(id string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryBox, id)
}
func GetQueryBoxParamsPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryParams)
}
func GetQueryBoxSearchPath(boxType string, name string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QuerySearch, boxType, name)
}
func GetQueryBoxsPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryList)
}
func QueryBoxParams(cliCtx context.CLIContext) ([]byte, error) {
	return cliCtx.QueryWithData(GetQueryBoxParamsPath(), nil)
}
func QueryBoxByName(boxType string, name string, cliCtx context.CLIContext) ([]byte, error) {
	return cliCtx.QueryWithData(GetQueryBoxSearchPath(boxType, name), nil)
}

func QueryBoxByID(id string, cliCtx context.CLIContext) ([]byte, error) {
	return cliCtx.QueryWithData(GetQueryBoxPath(id), nil)
}

func QueryBoxsList(params params.BoxQueryParams, cdc *codec.Codec, cliCtx context.CLIContext) ([]byte, error) {
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	return cliCtx.QueryWithData(GetQueryBoxsPath(), bz)
}

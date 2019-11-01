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
				csdt = types.CSDT{Owner: requestParams.Owner, CollateralDenom: requestParams.CollateralDenom, CollateralAmount: sdk.ZeroInt(), Debt: sdk.ZeroInt()}
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

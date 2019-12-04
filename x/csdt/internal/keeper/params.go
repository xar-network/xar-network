package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// ---------- Module Parameters ----------
// GetParams returns the params from the store
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramsSubspace.GetParamSet(ctx, &p)
	return p
}

// SetParams sets params on the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramsSubspace.SetParamSet(ctx, &params)
}

func (k Keeper) AddCollateralParam(ctx sdk.Context, nominee string, collateralParam types.CollateralParam) sdk.Error {
	if !k.IsNominee(ctx, nominee) {
		return sdk.ErrInternal(fmt.Sprintf("not a nominee: '%s'", nominee))
	}
	params := k.GetParams(ctx)
	if params.IsCollateralPresent(collateralParam.Denom) {
		return sdk.ErrInternal(fmt.Sprintf("param already exists: '%s'", collateralParam.String()))
	}
	params.CollateralParams = append(params.CollateralParams, collateralParam)
	k.SetParams(ctx, params)
	return nil
}

func (k Keeper) SetCollateralParam(ctx sdk.Context, nominee string, collateralParam types.CollateralParam) sdk.Error {
	if !k.IsNominee(ctx, nominee) {
		return sdk.ErrInternal(fmt.Sprintf("not a nominee: '%s'", nominee))
	}
	params := k.GetParams(ctx)
	if !params.IsCollateralPresent(collateralParam.Denom) {
		return sdk.ErrInternal(fmt.Sprintf("param doesnt exists: '%s'", collateralParam.String()))
	}
	for x, cp := range params.CollateralParams {
		if cp.Denom == collateralParam.Denom {
			params.CollateralParams[x] = collateralParam
		}
	}
	k.SetParams(ctx, params)
	return nil
}

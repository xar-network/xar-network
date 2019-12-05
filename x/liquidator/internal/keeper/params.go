package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

// GetParams returns the params for liquidator module
func (k Keeper) GetParams(ctx sdk.Context) types.LiquidatorParams {
	var params types.LiquidatorParams
	k.paramsSubspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets params for the liquidator module
func (k Keeper) SetParams(ctx sdk.Context, params types.LiquidatorParams) {
	k.paramsSubspace.SetParamSet(ctx, &params)
}

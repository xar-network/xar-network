package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

// GetParams gets params from the store
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.GetNomineeParams(ctx))
}

// SetParams updates params in the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetNomineeParams get nominee params from store
func (k Keeper) GetNomineeParams(ctx sdk.Context) []string {
	var nominees []string
	k.paramSubspace.Get(ctx, types.KeyNominees, &nominees)
	return nominees
}

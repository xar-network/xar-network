/**

Baseline from Kava Cosmos Module

**/

package liquidator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

// InitGenesis sets the genesis state in the keeper.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	return types.GenesisState{
		Params: params,
	}
}

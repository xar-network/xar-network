/**

Baseline from Kava Cosmos Module

**/

package liquidator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/liquidator/internal/keeper"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	LiquidatorModuleParams types.LiquidatorModuleParams `json:"params"`
}

// DefaultGenesisState returns a default genesis state
// TODO pick better values
func DefaultGenesisState() GenesisState {
	return GenesisState{
		types.LiquidatorModuleParams{
			DebtAuctionSize: sdk.NewInt(1000),
			CollateralParams: []types.CollateralParams{
				{
					Denom:       "btc",
					AuctionSize: sdk.NewInt(1),
				},
				{
					Denom:       "bnb",
					AuctionSize: sdk.NewInt(1000),
				},
				{
					Denom:       "eth",
					AuctionSize: sdk.NewInt(10),
				},
				{
					Denom:       "ftm",
					AuctionSize: sdk.NewInt(1000000),
				},
			},
		},
	}
}

// InitGenesis sets the genesis state in the keeper.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.LiquidatorModuleParams)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	// TODO implement this
	return DefaultGenesisState()
}

// ValidateGenesis performs basic validation of genesis data returning an error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	// TODO
	// check debt auction size > 0
	// validate denoms
	// check no repeated denoms
	// check collateral auction sizes > 0
	return nil
}

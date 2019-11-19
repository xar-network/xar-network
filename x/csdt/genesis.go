/**

Baseline from Kava Cosmos Module

**/

package csdt

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	CsdtModuleParams types.CsdtModuleParams `json:"params"`
	GlobalDebt       sdk.Int                `json:"global_debt"`
	// don't need to setup CollateralStates as they are created as needed
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		types.CsdtModuleParams{
			GlobalDebtLimit: sdk.NewInt(25000000000000),
			CollateralParams: []types.CollateralParams{
				{
					Denom:            "ubtc",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewInt(500000000000),
				},
				{
					Denom:            "ubnb",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewInt(500000000000),
				},
				{
					Denom:            "ueth",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewInt(500000000000),
				},
				{
					Denom:            "uftm",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewInt(500000000000),
				},
				{
					Denom:            "uzar",
					LiquidationRatio: sdk.MustNewDecFromStr("1.3"),
					DebtLimit:        sdk.NewInt(500000000000),
				},
			},
		},
		sdk.ZeroInt(),
	}
}

func NewGenesisState(csdtModuleParams types.CsdtModuleParams, globalDebt sdk.Int) GenesisState {
	return GenesisState{
		CsdtModuleParams: csdtModuleParams,
		GlobalDebt:       globalDebt,
	}
}

// InitGenesis sets the genesis state in the keeper.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	k.SetParams(ctx, data.CsdtModuleParams)
	k.SetGlobalDebt(ctx, data.GlobalDebt)
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	// TODO implement this
	// validate denoms
	// check collateral debt limits sum to global limit?
	// check limits are > 0
	// check ratios are > 1
	// check no repeated denoms

	// check global debt is zero - force the chain to always start with zero stable coin, otherwise collateralStatus's will need to be set up as well.
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	// TODO add csdt's to genesisState for export or track csdt in account space?
	params := k.GetParams(ctx)
	globalDebt := k.GetGlobalDebt(ctx)
	return NewGenesisState(params, globalDebt)
}

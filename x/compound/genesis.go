package compound

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zar-network/zar-network/x/compound/internal/keeper"
	"github.com/zar-network/zar-network/x/compound/internal/types"
)

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	CompoundModuleParams types.CompoundModuleParams `json:"params"`
	GlobalDebt           sdk.Int                    `json:"global_debt"`
	// don't need to setup CollateralStates as they are created as needed
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		types.CompoundModuleParams{
			GlobalDebtLimit: sdk.NewInt(1000000),
			CollateralParams: []types.CollateralParams{
				{
					Denom:            "btc",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewInt(500000),
				},
				{
					Denom:            "bnb",
					LiquidationRatio: sdk.MustNewDecFromStr("2.0"),
					DebtLimit:        sdk.NewInt(500000),
				},
				{
					Denom:            "eth",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewInt(500000),
				},
				{
					Denom:            "ftm",
					LiquidationRatio: sdk.MustNewDecFromStr("2.0"),
					DebtLimit:        sdk.NewInt(500000),
				},
			},
		},
		sdk.ZeroInt(),
	}
}

// InitGenesis sets the genesis state in the keeper.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	k.SetParams(ctx, data.CdpModuleParams)
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
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	// TODO implement this
	return DefaultGenesisState()
}

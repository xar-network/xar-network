package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type CompoundModuleParams struct {
	GlobalDebtLimit  sdk.Int
	CollateralParams []CollateralParams
}

type CollateralParams struct {
	Denom            string  // Coin name of collateral type
	LiquidationRatio sdk.Dec // The ratio (Collateral (priced in stable coin) / Debt) under which a CSDT will be liquidated
	DebtLimit        sdk.Int // Maximum amount of debt allowed to be drawn from this collateral type
	//DebtFloor        sdk.Int // used to prevent dust
}

var ModuleParamsKey = []byte("CompoundModuleParams")

func CreateParamsKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ModuleParamsKey, CompoundModuleParams{},
	)
}

// Implement fmt.Stringer interface for cli querying
func (p CompoundModuleParams) String() string {
	out := fmt.Sprintf(`Params:
	Global Debt Limit: %s
	Collateral Params:`,
		p.GlobalDebtLimit,
	)
	for _, cp := range p.CollateralParams {
		out += fmt.Sprintf(`
		%s
			Liquidation Ratio: %s
			Debt Limit:        %s`,
			cp.Denom,
			cp.LiquidationRatio,
			cp.DebtLimit,
		)
	}
	return out
}

// Helper methods to search the list of collateral params for a particular denom. Wouldn't be needed if amino supported maps.

func (p CompoundModuleParams) GetCollateralParams(collateralDenom string) CollateralParams {
	// search for matching denom, return
	for _, cp := range p.CollateralParams {
		if cp.Denom == collateralDenom {
			return cp
		}
	}
	// panic if not found, to be safe
	panic("collateral params not found in module params")
}
func (p CompoundModuleParams) IsCollateralPresent(collateralDenom string) bool {
	// search for matching denom, return
	for _, cp := range p.CollateralParams {
		if cp.Denom == collateralDenom {
			return true
		}
	}
	return false
}

package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryGetCdps   = "cdps"
	QueryGetParams = "params"
)

// StableDenom asset code of the dollar-denominated debt coin
const StableDenom = "ftg" // TODO allow to be changed
// GovDenom asset code of the governance coin
const GovDenom = "ftm"

// CDP is the state of a single Collateralized Debt Position.
type CDP struct {
	//ID             []byte                                    // removing IDs for now to make things simpler
	Owner            sdk.AccAddress `json:"owner"`             // Account that authorizes changes to the CDP
	CollateralDenom  string         `json:"collateral_denom"`  // Type of collateral stored in this CDP
	CollateralAmount sdk.Int        `json:"collateral_amount"` // Amount of collateral stored in this CDP
	Debt             sdk.Int        `json:"debt"`              // Amount of stable coin drawn from this CDP
}

func (cdp CDP) IsUnderCollateralized(price sdk.Dec, liquidationRatio sdk.Dec) bool {
	collateralValue := sdk.NewDecFromInt(cdp.CollateralAmount).Mul(price)
	minCollateralValue := liquidationRatio.Mul(sdk.NewDecFromInt(cdp.Debt))
	return collateralValue.LT(minCollateralValue) // TODO LT or LTE?
}

func (cdp CDP) String() string {
	return strings.TrimSpace(fmt.Sprintf(`CDP:
  Owner:      %s
  Collateral: %s
  Debt:       %s`,
		cdp.Owner,
		sdk.NewCoin(cdp.CollateralDenom, cdp.CollateralAmount),
		sdk.NewCoin(StableDenom, cdp.Debt),
	))
}

type CDPs []CDP

func (cdps CDPs) String() string {
	out := ""
	for _, cdp := range cdps {
		out += cdp.String() + "\n"
	}
	return out
}

type QueryCdpsParams struct {
	CollateralDenom       string         // get CDPs with this collateral denom
	Owner                 sdk.AccAddress // get CDPs belonging to this owner
	UnderCollateralizedAt sdk.Dec        // get CDPs that will be below the liquidation ratio when the collateral is at this price.
}

// byCollateralRatio is used to sort CDPs
type ByCollateralRatio CDPs

func (cdps ByCollateralRatio) Len() int      { return len(cdps) }
func (cdps ByCollateralRatio) Swap(i, j int) { cdps[i], cdps[j] = cdps[j], cdps[i] }
func (cdps ByCollateralRatio) Less(i, j int) bool {
	// Sort by "collateral ratio" ie collateralAmount/Debt
	// The comparison is: collat_i/debt_i < collat_j/debt_j
	// But to avoid division this can be rearranged to: collat_i*debt_j < collat_j*debt_i
	// Provided the values are positive, so check for positive values.
	if cdps[i].CollateralAmount.IsNegative() ||
		cdps[i].Debt.IsNegative() ||
		cdps[j].CollateralAmount.IsNegative() ||
		cdps[j].Debt.IsNegative() {
		panic("negative collateral and debt not supported in CDPs")
	}
	// TODO overflows could cause panics
	left := cdps[i].CollateralAmount.Mul(cdps[j].Debt)
	right := cdps[j].CollateralAmount.Mul(cdps[i].Debt)
	return left.LT(right)
}

// CollateralState stores global information tied to a particular collateral type.
type CollateralState struct {
	Denom     string  // Type of collateral
	TotalDebt sdk.Int // total debt collateralized by a this coin type
	//AccumulatedFees sdk.Int // Ignoring fees for now
}

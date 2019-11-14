package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CSDT is the state of a single account.
type CSDT struct {
	//ID             []byte                                    // removing IDs for now to make things simpler
	Owner            sdk.AccAddress `json:"owner"`             // Account that authorizes changes to the CSDT
	CollateralDenom  string         `json:"collateral_denom"`  // Type of collateral stored in this CSDT
	CollateralAmount sdk.Int        `json:"collateral_amount"` // Amount of collateral stored in this CSDT
	Debt             sdk.Int        `json:"debt"`              // Amount of stable coin drawn from this CSDT
}

func (csdt CSDT) IsUnderCollateralized(price sdk.Dec, liquidationRatio sdk.Dec) bool {
	collateralValue := sdk.NewDecFromInt(csdt.CollateralAmount).Mul(price)
	minCollateralValue := liquidationRatio.Mul(sdk.NewDecFromInt(csdt.Debt))
	return collateralValue.LT(minCollateralValue) // TODO LT or LTE?
}

func (csdt CSDT) String() string {
	return strings.TrimSpace(fmt.Sprintf(`CSDT:
  Owner:      %s
  Collateral: %s
  Debt:       %s`,
		csdt.Owner,
		sdk.NewCoin(csdt.CollateralDenom, csdt.CollateralAmount),
		sdk.NewCoin(StableDenom, csdt.Debt),
	))
}

type CSDTs []CSDT

func (csdts CSDTs) String() string {
	out := ""
	for _, csdt := range csdts {
		out += csdt.String() + "\n"
	}
	return out
}

type QueryCsdtsParams struct {
	CollateralDenom       string         // get CSDTs with this collateral denom
	Owner                 sdk.AccAddress // get CSDTs belonging to this owner
	UnderCollateralizedAt sdk.Dec        // get CSDTs that will be below the liquidation ratio when the collateral is at this price.
}

// byCollateralRatio is used to sort CSDTs
type ByCollateralRatio CSDTs

func (csdts ByCollateralRatio) Len() int      { return len(csdts) }
func (csdts ByCollateralRatio) Swap(i, j int) { csdts[i], csdts[j] = csdts[j], csdts[i] }
func (csdts ByCollateralRatio) Less(i, j int) bool {
	// Sort by "collateral ratio" ie collateralAmount/Debt
	// The comparison is: collat_i/debt_i < collat_j/debt_j
	// But to avoid division this can be rearranged to: collat_i*debt_j < collat_j*debt_i
	// Provided the values are positive, so check for positive values.
	if csdts[i].CollateralAmount.IsNegative() ||
		csdts[i].Debt.IsNegative() ||
		csdts[j].CollateralAmount.IsNegative() ||
		csdts[j].Debt.IsNegative() {
		panic("negative collateral and debt not supported in CSDTs")
	}
	// TODO overflows could cause panics
	left := csdts[i].CollateralAmount.Mul(csdts[j].Debt)
	right := csdts[j].CollateralAmount.Mul(csdts[i].Debt)
	return left.LT(right)
}

// CollateralState stores global information tied to a particular collateral type.
type CollateralState struct {
	Denom     string  // Type of collateral
	TotalDebt sdk.Int // total debt collateralized by a this coin type
	//AccumulatedFees sdk.Int // Ignoring fees for now
}

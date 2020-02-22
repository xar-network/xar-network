package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CSDT is the state of a single account.
type CSDT struct {
	// ID             []byte                                    // removing IDs for now to make things simpler
	Owner                sdk.AccAddress `json:"owner" yaml:"owner"`                         // Account that authorizes changes to the CSDT
	CollateralDenom      string         `json:"collateral_denom" yaml:"collateral_denom"`   // Type of collateral stored in this CSDT
	CollateralAmount     sdk.Coins      `json:"collateral_amount" yaml:"collateral_amount"` // Amount of collateral stored in this CSDT
	Debt                 sdk.Coins      `json:"debt" yaml:"debt"`
	DebtAccruedBlock     int64          `json:"debt_accrued_block" yaml:"debt_accrued_block"`
	Interest             sdk.Coins      `json:"interest" yaml:"interest"` // Amount of interest accumulated on collateral stored in this CSDT
	InterestAccruedBlock int64          `json:"interest_accrued_block" yaml:"interest_accrued_block"`
	AccumulatedFees      sdk.Coins      `json:"accumulated_fees" yaml:"accumulated_fees"`
	FeesUpdated          time.Time      `json:"fees_updated" yaml:"fees_updated"`
}

func (csdt CSDT) IsUnderCollateralized(price sdk.Dec, liquidationRatio sdk.Dec) bool {
	collateralValue := sdk.NewDecFromInt(csdt.CollateralAmount.AmountOf(csdt.CollateralDenom)).Mul(price)
	minCollateralValue := sdk.NewDec(0)
	for _, c := range csdt.Debt {
		minCollateralValue = minCollateralValue.Add(liquidationRatio.Mul(c.Amount.ToDec()))
	}
	return collateralValue.LT(minCollateralValue) // TODO LT or LTE?
}

func (csdt CSDT) String() string {
	return strings.TrimSpace(fmt.Sprintf(`CSDT:
  Owner:      %s
	Collateral Type: %s
	Collateral: %s
	Debt: %s
	Fees: %s
	Fees Last Updated: %s`,
		csdt.Owner,
		csdt.CollateralDenom,
		csdt.CollateralAmount,
		csdt.Debt,
		csdt.AccumulatedFees,
		csdt.FeesUpdated,
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

// byCollateralRatio is used to sort CSDTs
type ByCollateralRatio CSDTs

func (csdts ByCollateralRatio) Len() int      { return len(csdts) }
func (csdts ByCollateralRatio) Swap(i, j int) { csdts[i], csdts[j] = csdts[j], csdts[i] }
func (csdts ByCollateralRatio) Less(i, j int) bool {
	// Sort by "collateral ratio" ie collateralAmount/Debt
	// The comparison is: collat_i/debt_i < collat_j/debt_j
	// But to avoid division this can be rearranged to: collat_i*debt_j < collat_j*debt_i
	// Provided the values are positive, so check for positive values.
	if csdts[i].CollateralAmount.IsAnyNegative() ||
		csdts[i].Debt.IsAnyNegative() ||
		csdts[j].CollateralAmount.IsAnyNegative() ||
		csdts[j].Debt.IsAnyNegative() {
		panic("negative collateral and debt not supported in CSDTs")
	}
	// TODO overflows could cause panics
	left := csdts[i].CollateralAmount.AmountOf(csdts[i].CollateralDenom).Mul(csdts[j].Debt.AmountOf(StableDenom))
	right := csdts[j].CollateralAmount.AmountOf(csdts[j].CollateralDenom).Mul(csdts[i].Debt.AmountOf(StableDenom))
	return left.LT(right)
}

// CollateralState stores global information tied to a particular collateral type.
type CollateralState struct {
	Denom            string   // Type of collateral
	TotalDebt        sdk.Int  // total debt collateralized by a this coin type
	TotalCash        sdk.Uint // total cash supply
	TotalBorrows     sdk.Uint // total borrows balance
	Reserves         sdk.Uint // portion of accrued interest set aside as reserves
	LastAccrualBlock int64    // last block interest was accrued
	BorrowIndex      sdk.Uint // index for interest calculation
	// AccumulatedFees sdk.Int // Ignoring fees for now
}

type CoinU struct {
	Denom  string   `json:"denom"`
	Amount sdk.Uint `json:"amount"`
}

type CoinUs []CoinU

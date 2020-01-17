/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CSDT is the state of a single account.
type CSDT struct {
	//ID             []byte                                    // removing IDs for now to make things simpler
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`                         // Account that authorizes changes to the CSDT
	CollateralAmount sdk.Coins      `json:"collateral_amount" yaml:"collateral_amount"` // Amount of collateral stored in this CSDT
	Debt             sdk.Coins      `json:"debt" yaml:"debt"`
	AccumulatedFees  sdk.Coins      `json:"accumulated_fees" yaml:"accumulated_fees"`
	FeesUpdated      time.Time      `json:"fees_updated" yaml:"fees_updated"` // Amount of stable coin drawn from this CSDT
}

func (csdt CSDT) IsUnderCollateralized(price sdk.Dec, liquidationRatio sdk.Dec, denom string) bool {
	collateralValue := sdk.NewDecFromInt(csdt.CollateralAmount.AmountOf(denom)).Mul(price)
	minCollateralValue := sdk.NewDec(0)
	for _, c := range csdt.Debt {
		minCollateralValue = minCollateralValue.Add(liquidationRatio.Mul(c.Amount.ToDec()))
	}
	return collateralValue.LT(minCollateralValue) // TODO LT or LTE?
}

// will handle all the validation in the future updates
func (csdt CSDT) Validate(price sdk.Dec, liquidationRatio sdk.Dec, denom string) sdk.Error {
	isUnderCollateralized := csdt.IsUnderCollateralized(
		price,
		liquidationRatio,
		denom,
	)
	if isUnderCollateralized {
		return sdk.ErrInternal("Change to CSDT would put it below liquidation ratio")
	}
	return nil
}

func (csdt CSDT) String() string {
	return strings.TrimSpace(fmt.Sprintf(`CSDT:
  Owner:      %s
	Collateral: %s
	Debt: %s
	Fees: %s
	Fees Last Updated: %s`,
		csdt.Owner,
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

	ltCount := 0
	gtCount := 0
	for _, clt := range csdts[i].CollateralAmount {
		debtInt := csdts[j].Debt.AmountOf(clt.Denom)
		left := clt.Amount.Mul(debtInt)

		for _, cltR := range csdts[j].CollateralAmount {
			debtInt = csdts[i].Debt.AmountOf(cltR.Denom)
			right := cltR.Amount.Mul(debtInt)

			if left.LT(right) {
				ltCount++
			} else {
				gtCount++
			}
		}
	}
	//left := csdts[i].CollateralAmount.AmountOf(csdts[i].CollateralDenom).Mul(csdts[j].Debt.AmountOf(StableDenom))
	//right := csdts[j].CollateralAmount.AmountOf(csdts[j].CollateralDenom).Mul(csdts[i].Debt.AmountOf(StableDenom))

	return ltCount > gtCount
}

// CollateralState stores global information tied to a particular collateral type.
type CollateralState struct {
	Denom     string  // Type of collateral
	TotalDebt sdk.Int // total debt collateralized by a this coin type
	//AccumulatedFees sdk.Int // Ignoring fees for now
}

type PoolSnapValue struct {
	Limit 	PoolDecreaseLimitParam
	Val		sdk.Coin
}

type PoolSnapshot struct {
	ByLimits []PoolSnapValue
}

func (snap *PoolSnapshot) GetVal(limit PoolDecreaseLimitParam, denom string) *sdk.Coin {
	for _, v := range snap.ByLimits {
		if v.Val.Denom == denom {
			if v.Limit.IsEqual(limit) {
				return &v.Val
			}
		}
	}

	return nil
}

func (snap *PoolSnapshot) SetVal(limit PoolDecreaseLimitParam, val sdk.Coin) {
	for _, v := range snap.ByLimits {
		if v.Val.Denom == val.Denom {
			if v.Limit.IsEqual(limit) {
				v.Val = val
				return
			}
		}
	}

	snap.ByLimits = append(snap.ByLimits, PoolSnapValue{
		Limit: limit,
		Val:   val,
	})
}

type SignedCoin struct {
	sdk.Coin
	isNeg		bool
}

func NewSignedCoin(denom string, amount sdk.Int) SignedCoin {

	if amount.IsNegative() {
		return SignedCoin{
			Coin:  sdk.NewCoin(denom, amount.Neg()),
			isNeg: true,
		}
	}

	return SignedCoin{
		Coin:  sdk.NewCoin(denom, amount),
		isNeg: false,
	}
}

func NewSignedCoinFromCoin(coin sdk.Coin) SignedCoin {
	return SignedCoin{
		Coin:  coin,
		isNeg: false,
	}
}

func (c *SignedCoin) IsNegative() bool {
	return c.isNeg
}

func (c *SignedCoin) IsPositive() bool {
	return !c.isNeg
}

func (c SignedCoin) String() string {
	s := ""
	if c.isNeg {
		s = "-"
	}
	return s+c.Coin.String()
}
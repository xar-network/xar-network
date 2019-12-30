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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/types/fee"
)

/*
How this uses the sdk params module:
 - Put all the params for this module in one struct `CSDTModuleParams`
 - Store this in the keeper's paramSubspace under one key
 - Provide a function to load the param struct all at once `keeper.GetParams(ctx)`
It's possible to set individual key value pairs within a paramSubspace, but reading and setting them is awkward (an empty variable needs to be created, then Get writes the value into it)
This approach will be awkward if we ever need to write individual parameters (because they're stored all together). If this happens do as the sdk modules do - store parameters separately with custom get/set func for each.
*/

// Parameter keys
var (
	// ParamStoreKeyAuctionParams Param store key for auction params
	KeyGlobalDebtLimit      = []byte("GlobalDebtLimit")
	KeyCollateralParams     = []byte("CollateralParams")
	KeyDebtParams           = []byte("DebtParams")
	KeyCircuitBreaker       = []byte("CircuitBreaker")
	KeyNominees             = []byte("Nominees")
	KeyFee                  = []byte("Fee")
	DefaultGlobalDebt       = sdk.NewCoins(sdk.NewCoin(StableDenom, sdk.NewInt(500000000000)))
	DefaultCircuitBreaker   = false
	DefaultCollateralParams = CollateralParams{CollateralParam{
		Denom:            "uftm",
		LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
		DebtLimit:        sdk.NewCoins(sdk.NewCoin(StableDenom, sdk.NewInt(500000000000))),
	}}
	DefaultDebtParams = DebtParams{}
)

// Params governance parameters for cdp module
type Params struct {
	CollateralParams CollateralParams `json:"collateral_params" yaml:"collateral_params"`
	DebtParams       DebtParams       `json:"debt_params" yaml:"debt_params"`
	GlobalDebtLimit  sdk.Coins        `json:"global_debt_limit" yaml:"global_debt_limit"`
	CircuitBreaker   bool             `json:"circuit_breaker" yaml:"circuit_breaker"`
	Fee              fee.Fee          `json:"fee" yaml:"fee"`
	Nominees         []string         `json:"nominees" yaml:"nominees"`
}

func (cps Params) IsCollateralPresent(collateralDenom string) bool {
	// search for matching denom, return
	for _, cp := range cps.CollateralParams {
		if cp.Denom == collateralDenom {
			return true
		}
	}
	return false
}

func (cps Params) GetCollateralParam(collateralDenom string) CollateralParam {
	// search for matching denom, return
	for _, cp := range cps.CollateralParams {
		if cp.Denom == collateralDenom {
			return cp
		}
	}
	// panic if not found, to be safe
	panic("collateral params not found in module params")
}

// String implements fmt.Stringer
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	Global Debt Limit: %s
	Collateral Params: %s
	Debt Params: %s
	Nominees: %s
	Fee: %s
	Circuit Breaker: %t`,
		p.GlobalDebtLimit, p.CollateralParams, p.DebtParams, p.Nominees, p.Fee, p.CircuitBreaker,
	)
}

// NewParams returns a new params object
func NewParams(
	debtLimit sdk.Coins,
	collateralParams CollateralParams,
	debtParams DebtParams,
	breaker bool,
	nominees []string,
	fee fee.Fee,
) Params {
	return Params{
		GlobalDebtLimit:  debtLimit,
		CollateralParams: collateralParams,
		DebtParams:       debtParams,
		CircuitBreaker:   breaker,
		Nominees:         nominees,
		Fee:              fee,
	}
}

// DefaultParams returns default params for cdp module
func DefaultParams() Params {
	return NewParams(
		DefaultGlobalDebt,
		DefaultCollateralParams,
		DefaultDebtParams,
		DefaultCircuitBreaker,
		[]string{},
		fee.NewDefaultFee(),
	)
}

type CollateralParam struct {
	Denom            string    `json:"denom" yaml:"denom"`                         // Coin name of collateral type
	LiquidationRatio sdk.Dec   `json:"liquidation_ratio" yaml:"liquidation_ratio"` // The ratio (Collateral (priced in stable coin) / Debt) under which a CSDT will be liquidated
	DebtLimit        sdk.Coins `json:"debt_limit" yaml:"debt_limit"`               // Maximum amount of debt allowed to be drawn from this collateral type
	//DebtFloor        sdk.Int // used to prevent dust
}

// String implements fmt.Stringer
func (cp CollateralParam) String() string {
	return fmt.Sprintf(`Collateral:
	Denom: %s
	LiquidationRatio: %s
	DebtLimit: %s`, cp.Denom, cp.LiquidationRatio, cp.DebtLimit)
}

// CollateralParams array of CollateralParam
type CollateralParams []CollateralParam

// String implements fmt.Stringer
func (cps CollateralParams) String() string {
	out := "Collateral Params\n"
	for _, cp := range cps {
		out += fmt.Sprintf("%s\n", cp)
	}
	return out
}

// DebtParam governance params for debt assets
type DebtParam struct {
	Denom          string    `json:"denom" yaml:"denom"`
	ReferenceAsset string    `json:"reference_asset" yaml:"reference_asset"`
	DebtLimit      sdk.Coins `json:"debt_limit" yaml:"debt_limit"`
}

func (dp DebtParam) String() string {
	return fmt.Sprintf(`Debt:
	Denom: %s
	ReferenceAsset: %s
	DebtLimit: %s`, dp.Denom, dp.ReferenceAsset, dp.DebtLimit)
}

// DebtParams array of DebtParam
type DebtParams []DebtParam

// String implements fmt.Stringer
func (dps DebtParams) String() string {
	out := "Debt Params\n"
	for _, dp := range dps {
		out += fmt.Sprintf("%s\n", dp)
	}
	return out
}

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyGlobalDebtLimit, Value: &p.GlobalDebtLimit},
		{Key: KeyCollateralParams, Value: &p.CollateralParams},
		{Key: KeyDebtParams, Value: &p.DebtParams},
		{Key: KeyCircuitBreaker, Value: &p.CircuitBreaker},
		{Key: KeyNominees, Value: &p.Nominees},
		{Key: KeyFee, Value: &p.Fee},
	}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	debtDenoms := make(map[string]int)
	debtParamsDebtLimit := sdk.Coins{}
	for _, dp := range p.DebtParams {
		_, found := debtDenoms[dp.Denom]
		if found {
			return fmt.Errorf("duplicate debt denom: %s", dp.Denom)
		}
		debtDenoms[dp.Denom] = 1
		if dp.DebtLimit.IsAnyNegative() {
			return fmt.Errorf("debt limit for all debt tokens should be positive, is %s for %s", dp.DebtLimit, dp.Denom)
		}
		debtParamsDebtLimit = debtParamsDebtLimit.Add(dp.DebtLimit)
	}
	if debtParamsDebtLimit.IsAnyGT(p.GlobalDebtLimit) {
		return fmt.Errorf("debt limit exceeds global debt limit:\n\tglobal debt limit: %s\n\tdebt limits: %s",
			p.GlobalDebtLimit, debtParamsDebtLimit)
	}

	collateralDupMap := make(map[string]int)
	collateralParamsDebtLimit := sdk.Coins{}
	for _, cp := range p.CollateralParams {
		_, found := collateralDupMap[cp.Denom]
		if found {
			return fmt.Errorf("duplicate collateral denom: %s", cp.Denom)
		}
		collateralDupMap[cp.Denom] = 1

		if cp.DebtLimit.IsAnyNegative() {
			return fmt.Errorf("debt limit for all collaterals should be positive, is %s for %s", cp.DebtLimit, cp.Denom)
		}
		collateralParamsDebtLimit = collateralParamsDebtLimit.Add(cp.DebtLimit)
	}
	if collateralParamsDebtLimit.IsAnyGT(p.GlobalDebtLimit) {
		return fmt.Errorf("collateral debt limit exceeds global debt limit:\n\tglobal debt limit: %s\n\tcollateral debt limits: %s",
			p.GlobalDebtLimit, collateralParamsDebtLimit)
	}

	if p.GlobalDebtLimit.IsAnyNegative() {
		return fmt.Errorf("global debt limit should be positive for all debt tokens, is %s", p.GlobalDebtLimit)
	}
	return nil
}

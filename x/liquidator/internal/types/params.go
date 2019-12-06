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
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Parameter keys
var (
	KeyDebtAuctionSize  = []byte("DebtAuctionSize")
	KeyCollateralParams = []byte("CollateralParams")
)

// LiquidatorParams store params for the liquidator module
type LiquidatorParams struct {
	DebtAuctionSize sdk.Int `json:"debt_auction_size" yaml:"debt_auction_size"`
	//SurplusAuctionSize sdk.Int
	CollateralParams []CollateralParams `json:"collateral_params" yaml:"collateral_params"`
}

// NewLiquidatorParams returns a new params object for the liquidator module
func NewLiquidatorParams(debtAuctionSize sdk.Int, collateralParams []CollateralParams) LiquidatorParams {
	return LiquidatorParams{
		DebtAuctionSize:  debtAuctionSize,
		CollateralParams: collateralParams,
	}
}

// String implements fmt.Stringer
func (p LiquidatorParams) String() string {
	out := fmt.Sprintf(`Params:
		Debt Auction Size: %s
		Collateral Params: `,
		p.DebtAuctionSize,
	)
	for _, cp := range p.CollateralParams {
		out += fmt.Sprintf(`
		%s`, cp.String())
	}
	return out
}

// CollateralParams params storing information about each collateral for the liquidator module
type CollateralParams struct {
	Denom       string  `json:"denom" yaml:"denom"`
	AuctionSize sdk.Int `json:"auction_size" yaml:"auction_size"`
	// LiquidationPenalty
}

// String implements stringer interface
func (cp CollateralParams) String() string {
	return fmt.Sprintf(`
  Denom:        %s
  AuctionSize: %s`, cp.Denom, cp.AuctionSize)
}

// ParamKeyTable for the liquidator module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&LiquidatorParams{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of liquidator module's parameters.
// nolint
func (p *LiquidatorParams) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		subspace.NewParamSetPair(KeyDebtAuctionSize, &p.DebtAuctionSize),
		subspace.NewParamSetPair(KeyCollateralParams, &p.CollateralParams),
	}
}

// DefaultParams for the liquidator module
func DefaultParams() LiquidatorParams {
	return LiquidatorParams{
		DebtAuctionSize:  sdk.NewInt(1000),
		CollateralParams: []CollateralParams{},
	}
}

func (p LiquidatorParams) Validate() error {
	if p.DebtAuctionSize.IsNegative() {
		return fmt.Errorf("debt auction size should be positive, is %s", p.DebtAuctionSize)
	}
	denomDupMap := make(map[string]int)
	for _, cp := range p.CollateralParams {
		_, found := denomDupMap[cp.Denom]
		if found {
			return fmt.Errorf("duplicate denom: %s", cp.Denom)
		}
		denomDupMap[cp.Denom] = 1
		if cp.AuctionSize.IsNegative() {
			return fmt.Errorf(
				"auction size for each collateral should be positive, is %s for %s", cp.AuctionSize, cp.Denom,
			)
		}
	}
	return nil
}

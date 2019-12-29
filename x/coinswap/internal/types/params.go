/*

Copyright 2016 All in Bits, Inc
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
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName
)

// Parameter store keys
var (
	KeyNativeDenom = []byte("nativeDenom")
	KeyFee         = []byte("fee")
)

// Params defines the fee and native denomination for coinswap
type Params struct {
	NativeDenom string   `json:"native_denom"`
	Fee         FeeParam `json:"fee"`
}

func NewParams(nativeDenom string, fee FeeParam) Params {
	return Params{
		NativeDenom: nativeDenom,
		Fee:         fee,
	}
}

// FeeParam defines the numerator and denominator used in calculating the
// amount to be reserved as a liquidity fee.
// TODO: come up with a more descriptive name than Numerator/Denominator
// Fee = 1 - (Numerator / Denominator) TODO: move this to spec
type FeeParam struct {
	Numerator   sdk.Int `json:"fee_numerator"`
	Denominator sdk.Int `json:"fee_denominator"`
}

func NewFeeParam(numerator, denominator sdk.Int) FeeParam {
	return FeeParam{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

// ParamKeyTable returns the KeyTable for coinswap module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyNativeDenom, &p.NativeDenom},
		{KeyFee, &p.Fee},
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
Native Denom:	%s
Fee:			%s`, p.NativeDenom, p.Fee,
	)
}

// String returns a decimal representation of the parameters.
func (fp FeeParam) String() sdk.Dec {
	feeN := sdk.NewDecFromInt(fp.Numerator)
	feeD := sdk.NewDecFromInt(fp.Denominator)
	fee := sdk.OneDec().Sub((feeN.Quo(feeD)))
	return fee
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	feeParam := NewFeeParam(sdk.NewInt(997), sdk.NewInt(1000))

	return Params{
		NativeDenom: sdk.DefaultBondDenom,
		Fee:         feeParam,
	}
}

var (
	// Denominations can be 3 ~ 16 characters long.
	reDnmString = `[a-z][a-z0-9]{2,15}`
	reDnm       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reDnmString))
)

// ValidateParams validates a set of params
func ValidateParams(p Params) error {
	if strings.TrimSpace(p.NativeDenom) == "" {
		return fmt.Errorf("native denomination must not be empty")
	}
	if !reDnm.MatchString(p.NativeDenom) {
		return fmt.Errorf("invalid denom: %s", p.NativeDenom)
	}
	if !p.Fee.Numerator.IsPositive() {
		return fmt.Errorf("fee numerator is not positive: %v", p.Fee.Numerator)
	}
	if !p.Fee.Denominator.IsPositive() {
		return fmt.Errorf("fee denominator is not positive: %v", p.Fee.Denominator)
	}
	if p.Fee.Numerator.GTE(p.Fee.Denominator) {
		return fmt.Errorf("fee numerator is greater than or equal to fee numerator")
	}
	return nil
}

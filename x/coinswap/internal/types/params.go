/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultParamSpace for coinswap
	DefaultParamSpace = ModuleName
	MaxFeePrecision   = 10
)

// Parameter store keys
var (
	feeKey = []byte("fee")
)

// Params defines the fee and native denomination for coinswap
type Params struct {
	Fee sdk.Coins `json:"fee"`
}

// NewParams coinswap params constructor
func NewParams(fee sdk.Coins) Params {
	return Params{
		Fee: fee,
	}
}

// ParamTypeTable returns the TypeTable for coinswap module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Coinswap Params:
  Fee:			%s`, p.Fee.String(),
	)
}

// GetParamSpace Implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

// KeyValuePairs  Implements params.KeyValuePairs
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: feeKey, Value: &p.Fee},
	}
}

// Validate Implements params.Validate
func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(feeKey):
		fee, err := sdk.NewRatFromDecimal(value, MaxFeePrecision)
		if err != nil {
			return nil, err
		}
		if err := validateFee(fee); err != nil {
			return nil, err
		}
		return fee, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

// StringFromBytes Implements params.StringFromBytes
func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// ReadOnly Implements params.ReadOnly
func (p *Params) ReadOnly() bool {
	return false
}

// DefaultParams returns the default coinswap module parameters
func DefaultParams() Params {
	fee := sdk.NewRat(3, 1000)
	return Params{
		Fee: fee,
	}
}

// ValidateParams validates a set of params
func ValidateParams(p Params) error {
	return validateFee(p.Fee)
}

func validateFee(fee sdk.Rat) sdk.Error {
	if !fee.GT(sdk.ZeroRat()) {
		return sdk.ParseParamsErr(fmt.Errorf("fee is not positive: %s", fee.String()))
	}

	if !fee.LT(sdk.OneRat()) {
		return sdk.ParseParamsErr(fmt.Errorf("fee must be less than 1: %s", fee.String()))
	}
	return nil
}

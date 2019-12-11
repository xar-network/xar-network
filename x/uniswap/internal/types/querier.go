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
	"strings"
)

const (
	// Query endpoints supported by the coinswap querier
	QueryLiquidity  = "liquidity"
	QueryParameters = "parameters"

	ParamFee         = "fee"
	ParamNativeDenom = "nativeDenom"
)

var paramList = []string{ParamFee, ParamNativeDenom}

// defines the params for the following queries:
// - 'custom/coinswap/liquidity'
type QueryLiquidityParams struct {
	NonNativeDenom string
}

// Params used for querying liquidity
func NewQueryLiquidityParams(nonNativeDenom string) QueryLiquidityParams {
	return QueryLiquidityParams{
		NonNativeDenom: strings.TrimSpace(nonNativeDenom),
	}
}

func (l QueryLiquidityParams) String() string {
	return l.NonNativeDenom
}

// return if parameter is present at the paramList
func ParamIsValid(p string) bool {
	for _, v := range paramList {
		if v == p {
			return true
		}
	}
	return false
}

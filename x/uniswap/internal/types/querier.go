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
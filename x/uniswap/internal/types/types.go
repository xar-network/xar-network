package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "swap"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

type Pool struct {
	BalanceCoin  sdk.Coin `json:"balance_coin"` // intermediation for exchange tokens
	BalanceToken sdk.Coin `json:"balance_token"`
}

func NewPool(balanceCoin sdk.Coin, balanceToken sdk.Coin) Pool {
	return Pool{
		BalanceCoin:  balanceCoin,
		BalanceToken: balanceToken,
	}
}

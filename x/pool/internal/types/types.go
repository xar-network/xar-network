package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

	DefaultParamspace = ModuleName

	QueryTotalFunds = "funds"
	QueryReadFunds  = "get-funds"
)

type QueryFundsParams struct {
	Owner sdk.AccAddress
}

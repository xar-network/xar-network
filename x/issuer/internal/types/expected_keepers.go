package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type InterestKeeper interface {
	SetInterest(sdk.Context, sdk.Dec, string) sdk.Result
}

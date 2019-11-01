package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/interest/internal/types"
)

// NewQuerier returns an inflation Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {

		case types.QueryInterest:
			return queryInterest(ctx, k)

		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown inflation query endpoint: %s", path[0]))
		}
	}
}

func queryInterest(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	interestState := k.GetState(ctx)

	// TODO Introduce a more presentation-friendly response type
	res, err := types.ModuleCdc.MarshalJSON(interestState)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}

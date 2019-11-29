package market

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/market/types"
)

// NewHandler handles all oracle type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgCreateMarket:
			return k.CreateMarket(ctx, msg)
		default:
			return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized market message type: %T", msg)).Result()
		}
	}
}

package denominations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/denominations/internal/keeper"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgMint:
			return k.Mint(ctx, msg)
		case types.MsgBurn:
			return k.Burn(ctx, msg)
		default:
			return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized market message type: %T", msg)).Result()
		}
	}
}

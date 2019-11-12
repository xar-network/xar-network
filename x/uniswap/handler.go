package uniswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper, bk bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSwap:
			return handleMsgSwap(ctx, k, bk, msg)
		default:
			errMsg := "Unrecognized swap Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSwap(ctx sdk.Context, keeper Keeper, bk bank.Keeper, msg types.MsgSwap) sdk.Result {
	_, err := bk.SubtractCoins(ctx, msg.Sender, sdk.Coins{msg.Asset})
	if err != nil {
		return err.Result()
	}

	result, err := keeper.Swap(ctx, msg.Asset, msg.TargetDenom)
	if err != nil {
		return err.Result()
	}
	_, err = bk.AddCoins(ctx, msg.Sender, sdk.Coins{result})
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

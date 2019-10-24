package uniswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSwap:
			return handleMsgSwap(ctx, k, msg)
		default:
			errMsg := "Unrecognized swap Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSwap(ctx sdk.Context, keeper Keeper, msg MsgSwap) sdk.Result {
	_, _, err := keeper.bk.SubtractCoins(ctx, msg.Sender, sdk.Coins{msg.Asset})
	if err != nil {
		return err.Result()
	}

	result, tags, err := keeper.Swap(ctx, msg.Asset, msg.TargetDenom)
	if err != nil {
		return err.Result()
	}
	_, _, err = keeper.bk.AddCoins(ctx, msg.Sender, sdk.Coins{result})
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

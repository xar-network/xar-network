/**

Baseline from Kava Cosmos Module

**/

package pricefeed

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/zar-network/zar-network/x/pricefeed/internal/types"
)

// NewHandler handles all pricefeed type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgPostPrice:
			return HandleMsgPostPrice(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized pricefeed message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// price feed questions:
// do proposers need to post the round in the message? If not, how do we determine the round?

// HandleMsgPostPrice handles prices posted by oracles
func HandleMsgPostPrice(
	ctx sdk.Context,
	k Keeper,
	msg types.MsgPostPrice) sdk.Result {

	// TODO cleanup message validation and errors
	err := k.ValidatePostPrice(ctx, msg)
	if err != nil {
		return err.Result()
	}
	k.SetPrice(ctx, msg.From, msg.AssetCode, msg.Price, msg.Expiry)
	return sdk.Result{}
}

// EndBlocker updates the current pricefeed
func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
	// TODO val_state_change.go is relevant if we want to rotate the oracle set

	// Running in the end blocker ensures that prices will update at most once per block,
	// which seems preferable to having state storage values change in response to multiple transactions
	// which occur during a block
	//TODO use an iterator and update the prices for all assets in the store
	k.SetCurrentPrices(ctx)
	return []abci.ValidatorUpdate{}
}

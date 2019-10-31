/**

Baseline from Kava Cosmos Module

**/

package auction

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/auction/internal/keeper"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

// NewHandler returns a function to handle all "auction" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgPlaceBid:
			return handleMsgPlaceBid(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized auction msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgPlaceBid(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgPlaceBid) sdk.Result {

	err := keeper.PlaceBid(ctx, msg.AuctionID, msg.Bidder, msg.Bid, msg.Lot)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

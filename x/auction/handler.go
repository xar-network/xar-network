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

// EndBlocker runs at the end of every block.
func EndBlocker(ctx sdk.Context, k Keeper) sdk.Result {

	// get an iterator of expired auctions
	expiredAuctions := k.GetQueueIterator(ctx, types.EndTime(ctx.BlockHeight()))
	defer expiredAuctions.Close()

	// loop through and close them - distribute funds, delete from store (and queue)
	for ; expiredAuctions.Valid(); expiredAuctions.Next() {
		var auctionID types.ID
		ModuleCdc.MustUnmarshalBinaryLengthPrefixed(expiredAuctions.Value(), &auctionID)

		err := k.CloseAuction(ctx, auctionID)
		if err != nil {
			panic(err) // TODO how should errors be handled here?
		}
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

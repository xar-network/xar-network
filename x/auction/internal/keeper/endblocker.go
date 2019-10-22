package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zar-network/zar-network/x/auction/internal/types"
)

// EndBlocker runs at the end of every block.
func EndBlocker(ctx sdk.Context, k Keeper) sdk.Result {

	// get an iterator of expired auctions
	expiredAuctions := k.getQueueIterator(ctx, types.EndTime(ctx.BlockHeight()))
	defer expiredAuctions.Close()

	// loop through and close them - distribute funds, delete from store (and queue)
	for ; expiredAuctions.Valid(); expiredAuctions.Next() {
		var auctionID types.ID
		k.cdc.MustUnmarshalBinaryLengthPrefixed(expiredAuctions.Value(), &auctionID)

		err := k.CloseAuction(ctx, auctionID)
		if err != nil {
			panic(err) // TODO how should errors be handled here?
		}
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

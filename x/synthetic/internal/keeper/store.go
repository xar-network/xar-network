package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/synthetic/internal/types"
)

func (k Keeper) SetMarketBalance(ctx sdk.Context, mb types.MarketBalance) {
	balances := k.GetMarketBalances(ctx)
	balances.SetMarketBalance(mb)

	k.SetMarketBalances(ctx, balances)
}

func (k Keeper) GetMarketBalance(ctx sdk.Context, denom string) (mb types.MarketBalance, found bool) {
	balances := k.GetMarketBalances(ctx)

	return balances.GetMarketBalance(denom)
}

func (k Keeper) GetMarketBalances(ctx sdk.Context) (mb types.MarketBalances) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(getMarketBalanceStoreKey())
	if bz == nil {
		return mb
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &mb)
	return mb
}

func (k Keeper) SetMarketBalances(ctx sdk.Context, mb types.MarketBalances) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(mb)
	store.Set(getMarketBalanceStoreKey(), bz)
}

func (k Keeper) MustGetMarketBalance(ctx sdk.Context, denom string) types.MarketBalance {
	mb, found := k.GetMarketBalance(ctx, denom)
	if !found {
		// this is not going to happen
		panic("market balance not set")
	}

	return mb
}

func getMarketBalanceStoreKey() []byte {
	return []byte{0x01}
}

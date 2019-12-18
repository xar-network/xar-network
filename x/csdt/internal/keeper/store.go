package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	LastBlockKey = []byte{0x00} // key for the last interest accrual block
)

// GetLastAccrualBlock gets the last time of interest accrual
func (k Keeper) GetLastAccrualBlock(ctx sdk.Context) (lastBlock int64) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(LastBlockKey)
	if b == nil {
		panic("previous interest accrual block not set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &lastBlock)
	return
}

// SetLastAccrualBlock sets the last time of interest accrual
func (k Keeper) SetLastAccrualBlock(ctx sdk.Context, lastBlock int64) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(lastBlock)
	store.Set(LastBlockKey, b)
}

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

var reservePoolStorePrefix = []byte{0x01}

// creates reserve pool and returns it as a response
func (keeper Keeper) CreateReservePool(ctx sdk.Context, nativeDenom, nonNativeDenom string) types.ReservePool {
	poolName, err := keeper.GetPoolName(nativeDenom, nonNativeDenom)
	if err != nil {
		panic(err)
	}
	pool := types.NewReservePool(nativeDenom, nonNativeDenom, poolName)
	keeper.SetReservePool(ctx, pool)

	return pool
}

func (keeper Keeper) GetReservePool(ctx sdk.Context, poolName string) (rp types.ReservePool, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	key := reservePoolKey(poolName, reservePoolStorePrefix)
	value := store.Get(key)
	if value == nil {
		return
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(value, &rp)
	return rp, true
}


func (keeper Keeper) SetReservePool(ctx sdk.Context, pool types.ReservePool) {
	store := ctx.KVStore(keeper.storeKey)
	value := keeper.cdc.MustMarshalBinaryLengthPrefixed(pool)
	key := reservePoolKey(pool.GetName(), reservePoolStorePrefix)
	store.Set(key, value)
}

func reservePoolKey(poolName string, poolKeyPrefix []byte) []byte {
	return append(poolKeyPrefix, []byte(poolName)...)
}

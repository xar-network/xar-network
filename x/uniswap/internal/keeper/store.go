package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/uniswap/internal/types/pool"
)

var reservePoolStorePrefix = []byte{0x01}

// creates reserve pool and returns it as a response
func (keeper Keeper) CreateReservePool(ctx sdk.Context, nativeDenom, nonNativeDenom string) pool.ReservePool {
	poolName, err := keeper.GetPoolName(nativeDenom, nonNativeDenom)
	if err != nil {
		panic(err)
	}
	resPool := pool.NewReservePool(nativeDenom, nonNativeDenom, poolName)
	keeper.SetReservePool(ctx, resPool)

	return resPool
}

func (keeper Keeper) GetReservePool(ctx sdk.Context, poolName string) (rp pool.ReservePool, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	key := reservePoolKey(poolName, reservePoolStorePrefix)
	value := store.Get(key)
	if value == nil {
		return
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(value, &rp)
	return rp, true
}

func (keeper Keeper) CreateOrGetReservePool(ctx sdk.Context, nonNativeDenom string) pool.ReservePool {
	nativeDenom := keeper.GetNativeDenom(ctx)

	poolName := keeper.MustGetPoolName(nativeDenom, nonNativeDenom)
	rp, found := keeper.GetReservePool(ctx, poolName)
	if !found {
		rp = keeper.CreateReservePool(ctx, nativeDenom, nonNativeDenom)
	}
	return rp
}

func (keeper Keeper) SetReservePool(ctx sdk.Context, pool pool.ReservePool) {
	store := ctx.KVStore(keeper.storeKey)
	value := keeper.cdc.MustMarshalBinaryLengthPrefixed(pool)
	key := reservePoolKey(pool.GetName(), reservePoolStorePrefix)
	store.Set(key, value)
}

func reservePoolKey(poolName string, poolKeyPrefix []byte) []byte {
	return append(poolKeyPrefix, []byte(poolName)...)
}
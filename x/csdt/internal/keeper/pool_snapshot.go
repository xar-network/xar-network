package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sync"
)

type PoolSnapshot struct {
	store sdk.KVStore

	cache map[string]sdk.Int
	cacheMutex	sync.Mutex
}

func NewPoolSnapshot(ctx sdk.Context, denom string) *PoolSnapshot {
	return &PoolSnapshot{
		store: ctx.KVStore(sdk.NewKVStoreKey("PoolsSnapshots-" + denom)),
		cache: make(map[string]sdk.Int),
		cacheMutex: sync.Mutex{},
	}
}

func (ps *PoolSnapshot) periodKey(period string) []byte {
	return []byte(period)
}

func (ps *PoolSnapshot) SetPool(period string, pool sdk.Int) {
	ps.store.Set(ps.periodKey(period), []byte(pool.String()))

	ps.setCache(period, pool)
}

func (ps *PoolSnapshot) GetPool(period string) *sdk.Int {
	cachePool := ps.getCache(period)
	if cachePool != nil {
		return cachePool
	}

	valStr := ps.store.Get(ps.periodKey(period))
	if valStr == nil {
		return nil
	}

	valInt, ok := sdk.NewIntFromString(string(valStr))
	if ok {
		return &valInt
	}

	return nil
}

func (ps *PoolSnapshot) setCache(period string, pool sdk.Int) {
	ps.cacheMutex.Lock()
	defer ps.cacheMutex.Unlock()

	ps.cache[period] = pool
}

func (ps *PoolSnapshot) getCache(period string) *sdk.Int {
	ps.cacheMutex.Lock()
	defer ps.cacheMutex.Unlock()

	valInt, ok := ps.cache[period]
	if !ok {
		return nil
	}

	return &valInt
}


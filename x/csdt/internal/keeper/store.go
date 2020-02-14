package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	LastBlockKey = []byte{0x00} // key for the last interest accrual block
)

// Wrappers
func (k Keeper) getTotalBorrowsKey(collateralDenom string) []byte {
	return bytes.Join(
		[][]byte{
			[]byte("borrow"),
			[]byte(collateralDenom),
		},
		nil, // no separator
	)
}

func (k Keeper) getTotalCashKey(collateralDenom string) []byte {
	return bytes.Join(
		[][]byte{
			[]byte("cash"),
			[]byte(collateralDenom),
		},
		nil, // no separator
	)
}

func (k Keeper) getTotalReserveKey(collateralDenom string) []byte {
	return bytes.Join(
		[][]byte{
			[]byte("reserve"),
			[]byte(collateralDenom),
		},
		nil, // no separator
	)
}

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

// GetTotalBorrows gets the global borrows for a specific denomination
func (k Keeper) GetTotalBorrows(ctx sdk.Context, collateralDenom string) (sdk.Uint, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(k.getTotalBorrowsKey(collateralDenom))
	// unmarshal
	if bz == nil {
		return sdk.ZeroUint(), false
	}
	var cash sdk.Uint
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &cash)
	return cash, true
}

// GetTotalCash gets the global cash for a specific denomination
func (k Keeper) GetTotalCash(ctx sdk.Context, collateralDenom string) (sdk.Uint, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(k.getTotalCashKey(collateralDenom))
	// unmarshal
	if bz == nil {
		return sdk.ZeroUint(), false
	}
	var cash sdk.Uint
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &cash)
	return cash, true
}

// GetTotalReserve gets the global reserve value for a specific denomination
func (k Keeper) GetTotalReserve(ctx sdk.Context, collateralDenom string) (sdk.Uint, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(k.getTotalReserveKey(collateralDenom))
	// unmarshal
	if bz == nil {
		return sdk.ZeroUint(), false
	}
	var reserve sdk.Uint
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &reserve)
	return reserve, true
}

// SetLastAccrualBlock sets the last time of interest accrual
func (k Keeper) SetLastAccrualBlock(ctx sdk.Context, lastBlock int64) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(lastBlock)
	store.Set(LastBlockKey, b)
}

// SetTotalBorrows stores the global borrow value for a specific denomination
func (k Keeper) SetTotalBorrows(ctx sdk.Context, totalBorrows sdk.Uint, collateralDenom string) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalBorrows)
	store.Set(k.getTotalBorrowsKey(collateralDenom), bz)
}

// SetTotalCash stores the global cash value for a specific denomination
func (k Keeper) SetTotalCash(ctx sdk.Context, totalCash sdk.Uint, collateralDenom string) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalCash)
	store.Set(k.getTotalCashKey(collateralDenom), bz)
}

// SetTotalReserve stores the global reserve value for a specific denomination
func (k Keeper) SetTotalReserve(ctx sdk.Context, totalReserve sdk.Uint, collateralDenom string) {
	// get store
	store := ctx.KVStore(k.storeKey)
	// marshal and set
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(totalReserve)
	store.Set(k.getTotalReserveKey(collateralDenom), bz)
}

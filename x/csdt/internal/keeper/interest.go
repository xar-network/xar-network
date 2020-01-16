package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccrueInterest
func (k Keeper) AccrueInterest(ctx sdk.Context, lastAccruedBlock int64) {
	// logger := k.Logger(ctx)
	//
	// totalCash := sdk.NewInt(0) // k.GetTotalCash(ctx) // TODO: get balance

	// get the interest rate (that was in effect since the last update):
	// ic := types.NewCsdtInterest(sdk.OneUint())
	// totalBorrows := k.GetGlobalDebt(ctx)
	// totalReserves := k.GetGlobalReserves(ctx)
	// borrowRate := ic.GetBorrowRate(ctx, totalCash, totalBorrows, totalReserves)
	// blockDelta := uint64(ctx.BlockHeight() - lastAccruedBlock)
	// if blockDelta == 0 {
	// 	logger.Error("Failed to get accurate block delta")
	// 	return
	// }
	// simpleInterestFactor := borrowRate.MulUint64(blockDelta)

	// update borrowIndex :
	// borrowIndexNew = borrowIndex × (1 + simpleInterestFactor)

	// calculate the interest accrued:
	// interestAccumulated = totalBorrows × simpleInterestFactor

	// We update borrows and reserves:
	// totalBorrowsNew = totalBorrows + interestAccumulated
	// totalReservesNew = totalReserves + interestAccumulated × reserveFactor

	// We store the updates back to the blockchain:
	// Set accrualBlockNumber = getBlockNumber()
	// Set borrowIndex = borrowIndexNew
	// Set totalBorrows = totalBorrowsNew
	// Set totalReserves = totalReservesNew

	k.SetLastAccrualBlock(ctx, ctx.BlockHeight())
}

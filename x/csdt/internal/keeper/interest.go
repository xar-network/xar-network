package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

func (k Keeper) getAccrualBlockAndIndex(ctx sdk.Context, collateralDenom string) (lastAccrualBlock int64, borrowIndex sdk.Uint) {
	lastAccrualBlock, ok := k.GetLastAccrualBlock(ctx, collateralDenom)
	if !ok {
		panic(fmt.Sprintf("failed to get last accrual block for '%v'", collateralDenom)) // This should exist already
	}

	borrowIndex, ok = k.GetBorrowIndex(ctx, collateralDenom)
	if !ok {
		panic(fmt.Sprintf("failed to get borrow index for '%v'", collateralDenom)) // This should exist already
	}
	return
}

func (k Keeper) getInterestRateModel(ctx sdk.Context, collateralDenom string) types.InterestRateModel {
	p := k.GetParams(ctx).GetCollateralParam(collateralDenom)
	return p.InterestModel
}

func (k Keeper) getTotals(ctx sdk.Context, collateralDenom string) (
	totalBorrows sdk.Uint, totalCash sdk.Uint, globalReserves sdk.Uint) {

	totalCash, ok := k.GetTotalCash(ctx, collateralDenom)
	if !ok {
		panic(fmt.Sprintf("failed to get global cash value for '%v'", collateralDenom)) // This should exist already
	}

	totalBorrows, ok = k.GetTotalBorrows(ctx, collateralDenom)
	if !ok {
		panic(fmt.Sprintf("failed to get global borrow value for '%v'", collateralDenom)) // This should exist already
	}

	globalReserves, ok = k.GetTotalReserve(ctx, collateralDenom)
	if !ok {
		panic(fmt.Sprintf("failed to get global reserve value for '%v'", collateralDenom)) // This should exist already
	}

	return totalCash, totalBorrows, globalReserves
}

// AccrueInterest accrues interest and updates the borrow index on every operation.
// This increases compounding, approaching the true value, regardless of whether the rest of the operation succeeds or not
func (k Keeper) AccrueInterest(ctx sdk.Context, collateralDenom string, reserveFactor sdk.Uint) {
	logger := k.Logger(ctx)

	lastAccruedBlock, borrowIndex := k.getAccrualBlockAndIndex(ctx, collateralDenom)

	// get the interest rate (that was in effect since the last update):
	ic := k.getInterestRateModel(ctx, collateralDenom)
	lastTotalSupply, totalBorrows, totalReserves := k.getTotals(ctx, collateralDenom)
	borrowRate := ic.GetBorrowRate(lastTotalSupply, totalBorrows, totalReserves)

	blockDelta := uint64(ctx.BlockHeight() - lastAccruedBlock)
	if blockDelta == 0 {
		logger.Error("Failed to get accurate block delta")
		return
	}
	simpleInterestFactor := borrowRate.MulUint64(blockDelta)

	// update borrowIndex:
	// borrowIndexNew = borrowIndex × (1 + simpleInterestFactor)
	borrowIndexNew := borrowIndex.Mul(sdk.OneUint().Add(simpleInterestFactor))

	// calculate the interest accrued:
	// interestAccumulated = totalBorrows × simpleInterestFactor
	interestAccumulated := totalBorrows.Mul(simpleInterestFactor)

	// update borrows and reserves:
	// totalBorrowsNew = totalBorrows + interestAccumulated
	// totalReservesNew = totalReserves + interestAccumulated × reserveFactor
	totalBorrowsNew := totalBorrows.Add(interestAccumulated)
	totalReservesNew := totalReserves.Add(interestAccumulated.Mul(reserveFactor))

	// store the updates back to the blockchain:
	k.SetLastAccrualBlock(ctx, ctx.BlockHeight(), collateralDenom)
	k.SetBorrowIndex(ctx, borrowIndexNew, collateralDenom)
	k.SetTotalBorrows(ctx, totalBorrowsNew, collateralDenom)
	k.SetTotalReserve(ctx, totalReservesNew, collateralDenom)
}

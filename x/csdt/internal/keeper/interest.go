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

func (k Keeper) getReserveFactor(ctx sdk.Context, collateralDenom string) sdk.Uint {
	p := k.GetParams(ctx).GetCollateralParam(collateralDenom)
	return p.ReserveFactor
}

func (k Keeper) getTotals(ctx sdk.Context, collateralDenom string) (
	totalCash sdk.Uint, totalBorrows sdk.Uint, globalReserves sdk.Uint) {

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
func (k Keeper) AccrueInterest(ctx sdk.Context, collateralDenom string) {
	logger := k.Logger(ctx)

	reserveFactorMantissa := k.getReserveFactor(ctx, collateralDenom)
	currentBlockNumber := ctx.BlockHeight()
	lastAccruedBlock, borrowIndex := k.getAccrualBlockAndIndex(ctx, collateralDenom)

	// get the interest rate (that was in effect since the last update):
	ic := k.getInterestRateModel(ctx, collateralDenom)
	lastTotalSupply, totalBorrows, totalReserves := k.getTotals(ctx, collateralDenom)
	borrowRateMantissa := ic.GetBorrowRate(lastTotalSupply, totalBorrows, totalReserves)
	if borrowRateMantissa.GT(types.BorrowRateMaxMantissa()) {
		logger.Error("borrow rate is absurdly high")
	}

	// Calculate the number of blocks elapsed since the last accrual
	blockDelta := uint64(currentBlockNumber - lastAccruedBlock)
	if blockDelta == 0 {
		logger.Error("Failed to get accurate block delta")
		return
	}
	simpleInterestFactor := types.NewExp(borrowRateMantissa).MultiplyScalarUint64(blockDelta)

	// Calculate the interest accumulated into borrows and reserves and the new index:
	// interestAccumulated = totalBorrows × simpleInterestFactor
	interestAccumulated := simpleInterestFactor.MultiplyScalarTruncate(totalBorrows)

	// update borrows and reserves:
	// totalBorrowsNew = totalBorrows + interestAccumulated
	// totalReservesNew = interestAccumulated × reserveFactor + totalReserves
	totalBorrowsNew := totalBorrows.Add(interestAccumulated)
	totalReservesNew := types.NewExp(reserveFactorMantissa).MultiplyScalarTruncateAddUInt(
		interestAccumulated, totalReserves)

	// update borrowIndex:
	// borrowIndexNew = simpleInterestFactor * borrowIndex + borrowIndex
	borrowIndexNew := simpleInterestFactor.MultiplyScalarTruncateAddUInt(borrowIndex, borrowIndex)

	// store the updates back to the blockchain:
	k.SetLastAccrualBlock(ctx, currentBlockNumber, collateralDenom)
	k.SetBorrowIndex(ctx, borrowIndexNew, collateralDenom)
	k.SetTotalBorrows(ctx, totalBorrowsNew, collateralDenom)
	k.SetTotalReserve(ctx, totalReservesNew, collateralDenom)
}

// map collateral denomination to array of [borrow, supply] rates
type borrowSupplyArray = [2]sdk.Uint
type ratesBorrowSupplyMap = map[string]borrowSupplyArray

const (
	keyBorrowRate = 0
	keySupplyRate = 1
)

func (k Keeper) getBorrowSupplyRates(ctx sdk.Context) ratesBorrowSupplyMap {
	params := k.GetParams(ctx).CollateralParams
	rates := make(ratesBorrowSupplyMap, len(params))
	for _, c := range params {
		denom := c.Denom
		irm := k.getInterestRateModel(ctx, denom)
		cash, borrows, reserves := k.getTotals(ctx, denom)
		borrowRate := irm.GetBorrowRate(cash, borrows, reserves)
		supplyRate := irm.GetSupplyRate(cash, borrows, reserves, k.getReserveFactor(ctx, denom))
		rates[denom] = borrowSupplyArray{borrowRate, supplyRate}
	}
	return rates
}

func (k Keeper) adjustCsdtBalances(ctx sdk.Context, csdts types.CSDTs, interestRates ratesBorrowSupplyMap) {
	logger := k.Logger(ctx)
	currentBlock := ctx.BlockHeight()
	for _, csdt := range csdts {
		denom := csdt.CollateralDenom
		borrowRate := interestRates[denom][keyBorrowRate]
		supplyRate := interestRates[denom][keySupplyRate]

		// Adjust for borrows
		if currentBlock > csdt.DebtAccruedBlock {
			debtBalance := sdk.NewUintFromString(csdt.Debt.AmountOf(denom).String())
			if !csdt.Debt.IsZero() {
				interestUint := types.NewExp(borrowRate).MultiplyScalarTruncate(debtBalance)
				interest, ok := sdk.NewIntFromString(interestUint.String())
				if ok {
					csdt.Debt.Add(sdk.NewCoins(sdk.NewCoin(denom, interest)))
					csdt.DebtAccruedBlock = currentBlock // block when debt accrued interest
				} else {
					logger.Error("failed to add borrow interest to csdt for owner: %s. '%s' could not be converted",
						csdt.Owner.String(), interestUint.String())
				}
			}
		}

		// Adjust for collateral
		if currentBlock > csdt.InterestAccruedBlock {
			collateralBalance := sdk.NewUintFromString(csdt.CollateralAmount.AmountOf(denom).String())
			if !csdt.CollateralAmount.IsZero() {
				interestUint := types.NewExp(supplyRate).MultiplyScalarTruncate(collateralBalance)
				interest, ok := sdk.NewIntFromString(interestUint.String())
				if ok {
					csdt.CollateralAmount.Add(sdk.NewCoins(sdk.NewCoin(denom, interest)))
					csdt.Interest.Add(sdk.NewCoins(sdk.NewCoin(denom, interest))) // lifetime accumulated interest
					csdt.InterestAccruedBlock = currentBlock
				} else {
					logger.Error("failed to add supply interest to csdt for owner: %s. '%s' could not be converted",
						csdt.Owner.String(), interestUint.String())
				}
			}
		}
	}
}

func (k Keeper) AdjustBalances(ctx sdk.Context) {
	logger := k.Logger(ctx)
	interestRates := k.getBorrowSupplyRates(ctx)

	csdts, err := k.GetCSDTs(ctx, "", sdk.Dec{})
	if err != nil {
		logger.Error("Failed to get CSDTs to adjust balances because: %v", err)
	}

	k.adjustCsdtBalances(ctx, csdts, interestRates)
}

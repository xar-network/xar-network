package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InterestRateModel describes the contract needed to calculate the Interest Index.
// Track the growth of principal of an arbitrary account over time.
// Using the ratio of that account's interest versus initial principal
// to calculate the growth of any given account's interest over a subset of that time interval. The
// interest model contract specifies the simple interest rate at any moment (which, when
// compounded for each transaction becomes compound interest). We force this interest rate
// model to be a pure function over the cash, borrows, and reserves of an asset
type InterestRateModel interface {
	/**
	 * @notice Calculates the current borrow interest rate per block
	 * @param ctx The context
	 * @param totalCash The total amount of cash available
	 * @param totalBorrows The total amount of borrows outstanding
	 * @param totalReserves The total amount of reserves available
	 * @return The borrow rate per block (as a percentage, and scaled by 1e18)
	 */
	GetBorrowRate(totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint) sdk.Uint

	/**
	 * @notice Calculates the current supply interest rate per block
	 * @param ctx The context
	 * @param totalCash The total amount of cash available
	 * @param totalBorrows The total amount of borrows outstanding
	 * @param totalReserves The total amount of reserves available
	 * @param reserveFactorMantissa The current reserve factor for the market
	 * @return The supply rate per block (as a percentage, and scaled by 1e18)
	 */
	GetSupplyRate(totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint, reserveFactorMantissa sdk.Uint) sdk.Uint
}

// CsdtInterest definition
type CsdtInterest struct {
	MultiplierPerBlock sdk.Uint
	BaseRatePerBlock   sdk.Uint
}

// NewCsdtInterest creates a CSDT Interest model
func NewCsdtInterest(baseRatePerYear sdk.Uint, multiplierPerBlock sdk.Uint) CsdtInterest {
	return CsdtInterest{BaseRatePerBlock: baseRatePerYear, MultiplierPerBlock: multiplierPerBlock}
}

// // InitialExchangeRate returns the initial exchange rate
// func (CsdtInterest) InitialExchangeRate() sdk.Uint { return sdk.NewUint(1) }

// BlocksPerYear is the approximate number of blocks per year that is assumed by the interest rate model
func (CsdtInterest) BlocksPerYear() sdk.Uint { return sdk.NewUint(2102400) }

// func (ci CsdtInterest) ExchangeRateStored(totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint) sdk.Uint {
// 	/*
// 		● Note: we do not assert that the market is up to date.
// 		● If there are no tokens minted:
// 			○ exchangeRate = initial exchange rate
// 		● Otherwise:
// 			○ totalCash = invoke getCash()
// 				■ Note: likely makes an external call
// 			○ exchangeRate = (totalCash + totalBorrows − totalReserves) / totalSupply
// 		● Return exchangeRate
// 	*/
// 	totalSupply := ci.GetSupply().GetSupply(ctx).GetTotal()
// 	if totalSupply.IsZero() {
// 		return ci.InitialExchangeRate()
// 	}
//
// 	totalStableSupplyUint := sdk.NewUintFromBigInt(totalSupply.AmountOf(StableDenom).BigInt())
// 	return totalCash.Add(totalBorrows).Sub(totalReserves).Mul(totalStableSupplyUint)
// }

func pow18() sdk.Uint {
	z := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	return sdk.NewUintFromBigInt(z)
}

/**
 * @notice Calculates the utilization rate of the market: `borrows / (cash + borrows - reserves)`
 * @param cash The amount of cash in the market
 * @param borrows The amount of borrows in the market
 * @param reserves The amount of reserves in the market (amount currently unused)
 * @return The utilization rate as a mantissa between [0, 1e18]
 */
func (CsdtInterest) utilizationRate(totalCash sdk.Uint, totalBorrows sdk.Uint, reserves sdk.Uint) sdk.Uint {
	// Utilization rate is 0 when there are no borrows
	if totalBorrows.IsZero() {
		return sdk.ZeroUint()
	}

	return totalBorrows.Mul(pow18()).Quo(totalCash.Add(totalBorrows).Sub(reserves))
}

// GetBorrowRate calculates the current borrow interest rate per block
func (ci CsdtInterest) GetBorrowRate(totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint) sdk.Uint {
	ur := ci.utilizationRate(totalCash, totalBorrows, totalReserves)
	return ur.Mul(ci.MultiplierPerBlock).Quo(pow18().Add(ci.BaseRatePerBlock))
}

// GetSupplyRate calculates the current supply interest rate per block
func (ci CsdtInterest) GetSupplyRate(totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint,
	reserveFactorMantissa sdk.Uint) sdk.Uint {

	oneMinusReserveFactor := pow18().Sub(reserveFactorMantissa)
	borrowRate := ci.GetBorrowRate(totalCash, totalBorrows, totalReserves)
	rateToPool := borrowRate.Mul(oneMinusReserveFactor).Quo(pow18())
	return ci.utilizationRate(totalCash, totalBorrows, totalReserves).Mul(rateToPool).Quo(pow18())
}

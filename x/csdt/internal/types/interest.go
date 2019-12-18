package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type InterestRateModel interface {
	/**
	 * @notice Calculates the current borrow interest rate per block
	 * @param totalCash The total amount of cash available
	 * @param totalBorrows The total amount of borrows outstanding
	 * @param totalReserves The total amount of reserves available
	 * @return The borrow rate per block (as a percentage, and scaled by 1e18)
	 */
	GetBorrowRate(totalCash sdk.Int, totalBorrows sdk.Int, totalReserves sdk.Int) sdk.Int

	/**
	 * @notice Calculates the current supply interest rate per block
	 * @param totalCash The total amount of cash available
	 * @param totalBorrows The total amount of borrows outstanding
	 * @param totalReserves The total amount of reserves available
	 * @param reserveFactorMantissa The current reserve factor available
	 * @return The supply rate per block (as a percentage, and scaled by 1e18)
	 */
	GetSupplyRate(totalCash sdk.Int, totalBorrows sdk.Int, totalReserves sdk.Int, reserveFactorMantissa sdk.Int) sdk.Int
}

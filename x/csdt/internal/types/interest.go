package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
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
	GetBorrowRate(ctx sdk.Context, totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint) sdk.Uint

	/**
	 * @notice Calculates the current supply interest rate per block
	 * @param ctx The context
	 * @param totalCash The total amount of cash available
	 * @param totalBorrows The total amount of borrows outstanding
	 * @param totalReserves The total amount of reserves available
	 * @param reserveFactor The current reserve factor available
	 * @return The supply rate per block (as a percentage, and scaled by 1e18)
	 */
	GetSupplyRate(ctx sdk.Context, totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint, reserveFactor sdk.Uint) sdk.Uint
}

// CsdtInterest definition
type CsdtInterest struct {
	*keeper.Keeper
	BorrowRate sdk.Uint
}

// NewCsdtInterest creates a CSDT Interest model
func NewCsdtInterest(keeper *keeper.Keeper, borrowRate sdk.Uint) *CsdtInterest {
	return &CsdtInterest{Keeper: keeper, BorrowRate: borrowRate}
}

// InitialExchangeRate returns the initial exchange rate
func (CsdtInterest) InitialExchangeRate() sdk.Uint { return sdk.NewUint(1) }

// BlocksPerYear is the approximate number of blocks per year that is assumed by the interest rate model
func (CsdtInterest) BlocksPerYear() sdk.Uint { return sdk.NewUint(2102400) }

func (ci CsdtInterest) ExchangeRateStored(ctx sdk.Context, totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint) sdk.Uint {
	/*
		● Note: we do not assert that the market is up to date.
		● If there are no tokens minted:
			○ exchangeRate = initial exchange rate
		● Otherwise:
			○ totalCash = invoke getCash()
				■ Note: likely makes an external call
			○ exchangeRate = (totalCash + totalBorrows − totalReserves) / totalSupply
		● Return exchangeRate
	*/
	totalSupply := ci.GetSupply().GetSupply(ctx).GetTotal()
	if totalSupply.IsZero() {
		return ci.InitialExchangeRate()
	}

	totalStableSupplyUint := sdk.NewUintFromBigInt(totalSupply.AmountOf(StableDenom).BigInt())
	return totalCash.Add(totalBorrows).Sub(totalReserves).Mul(totalStableSupplyUint)
}

// GetBorrowRate calculates the current borrow interest rate per block
func (ci CsdtInterest) GetBorrowRate(ctx sdk.Context, totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint) sdk.Uint {
	return ci.BorrowRate
}

// GetSupplyRate calculates the current supply interest rate per block
func (ci CsdtInterest) GetSupplyRate(ctx sdk.Context, totalCash sdk.Uint, totalBorrows sdk.Uint, totalReserves sdk.Uint,
	reserveFactor sdk.Uint) sdk.Uint {
	// TODO: - change to:
	// underlying = totalSupply × exchangeRate
	// borrowsPer = totalBorrows ÷ underlying
	// supplyRate = borrowRate × (1 − reserveFactor) × borrowsPer
	return sdk.OneUint().Sub(reserveFactor).Mul(ci.BorrowRate)
}

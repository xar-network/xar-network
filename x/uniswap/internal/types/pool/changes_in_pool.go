package pool

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// as far as reservePool does not handle transfer coin operations
// and almost all setters assume that incoming amounts can be changed in order to complete operation correct, a "change in pool" variable is made.
// use it for a transfer operations outside the reserve pool.
type ChangeInPool struct {
	NativeCoins    sdk.Coin `json:"native_coins" yaml:"native_coins"`
	NonNativeCoins sdk.Coin `json:"non_native_coins" yaml:"non_native_coins"`
	LiquidityCoins sdk.Coin `json:"liquidity_coins" yaml:"liquidity_coins"`
}

func NewChangeInPool(nc, nnc, lc sdk.Coin) ChangeInPool {
	return ChangeInPool{
		nc,
		nnc,
		lc,
	}
}

// returns a copy of a ChangeInPool but with absolute values
func (c ChangeInPool) ToAbsolute() ChangeInPool {
	c.NativeCoins.Amount = c.amountToAbsolute(c.NativeCoins.Amount)
	c.NonNativeCoins.Amount = c.amountToAbsolute(c.NonNativeCoins.Amount)
	c.LiquidityCoins.Amount = c.amountToAbsolute(c.LiquidityCoins.Amount)
	return c
}

func (c ChangeInPool) amountToAbsolute(amt sdk.Int) sdk.Int {
	if amt.IsNegative() {
		return amt.Neg()
	}

	return amt
}

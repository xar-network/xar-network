package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ReservePool struct {
	NativeCoins    sdk.Coin `json:"native_coins" yaml:"native_coins"`
	NonNativeCoins sdk.Coin `json:"non_native_coins" yaml:"non_native_coins"`
	LiquidityCoins sdk.Coin `json:"liquidity_coins" yaml:"liquidity_coins"`
}

func NewReservePool(nativeDenom, nonNativeDenom, poolName string) ReservePool {
	nativeCoins := sdk.NewCoin(nativeDenom, sdk.ZeroInt())
	nonNativeCoins := sdk.NewCoin(nonNativeDenom, sdk.ZeroInt())

	return ReservePool{nativeCoins, nonNativeCoins, sdk.NewCoin(poolName, sdk.ZeroInt())}
}

func validateNewPool(nativeCoins, nonNativeCoins sdk.Coin, poolName string) {

}

// returns a name of an exchange
func (r ReservePool) GetName() string {
	return r.LiquidityCoins.Denom
}

// returns liquidityVouchers and error as a response
func (r *ReservePool) AddLiquidity(nativeCoins, nonNativeCoins sdk.Coin) (sdk.Coin, sdk.Error) {
	if r.NativeCoins.Amount.Equal(sdk.ZeroInt()) {
		return r.addInitialLiquidity(nativeCoins, nonNativeCoins)
	}

	return r.addLiquidity(nativeCoins, nonNativeCoins)
}

func (r *ReservePool) addLiquidity(nativeCoins, nonNativeCoins sdk.Coin) (sdk.Coin, sdk.Error) {
	err := r.validateLiquidityParams(nativeCoins, nonNativeCoins)
	if err != nil {
		return sdk.Coin{}, err
	}

	liquidityCoinAmt := r.LiquidityCoins.Amount
	nativeInPool := r.NativeCoins.Amount
	nativeIncoming := nativeCoins.Amount

	amtToMint := (liquidityCoinAmt.Mul(nativeIncoming)).Quo(nativeInPool)
	coinAmountDeposited := (liquidityCoinAmt.Mul(nativeIncoming)).Quo(nativeInPool)
	nonNativeCoins.Amount = coinAmountDeposited

	r.NativeCoins = r.NativeCoins.Add(nativeCoins)
	r.NonNativeCoins = r.NonNativeCoins.Add(nonNativeCoins)
	r.LiquidityCoins.Amount = r.LiquidityCoins.Amount.Add(amtToMint)
	return sdk.NewCoin(r.LiquidityCoins.Denom, amtToMint), nil
}

func (r *ReservePool) addInitialLiquidity(nativeCoins, nonNativeCoins sdk.Coin) (sdk.Coin, sdk.Error) {
	err := r.validateLiquidityParams(nativeCoins, nonNativeCoins)
	if err != nil {
		return sdk.Coin{}, err
	}

	coinProduct := nativeCoins.Amount.Mul(nonNativeCoins.Amount)
	mintAmtBigint := coinProduct.BigInt().Sqrt(nonNativeCoins.Amount.BigInt())

	amtToMint := sdk.NewIntFromBigInt(mintAmtBigint)
	r.NativeCoins = nativeCoins
	r.NonNativeCoins = nonNativeCoins
	r.LiquidityCoins.Amount = r.LiquidityCoins.Amount.Add(amtToMint)

	return r.LiquidityCoins, nil
}

func (r ReservePool) validateLiquidityParams(nativeCoins, nonNativeCoins sdk.Coin) sdk.Error {
	if nativeCoins.Denom != r.NativeCoins.Denom {
		return ErrIncorrectNativeDenom
	}

	if nonNativeCoins.Denom != r.NonNativeCoins.Denom {
		return ErrIncorrectNonNativeDenom
	}

	if nativeCoins.Amount.Equal(sdk.ZeroInt()) {
		return ErrIncorrectNativeAmount("native coins amount cannot be zero")
	}

	if nativeCoins.Amount.Equal(sdk.ZeroInt()) {
		return ErrIncorrectNonNativeAmount("non-native coins amount cannot be zero")
	}
	return nil
}

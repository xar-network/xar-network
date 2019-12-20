package pool

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ReservePool struct {
	nativeCoins    sdk.Coin
	nonNativeCoins sdk.Coin
	liquidityCoins sdk.Coin
}

func NewReservePool(nativeDenom, nonNativeDenom, poolName string) ReservePool {
	nativeCoins := sdk.NewCoin(nativeDenom, sdk.ZeroInt())
	nonNativeCoins := sdk.NewCoin(nonNativeDenom, sdk.ZeroInt())

	return ReservePool{nativeCoins, nonNativeCoins, sdk.NewCoin(poolName, sdk.ZeroInt())}
}

// returns a name of an exchange
func (r ReservePool) GetName() string {
	return r.liquidityCoins.Denom
}

// returns liquidityVouchers and error as a response
func (r *ReservePool) AddLiquidity(nativeCoins, nonNativeCoins sdk.Coin) (ChangeInPool, sdk.Error) {
	if r.nativeCoins.Amount.Equal(sdk.ZeroInt()) {
		return r.addInitialLiquidity(nativeCoins, nonNativeCoins)
	}

	return r.addLiquidity(nativeCoins, nonNativeCoins)
}

func (r *ReservePool) addLiquidity(nativeCoins, nonNativeCoins sdk.Coin) (ChangeInPool, sdk.Error) {
	err := r.validateLiquidityParams(nativeCoins, nonNativeCoins)
	if err != nil {
		return ChangeInPool{}, err
	}

	liquidityCoinAmt := r.liquidityCoins.Amount
	nativeInPool := r.nativeCoins.Amount
	nativeIncoming := nativeCoins.Amount

	amtToMint := (liquidityCoinAmt.Mul(nativeIncoming)).Quo(nativeInPool)
	coinAmountDeposited := (liquidityCoinAmt.Mul(nativeIncoming)).Quo(nativeInPool)
	nonNativeCoins.Amount = coinAmountDeposited

	changeInPool := NewChangeInPool(nativeCoins, nonNativeCoins, sdk.NewCoin(r.GetName(), amtToMint))

	return changeInPool, r.applyChanges(changeInPool)
}

func (r *ReservePool) addInitialLiquidity(nativeCoins, nonNativeCoins sdk.Coin) (ChangeInPool, sdk.Error) {
	err := r.validateLiquidityParams(nativeCoins, nonNativeCoins)
	if err != nil {
		return ChangeInPool{}, err
	}

	coinProduct := nativeCoins.Amount.Mul(nonNativeCoins.Amount)
	mintAmtBigint := coinProduct.BigInt().Sqrt(nonNativeCoins.Amount.BigInt())

	amtToMint := sdk.NewIntFromBigInt(mintAmtBigint)
	changeInPool := NewChangeInPool(nativeCoins, nonNativeCoins, sdk.NewCoin(r.GetName(), amtToMint))

	return changeInPool, r.applyChanges(changeInPool)
}

func (r *ReservePool) RemoveLiquidity(nativeCoins, nonNativeCoins sdk.Coin) (ChangeInPool, sdk.Error) {
	err := r.validateRemoveLiquidity(nativeCoins, nonNativeCoins)
	if err != nil {
		return ChangeInPool{}, err
	}

	liquidityCoinAmt := r.liquidityCoins.Amount
	nativeInPool := r.nativeCoins.Amount
	nativeIncoming := nativeCoins.Amount

	coinsToWithdraw := r.nonNativeCoins.Amount.Mul(nativeIncoming).Quo(nativeInPool)
	exchangeCoin := sdk.NewCoin(nonNativeCoins.Denom, coinsToWithdraw)
	amtToBurn := (liquidityCoinAmt.Mul(nativeIncoming)).Quo(nativeInPool)

	nativeCoinsNeg := r.coinToNegative(nativeCoins)
	nonNativeCoinsNeg := r.coinToNegative(exchangeCoin)
	liquidityCoinsNeg := sdk.NewCoin(r.GetName(), amtToBurn.Neg())

	changeInPool := NewChangeInPool(nativeCoinsNeg, nonNativeCoinsNeg, liquidityCoinsNeg)

	return changeInPool, r.applyChanges(changeInPool)
}

func (r ReservePool) AmountOf(denom string) sdk.Int {
	if r.nativeCoins.Denom == denom {
		return r.nativeCoins.Amount
	}

	if r.nonNativeCoins.Denom == denom {
		return r.nonNativeCoins.Amount
	}

	return sdk.ZeroInt()
}

func (r ReservePool) coinToNegative(coin sdk.Coin) sdk.Coin {
	coin.Amount = coin.Amount.Neg()
	return coin
}

func (r ReservePool) validateRemoveLiquidity(nativeCoins, nonNativeCoins sdk.Coin) sdk.Error {
	return nil
}

func (r ReservePool) validateLiquidityParams(nativeCoins, nonNativeCoins sdk.Coin) sdk.Error {
	if nativeCoins.Denom != r.nativeCoins.Denom {
		return ErrIncorrectNativeDenom
	}

	if nonNativeCoins.Denom != r.nonNativeCoins.Denom {
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

// changes a state of a reserve pool. in a future updates it is possible to add an additional logic here.
// for example:
// as far as reserve pool fields are public (for a purpose of a json marshal/unmarshal logic) we can prevent manual changes of a state by taking a hash of fields
func (r *ReservePool) applyChanges(change ChangeInPool) sdk.Error {
	r.nativeCoins.Add(change.NativeCoins)
	r.nonNativeCoins.Add(change.NonNativeCoins)
	r.liquidityCoins.Add(change.LiquidityCoins)
	return nil
}

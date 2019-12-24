package pool

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
	"strings"
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
	if r.IsEmpty() {
		return r.addInitialLiquidity(nativeCoins, nonNativeCoins)
	}

	return r.addLiquidity(nativeCoins, nonNativeCoins)
}

func (r ReservePool) IsEmpty() bool {
	return r.nativeCoins.Amount.Equal(r.nonNativeCoins.Amount) && r.nativeCoins.Amount.Equal(sdk.ZeroInt())
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

// TODO: a swap logic is currently placed inside msgSwap handler. it should be considered to put it below
func (r *ReservePool) Swap(inputCoin, outputCoin sdk.Coin) (changeInPool ChangeInPool, err sdk.Error) {
	err = r.validateSwap(inputCoin, outputCoin)
	if err != nil {
		return ChangeInPool{}, err
	}

	if r.nativeCoins.Denom == inputCoin.Denom {
		changeInPool = NewChangeInPool(inputCoin, r.coinToNegative(outputCoin), r.EmptyLiquidityCoin())
	} else {
		changeInPool = NewChangeInPool(r.coinToNegative(outputCoin), inputCoin, r.EmptyLiquidityCoin())
	}

	return changeInPool, r.applyChanges(changeInPool)
}

func (r ReservePool) EmptyLiquidityCoin() sdk.Coin {
	return sdk.NewCoin(r.GetName(), sdk.ZeroInt())
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

	if r.liquidityCoins.Denom == denom {
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

func (r ReservePool) validateSwap(inputCoin, outputCoin sdk.Coin) sdk.Error {
	if !r.containsNative(sdk.Coins{inputCoin, outputCoin}) {
		return ErrNoNativeDenomPresent
	}

	if !r.allDenomsInPool(sdk.Coins{inputCoin, outputCoin}) {
		return ErrNotAllDenomsAreInPool
	}

	return nil
}

// returns true if a coin with a native denom is present in an array
func (r ReservePool) containsNative(cns sdk.Coins) bool {
	for _, coin := range cns {
		if coin.Denom == r.nativeCoins.Denom {
			return true
		}
	}
	return false
}

// returns true if a all denoms of coins in array are present in a reserve pool
func (r ReservePool) allDenomsInPool(cns sdk.Coins) bool {
	for _, coin := range cns {
		if coin.Denom == r.nativeCoins.Denom {
			continue
		}
		if coin.Denom == r.nonNativeCoins.Denom {
			continue
		}
		if coin.Denom == r.liquidityCoins.Denom {
			continue
		}
		return false
	}
	return true
}

// changes a state of a reserve pool. in a future updates it is possible to add an additional logic here.
// for example:
// as far as reserve pool fields are public (for a purpose of a json marshal/unmarshal logic) we can prevent manual changes of a state by taking a hash of fields
func (r *ReservePool) applyChanges(change ChangeInPool) sdk.Error {
	err := r.validateChanges(change)
	if err != nil {
		return err
	}

	r.nativeCoins = r.nativeCoins.Add(change.NativeCoins)
	r.nonNativeCoins = r.nonNativeCoins.Add(change.NonNativeCoins)
	r.liquidityCoins = r.liquidityCoins.Add(change.LiquidityCoins)
	return nil
}

// place err to errors.go
func (r *ReservePool) validateChanges(change ChangeInPool) sdk.Error {
	zero := sdk.ZeroInt()
	if !change.NativeCoins.Amount.Equal(zero) {
		if r.nativeCoins.Denom != change.NativeCoins.Denom {
			return sdk.NewError(types.DefaultCodespace, 200, "")
		}
	}

	if !change.NonNativeCoins.Amount.Equal(zero) {
		if r.nonNativeCoins.Denom != change.NonNativeCoins.Denom {
			return sdk.NewError(types.DefaultCodespace, 201, "")
		}
	}

	if !change.LiquidityCoins.Amount.Equal(zero) {
		if r.liquidityCoins.Denom != change.LiquidityCoins.Denom {
			return sdk.NewError(types.DefaultCodespace, 202, "")
		}
	}
	return nil
}

func (r ReservePool) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ReservePool:
    r.nativeCoins:      %s
	r.nonNativeCoins: %s
	r.liquidityCoins: %s`,
		r.nativeCoins,
		r.nonNativeCoins,
		r.liquidityCoins,
	))
}

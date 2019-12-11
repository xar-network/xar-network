/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package keeper

import (
	"fmt"
	"log"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

func (keeper Keeper) SwapCoins(ctx sdk.Context, sender, recipient sdk.AccAddress, coinSold, coinBought sdk.Coin) sdk.Error {
	if !keeper.HasCoins(ctx, sender, coinSold) {
		cns := keeper.bk.GetCoins(ctx, sender)
		log.Println(coinSold)
		log.Println(coinBought)
		for _, v := range cns {
			log.Println(v)
		}
		return sdk.ErrInsufficientCoins(fmt.Sprintf("sender account does not have sufficient amount of %s to fulfill the swap order", coinSold.Denom))
	}

	moduleName, err := keeper.GetModuleName(coinSold.Denom, coinBought.Denom)
	if err != nil {
		return err
	}

	mAcc := keeper.ModuleAccountFromName(ctx, moduleName)
	err = keeper.SendCoins(ctx, sender, mAcc.Address, coinSold)
	if err != nil {
		return err
	}

	err = keeper.SendCoins(ctx, mAcc.Address, recipient, coinBought)
	if err != nil {
		return err
	}

	return nil
}

// GetInputAmount returns the amount of coins sold (calculated) given the output amount being bought (exact)
// The fee is included in the output coins being bought
// https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf
// TODO: continue using numerator/denominator -> open issue for eventually changing to sdk.Dec
func (keeper Keeper) GetInputAmount(ctx sdk.Context, outputAmt sdk.Int, inputDenom, outputDenom string) sdk.Int {
	moduleName, err := keeper.GetModuleName(inputDenom, outputDenom)
	if err != nil {
		panic(err)
	}
	reservePool, found := keeper.GetReservePool(ctx, moduleName)
	if !found {
		panic(fmt.Sprintf("reserve pool for %s not found", moduleName))
	}
	inputBalance := reservePool.AmountOf(inputDenom)
	outputBalance := reservePool.AmountOf(outputDenom)
	fee := keeper.GetFeeParam(ctx)

	numerator := inputBalance.Mul(outputAmt).Mul(fee.Denominator)
	denominator := (outputBalance.Sub(outputAmt)).Mul(fee.Numerator)
	return numerator.Quo(denominator).Add(sdk.OneInt())
}

// GetOutputAmount returns the amount of coins bought (calculated) given the input amount being sold (exact)
// The fee is included in the input coins being bought
// https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf
// TODO: continue using numerator/denominator -> open issue for eventually changing to sdk.Dec
func (keeper Keeper) GetOutputAmount(ctx sdk.Context, inputAmt sdk.Int, inputDenom, outputDenom string) sdk.Int {
	moduleName, err := keeper.GetModuleName(inputDenom, outputDenom)
	if err != nil {
		panic(err)
	}
	reservePool, found := keeper.GetReservePool(ctx, moduleName)
	if !found {
		panic(fmt.Sprintf("reserve pool for %s not found", moduleName))
	}

	inputBalance := reservePool.AmountOf(inputDenom)   // coin
	outputBalance := reservePool.AmountOf(outputDenom) // native
	fee := keeper.GetFeeParam(ctx)

	inputAmtWithFee := inputAmt.Mul(fee.Numerator)
	numerator := inputAmtWithFee.Mul(outputBalance)
	denominator := inputBalance.Mul(fee.Denominator).Add(inputAmtWithFee)
	return numerator.Quo(denominator)
}

// IsDoubleSwap returns true if the trade requires a double swap.
func (keeper Keeper) IsDoubleSwap(ctx sdk.Context, denom1, denom2 string) bool {
	nativeDenom := keeper.GetNativeDenom(ctx)
	return denom1 != nativeDenom && denom2 != nativeDenom
}

func (keeper Keeper) GetTradingAmount(ctx sdk.Context, input sdk.Coin, outputDenom string, buyOrder bool) sdk.Int {
	if buyOrder {
		return keeper.InputAmount(ctx, input, outputDenom)
	}

	return keeper.OutputAmount(ctx, input, outputDenom)
}

// TODO: replace GetInputAmount
func (keeper Keeper) InputAmount(ctx sdk.Context, outputCoins sdk.Coin, inputDenom string) sdk.Int {
	keeper.ValidateSwap(ctx, outputCoins.Denom, inputDenom)

	outputAmount := outputCoins.Amount
	outputReserve, inputReserve, fee := keeper.getSwapBalances(ctx, outputCoins.Denom, inputDenom)

	numerator := outputAmount.Mul(inputReserve).Mul(fee.Denominator)
	denominator := (outputReserve.Sub(outputAmount)).Mul(fee.Numerator)
	return numerator.Quo(denominator).Add(sdk.OneInt())
}

func (keeper Keeper) DoubleSwapOutputAmount(ctx sdk.Context, coinA, coinB sdk.Coin) (nativeAmount, nonNativeAmount sdk.Int) {
	if coinA.Denom == coinB.Denom {
		panic("denoms are equal")
	}

	fee := keeper.GetFeeParam(ctx)
	nativeDenom := keeper.GetNativeDenom(ctx)
	moduleNameA := keeper.MustGetModuleName(coinA.Denom, nativeDenom)
	moduleNameB := keeper.MustGetModuleName(coinB.Denom, nativeDenom)

	reservePoolA, found := keeper.GetReservePool(ctx, moduleNameA)
	if !found {
		panic("reserve pool not found")
	}

	reservePoolB, found := keeper.GetReservePool(ctx, moduleNameB)
	if !found {
		panic("reserve pool not found")
	}

	// nonNativeA to native conversion
	inputAmountA := coinA.Amount
	inputReserveA := reservePoolA.AmountOf(coinA.Denom)
	outputReserveA := reservePoolA.AmountOf(nativeDenom)

	inputAWithoutFee := inputAmountA.Mul(fee.Numerator)
	numeratorA := inputAWithoutFee.Mul(outputReserveA)
	denominatorA := (inputReserveA.Mul(fee.Denominator)).Add(inputAWithoutFee)
	outputAmountA := numeratorA.Quo(denominatorA)

	// native to nonNativeB conversion
	inputAmountB := outputAmountA
	inputReserveB := reservePoolB.AmountOf(nativeDenom)
	outputReserveB := reservePoolB.AmountOf(coinB.Denom)

	outputBWithoutFee := inputAmountB.Mul(fee.Numerator)
	numeratorB := outputBWithoutFee.Mul(outputReserveB)
	denominatorB := (inputReserveB.Mul(fee.Denominator)).Add(outputBWithoutFee)
	outputAmountB := numeratorB.Quo(denominatorB)

	return outputAmountA, outputAmountB
}

func (keeper Keeper) DoubleSwapInputAmount(ctx sdk.Context, outputCoinsA, outputCoinsB sdk.Coin) (inputAmountA, inputAmountB sdk.Int) {
	if outputCoinsA.Denom == outputCoinsB.Denom {
		panic("denoms are equal")
	}

	fee := keeper.GetFeeParam(ctx)
	nativeDenom := keeper.GetNativeDenom(ctx)
	moduleNameA := keeper.MustGetModuleName(outputCoinsA.Denom, nativeDenom)
	moduleNameB := keeper.MustGetModuleName(outputCoinsB.Denom, nativeDenom)

	reservePoolA, found := keeper.GetReservePool(ctx, moduleNameA)
	if !found {
		panic("reserve pool not found")
	}

	reservePoolB, found := keeper.GetReservePool(ctx, moduleNameB)
	if !found {
		panic("reserve pool not found")
	}

	outputAmountB := outputCoinsB.Amount
	inputReserveB := reservePoolB.AmountOf(nativeDenom)
	outputReserveB := reservePoolB.AmountOf(outputCoinsB.Denom)

	numeratorB := outputAmountB.Mul(inputReserveB).Mul(fee.Denominator)
	denominatorB := (outputReserveB.Sub(outputAmountB)).Mul(fee.Numerator)
	inputAmountB = numeratorB.Quo(denominatorB).Add(sdk.OneInt())

	outputAmountA := inputAmountB
	inputReserveA := reservePoolA.AmountOf(outputCoinsA.Denom)
	outputReserveA := reservePoolA.AmountOf(nativeDenom)

	numeratorA := outputAmountA.Mul(inputReserveA).Mul(fee.Denominator)
	denominatorA := (outputReserveA.Sub(outputAmountA)).Mul(fee.Numerator)
	inputAmountA = numeratorA.Quo(denominatorA).Add(sdk.OneInt())

	return inputAmountA, inputAmountB
}

// TODO: replace GetInputAmount
func (keeper Keeper) OutputAmount(ctx sdk.Context, inputCoins sdk.Coin, outputDenom string) sdk.Int {
	keeper.ValidateSwap(ctx, inputCoins.Denom, outputDenom)

	inputAmount := inputCoins.Amount
	inputReserve, outputReserve, fee := keeper.getSwapBalances(ctx, inputCoins.Denom, outputDenom)

	inputMinusFee := inputAmount.Mul(fee.Numerator)
	numerator := inputMinusFee.Mul(outputReserve)
	denominator := (inputReserve.Mul(fee.Denominator)).Add(inputMinusFee)
	return numerator.Quo(denominator)
}

func (keeper Keeper) ValidateSwap(ctx sdk.Context, denom1, denom2 string) {
	if keeper.IsDoubleSwap(ctx, denom1, denom2) {
		panic("cannot commit double swap straightly. use native denomination as middleware")
	}

	if denom1 == denom2 {
		panic("cannot swap equal denoms")
	}
}

// returns balances from reservePool and fees from genesis
// panics if some of them is zero or uninitialized.
func (keeper Keeper) getSwapBalances(ctx sdk.Context, denom1, denom2 string) (balance1, balance2 sdk.Int, fee types.FeeParam) {
	moduleName, err := keeper.GetModuleName(denom1, denom2)
	if err != nil {
		panic(err)
	}

	reservePool, found := keeper.GetReservePool(ctx, moduleName)
	if !found {
		panic(fmt.Sprintf("reserve pool for %s not found", moduleName))
	}

	balance1 = reservePool.AmountOf(denom1)
	if balance1.IsZero() {
		panic("non native coin has zero value")
	}

	balance2 = reservePool.AmountOf(denom2)
	if balance2.IsZero() {
		panic("native coin has zero value")
	}

	fee = keeper.GetFeeParam(ctx)
	if fee.Denominator.IsZero() {
		panic("fee denominator is zero")
	}

	if fee.Numerator.IsZero() {
		panic("fee numerator is zero")
	}

	return
}

func (keeper Keeper) DenomIsNative(ctx sdk.Context, denom string) bool {
	return keeper.GetNativeDenom(ctx) == denom
}

// GetModuleName returns the ModuleAccount name for the provided denominations.
// The module name is in the format of 'swap:denom:denom' where the denominations
// are sorted alphabetically.
func (keeper Keeper) GetModuleName(denom1, denom2 string) (string, sdk.Error) {
	// replaced ':' with digits to pass a regex check inside 'AddCoins'
	// though punctuation is not suitable, it is possible to use digits as a trailing symbols
	switch strings.Compare(denom1, denom2) {
	case -1:
		return "swap1" + denom1 + "2" + denom2, nil
	case 1:
		return "swap1" + denom2 + "2" + denom1, nil
	default:
		return "", types.ErrEqualDenom(types.DefaultCodespace, "denomnations for forming module name are equal")
	}
}

func (keeper Keeper) MustGetModuleName(denom1, denom2 string) string {
	moduleName, err := keeper.GetModuleName(denom1, denom2)
	if err != nil {
		panic(err)
	}

	return moduleName
}

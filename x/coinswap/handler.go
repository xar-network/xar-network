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

package coinswap

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/coinswap/internal/types"
)

// NewHandler returns a handler for "coinswap" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgSwapOrder:
			return HandleMsgSwapOrder(ctx, msg, k)
		case MsgAddLiquidity:
			return HandleMsgAddLiquidity(ctx, msg, k)
		case MsgRemoveLiquidity:
			return HandleMsgRemoveLiquidity(ctx, msg, k)
		case MsgTransactionOrder:
			return HandleMsgTransactionOrder(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized coinswap message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// HandleMsgSwapOrder.
func HandleMsgSwapOrder(ctx sdk.Context, msg MsgSwapOrder, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(msg.Deadline) {
		return ErrInvalidDeadline(DefaultCodespace, "deadline has passed for MsgSwapOrder").Result()
	}

	if msg.IsDoubleSwap(k.GetNativeDenom(ctx)) {
		return DoubleSwap(ctx, k, msg)
	}
	return ValidateAndSwap(&k, ctx, &msg)
}

// just a wrapper around "swap"
// the only difference is that msgTransactionOrder has an additional check.
// it is necessary for sender and recipient to be different addresses.
// it just reproduces logic of the original "transactionOrder" message from coinswap
func HandleMsgTransactionOrder(ctx sdk.Context, msg MsgTransactionOrder, k Keeper) sdk.Result {
	m := MsgSwapOrder{
		msg.Input,
		msg.Output,
		msg.Deadline,
		msg.Sender,
		msg.Recipient,
		msg.IsBuyOrder,
	}
	return HandleMsgSwapOrder(ctx, m, k)
}

func ValidateAndSwap(keeper *Keeper, ctx sdk.Context, msg *MsgSwapOrder) sdk.Result {
	var calculatedAmount sdk.Int
	if msg.IsBuyOrder {
		calculatedAmount = keeper.GetInputAmount(ctx, msg.Output.Amount, msg.Input.Denom, msg.Output.Denom)
	} else {
		calculatedAmount = keeper.GetOutputAmount(ctx, msg.Input.Amount, msg.Input.Denom, msg.Output.Denom)
	}

	err := validateSwapMsg(msg, calculatedAmount)
	if err != nil {
		return err.Result()
	}

	err = makeSwap(keeper, ctx, msg, calculatedAmount)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// function is too long
// TODO: decompose function
func DoubleSwap(ctx sdk.Context, keeper Keeper, msg types.MsgSwapOrder) sdk.Result {
	nativeDenom := keeper.GetNativeDenom(ctx)
	var nativeMediatorAmt sdk.Int
	var outputAmt sdk.Int
	var inputAmt sdk.Int

	if msg.IsBuyOrder {
		nAmt, oAmt := keeper.DoubleSwapOutputAmount(ctx, msg.Input, msg.Output)
		if msg.Input.Amount.LT(oAmt) {
			return sdk.ErrInsufficientCoins("output amount ").Result()
		}

		nativeMediatorAmt = nAmt
		inputAmt = oAmt
		outputAmt = msg.Output.Amount
	} else {
		iAmt, nAmt := keeper.DoubleSwapInputAmount(ctx, msg.Input, msg.Output)
		if msg.Input.Amount.LT(iAmt) {
			return sdk.ErrInsufficientCoins("output amount ").Result()
		}
		nativeMediatorAmt = nAmt
		inputAmt = msg.Input.Amount
		outputAmt = iAmt
	}

	inputCoin := sdk.NewCoin(msg.Input.Denom, inputAmt)
	outputCoin := sdk.NewCoin(msg.Output.Denom, outputAmt)
	nativeMideator := sdk.NewCoin(nativeDenom, nativeMediatorAmt)

	rpInput, found := keeper.GetReservePool(ctx, inputCoin.Denom)
	if !found {
		return types.ErrReservePoolNotFound(DefaultCodespace, inputCoin.Denom).Result()
	}

	rpOutput, found := keeper.GetReservePool(ctx, outputCoin.Denom)
	if !found {
		return types.ErrReservePoolNotFound(DefaultCodespace, inputCoin.Denom).Result()
	}

	_, err := rpInput.Swap(inputCoin, nativeMideator)
	if err != nil {
		return err.Result()
	}

	_, err = rpOutput.Swap(nativeMideator, outputCoin)
	if err != nil {
		return err.Result()
	}

	err = keeper.HandleCoinSwap(ctx, msg.Sender, msg.Recipient, inputCoin, outputCoin)
	if err != nil {
		return err.Result()
	}

	keeper.SetReservePool(ctx, rpInput)
	keeper.SetReservePool(ctx, rpOutput)

	return sdk.Result{}
}

// TODO: replace
func DoubleSwapCoins(keeper *Keeper, ctx sdk.Context, msg MsgSwapOrder) sdk.Result {
	nativeDenom := keeper.GetNativeDenom(ctx)
	calculatedAmount := keeper.GetOutputAmount(ctx, msg.Input.Amount, msg.Input.Denom, nativeDenom)
	nativeCoinMideator := sdk.NewCoin(nativeDenom, calculatedAmount)
	firstSwapMsg := msg
	firstSwapMsg.Output = nativeCoinMideator
	result := ValidateAndSwap(keeper, ctx, &firstSwapMsg)
	if !result.IsOK() {
		return result
	}

	secondSwapMsg := msg
	secondSwapMsg.Input = nativeCoinMideator
	return ValidateAndSwap(keeper, ctx, &secondSwapMsg)
}

func makeSwap(keeper *Keeper, ctx sdk.Context, msg *MsgSwapOrder, calculatedAmount sdk.Int) sdk.Error {
	if msg.IsBuyOrder {
		return keeper.SwapCoins(ctx, msg.Sender, msg.Recipient, sdk.NewCoin(msg.Input.Denom, calculatedAmount), msg.Output)
	}

	return keeper.SwapCoins(ctx, msg.Sender, msg.Recipient, msg.Input, sdk.NewCoin(msg.Output.Denom, calculatedAmount))
}

func validateSwapMsg(msg *MsgSwapOrder, calculatedAmount sdk.Int) sdk.Error {
	if msg.IsBuyOrder {
		if !calculatedAmount.LTE(msg.Input.Amount) {
			return ErrConstraintNotMet(DefaultCodespace, fmt.Sprintf("maximum amount (%d) to be sold was exceeded (%d)", msg.Input.Amount, calculatedAmount))
		}
		return nil
	}

	if !calculatedAmount.GTE(msg.Output.Amount) {
		return ErrConstraintNotMet(DefaultCodespace, fmt.Sprintf("minimum amount (%d) to be sold was not met (%d)", msg.Output.Amount, calculatedAmount))
	}

	return nil
}

// HandleMsgAddLiquidity. If the reserve pool does not exist, it will be
// created. The first liquidity provider sets the exchange rate.
// TODO create the initial setting liquidity, additional liquidity does not have to be in the same ratio
func HandleMsgAddLiquidity(ctx sdk.Context, msg MsgAddLiquidity, keeper Keeper) sdk.Result {
	nativeCoins := sdk.NewCoin(keeper.GetNativeDenom(ctx), msg.DepositAmount)

	rp := keeper.CreateOrGetReservePool(ctx, msg.Deposit.Denom)

	changesInPool, err := rp.AddLiquidity(nativeCoins, msg.Deposit)
	if err != nil {
		return err.Result()
	}

	coins := sdk.NewCoins(changesInPool.NonNativeCoins, changesInPool.NativeCoins)
	liquidityVouchers := changesInPool.LiquidityCoins

	err = keeper.AddLiquidityTransfer(ctx, msg.Sender, coins, liquidityVouchers)
	if err != nil {
		return err.Result()
	}

	keeper.SetReservePool(ctx, rp)

	return sdk.Result{}
}

// HandleMsgRemoveLiquidity handler for MsgRemoveLiquidity
func HandleMsgRemoveLiquidity(ctx sdk.Context, msg MsgRemoveLiquidity, keeper Keeper) sdk.Result {
	// check that deadline has not passed
	nativeCoins := sdk.NewCoin(keeper.GetNativeDenom(ctx), msg.WithdrawAmount)

	rp := keeper.CreateOrGetReservePool(ctx, msg.Withdraw.Denom)

	rawChangesInPool, err := rp.RemoveLiquidity(nativeCoins, msg.Withdraw)
	if err != nil {
		return err.Result()
	}
	changesInPool := rawChangesInPool.ToAbsolute()

	coins := sdk.NewCoins(changesInPool.NonNativeCoins, changesInPool.NativeCoins)
	liquidityVouchers := changesInPool.LiquidityCoins

	err = keeper.RemoveLiquidityTransfer(ctx, msg.Sender, coins, liquidityVouchers)
	if err != nil {
		return err.Result()
	}

	keeper.SetReservePool(ctx, rp)
	return sdk.Result{}
}

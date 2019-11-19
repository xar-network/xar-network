package uniswap

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	var calculatedAmount, nativeAmount sdk.Int
	doubleSwap := k.IsDoubleSwap(ctx, msg.Input.Denom, msg.Output.Denom)
	nativeDenom := k.GetNativeDenom(ctx)

	if msg.IsBuyOrder {
		if doubleSwap {
			return ErrNotSupported(DefaultCodespace, "double swaps are currently not supported").Result()
			/*nativeAmount = k.GetInputAmount(ctx, msg.Output.Amount, nativeDenom, msg.Output.Denom)
			calculatedAmount = k.GetInputAmount(ctx, nativeAmount, msg.Input.Denom, nativeDenom)
			nativeCoin := sdk.NewCoin(nativeDenom, nativeAmount)
			k.SwapCoins(ctx, msg.Sender, sdk.NewCoin(msg.Input.Denom, calculatedAmount), nativeCoin)
			k.SwapCoins(ctx, msg.Sender, nativeCoin, msg.Output)*/
		} else {
			calculatedAmount = k.GetInputAmount(ctx, msg.Output.Amount, msg.Input.Denom, msg.Output.Denom)
			// assert that the calculated amount is less than or equal to the
			// maximum amount the buyer is willing to pay.
			if calculatedAmount.LTE(msg.Input.Amount) {
				k.SwapCoins(ctx, msg.Sender, sdk.NewCoin(msg.Input.Denom, calculatedAmount), msg.Output)
			} else {
				return ErrConstraintNotMet(DefaultCodespace, fmt.Sprintf("maximum amount (%d) to be sold was exceeded (%d)", msg.Input.Amount, calculatedAmount)).Result()
			}
		}
	} else {
		if doubleSwap {
			return ErrNotSupported(DefaultCodespace, "double swaps are currently not supported").Result()
			/*nativeAmount = k.GetOutputAmount(ctx, msg.Input.Amount, msg.Input.Denom, nativeDenom)
			calculatedAmount = k.GetOutputAmount(ctx, nativeAmount, nativeDenom, msg.Output.Denom)
			nativeCoin := sdk.NewCoin(nativeDenom, nativeAmount)
			k.SwapCoins(ctx, msg.Sender, msg.Input, nativeCoin)
			k.SwapCoins(ctx, msg.Sender, nativeCoin, sdk.NewCoin(msg.Output.Denom, calculatedAmount))*/
		} else {
			calculatedAmount = k.GetOutputAmount(ctx, msg.Input.Amount, msg.Input.Denom, msg.Output.Denom)

			// assert that the calculated amount is greater than or equal to the
			// minimum amount the sender is willing to buy.
			if !calculatedAmount.GTE(msg.Output.Amount) {
				return ErrConstraintNotMet(DefaultCodespace, fmt.Sprintf("minimum amount (%d) to be sold was not met (%d)", msg.Output.Amount, calculatedAmount)).Result()
			} else {
				k.SwapCoins(ctx, msg.Sender, msg.Input, sdk.NewCoin(msg.Output.Denom, calculatedAmount))
			}
		}
	}
	return sdk.Result{}
}

// HandleMsgAddLiquidity. If the reserve pool does not exist, it will be
// created. The first liquidity provider sets the exchange rate.
// TODO create the initial setting liquidity, additional liquidity does not have to be in the same ratio
func HandleMsgAddLiquidity(ctx sdk.Context, msg MsgAddLiquidity, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(msg.Deadline) {
		return ErrInvalidDeadline(DefaultCodespace, "deadline has passed for MsgAddLiquidity").Result()
	}

	nativeDenom := k.GetNativeDenom(ctx)
	//swap:nativeDenom:msg.Deposit.Denom
	moduleName, err := k.GetModuleName(nativeDenom, msg.Deposit.Denom)
	if err != nil {
		return err.Result()
	}
	// create reserve pool if it does not exist
	reservePool, found := k.GetReservePool(ctx, moduleName)
	if !found {
		k.CreateReservePool(ctx, moduleName)
		reservePool, found = k.GetReservePool(ctx, moduleName)
		if !found {
			// return pool not found error
		}
	}

	// New pool would have 0 nativeBalance, coinBalance, and liquidityCoinBalance
	nativeBalance := reservePool.AmountOf(nativeDenom)       //AmountOf FTM
	coinBalance := reservePool.AmountOf(msg.Deposit.Denom)   //AmountOf DAI
	liquidityCoinBalance := reservePool.AmountOf(moduleName) //AmountOf swap:DAI:FTM

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	// TODO: verify

	amtToMint := (liquidityCoinBalance.Mul(msg.DepositAmount)).Quo(nativeBalance)           //0 on new pool
	coinAmountDeposited := (liquidityCoinBalance.Mul(msg.DepositAmount)).Quo(nativeBalance) //calculate amount of DAI based on pool, this would mean ratio is ever depleting
	// Should be able to add any amount and not calculated amounts since withdrawals can skew these numbers
	nativeCoinDeposited := sdk.NewCoin(nativeDenom, msg.DepositAmount)
	coinDeposited := sdk.NewCoin(msg.Deposit.Denom, coinAmountDeposited)

	if !k.HasCoins(ctx, msg.Sender, nativeCoinDeposited, coinDeposited) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to add liquidity").Result()
	}

	// transfer deposited liquidity into uniswap ModuleAccount
	err = k.SendCoins(ctx, msg.Sender, moduleName, nativeCoinDeposited, coinDeposited)
	if err != nil {
		return err.Result()
	}

	// mint liquidity vouchers for sender
	k.MintCoins(ctx, moduleName, amtToMint)
	k.RecieveCoins(ctx, msg.Sender, moduleName, sdk.NewCoin(moduleName, amtToMint))

	return sdk.Result{}
}

// HandleMsgRemoveLiquidity handler for MsgRemoveLiquidity
func HandleMsgRemoveLiquidity(ctx sdk.Context, msg MsgRemoveLiquidity, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(msg.Deadline) {
		return ErrInvalidDeadline(DefaultCodespace, "deadline has passed for MsgRemoveLiquidity")
	}

	nativeDenom := k.GetNativeDenom(ctx)
	moduleName := k.GetModuleName(ctx, nativeDenom, msg.Deposit.Denom)

	// check if reserve pool exists
	reservePool, found := k.GetReservePool(ctx, msg.Withdraw.Denom)
	if !found {
		panic(fmt.Sprintf("error retrieving reserve pool for ModuleAccoint name: %s", moduleName))
	}

	nativeBalance := reservePool.AmountOf(nativeDenom)
	coinBalance := reservePool.AmountOf(msg.Withdraw.Denom)
	liquidityCoinBalance := reservePool.AmountOf(moduleName)

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	// TODO: verify, add amt burned
	nativeWithdrawn := msg.WithdrawAmount.Mul(nativeBalance).Quo(liquidityCoinBalance)
	coinWithdrawn := msg.WithdrawAmount.Mul(coinBalance).Quo(liquidityCoinBalance)
	nativeCoin := sdk.NewCoin(nativeDenom, nativeWithdrawn)
	exchangeCoin = sdk.NewCoin(msg.Withdraw.Denom, coinWithdrawn)

	if !k.HasCoins(ctx, msg.Sender, sdk.NewCoin(moduleName, amtBurned)) {
		return sdk.ErrInsufficientCoins("sender does not have sufficient funds to remove liquidity").Result()
	}

	// burn liquidity vouchers
	k.SendCoins(ctx, msg.Sender, moduleName, sdk.NewCoin(moduleName, amtBurned))
	k.BurnCoins(ctx, moduleName, amtBurned)

	// transfer withdrawn liquidity from coinswaps ModuleAccount to sender's account
	k.RecieveCoins(ctx, msg.Sender, moduleName, nativeCoin, coinDeposited)

	return sdk.Result{}
}

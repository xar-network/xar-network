package uniswap

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
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
// it just reproduces logic of the original "transactionOrder" message from uniswap
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

	moduleNameA := keeper.MustGetModuleName(nativeDenom, inputCoin.Denom)
	mAccA := keeper.ModuleAccountFromName(ctx, moduleNameA)

	moduleNameB := keeper.MustGetModuleName(nativeDenom, outputCoin.Denom)
	mAccB := keeper.ModuleAccountFromName(ctx, moduleNameB)

	err := keeper.SendCoins(ctx, msg.Sender, mAccA.Address, inputCoin)
	if err != nil {
		return err.Result()
	}

	err = keeper.SendCoins(ctx, mAccA.Address, mAccB.Address, nativeMideator)
	if err != nil {
		return err.Result()
	}

	err = keeper.SendCoins(ctx, mAccB.Address, msg.Recipient, outputCoin)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// incorrect. will be removed in the upcoming commit
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
	if ctx.BlockHeader().Time.After(msg.Deadline) {
		return ErrInvalidDeadline(DefaultCodespace, types.LiquidityAddDeadLineHasPassed).Result()
	}

	nativeDenom := keeper.GetNativeDenom(ctx) // swap:nativeDenom:msg.Deposit.Denom
	moduleName, err := keeper.GetModuleName(nativeDenom, msg.Deposit.Denom)
	if err != nil {
		return err.Result()
	}

	// create reserve pool if it does not exist
	reservePool, found := keeper.GetReservePool(ctx, moduleName)
	if !found {
		err := newReservePool(ctx, moduleName, keeper)
		if err != nil {
			return err.Result()
		}

		return keeper.AddInitialLiquidity(ctx, &msg)
		//addInitialLiquidity(ctx, keeper, &msg)
	}

	return keeper.AddLiquidity(ctx, &msg, reservePool)
	//addLiquidity(ctx, keeper, &msg, reservePool)
}

// creates new reserve pool and verifies it was created successfully
func newReservePool(ctx sdk.Context, moduleName string, keeper Keeper) sdk.Error {
	keeper.CreateReservePool(ctx, moduleName)

	if _, found := keeper.GetReservePool(ctx, moduleName); !found {
		return ErrCannotCreateReservePool(DefaultCodespace)
	}

	return nil
}

func addInitialLiquidity(ctx sdk.Context, k Keeper, msg *MsgAddLiquidity) sdk.Result {
	nativeDenom, _, moduleName := k.MustGetAllDenoms(ctx, msg)

	coinAmount := msg.Deposit.Amount.BigInt()
	nativeAmount := msg.DepositAmount.BigInt()
	mintAmtBigint := (coinAmount.Mul(coinAmount, nativeAmount)).Sqrt(coinAmount)
	nativeCoinDeposited := sdk.NewCoin(nativeDenom, msg.DepositAmount)
	amtToMint := sdk.NewIntFromBigInt(mintAmtBigint)

	if !k.HasCoins(ctx, msg.Sender, nativeCoinDeposited, msg.Deposit) {
		return sdk.ErrInsufficientCoins(types.InsufficientCoins).Result()
	}

	mAcc := k.ModuleAccountFromName(ctx, moduleName)
	err := transferLiquidityCoins(ctx, msg, k, msg.Deposit, nativeCoinDeposited, amtToMint, mAcc)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func addLiquidity(ctx sdk.Context, k Keeper, msg *MsgAddLiquidity, reservePool sdk.Coins) sdk.Result {
	nativeDenom, _, moduleName := k.MustGetAllDenoms(ctx, msg)
	nativeBalance := reservePool.AmountOf(nativeDenom)
	liquidityCoinBalance := reservePool.AmountOf(moduleName)
	if liquidityCoinBalance.LTE(sdk.NewInt(0)) {
		return types.ErrInsufficientLiquidityAmount(DefaultCodespace).Result()
	}

	amtToMint := (liquidityCoinBalance.Mul(msg.DepositAmount)).Quo(nativeBalance)
	coinAmountDeposited := (liquidityCoinBalance.Mul(msg.DepositAmount)).Quo(nativeBalance)
	nativeCoinDeposited := sdk.NewCoin(nativeDenom, msg.DepositAmount)
	coinDeposited := sdk.NewCoin(msg.Deposit.Denom, coinAmountDeposited)

	if !k.HasCoins(ctx, msg.Sender, nativeCoinDeposited, coinDeposited) {
		return sdk.ErrInsufficientCoins(types.InsufficientCoins).Result()
	}

	mAcc := k.ModuleAccountFromName(ctx, moduleName)
	err := transferLiquidityCoins(ctx, msg, k, nativeCoinDeposited, coinDeposited, amtToMint, mAcc)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func transferLiquidityCoins(ctx sdk.Context, msg *MsgAddLiquidity, k Keeper, nativeCoin, coin sdk.Coin, amtToMint sdk.Int, moduleAcc *supply.ModuleAccount) sdk.Error {
	err := k.SendCoins(ctx, msg.Sender, moduleAcc.Address, nativeCoin, coin)
	if err != nil {
		return err
	}

	// mint liquidity vouchers for sender
	mintCoins, err := k.MintCoins(ctx, moduleAcc.Address, amtToMint)
	if err != nil {
		return err
	}

	return k.RecieveCoins(ctx, msg.Sender, mintCoins...)
}

// HandleMsgRemoveLiquidity handler for MsgRemoveLiquidity
func HandleMsgRemoveLiquidity(ctx sdk.Context, msg MsgRemoveLiquidity, k Keeper) sdk.Result {
	// check that deadline has not passed
	if ctx.BlockHeader().Time.After(msg.Deadline) {
		return ErrInvalidDeadline(DefaultCodespace, types.LiquidityRemoveDeadLineHasPassed).Result()
	}

	nativeDenom := k.GetNativeDenom(ctx)
	moduleName, err := k.GetModuleName(nativeDenom, msg.Withdraw.Denom)
	if err != nil {
		return err.Result()
	}

	// check if reserve pool exists
	reservePool, found := k.GetReservePool(ctx, moduleName)
	if !found {
		return types.ErrReservePoolNotFound(DefaultCodespace, moduleName).Result()
	}

	nativeBalance := reservePool.AmountOf(nativeDenom)       // n coin
	coinBalance := reservePool.AmountOf(msg.Withdraw.Denom)  // erc20
	liquidityCoinBalance := reservePool.AmountOf(moduleName) //total amount of liquidity tokens

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	amtCoinWithdrawn := coinBalance.Mul(msg.WithdrawAmount).Quo(nativeBalance) // msg.WithdrawAmount.Mul(coinBalance).Quo(liquidityCoinBalance)
	nativeCoin := sdk.NewCoin(nativeDenom, msg.WithdrawAmount)
	exchangeCoin := sdk.NewCoin(msg.Withdraw.Denom, amtCoinWithdrawn)
	amtToBurn := (liquidityCoinBalance.Mul(msg.WithdrawAmount)).Quo(nativeBalance)

	if !k.HasCoins(ctx, msg.Sender, sdk.NewCoin(moduleName, amtToBurn)) {
		return sdk.ErrInsufficientCoins(types.InsufficientCoins).Result()
	}

	// burn liquidity vouchers
	mAcc := k.ModuleAccountFromName(ctx, moduleName)
	err = k.SendCoins(ctx, msg.Sender, mAcc.Address, sdk.NewCoin(moduleName, amtToBurn))
	if err != nil {
		return err.Result()
	}

	err = k.BurnCoins(ctx, moduleName, amtToBurn)
	if err != nil {
		return err.Result()
	}

	err = k.SendCoins(ctx, mAcc.Address, msg.Sender, nativeCoin, exchangeCoin)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

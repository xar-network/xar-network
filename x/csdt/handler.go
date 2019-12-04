/**

Baseline from Kava Cosmos Module

**/

package csdt

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// Handle all csdt messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgCreateOrModifyCSDT:
			return handleMsgCreateOrModifyCSDT(ctx, keeper, msg)
		case types.MsgDepositCollateral:
			return handleMsgDepositCollateral(ctx, keeper, msg)
		case types.MsgWithdrawCollateral:
			return handleMsgWithdrawCollateral(ctx, keeper, msg)
		case types.MsgSettleDebt:
			return handleMsgSettleDebt(ctx, keeper, msg)
		case types.MsgWithdrawDebt:
			return handleMsgWithdrawDebt(ctx, keeper, msg)
		case types.MsgSetCollateralParam:
			return handleMsgSetCollateralParam(ctx, keeper, msg)
		case types.MsgAddCollateralParam:
			return handleMsgAddCollateralParam(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized csdt msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateOrModifyCSDT(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgCreateOrModifyCSDT) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = keeper.ModifyCSDT(ctx, msg.Sender, msg.CollateralDenom, msg.CollateralChange, msg.DebtChange)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgDepositCollateral(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDepositCollateral) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = keeper.ModifyCSDT(ctx, msg.Sender, msg.CollateralDenom, msg.CollateralChange, sdk.NewInt(0))
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdrawCollateral(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgWithdrawCollateral) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = keeper.ModifyCSDT(ctx, msg.Sender, msg.CollateralDenom, msg.CollateralChange.Neg(), sdk.NewInt(0))
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSettleDebt(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSettleDebt) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = keeper.ModifyCSDT(ctx, msg.Sender, msg.CollateralDenom, sdk.NewInt(0), msg.DebtChange.Neg())
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdrawDebt(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgWithdrawDebt) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = keeper.ModifyCSDT(ctx, msg.Sender, msg.CollateralDenom, sdk.NewInt(0), msg.DebtChange)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSetCollateralParam(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgSetCollateralParam) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	params := types.CollateralParam{
		Denom:            msg.CollateralDenom,
		LiquidationRatio: msg.LiquidationRatio,
		DebtLimit:        msg.DebtLimit,
	}

	err = keeper.SetCollateralParam(ctx, msg.Nominee.String(), params)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgAddCollateralParam(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgAddCollateralParam) sdk.Result {

	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	params := types.CollateralParam{
		Denom:            msg.CollateralDenom,
		LiquidationRatio: msg.LiquidationRatio,
		DebtLimit:        msg.DebtLimit,
	}

	err = keeper.AddCollateralParam(ctx, msg.Nominee.String(), params)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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

	collateral := sdk.NewCoin(msg.CollateralDenom, msg.CollateralChange)
	debt := sdk.NewCoin(msg.DebtDenom, msg.DebtChange)
	err = keeper.ModifyCSDT(ctx, msg.Sender, collateral, debt)
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

	collateral := sdk.NewCoin(msg.CollateralDenom, msg.CollateralChange)
	debt := sdk.NewCoin("", sdk.NewInt(0))
	err = keeper.ModifyCSDT(ctx, msg.Sender, collateral, debt)
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

	collateral := sdk.NewCoin(msg.CollateralDenom, msg.CollateralChange.Neg())
	debt := sdk.NewCoin("", sdk.NewInt(0))
	err = keeper.ModifyCSDT(ctx, msg.Sender, collateral, debt)
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

	collateral := sdk.NewCoin(msg.CollateralDenom, sdk.NewInt(0))
	debt := sdk.NewCoin(msg.DebtDenom, msg.DebtChange.Neg())
	err = keeper.ModifyCSDT(ctx, msg.Sender, collateral, debt)
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

	collateral := sdk.NewCoin(msg.CollateralDenom, sdk.NewInt(0))
	debt := sdk.NewCoin(msg.DebtDenom, msg.DebtChange)
	err = keeper.ModifyCSDT(ctx, msg.Sender, collateral, debt)
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

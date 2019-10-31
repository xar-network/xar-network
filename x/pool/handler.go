package pool

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all pool messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDepositFund:
			return handleMsgDepositFund(ctx, keeper, msg)
		case MsgWithdrawFund:
			return handleMsgWithdrawFund(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized pool msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handles the message that allows a user to deposit funds into the pool
func handleMsgDepositFund(ctx sdk.Context, keeper Keeper, msg MsgDepositFund) sdk.Result {

	err := keeper.DepositFundFromAddress(ctx, msg.Sender, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// handles the message that allows a user to withdraw funds from the pool
func handleMsgWithdrawFund(ctx sdk.Context, keeper Keeper, msg MsgWithdrawFund) sdk.Result {

	err := keeper.WithdrawFundToAddress(ctx, msg.Amount, msg.Sender)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// EndBlocker distributes the rewards
func EndBlocker(ctx sdk.Context, k Keeper) sdk.Tags {

	// Running in the end blocker ensures that rewards will update at most once per block
	err := k.DistributeReward(ctx)
	if err != nil {
		panic(err)
	}

	return sdk.Tags{}
}

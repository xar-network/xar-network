/**

Baseline from Hashgard Cosmos Module

**/

package issue

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/issue/internal/keeper"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

// NewHandler creates an sdk.Handler for all the issue type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgIssue:
			return handleMsgIssue(ctx, k, msg)
		case types.MsgIssueTransferOwnership:
			return handleMsgIssueTransferOwnership(ctx, k, msg)
		case types.MsgIssueDescription:
			return handleMsgIssueDescription(ctx, k, msg)
		case types.MsgIssueMint:
			return handleMsgIssueMint(ctx, k, msg)
		case types.MsgIssueBurnOwner:
			return handleMsgIssueBurnOwner(ctx, k, msg)
		case types.MsgIssueBurnHolder:
			return handleMsgIssueBurnHolder(ctx, k, msg)
		case types.MsgIssueBurnFrom:
			return handleMsgIssueBurnFrom(ctx, k, msg)
		case types.MsgIssueDisableFeature:
			return handleMsgIssueDisableFeature(ctx, k, msg)
		case types.MsgIssueApprove:
			return handleMsgIssueApprove(ctx, k, msg)
		case types.MsgIssueSendFrom:
			return handleMsgIssueSendFrom(ctx, k, msg)
		case types.MsgIssueIncreaseApproval:
			return handleMsgIssueIncreaseApproval(ctx, k, msg)
		case types.MsgIssueDecreaseApproval:
			return handleMsgIssueDecreaseApproval(ctx, k, msg)
		case types.MsgIssueFreeze:
			return handleMsgIssueFreeze(ctx, k, msg)
		case types.MsgIssueUnFreeze:
			return handleMsgIssueUnFreeze(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized gov msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

//Handle handleMsgIssueDecreaseApproval
func handleMsgIssueDecreaseApproval(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueDecreaseApproval) sdk.Result {

	if err := k.DecreaseApproval(ctx, msg.FromAddress, msg.ToAddress, msg.IssueId, msg.Amount); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueIncreaseApproval
func handleMsgIssueIncreaseApproval(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssueIncreaseApproval) sdk.Result {

	if err := keeper.IncreaseApproval(ctx, msg.FromAddress, msg.ToAddress, msg.IssueId, msg.Amount); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueApprove
func handleMsgIssueApprove(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgIssueApprove) sdk.Result {

	if err := keeper.Approve(ctx, msg.FromAddress, msg.ToAddress, msg.IssueId, msg.Amount); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueBurnFrom
func handleMsgIssueBurnFrom(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueBurnFrom) sdk.Result {
	fee := k.GetParams(ctx).BurnFromFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}

	_, err := k.BurnFrom(ctx, msg.IssueId, msg.Amount, msg.FromAddress, msg.ToAddress)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueBurnHolder
func handleMsgIssueBurnHolder(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueBurnHolder) sdk.Result {
	_, err := k.BurnHolder(ctx, msg.IssueId, msg.Amount, msg.FromAddress)

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueBurnOwner
func handleMsgIssueBurnOwner(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueBurnOwner) sdk.Result {
	fee := k.GetParams(ctx).BurnFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}
	_, err := k.BurnOwner(ctx, msg.IssueId, msg.Amount, msg.FromAddress)

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueDescription
func handleMsgIssueDescription(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueDescription) sdk.Result {
	fee := k.GetParams(ctx).DescribeFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}
	if err := k.SetIssueDescription(ctx, msg.IssueId, msg.FromAddress, msg.Description); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueDisableFeature
func handleMsgIssueDisableFeature(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueDisableFeature) sdk.Result {
	if err := k.DisableFeature(ctx, msg.FromAddress, msg.IssueId, msg.Feature); err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueFreeze
func handleMsgIssueFreeze(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueFreeze) sdk.Result {
	fee := k.GetParams(ctx).FreezeFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}
	if err := k.Freeze(ctx, msg.IssueId, msg.FromAddress, msg.ToAddress, msg.FreezeType); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueMint
func handleMsgIssueMint(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueMint) sdk.Result {
	fee := k.GetParams(ctx).MintFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}
	_, err := k.Mint(ctx, msg.IssueId, msg.Amount, msg.FromAddress, msg.ToAddress)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle MsgIssueSendFrom
func handleMsgIssueSendFrom(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueSendFrom) sdk.Result {

	if err := k.SendFrom(ctx, msg.FromAddress, msg.From, msg.ToAddress, msg.IssueId, msg.Amount); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueTransferOwnership
func handleMsgIssueTransferOwnership(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueTransferOwnership) sdk.Result {
	fee := k.GetParams(ctx).TransferOwnerFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}
	if err := k.TransferOwnership(ctx, msg.IssueId, msg.FromAddress, msg.ToAddress); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssueUnFreeze
func handleMsgIssueUnFreeze(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssueUnFreeze) sdk.Result {
	fee := k.GetParams(ctx).UnFreezeFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}
	if err := k.UnFreeze(ctx, msg.IssueId, msg.FromAddress, msg.ToAddress, msg.FreezeType); err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//Handle handleMsgIssue
func handleMsgIssue(ctx sdk.Context, k keeper.Keeper, msg types.MsgIssue) sdk.Result {
	fee := k.GetParams(ctx).IssueFee
	if err := k.Fee(ctx, msg.FromAddress, fee); err != nil {
		return err.Result()
	}

	coinIssueInfo := types.CoinIssueInfo{
		Owner:              msg.FromAddress,
		Issuer:             msg.FromAddress,
		Name:               msg.Name,
		Symbol:             strings.ToUpper(msg.Symbol),
		TotalSupply:        msg.TotalSupply,
		Decimals:           msg.Decimals,
		Description:        msg.Description,
		BurnOwnerDisabled:  msg.BurnOwnerDisabled,
		BurnHolderDisabled: msg.BurnHolderDisabled,
		BurnFromDisabled:   msg.BurnFromDisabled,
		MintingFinished:    msg.MintingFinished,
		FreezeDisabled:     msg.FreezeDisabled,
	}

	_, err := k.CreateIssue(ctx, &coinIssueInfo)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute("issue-id", coinIssueInfo.IssueId),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

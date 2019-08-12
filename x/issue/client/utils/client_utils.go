package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	issuequeriers "github.com/zar-network/zar-network/x/issue/client/queriers"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

func burnCheck(sender auth.Account, burnFrom sdk.AccAddress, issueInfo types.Issue, amount sdk.Int, burnType string, cli bool) error {
	coins := sender.GetCoins()
	switch burnType {
	case types.BurnOwner:
		{
			if !sender.GetAddress().Equals(issueInfo.GetOwner()) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			if !sender.GetAddress().Equals(burnFrom) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			if issueInfo.IsBurnOwnerDisabled() {
				return types.Errorf(types.ErrCanNotBurn())
			}
			break
		}
	case types.BurnHolder:
		{
			if issueInfo.IsBurnHolderDisabled() {
				return types.Errorf(types.ErrCanNotBurn())
			}
			if !sender.GetAddress().Equals(burnFrom) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			break
		}
	case types.BurnFrom:
		{
			if !sender.GetAddress().Equals(issueInfo.GetOwner()) {
				return types.Errorf(types.ErrOwnerMismatch())
			}
			if issueInfo.IsBurnFromDisabled() {
				return types.Errorf(types.ErrCanNotBurn())
			}
			if issueInfo.GetOwner().Equals(burnFrom) {
				//burnFrom
				if issueInfo.IsBurnOwnerDisabled() {
					return types.Errorf(types.ErrCanNotBurn())
				}
			}
			break
		}
	default:
		{
			panic("not support")
		}

	}
	if cli {
		amount = types.MulDecimals(amount, issueInfo.GetDecimals())
	}
	// ensure account has enough coins
	if !coins.IsAllGTE(sdk.NewCoins(sdk.NewCoin(issueInfo.GetIssueId(), amount))) {
		return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", sender.GetAddress())
	}
	return nil
}

func GetBurnMsg(cdc *codec.Codec, cliCtx context.CLIContext, sender auth.Account,
	burnFrom sdk.AccAddress, issueID string, amount sdk.Int, burnFromType string, cli bool) (sdk.Msg, error) {
	issueInfo, err := GetIssueByID(cdc, cliCtx, issueID)
	if err != nil {
		return nil, err
	}
	if types.BurnHolder == burnFromType {
		if issueInfo.GetOwner().Equals(sender.GetAddress()) {
			burnFromType = types.BurnOwner
		}
	}
	err = burnCheck(sender, burnFrom, issueInfo, amount, burnFromType, cli)
	if err != nil {
		return nil, err
	}
	if cli {
		amount = types.MulDecimals(amount, issueInfo.GetDecimals())
	}
	var msg sdk.Msg
	switch burnFromType {

	case types.BurnOwner:
		msg = types.NewMsgIssueBurnOwner(issueID, sender.GetAddress(), amount)
		break
	case types.BurnHolder:
		msg = types.NewMsgIssueBurnHolder(issueID, sender.GetAddress(), amount)
		break
	case types.BurnFrom:
		msg = types.NewMsgIssueBurnFrom(issueID, sender.GetAddress(), burnFrom, amount)
		break
	default:
		return nil, types.ErrCanNotBurn()
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, types.Errorf(err)
	}
	return msg, nil
}

func GetIssueFreezeMsg(cdc *codec.Codec, cliCtx context.CLIContext, account auth.Account, freezeType string, issueID string, address string, endTime string, freeze bool) (sdk.Msg, error) {
	_, ok := types.FreezeTypes[freezeType]
	if !ok {
		return nil, types.Errorf(types.ErrUnknownFreezeType())
	}
	if err := types.CheckIssueId(issueID); err != nil {
		return nil, types.Errorf(err)
	}
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	issueInfo, err := IssueOwnerCheck(cdc, cliCtx, account, issueID)
	if err != nil {
		return nil, err
	}
	if freeze {
		if issueInfo.IsFreezeDisabled() {
			return nil, types.ErrCanNotFreeze()
		}
		freezeEndTime, err := strconv.ParseInt(endTime, 10, 64)
		if err != nil {
			return nil, types.Errorf(types.ErrFreezeEndTimestampNotValid())
		}
		msg := types.NewMsgIssueFreeze(issueID, account.GetAddress(), accAddress, freezeType, freezeEndTime)
		if err := msg.ValidateService(); err != nil {
			return msg, types.Errorf(err)
		}
		return msg, nil
	}
	msg := types.NewMsgIssueUnFreeze(issueID, account.GetAddress(), accAddress, freezeType)
	if err := msg.ValidateBasic(); err != nil {
		return msg, types.Errorf(err)
	}
	return msg, nil
}

func GetIssueApproveMsg(cdc *codec.Codec, cliCtx context.CLIContext, issueID string, account auth.Account, accAddress sdk.AccAddress, approveType string, amount sdk.Int, cli bool) (sdk.Msg, error) {
	if err := types.CheckIssueId(issueID); err != nil {
		return nil, types.Errorf(err)
	}
	issueInfo, err := GetIssueByID(cdc, cliCtx, issueID)
	if err != nil {
		return nil, err
	}
	if cli {
		amount = types.MulDecimals(amount, issueInfo.GetDecimals())
	}
	var msg sdk.Msg
	switch approveType {
	case types.Approve:
		msg = types.NewMsgIssueApprove(issueID, account.GetAddress(), accAddress, amount)
		break
	case types.IncreaseApproval:
		msg = types.NewMsgIssueIncreaseApproval(issueID, account.GetAddress(), accAddress, amount)
		break
	case types.DecreaseApproval:
		msg = types.NewMsgIssueDecreaseApproval(issueID, account.GetAddress(), accAddress, amount)
		break
	default:
		return nil, sdk.ErrInternal("not support")
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, types.Errorf(err)
	}
	return msg, nil
}

func CheckAllowance(cdc *codec.Codec, cliCtx context.CLIContext, issueID string, owner sdk.AccAddress, spender sdk.AccAddress, amount sdk.Int) error {
	res, err := issuequeriers.QueryIssueAllowance(issueID, owner, spender, cliCtx)
	if err != nil {
		return err
	}
	var approval types.Approval
	cdc.MustUnmarshalJSON(res, &approval)

	if approval.Amount.LT(amount) {
		return types.Errorf(types.ErrNotEnoughAmountToTransfer())
	}
	return nil

}

func GetIssueByID(cdc *codec.Codec, cliCtx context.CLIContext, issueID string) (types.Issue, error) {
	var issueInfo types.Issue
	// Query the issue
	res, err := issuequeriers.QueryIssueByID(issueID, cliCtx)
	if err != nil {
		return nil, err
	}
	cdc.MustUnmarshalJSON(res, &issueInfo)
	return issueInfo, nil
}

func IssueOwnerCheck(cliCtx context.CLIContext, sender auth.Account, issueID string) (types.Issue, error) {
	var issueInfo types.Issue
	// Query the issue
	res, err := issuequeriers.QueryIssueByID(issueID, cliCtx)
	if err != nil {
		return nil, err
	}
	cdc.MustUnmarshalJSON(res, &issueInfo)

	if !sender.GetAddress().Equals(issueInfo.GetOwner()) {
		return nil, types.Errorf(types.ErrOwnerMismatch(issueID))
	}
	return issueInfo, nil
}

func checkFreezeByOut(issueID string, freeze types.IssueFreeze, from sdk.AccAddress) sdk.Error {
	if freeze.OutEndTime > 0 && time.Unix(freeze.OutEndTime, 0).After(time.Now()) {
		return types.ErrCanNotTransferOut()
	}
	return nil
}

func checkFreezeByIn(issueID string, freeze types.IssueFreeze, to sdk.AccAddress) sdk.Error {
	if freeze.InEndTime > 0 && time.Unix(freeze.InEndTime, 0).After(time.Now()) {
		return types.ErrCanNotTransferIn()
	}
	return nil
}

func CheckFreeze(cdc *codec.Codec, cliCtx context.CLIContext, issueID string, from sdk.AccAddress, to sdk.AccAddress) error {
	res, err := issuequeriers.QueryIssueFreeze(issueID, from, cliCtx)
	if err != nil {
		return err
	}
	var freeze types.IssueFreeze
	cdc.MustUnmarshalJSON(res, &freeze)

	if checkErr := checkFreezeByOut(issueID, freeze, from); checkErr != nil {
		return types.Errorf(checkErr)
	}
	res, err = issuequeriers.QueryIssueFreeze(issueID, to, cliCtx)
	if err != nil {
		return err
	}
	cdc.MustUnmarshalJSON(res, &freeze)
	if checkErr := checkFreezeByIn(issueID, freeze, to); checkErr != nil {
		return types.Errorf(checkErr)
	}
	return nil
}

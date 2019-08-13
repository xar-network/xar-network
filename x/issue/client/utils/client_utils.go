package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

func GetQueryIssuePath(issueID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryIssue, issueID)
}
func GetQueryParamsPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryParams)
}
func GetQueryIssueAllowancePath(issueID string, owner sdk.AccAddress, spender sdk.AccAddress) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryAllowance, issueID, owner.String(), spender.String())
}
func GetQueryIssueFreezePath(issueID string, accAddress sdk.AccAddress) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryFreeze, issueID, accAddress.String())
}
func GetQueryIssueFreezesPath(issueID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryFreezes, issueID)
}
func GetQueryIssueSearchPath(symbol string) string {
	return fmt.Sprintf("%s/%s/%s/%s", types.Custom, types.QuerierRoute, types.QuerySearch, symbol)
}
func GetQueryIssuesPath() string {
	return fmt.Sprintf("%s/%s/%s", types.Custom, types.QuerierRoute, types.QueryIssues)
}

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

func GetBurnMsg(
	cliCtx context.CLIContext,
	sender auth.Account,
	burnFrom sdk.AccAddress,
	issueID string,
	amount sdk.Int,
	burnFromType string,
	cli bool,
) (sdk.Msg, error) {
	issueInfo, err := GetIssueByID(cliCtx, issueID)
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

func GetIssueFreezeMsg(
	cliCtx context.CLIContext,
	account auth.Account,
	freezeType string,
	issueID string,
	address string,
	freeze bool,
) (sdk.Msg, error) {
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
	issueInfo, err := IssueOwnerCheck(cliCtx, account, issueID)
	if err != nil {
		return nil, err
	}
	if freeze {
		if issueInfo.IsFreezeDisabled() {
			return nil, types.ErrCanNotFreeze()
		}
		msg := types.NewMsgIssueFreeze(issueID, account.GetAddress(), accAddress, freezeType)
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

func GetIssueApproveMsg(
	cliCtx context.CLIContext,
	issueID string,
	account auth.Account,
	accAddress sdk.AccAddress,
	approveType string,
	amount sdk.Int,
	cli bool,
) (sdk.Msg, error) {
	if err := types.CheckIssueId(issueID); err != nil {
		return nil, types.Errorf(err)
	}
	issueInfo, err := GetIssueByID(cliCtx, issueID)
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

func CheckAllowance(
	cliCtx context.CLIContext,
	issueID string,
	owner sdk.AccAddress,
	spender sdk.AccAddress,
	amount sdk.Int,
) error {
	res, _, err := cliCtx.QueryWithData(GetQueryIssueAllowancePath(issueID, owner, spender), nil)
	if err != nil {
		return err
	}
	var approval types.Approval
	cliCtx.Codec.MustUnmarshalJSON(res, &approval)

	if approval.Amount.LT(amount) {
		return types.Errorf(types.ErrNotEnoughAmountToTransfer())
	}
	return nil
}

func GetIssueByID(cliCtx context.CLIContext, issueID string) (types.Issue, error) {
	var issueInfo types.Issue
	// Query the issue
	res, _, err := cliCtx.QueryWithData(GetQueryIssuePath(issueID), nil)
	if err != nil {
		return nil, err
	}
	cliCtx.Codec.MustUnmarshalJSON(res, &issueInfo)
	return issueInfo, nil
}

func IssueOwnerCheck(cliCtx context.CLIContext, sender sdk.AccAddress, issueID string) (types.Issue, error) {
	var issueInfo types.Issue
	// Query the issue
	res, _, err := cliCtx.QueryWithData(GetQueryIssuePath(issueID), nil)
	if err != nil {
		return nil, err
	}
	cliCtx.Codec.MustUnmarshalJSON(res, &issueInfo)

	if !sender.Equals(issueInfo.GetOwner()) {
		return nil, types.Errorf(types.ErrOwnerMismatch())
	}
	return issueInfo, nil
}

func CheckFreeze(cdc *codec.Codec, cliCtx context.CLIContext, issueID string, from sdk.AccAddress, to sdk.AccAddress) error {
	res, _, err := cliCtx.QueryWithData(GetQueryIssueFreezePath(issueID, from), nil)
	if err != nil {
		return err
	}
	var freeze types.IssueFreeze
	cliCtx.Codec.MustUnmarshalJSON(res, &freeze)

	res, _, err = cliCtx.QueryWithData(GetQueryIssueFreezePath(issueID, to), nil)
	if err != nil {
		return err
	}
	cliCtx.Codec.MustUnmarshalJSON(res, &freeze)
	return nil
}

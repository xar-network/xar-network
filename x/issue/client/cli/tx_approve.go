package cli

import (
	"github.com/zar-network/zar-network/x/issue/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"
	issueutils "github.com/zar-network/zar-network/x/issue/utils"

	"github.com/zar-network/zar-network/x/issue/errors"
	"github.com/zar-network/zar-network/x/issue/msgs"
)

// GetCmdIssueSendFrom implements send from a token transaction command.
func GetCmdIssueSendFrom(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-from [issue-id] [from_address] [to_address] [amount]",
		Args:  cobra.ExactArgs(4),
		Short: "Send tokens from one address to another",
		Long:  "Send tokens from one address to another by allowance",
		Example: "$ zar-networkcli issue send-from coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n gard1vud9ptwagudgq7yht53cwuf8qfmgkd0qcej0ah " +
			"88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[0]
			if err := issueutils.CheckIssueId(issueID); err != nil {
				return errors.Errorf(err)
			}
			fromAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			toAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.Errorf(errors.ErrAmountNotValid(args[3]))
			}

			txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
			if err != nil {
				return err
			}

			if err := clientutils.CheckAllowance(cdc, cliCtx, issueID, fromAddress, account.GetAddress(), amount); err != nil {
				return err
			}

			if err = clientutils.CheckFreeze(cdc, cliCtx, issueID, fromAddress, toAddress); err != nil {
				return err
			}

			issueInfo, err := clientutils.GetIssueByID(cdc, cliCtx, issueID)
			if err != nil {
				return err
			}
			amount = issueutils.MulDecimals(amount, issueInfo.GetDecimals())

			msg := msgs.NewMsgIssueSendFrom(issueID, account.GetAddress(), fromAddress, toAddress, amount)

			validateErr := msg.ValidateBasic()
			if validateErr != nil {
				return errors.Errorf(validateErr)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
	return cmd
}

// GetCmdIssueApprove implements approve a token transaction command.
func GetCmdIssueApprove(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approve [issue-id] [address] [amount]",
		Args:    cobra.ExactArgs(3),
		Short:   "Approve spend tokens on behalf of sender",
		Long:    "Approve the passed address to spend the specified amount of tokens on behalf of sender",
		Example: "$ zar-networkcli issue approve coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueApprove(cdc, args, types.Approve)
		},
	}
	return cmd
}

// GetCmdIssueIncreaseApproval implements increase approval a token transaction command.
func GetCmdIssueIncreaseApproval(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "increase-approval [issue-id] [address] [amount]",
		Args:    cobra.ExactArgs(3),
		Short:   "Increase approve spend tokens on behalf of sender",
		Long:    "Increase approve the passed address to spend the specified amount of tokens on behalf of sender",
		Example: "$ zar-networkcli issue increase-approval coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueApprove(cdc, args, types.IncreaseApproval)
		},
	}
	return cmd
}

// GetCmdIssueDecreaseApproval implements decrease approval a token transaction command.
func GetCmdIssueDecreaseApproval(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrease-approval [issue-id] [address] [amount]",
		Args:    cobra.ExactArgs(3),
		Short:   "Decrease approve spend tokens on behalf of sender",
		Long:    "Decrease approve the passed address to spend the specified amount of tokens on behalf of sender",
		Example: "$ zar-networkcli issue increase-approval coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueApprove(cdc, args, types.DecreaseApproval)
		},
	}
	return cmd
}
func issueApprove(cdc *codec.Codec, args []string, approveType string) error {
	issueID := args[0]
	accAddress, err := sdk.AccAddressFromBech32(args[1])
	if err != nil {
		return err
	}
	amount, ok := sdk.NewIntFromString(args[2])
	if !ok {
		return errors.Errorf(errors.ErrAmountNotValid(args[2]))
	}
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	msg, err := clientutils.GetIssueApproveMsg(cdc, cliCtx, issueID, account, accAddress, approveType, amount, true)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

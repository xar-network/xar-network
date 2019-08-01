package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"
	issueutils "github.com/zar-network/zar-network/x/issue/utils"

	"github.com/zar-network/zar-network/x/issue/errors"
	"github.com/zar-network/zar-network/x/issue/types"
)

// GetCmdIssueBurnFrom implements burn a coinIssue transaction command.
func GetCmdIssueBurn(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "burn [issue-id] [amount]",
		Args:    cobra.ExactArgs(2),
		Short:   "Token holder burn the token",
		Long:    "Token holder or the Owner burn the token he holds (the Owner can burn if 'burning_owner_disabled' is false, the holder can burn if 'burning_holder_disabled' is false)",
		Example: "$ zar-networkcli issue burn coin174876e800 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueBurnFrom(cdc, args, types.BurnHolder)
		},
	}
	return cmd
}

// GetCmdIssueBurnFrom implements burn a coinIssue transaction command.
func GetCmdIssueBurnFrom(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "burn-from [issue-id] [from-address] [amount]",
		Args:    cobra.ExactArgs(3),
		Short:   "Token owner burn the token",
		Long:    "Token Owner burn the token from any holder (the Owner can burn if 'burning_any_disabled' is false)",
		Example: "$ zar-networkcli issue burn-from coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 88888 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueBurnFrom(cdc, args, types.BurnFrom)
		},
	}
	return cmd
}

func issueBurnFrom(cdc *codec.Codec, args []string, burnFromType string) error {
	issueID := args[0]
	if err := issueutils.CheckIssueId(issueID); err != nil {
		return errors.Errorf(err)
	}
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	amountStr := ""
	//burn sender
	accAddress := account.GetAddress()

	if types.BurnFrom == burnFromType {
		accAddress, err = sdk.AccAddressFromBech32(args[1])
		if err != nil {
			return err
		}
		amountStr = args[2]
	} else {
		amountStr = args[1]
	}
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		return errors.Errorf(errors.ErrAmountNotValid(amountStr))
	}
	msg, err := clientutils.GetBurnMsg(cdc, cliCtx, account, accAddress, issueID, amount, burnFromType, true)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

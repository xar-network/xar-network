package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"

	"github.com/zar-network/zar-network/x/issue/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetCmdIssueUnFreeze implements freeze a token transaction command.
func GetCmdIssueFreeze(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze [freeze-type] [issue-id] [acc-address] [end-time]",
		Args:  cobra.ExactArgs(4),
		Short: "Freeze the transfer from a address",
		Long: fmt.Sprintf("Token owner freeze the transfer from a address:\n\n"+
			"%s:The address can not transfer in before the end time\n"+
			"%s:The address can not transfer out before the end time\n"+
			"%s:The address not can transfer in and out before the end time\n\n", types.FreezeIn, types.FreezeOut, types.FreezeInAndOut) +
			"Note:The end-time is unix timestamp.\nExample:date -d \"2020-01-01 10:30:00\" +%s",
		Example: "$ zar-networkcli issue freeze in coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 1577845800 --from foo\n" +
			"$ zar-networkcli issue freeze out coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 1577845800 --from foo\n" +
			"$ zar-networkcli issue freeze in-out coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 1577845800 --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueFreeze(cdc, args, true)
		},
	}
	return cmd
}

// GetCmdIssueUnFreeze implements un freeze  a token transaction command.
func GetCmdIssueUnFreeze(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfreeze [freeze-type] [issue-id] [acc-address]",
		Args:  cobra.ExactArgs(3),
		Short: "UnFreeze the transfer from a address",
		Long: fmt.Sprintf("Token owner unFreeze the transfer from a address:\n\n"+
			"%s:The address can transfer in\n"+
			"%s:The address can transfer out\n"+
			"%s:The address can transfer in and out", types.FreezeIn, types.FreezeOut, types.FreezeInAndOut),
		Example: "$ zar-networkcli issue unfreeze in coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n --from foo\n" +
			"$ zar-networkcli issue unfreeze out coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n --from foo\n" +
			"$ zar-networkcli issue unfreeze in-out coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueFreeze(cdc, args, false)
		},
	}
	return cmd
}

func issueFreeze(cdc *codec.Codec, args []string, freeze bool) error {
	txBldr, cliCtx, account, err := clientutils.GetCliContext(cdc)
	if err != nil {
		return err
	}
	endTime := ""
	if freeze {
		endTime = args[3]
	}
	msg, err := clientutils.GetIssueFreezeMsg(cdc, cliCtx, account, args[0], args[1], args[2], endTime, freeze)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
}

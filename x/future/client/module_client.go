package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/hashgard/hashgard/x/box/types"
	"github.com/hashgard/hashgard/x/future/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	cdc *amino.Codec
}

//New ModuleClient Instance
func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{cdc}
}

// GetCmd returns the future box commands for this module
func (mc ModuleClient) GetCmd() *cobra.Command {
	boxCmd := &cobra.Command{
		Use:   types.Future,
		Short: "Future payment protocol (HRC13) subcommands",
		Long:  "FutureBox is a system native special payment box that can be set to pay different amounts of tokens for different users at multiple times. And you can set whether the payment certificate has a trading function. User deposits ones own token to the FutureBox and sets the account address, amounts to be paid, payment time and whether to support the receivable certificate trading function. After the setting is completed, the receiving account will get a receivable certificate. It can be traded according to the issuer's settings. Upon expiration, the system will automatically convert the user's receivable certificate into a 1:1 spot token. Forward payment protocol can be used in financial areas such as bonds, checks, futures and other application scenarios.",
	}
	boxCmd.AddCommand(
		client.GetCommands(
			cli.GetQueryParamsCmd(mc.cdc),
			cli.GetQueryCmd(mc.cdc),
			cli.GetListCmd(mc.cdc),
			cli.GetSearchCmd(mc.cdc),
		)...)
	boxCmd.AddCommand(client.LineBreak)

	txCmd := client.PostCommands(
		cli.GetCreateCmd(mc.cdc),
		cli.GetInjectCmd(mc.cdc),
		cli.GetCancelDepositCmd(mc.cdc),
		cli.GetDescriptionCmd(mc.cdc),
		cli.GetDisableFeatureCmd(mc.cdc),
	)

	for _, cmd := range txCmd {
		_ = cmd.MarkFlagRequired(client.FlagFrom)
		boxCmd.AddCommand(cmd)
	}

	return boxCmd
}

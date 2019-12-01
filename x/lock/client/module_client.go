package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/xar-network/xar-network/x/escrow/types"
	"github.com/xar-network/xar-network/x/lock/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	cdc *amino.Codec
}

//New ModuleClient Instance
func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{cdc}
}

// GetLockCmd returns the box commands for this module
func (mc ModuleClient) GetCmd() *cobra.Command {
	boxCmd := &cobra.Command{
		Use:   types.Lock,
		Short: "Token lock-up subcommands",
		Long: "Xar supports the token lock-up function, financial institutions and users can lock up tokens for a limited period of time in financial activities. It has following features:" +
			"		\n1.Supports this native function from the bottom of the blockchain. It is simple, secure and efficient. Users need to set the token type, amount and lock-up time of the token." +
			"		\n2.Information is transparent and available for inquiry.",
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
		cli.GetDescriptionCmd(mc.cdc),
	)

	for _, cmd := range txCmd {
		_ = cmd.MarkFlagRequired(client.FlagFrom)
		boxCmd.AddCommand(cmd)
	}

	return boxCmd
}

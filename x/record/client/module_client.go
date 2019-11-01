package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/xar-network/xar-network/x/record/client/cli"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	cdc *amino.Codec
}

//New ModuleClient Instance
func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{cdc}
}

// GetIssueCmd returns the record commands for this module
func (mc ModuleClient) GetCmd() *cobra.Command {
	recordCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Xar native recording service subcommands",
		Long:  "Record a 64 characters hash of any user data into the chain.",
	}
	recordCmd.AddCommand(
		client.GetCommands(
			cli.GetCmdQueryRecord(mc.cdc),
			cli.GetCmdQueryList(mc.cdc),
		)...)
	recordCmd.AddCommand(client.LineBreak)

	txCmd := client.PostCommands(
		cli.GetCmdRecordCreate(mc.cdc),
	)

	for _, cmd := range txCmd {
		_ = cmd.MarkFlagRequired(client.FlagFrom)
		recordCmd.AddCommand(cmd)
	}

	return recordCmd
}

package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/zar-network/zar-network/x/pool/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

// NewModuleClient creates client for the module
func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "pool",
		Short: "Querying commands for the pool module",
	}

	queryCmd.AddCommand(client.GetCommands(
		cli.GetCmdGetFunds(mc.storeKey, mc.cdc),
		cli.GetCmdGetAllFunds(mc.storeKey, mc.cdc),
	)...)

	return queryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "pool",
		Short: "Pool transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		cli.GetCmdDepositFunds(mc.cdc),
		cli.GetCmdWithdrawFunds(mc.cdc),
	)...)

	return txCmd
}

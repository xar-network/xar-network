package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	csdtcmd "github.com/xar-network/xar-network/x/csdt/client/cli"
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
	// Group nameservice queries under a subcommand
	csdtQueryCmd := &cobra.Command{
		Use:   "csdt",
		Short: "Querying commands for the csdt module",
	}

	csdtQueryCmd.AddCommand(client.GetCommands(
		csdtcmd.GetCmd_GetCsdt(mc.storeKey, mc.cdc),
		csdtcmd.GetCmd_GetCsdts(mc.storeKey, mc.cdc),
		csdtcmd.GetCmd_GetUnderCollateralizedCsdts(mc.storeKey, mc.cdc),
		csdtcmd.GetCmd_GetParams(mc.storeKey, mc.cdc),
	)...)

	return csdtQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	csdtTxCmd := &cobra.Command{
		Use:   "csdt",
		Short: "csdt transactions subcommands",
	}

	csdtTxCmd.AddCommand(client.PostCommands(
		csdtcmd.GetCmdModifyCsdt(mc.cdc),
	)...)

	return csdtTxCmd
}

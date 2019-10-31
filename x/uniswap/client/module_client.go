package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/xar-network/xar-network/x/uniswap/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	stakingQueryCmd := &cobra.Command{
		Use: "uniswap",
	}

	return stakingQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	swapTxCmd := &cobra.Command{
		Use: "uniswap",
	}

	swapTxCmd.AddCommand(client.PostCommands(
		cli.GetSwapTxCmd(mc.cdc),
	)...)

	return swapTxCmd
}

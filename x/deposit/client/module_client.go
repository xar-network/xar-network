package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/hashgard/hashgard/x/box/types"
	"github.com/hashgard/hashgard/x/deposit/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	cdc *amino.Codec
}

//New ModuleClient Instance
func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{cdc}
}

// GetLockCmd returns the deposit box commands for this module
func (mc ModuleClient) GetCmd() *cobra.Command {
	boxCmd := &cobra.Command{
		Use:   types.Deposit,
		Short: "Deposit protocol (HRC12) subcommands",
		Long: "In Hashgard, users (called box owners) can create a \"deposit box\" acting as a timed deposit " +
			"service in a bank. Other users (called box \"investors\") can add tokens of a specified type and " +
			"limited number into the box to obtain certificates, using which the investors can receive the " +
			"principal back with corresponding interest upon deposit maturity. Issuer can set whether the " +
			"certificate can be traded and transferred and whether the principal and interest are the same token type.",
	}
	boxCmd.AddCommand(
		client.GetCommands(
			cli.GetQueryParamsCmd(mc.cdc),
			cli.GetQueryCmd(mc.cdc),
			cli.GetListCmd(mc.cdc),
			cli.GetSearchCmd(mc.cdc),
		)...)
	boxCmd.AddCommand(client.LineBreak)

	cmdCreate := cli.GetCreateCmd(mc.cdc)
	cli.MarkCmdDepositBoxCreateFlagRequired(cmdCreate)

	txCmd := client.PostCommands(
		cmdCreate,
		cli.GetInterestInjectCmd(mc.cdc),
		cli.GetInterestCancelCmd(mc.cdc),
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

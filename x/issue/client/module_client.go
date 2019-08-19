package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	issueCli "github.com/zar-network/zar-network/x/issue/client/cli"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	cdc *amino.Codec
}

//New ModuleClient Instance
func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{cdc}
}

// GetIssueCmd returns the issue commands for this module
func (mc ModuleClient) GetCmd() *cobra.Command {
	issueCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Fungible Token Issuance Protocol (HRC10) subcommands",
		Long: "Hashgard supports the issuance of fungible utility tokens with similar functions to the well-known ERC20 token of Ethereum, but with following features." +
			"		\n1.Token issuance is supported from the bottom level of the blockchain, rather than using smart contracts. Users only need to call the standard system interface to issue tokens instead of using solidity or other languages. Given the many security vulnerabilities in Ethereum smart contracts, we realize that user-customized smart contracts contain high potential security risks. To reduce risks, user programming for standard functional components should be minimized." +
			"		\n2.The creator of a token is called its owner. This user has ownership of the token and can modify and configure its parameters.",
	}
	issueCmd.AddCommand(
		client.GetCommands(
			issueCli.GetQueryParamsCmd(mc.cdc),
			issueCli.GetCmdQueryIssues(mc.cdc),
			issueCli.GetCmdQueryFreezes(mc.cdc),
			issueCli.GetCmdQueryIssue(mc.cdc),
			issueCli.GetCmdQueryAllowance(mc.cdc),
			issueCli.GetCmdQueryFreeze(mc.cdc),
			issueCli.GetCmdSearchIssues(mc.cdc),
		)...)
	issueCmd.AddCommand(client.LineBreak)

	txCmd := client.PostCommands(
		issueCli.GetCmdIssueApprove(mc.cdc),
		issueCli.GetCmdIssueBurn(mc.cdc),
		issueCli.GetCmdIssueBurnFrom(mc.cdc),
		issueCli.IssueCreateCmd(mc.cdc),
		issueCli.IssueDescriptionCmd(mc.cdc),
		issueCli.GetCmdIssueDecreaseApproval(mc.cdc),
		issueCli.GetCmdIssueFreeze(mc.cdc),
		issueCli.GetCmdIssueUnFreeze(mc.cdc),
		issueCli.GetCmdIssueIncreaseApproval(mc.cdc),
		issueCli.IssueMintCmd(mc.cdc),
		issueCli.GetCmdIssueSendFrom(mc.cdc),
		issueCli.IssueTransferOwnershipCmd(mc.cdc),
		client.LineBreak,
		issueCli.IssueDisableFeatureCmd(mc.cdc),
	)

	for _, cmd := range txCmd {
		_ = cmd.MarkFlagRequired(client.FlagFrom)
		issueCmd.AddCommand(cmd)
	}

	return issueCmd
}

package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	boxcli "github.com/hashgard/hashgard/x/box/client/cli"
	"github.com/hashgard/hashgard/x/box/types"
	"github.com/spf13/cobra"
)

// GetQueryParamsCmd implements the query params command.
func GetQueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "params",
		Short:   "Query the parameters of the lock process",
		Long:    "Query the all the parameters",
		Example: "$ hashgardcli lock params",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessQueryBoxParamsCmd(cdc, types.Lock)
		},
	}
}

// GetQueryCmd implements the query box command.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "query [id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query a single lock",
		Long:    "Query details for a lock. You can find the id by running hashgardcli lock list",
		Example: "$ hashgardcli lock query boxab3jlxpt2ps",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessQueryBoxCmd(cdc, types.Lock, args[0])
		},
	}
}

// GetListCmd implements the query box command.
func GetListCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Query lock list",
		Long:    "Query all or one of the account lock list, the limit default is 30",
		Example: "$ hashgardcli lock list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessListBoxCmd(cdc, types.Lock)
		},
	}

	cmd.Flags().String(boxcli.FlagAddress, "", "Lock owner address")
	cmd.Flags().String(boxcli.FlagStartId, "", "Start id of lock results")
	cmd.Flags().Int32(boxcli.FlagLimit, 30, "Query number of lock results per page returned")

	return cmd
}

// GetSearchCmd implements the query box command.
func GetSearchCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search [name]",
		Args:    cobra.ExactArgs(1),
		Short:   "Search lock",
		Long:    "Search lock based on name",
		Example: "$ hashgardcli lock search fo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessSearchBoxsCmd(cdc, types.Lock, args[0])
		},
	}
	return cmd
}

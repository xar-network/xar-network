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
		Short:   "Query the parameters of the future box process",
		Long:    "Query the all the parameters",
		Example: "$ hashgardcli future params",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessQueryBoxParamsCmd(cdc, types.Future)
		},
	}
}

// GetQueryCmd implements the query box command.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "query [id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query a single future box",
		Long:    "Query details for a future box. You can find the id by running hashgardcli future box list",
		Example: "$ hashgardcli future box query boxab3jlxpt2ps",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessQueryBoxCmd(cdc, types.Future, args[0])
		},
	}
}

// GetListCmd implements the query box command.
func GetListCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Query future box list",
		Long:    "Query all or one of the account future box list, the limit default is 30",
		Example: "$ hashgardcli future box list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessListBoxCmd(cdc, types.Future)
		},
	}

	cmd.Flags().String(boxcli.FlagAddress, "", "Future box owner address")
	cmd.Flags().String(boxcli.FlagStartId, "", "Start id of future box results")
	cmd.Flags().Int32(boxcli.FlagLimit, 30, "Query number of future box results per page returned")

	return cmd
}

// GetSearchCmd implements the query box command.
func GetSearchCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search [name]",
		Args:    cobra.ExactArgs(1),
		Short:   "Search future box",
		Long:    "Search future box based on name",
		Example: "$ hashgardcli future box search fo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessSearchBoxsCmd(cdc, types.Future, args[0])
		},
	}
	return cmd
}

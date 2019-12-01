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
		Short:   "Query the parameters of the deposit box process",
		Long:    "Query the all the parameters",
		Example: "$ hashgardcli deposit params",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessQueryBoxParamsCmd(cdc, types.Deposit)
		},
	}
}

// GetQueryCmd implements the query box command.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "query [id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query a single deposit box",
		Long:    "Query details for a deposit box. You can find the id by running hashgardcli deposit box list",
		Example: "$ hashgardcli deposit box query boxab3jlxpt2ps",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessQueryBoxCmd(cdc, types.Deposit, args[0])
		},
	}
}

// GetListCmd implements the query box command.
func GetListCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Query deposit box list",
		Long:    "Query all or one of the account deposit box list, the limit default is 30",
		Example: "$ hashgardcli deposit box list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessListBoxCmd(cdc, types.Deposit)
		},
	}

	cmd.Flags().String(boxcli.FlagAddress, "", "Deposit box owner address")
	cmd.Flags().String(boxcli.FlagStartId, "", "Start id of deposit box results")
	cmd.Flags().Int32(boxcli.FlagLimit, 30, "Query number of deposit box results per page returned")

	return cmd
}

// GetSearchCmd implements the query box command.
func GetSearchCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search [name]",
		Args:    cobra.ExactArgs(1),
		Short:   "Search deposit box",
		Long:    "Search deposit box based on name",
		Example: "$ hashgardcli deposit box search fo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxcli.ProcessSearchBoxsCmd(cdc, types.Deposit, args[0])
		},
	}
	return cmd
}

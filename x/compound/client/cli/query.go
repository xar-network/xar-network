package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/xar-network/xar-network/x/compound/internal/types"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the compound module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(client.GetCommands(
		GetCmdMarketInfo(storeKey, cdc),
		GetCmdMarketPosition(storeKey, cdc),
	)...)
	return queryCmd
}

// GetCmdWhois queries information about a domain
func GetCmdMarketInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "marketinfo [name]",
		Short: "Query market info of name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/compound/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve market info - %s \n", name)
				return nil
			}

			var out types.Compound
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetMarketPosition queries information about a domain
func GetCmdMarketPosition(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "marketposition [account]",
		Short: "Query market position for an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			owner := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/compound/%s", queryRoute, owner), nil)
			if err != nil {
				fmt.Printf("could not resolve market position - %s \n", owner)
				return nil
			}

			var out types.CompoundPosition
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

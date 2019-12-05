package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/xar-network/xar-network/x/denominations/internal/keeper"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the asset management module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(client.GetCommands(
		GetCmdFindToken(storeKey, cdc),
		GetCmdSymbols(storeKey, cdc),
	)...)
	return queryCmd
}

// GetCmdFindToken queries information about a token through its unique symbol
func GetCmdFindToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "find [symbol]",
		Short: "find symbol",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			symbol := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryToken, symbol), nil)
			if err != nil {
				fmt.Printf("could not find symbol - '%s'. reason: '%s'\n", symbol, err)
				return nil
			}

			var out types.Token
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdSymbols queries a list of all symbols
func GetCmdSymbols(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "symbols",
		Short: "symbols",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QuerySymbols), nil)
			if err != nil {
				fmt.Printf("could not get query symbols. reason: '%s'\n", err)
				return nil
			}

			var out types.QueryResultSymbol
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

const (
	nativeDenom = "atom"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	coinswapQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the coinswap module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	coinswapQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryLiquidity(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc))...)

	return coinswapQueryCmd
}

// GetCmdQueryLiquidity implements the liquidity query command
func GetCmdQueryLiquidity(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "liquidity [denom]",
		Short: "Query the current liquidity values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the liquidity of a specific trading pair stored in the reserve pool.

Example:
$ %s query coinswap liquidity btc
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Added a check to ensure that input provided is not a native denom
			if strings.Compare(strings.TrimSpace(args[0]), nativeDenom) == 0 {
				return fmt.Errorf("%s is not a valid denom, please input a valid denom", args[0])
			}

			bz, err := cdc.MarshalJSON(types.NewQueryLiquidityParams(args[0]))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryLiquidity)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var liquidity sdk.Coins
			if err := cdc.UnmarshalJSON(res, &liquidity); err != nil {
				return err
			}
			return cliCtx.PrintOutput(liquidity)
		},
	}
}

// GetCmdQueryParams implements the params query command
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters involved in the coinswap process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the parameters for the coinswap process.

Example:
$ %s query coinswap params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(bz, &params); err != nil {
				return err
			}
			return cliCtx.PrintOutput(params)
		},
	}
}

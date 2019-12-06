/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
)

// GetCmd_GetOutstandingDebt queries for the remaining available debt in the liquidator module after settlement with the module's stablecoin balance.
func GetCmd_GetOutstandingDebt(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "debt",
		Short: "get the outstanding seized debt",
		Long:  "Get the remaining available debt after settlement with the liquidator's stable coin balance.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetOutstandingDebt), nil)
			if err != nil {
				return err
			}

			var outstandingDebt sdk.Int
			cdc.MustUnmarshalJSON(res, &outstandingDebt)
			return cliCtx.PrintOutput(outstandingDebt)
		},
	}
}

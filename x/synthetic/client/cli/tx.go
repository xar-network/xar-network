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
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/xar-network/xar-network/x/synthetic/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	csdtTxCmd := &cobra.Command{
		Use:   "csdt",
		Short: "CSDT transactions subcommands",
	}

	csdtTxCmd.AddCommand(
		GetCmdBuySynthetic(cdc),
		GetCmdSellSynthetic(cdc),
	)

	return csdtTxCmd
}

// GetCmdBuySynthetic cli command for buying synthetics
func GetCmdBuySynthetic(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [from_key_or_addres] [denom] [amount]",
		Short: "buy a synthetic asset",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				fmt.Printf("invalid amount - %s \n", string(args[2]))
				return nil
			}
			coin := sdk.NewCoin(args[1], amount)
			if !coin.IsValid() {
				fmt.Printf("invalid coin - %s \n", string(args[1]))
			}
			msg := types.NewMsgBuySynthetic(cliCtx.GetFromAddress(), coin)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// GetCmdSellSynthetic cli command for selling synthetics
func GetCmdSellSynthetic(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sell [from_key_or_addres] [denom] [amount]",
		Short: "sell a synthetic asset",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				fmt.Printf("invalid amount - %s \n", string(args[2]))
				return nil
			}
			coin := sdk.NewCoin(args[1], amount)
			if !coin.IsValid() {
				fmt.Printf("invalid coin - %s \n", string(args[1]))
			}
			msg := types.NewMsgSellSynthetic(cliCtx.GetFromAddress(), coin)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

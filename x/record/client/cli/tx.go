/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
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

	"github.com/xar-network/xar-network/x/record/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagSender        = "sender"
	flagRecordType    = "record-type"
	flagAuthor        = "author"
	flagRecordNo      = "record-number"
	flagDescription   = "description"
	flagStartRecordId = "start-record-id"
	flagLimit         = "limit"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Record transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		GetCmdRecordCreate(cdc),
	)

	return cmd
}

// GetCmdIssue implements record a coin transaction command.
func GetCmdRecordCreate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [name] [hash]",
		Args:    cobra.ExactArgs(2),
		Short:   "Create a new record",
		Long:    "Create a new record",
		Example: "$ xarcli record create contractAEE321 BC38CAEE32149BEF4CCFAEAB518EC9A5FBC85AE6AC8D5A9F6CD710FAF5E4A2B8 --from live_key",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			para := types.RecordParams{
				Name:        args[0],
				Hash:        args[1],
				RecordType:  viper.GetString(flagRecordType),
				Author:      viper.GetString(flagAuthor),
				RecordNo:    viper.GetString(flagRecordNo),
				Description: viper.GetString(flagDescription),
			}
			msg := types.NewMsgRecord(cliCtx.GetFromAddress(), &para)

			validateErr := msg.ValidateBasic()

			if validateErr != nil {
				return types.Errorf(validateErr)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagRecordType, "", "customized record-type")
	cmd.Flags().String(flagAuthor, "", "author of the record data")
	cmd.Flags().String(flagRecordNo, "", "customized record-number")
	cmd.Flags().String(flagDescription, "", "customized description of the record")

	return cmd
}

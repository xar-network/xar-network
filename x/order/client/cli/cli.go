/*

Copyright 2019 All in Bits, Inc
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
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetQueryCmd(sk string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "order",
		Short: "queries orders",
	}
	queryCmd.AddCommand(client.GetCommands(
		GetCmdListOrders(sk, cdc),
	)...)
	return queryCmd
}

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "order",
		Short: "manages orders",
	}
	txCmd.AddCommand(client.PostCommands(
		GetCmdPost(cdc),
		GetCmdCancel(cdc),
	)...)
	return txCmd
}

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

package client

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/xar-network/xar-network/embedded/price/client/cli"
)

type ModuleClient struct {
	cdc *amino.Codec
}

func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{
		cdc: cdc,
	}
}

func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	priceQueryCmd := &cobra.Command{
		Use:   "price",
		Short: "queries price data",
	}
	priceQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdHistory(mc.cdc),
	)...)
	return priceQueryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	priceTxCmd := &cobra.Command{
		Use:   "price",
		Short: "manages price data",
	}
	return priceTxCmd
}

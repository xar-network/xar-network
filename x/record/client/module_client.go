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

package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/xar-network/xar-network/x/record/client/cli"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	cdc *amino.Codec
}

//New ModuleClient Instance
func NewModuleClient(cdc *amino.Codec) ModuleClient {
	return ModuleClient{cdc}
}

// GetIssueCmd returns the record commands for this module
func (mc ModuleClient) GetCmd() *cobra.Command {
	recordCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Xar native recording service subcommands",
		Long:  "Record a 64 characters hash of any user data into the chain.",
	}
	recordCmd.AddCommand(
		client.GetCommands(
			cli.GetCmdQueryRecord(mc.cdc),
			cli.GetCmdQueryList(mc.cdc),
		)...)
	recordCmd.AddCommand(client.LineBreak)

	txCmd := client.PostCommands(
		cli.GetCmdRecordCreate(mc.cdc),
	)

	for _, cmd := range txCmd {
		_ = cmd.MarkFlagRequired(client.FlagFrom)
		recordCmd.AddCommand(cmd)
	}

	return recordCmd
}

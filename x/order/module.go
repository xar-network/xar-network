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

package order

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/xar-network/xar-network/x/order/client/cli"
	types3 "github.com/xar-network/xar-network/x/order/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return types3.ModuleName
}

func (a AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types3.RegisterCodec(cdc)
}

func (a AppModuleBasic) DefaultGenesis() json.RawMessage {
	return []byte("{}")
}

func (a AppModuleBasic) ValidateGenesis(json.RawMessage) error {
	return nil
}

func (a AppModuleBasic) RegisterRESTRoutes(context.CLIContext, *mux.Router) {
}

func (a AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

func (a AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(types3.StoreKey, cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

func (a AppModule) InitGenesis(types.Context, json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (a AppModule) ExportGenesis(types.Context) json.RawMessage {
	return []byte("{}")
}

func (a AppModule) RegisterInvariants(types.InvariantRegistry) {
}

func (a AppModule) Route() string {
	return types3.RouterKey
}

func (a AppModule) NewHandler() types.Handler {
	return NewHandler(a.keeper)
}

func (a AppModule) QuerierRoute() string {
	return types3.RouterKey
}

func (a AppModule) NewQuerierHandler() types.Querier {
	return NewQuerier(a.keeper)
}

func (a AppModule) BeginBlock(types.Context, abci.RequestBeginBlock) {
}

func (a AppModule) EndBlock(types.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

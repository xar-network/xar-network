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

package synthetic

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/synthetic/client/cli"
	"github.com/xar-network/xar-network/x/synthetic/internal/keeper"
	"github.com/xar-network/xar-network/x/synthetic/internal/types"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/xar-network/xar-network/x/synthetic/client/rest"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic app module basics object
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name get module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the bank module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the bank module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the bank module.
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command { return nil }

// AppModule app module type
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name module name
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants register module invariants
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route module message route name
func (AppModule) Route() string {
	return ModuleName
}

// NewHandler module handler
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute module querier route name
func (AppModule) QuerierRoute() string {
	return ModuleName
}

// NewQuerierHandler module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// InitGenesis module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock performs a no-op.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the bank module. It returns no validator
// updates.

// Fees should be accrued here, need a global fee that can be used for buybacks or liquidation
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	balances := am.keeper.GetMarketBalances(ctx)
	for _, v := range balances {
		v.HandleBlockEvent(ctx.BlockHeader())
	}

	am.keeper.SetMarketBalances(ctx, balances)

	return []abci.ValidatorUpdate{}
}

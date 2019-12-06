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

package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/oracle/internal/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {

	// Set the assets and oracles from params
	keeper.SetParams(ctx, data.Params)

	// Iterate through the posted prices and set them in the store
	for _, pp := range data.PostedPrices {
		_, err := keeper.SetPrice(ctx, pp.OracleAddress, pp.AssetCode, pp.Price, pp.Expiry)
		if err != nil {
			panic(err)
		}
	}

	_ = keeper.SetCurrentPrices(ctx)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {

	// Get the params for assets and oracles
	params := keeper.GetParams(ctx)

	var postedPrices []PostedPrice
	for _, asset := range keeper.GetAssetParams(ctx) {
		pp := keeper.GetRawPrices(ctx, asset.AssetCode)
		postedPrices = append(postedPrices, pp...)
	}

	return types.GenesisState{
		Params:       params,
		PostedPrices: postedPrices,
	}
}

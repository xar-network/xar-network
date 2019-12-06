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

package auction

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/auction/internal/keeper"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

// InitGenesis - initializes the store state from genesis data
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.AuctionParams)

	for _, a := range data.Auctions {
		keeper.SetAuction(ctx, a)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)

	var genAuctions types.GenesisAuctions
	iterator := keeper.GetAuctionIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {

		auction := keeper.DecodeAuction(ctx, iterator.Value())
		genAuctions = append(genAuctions, auction)

	}
	return NewGenesisState(params, genAuctions)
}

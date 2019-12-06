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

package keeper

// DONTCOVER

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/nft/internal/types"
)

// RegisterInvariants registers all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(
		types.ModuleName, "supply",
		SupplyInvariant(k),
	)
}

// AllInvariants runs all invariants of the nfts module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		return SupplyInvariant(k)(ctx)
	}
}

// SupplyInvariant checks that the total amount of nfts on collections matches the total amount owned by addresses
func SupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		collectionsSupply := make(map[string]int)
		ownersCollectionsSupply := make(map[string]int)
		var msg string
		count := 0

		k.IterateCollections(ctx, func(collection types.Collection) bool {
			collectionsSupply[collection.Denom] = collection.Supply()
			return false
		})

		for _, owner := range k.GetOwners(ctx) {
			for _, idCollection := range owner.IDCollections {
				ownersCollectionsSupply[idCollection.Denom] += idCollection.Supply()
			}
		}

		for denom, supply := range collectionsSupply {
			if supply != ownersCollectionsSupply[denom] {
				count++
				msg += fmt.Sprintf("total %s NFTs supply invariance:\n"+
					"\ttotal %s NFTs supply: %d\n"+
					"\tsum of %s NFTs by owner: %d\n", denom, denom, supply, denom, ownersCollectionsSupply[denom])
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "supply", fmt.Sprintf(
			"%d NFT supply invariants found\n%s", count, msg)), broken
	}
}

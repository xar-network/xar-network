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

package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xar-network/xar-network/x/nft/internal/keeper"
	"github.com/xar-network/xar-network/x/nft/internal/types"
)

func TestSetCollection(t *testing.T) {
	app, ctx := createTestApp(false)

	// create a new nft with id = "id" and owner = "address"
	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// collection should exist
	collection, exists := app.NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)

	// create a new NFT and add it to the collection created with the NFT mint
	nft2 := types.NewBaseNFT(id2, address, tokenURI)
	collection2, err2 := collection.AddNFT(&nft2)
	require.NoError(t, err2)
	app.NFTKeeper.SetCollection(ctx, denom, collection2)

	collection2, exists = app.NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)
	require.Len(t, collection2.NFTs, 2)

	// reset collection for invariant sanity
	app.NFTKeeper.SetCollection(ctx, denom, collection)

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
func TestGetCollection(t *testing.T) {
	app, ctx := createTestApp(false)

	// collection shouldn't exist
	collection, exists := app.NFTKeeper.GetCollection(ctx, denom)
	require.Empty(t, collection)
	require.False(t, exists)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// collection should exist
	collection, exists = app.NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)
	require.NotEmpty(t, collection)

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
func TestGetCollections(t *testing.T) {
	app, ctx := createTestApp(false)

	// collections should be empty
	collections := app.NFTKeeper.GetCollections(ctx)
	require.Empty(t, collections)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// collections should equal 1
	collections = app.NFTKeeper.GetCollections(ctx)
	require.NotEmpty(t, collections)
	require.Equal(t, len(collections), 1)

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

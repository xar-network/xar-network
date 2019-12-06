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

func TestMintNFT(t *testing.T) {
	app, ctx := createTestApp(false)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// MintNFT shouldn't fail when collection exists
	nft2 := types.NewBaseNFT(id2, address, tokenURI)
	err = app.NFTKeeper.MintNFT(ctx, denom, &nft2)
	require.NoError(t, err)
}

func TestGetNFT(t *testing.T) {
	app, ctx := createTestApp(false)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// GetNFT should get the NFT
	receivedNFT, err := app.NFTKeeper.GetNFT(ctx, denom, id)
	require.NoError(t, err)
	require.Equal(t, receivedNFT.GetID(), id)
	require.True(t, receivedNFT.GetOwner().Equals(address))
	require.Equal(t, receivedNFT.GetTokenURI(), tokenURI)

	// MintNFT shouldn't fail when collection exists
	nft2 := types.NewBaseNFT(id2, address, tokenURI)
	err = app.NFTKeeper.MintNFT(ctx, denom, &nft2)
	require.NoError(t, err)

	// GetNFT should get the NFT when collection exists
	receivedNFT2, err := app.NFTKeeper.GetNFT(ctx, denom, id2)
	require.NoError(t, err)
	require.Equal(t, receivedNFT2.GetID(), id2)
	require.True(t, receivedNFT2.GetOwner().Equals(address))
	require.Equal(t, receivedNFT2.GetTokenURI(), tokenURI)

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

func TestUpdateNFT(t *testing.T) {
	app, ctx := createTestApp(false)

	nft := types.NewBaseNFT(id, address, tokenURI)

	// UpdateNFT should fail when NFT doesn't exists
	err := app.NFTKeeper.UpdateNFT(ctx, denom, &nft)
	require.Error(t, err)

	// MintNFT shouldn't fail when collection does not exist
	err = app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	nonnft := types.NewBaseNFT(id2, address, tokenURI)
	// UpdateNFT should fail when NFT doesn't exists
	err = app.NFTKeeper.UpdateNFT(ctx, denom, &nonnft)
	require.Error(t, err)

	// UpdateNFT shouldn't fail when NFT exists
	nft2 := types.NewBaseNFT(id, address, tokenURI2)
	err = app.NFTKeeper.UpdateNFT(ctx, denom, &nft2)
	require.NoError(t, err)

	// UpdateNFT shouldn't fail when NFT exists
	nft2 = types.NewBaseNFT(id, address2, tokenURI2)
	err = app.NFTKeeper.UpdateNFT(ctx, denom, &nft2)
	require.NoError(t, err)

	// GetNFT should get the NFT with new tokenURI
	receivedNFT, err := app.NFTKeeper.GetNFT(ctx, denom, id)
	require.NoError(t, err)
	require.Equal(t, receivedNFT.GetTokenURI(), tokenURI2)
}

func TestDeleteNFT(t *testing.T) {
	app, ctx := createTestApp(false)

	// DeleteNFT should fail when NFT doesn't exist and collection doesn't exist
	err := app.NFTKeeper.DeleteNFT(ctx, denom, id)
	require.Error(t, err)

	// MintNFT should not fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err = app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// DeleteNFT should fail when NFT doesn't exist but collection does exist
	err = app.NFTKeeper.DeleteNFT(ctx, denom, id2)
	require.Error(t, err)

	// DeleteNFT should not fail when NFT and collection exist
	err = app.NFTKeeper.DeleteNFT(ctx, denom, id)
	require.NoError(t, err)

	// NFT should no longer exist
	isNFT := app.NFTKeeper.IsNFT(ctx, denom, id)
	require.False(t, isNFT)

	owner := app.NFTKeeper.GetOwner(ctx, address)
	require.Equal(t, 0, owner.Supply())

	msg, fail := keeper.SupplyInvariant(app.NFTKeeper)(ctx)
	require.False(t, fail, msg)
}

func TestIsNFT(t *testing.T) {
	app, ctx := createTestApp(false)

	// IsNFT should return false
	isNFT := app.NFTKeeper.IsNFT(ctx, denom, id)
	require.False(t, isNFT)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, tokenURI)
	err := app.NFTKeeper.MintNFT(ctx, denom, &nft)
	require.NoError(t, err)

	// IsNFT should return true
	isNFT = app.NFTKeeper.IsNFT(ctx, denom, id)
	require.True(t, isNFT)
}

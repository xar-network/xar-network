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

package nft_test

import (
	"testing"
	"encoding/json"

	"github.com/stretchr/testify/require"

	"github.com/xar-network/xar-network/x/nft"
)

type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func TestInitGenesis(t *testing.T) {
	app, ctx := createTestApp(false)
	genesisState := nft.DefaultGenesisState()
	require.Equal(t, 0, len(genesisState.Owners))
	require.Equal(t, 0, len(genesisState.Collections))

	ids := []string{id, id2, id3}
	idCollection := nft.NewIDCollection(denom, ids)
	idCollection2 := nft.NewIDCollection(denom2, ids)
	owner := nft.NewOwner(address, idCollection)

	owner2 := nft.NewOwner(address2, idCollection2)

	owners := []nft.Owner{owner, owner2}

	nft1 := nft.NewBaseNFT(id, address, tokenURI1)
	nft2 := nft.NewBaseNFT(id2, address, tokenURI1)
	nft3 := nft.NewBaseNFT(id3, address, tokenURI1)
	nfts := nft.NewNFTs(&nft1, &nft2, &nft3)
	collection := nft.NewCollection(denom, nfts)

	nftx := nft.NewBaseNFT(id, address2, tokenURI1)
	nft2x := nft.NewBaseNFT(id2, address2, tokenURI1)
	nft3x := nft.NewBaseNFT(id3, address2, tokenURI1)
	nftsx := nft.NewNFTs(&nftx, &nft2x, &nft3x)
	collection2 := nft.NewCollection(denom2, nftsx)

	collections := nft.NewCollections(collection, collection2)

	genesisState = nft.NewGenesisState(owners, collections)

	nft.InitGenesis(ctx, app.NFTKeeper, genesisState)

	returnedOwners := app.NFTKeeper.GetOwners(ctx)
	require.Equal(t, 2, len(owners))
	require.Equal(t, returnedOwners[0].String(), owners[0].String())
	require.Equal(t, returnedOwners[1].String(), owners[1].String())

	returnedCollections := app.NFTKeeper.GetCollections(ctx)
	require.Equal(t, 2, len(returnedCollections))
	require.Equal(t, returnedCollections[0].String(), collections[0].String())
	require.Equal(t, returnedCollections[1].String(), collections[1].String())

	exportedGenesisState := nft.ExportGenesis(ctx, app.NFTKeeper)
	require.Equal(t, len(genesisState.Owners), len(exportedGenesisState.Owners))
	require.Equal(t, genesisState.Owners[0].String(), exportedGenesisState.Owners[0].String())
	require.Equal(t, genesisState.Owners[1].String(), exportedGenesisState.Owners[1].String())

	require.Equal(t, len(genesisState.Collections), len(exportedGenesisState.Collections))
	require.Equal(t, genesisState.Collections[0].String(), exportedGenesisState.Collections[0].String())
	require.Equal(t, genesisState.Collections[1].String(), exportedGenesisState.Collections[1].String())
}

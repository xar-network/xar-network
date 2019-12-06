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

package market_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/market"
	"github.com/xar-network/xar-network/x/market/types"

	cstore "github.com/cosmos/cosmos-sdk/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/xar-network/xar-network/types/store"
)

func TestKeeperCoverage(t *testing.T) {

	cdc := makeTestCodec()

	logger := log.NewNopLogger() // Default
	//logger = log.NewTMLogger(os.Stdout) // Override to see output

	var (
		keyParams  = sdk.NewKVStoreKey(params.StoreKey)
		keyMarket  = sdk.NewKVStoreKey(market.StoreKey)
		tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	)

	db := dbm.NewMemDB()
	ms := cstore.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "xar-chain"}, true, logger)

	var (
		pk = params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
		mk = market.NewKeeper(keyMarket, cdc, pk.Subspace(market.DefaultParamspace), market.DefaultCodespace)
	)
	mk.SetParams(ctx, types.NewParams(market.DefaultGenesisState().Markets, []string{"cosmos1wdhk6e2wv9kk2j88d92"}))

	// Get market with ID 1
	market, err := mk.Get(ctx, store.NewEntityID(1))
	require.Nil(t, err)
	require.Equal(t, "1", market.ID.String())
	require.Equal(t, "uftm", market.BaseAssetDenom)

	// Get pair for market 1
	pair, err := mk.Pair(ctx, store.NewEntityID(1))
	require.Nil(t, err)
	require.Equal(t, "1", market.ID.String())
	require.Equal(t, "uftm/uzar", pair)

	// Create market as a nominee
	addr := sdk.AccAddress([]byte("someName"))
	msg := types.NewMsgCreateMarket(addr, "new1", "new2")
	mkt, err := mk.CreateMarket(ctx, msg.Nominee.String(), msg.BaseAsset, msg.QuoteAsset)
	require.Nil(t, err)
	require.Equal(t, mkt.BaseAssetDenom, msg.BaseAsset)

	// Create market as a nominee
	addr = sdk.AccAddress([]byte("someInvalidName"))
	msg = types.NewMsgCreateMarket(addr, "new1", "new2")
	mkt, err = mk.CreateMarket(ctx, msg.Nominee.String(), msg.BaseAsset, msg.QuoteAsset)
	assert.Error(t, err)
	require.Equal(t, "", mkt.BaseAssetDenom)
}

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()

	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return
}

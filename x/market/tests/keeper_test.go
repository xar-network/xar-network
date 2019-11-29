package tests

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/market"
	"github.com/xar-network/xar-network/x/market/types"

	cstore "github.com/cosmos/cosmos-sdk/store"
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
		mk = market.NewKeeper(keyMarket, cdc, pk.Subspace(market.DefaultParamspace))
	)
	mk.SetParams(ctx, types.NewParams(market.DefaultGenesisState().Markets))
	market, err := mk.Get(ctx, store.NewEntityID(1))
	require.Nil(t, err)
	require.Equal(t, "ueur", market.BaseAssetDenom)
	pair, err := mk.Pair(ctx, store.NewEntityID(1))
	require.Nil(t, err)
	require.Equal(t, "ueur/uzar", pair)
}

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()

	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return
}

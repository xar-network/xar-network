package tests

import (
	"testing"

	"github.com/xar-network/xar-network/x/record/params"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"

	"github.com/xar-network/xar-network/x/record"
	"github.com/xar-network/xar-network/x/record/msgs"
	"github.com/xar-network/xar-network/x/record/types"

	"github.com/xar-network/xar-network/x/record/keeper"
)

var (
	ReceiverCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("receiverCoins")))
	SenderAccAddr        sdk.AccAddress

	RecordParams = params.RecordParams{
		Hash:        "BC38CAEE32149BEF4CCFAEAB518EC9A5FBC85AE6AC8D5A9F6CD710FAF5E4A2B8",
		Name:        "testRecord",
		RecordType:  "image-hash",
		Description: "{}",
		Author:      "TEST",
		RecordNo:    "test-008"}

	RecordInfo = types.RecordInfo{
		Sender:      SenderAccAddr,
		Hash:        "BC38CAEE32149BEF4CCFAEAB518EC9A5FBC85AE6AC8D5A9F6CD710FAF5E4A2B8",
		Name:        "testRecord",
		RecordType:  "image-hash",
		Description: "{}",
		Author:      "TEST",
		RecordNo:    "test-008"}

	RecordQueryParams = params.RecordQueryParams{
		Limit:  30,
		Sender: SenderAccAddr}
)

// initialize the mock application for this module
func getMockApp(t *testing.T, genState record.GenesisState, genAccs []auth.Account) (
	mapp *mock.App, keeper keeper.Keeper, addrs []sdk.AccAddress,
	pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {
	mapp = mock.NewApp()
	msgs.RegisterCodec(mapp.Cdc)
	keyRecord := sdk.NewKVStoreKey(types.StoreKey)

	pk := mapp.ParamsKeeper

	keeper = record.NewKeeper(mapp.Cdc, keyRecord, pk, pk.Subspace("testrecord"), types.DefaultCodespace)

	mapp.Router().AddRoute(types.RouterKey, record.NewHandler(keeper))
	mapp.QueryRouter().AddRoute(types.QuerierRoute, record.NewQuerier(keeper))
	//mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper, genState))

	require.NoError(t, mapp.CompleteSetup(keyRecord))

	valTokens := sdk.TokensFromTendermintPower(1000000000000)
	if len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(2,
			sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens)))
	}
	SenderAccAddr = genAccs[0].GetAddress()

	RecordInfo.Sender = SenderAccAddr

	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, addrs, pubKeys, privKeys
}
func getInitChainer(mapp *mock.App, keeper keeper.Keeper, genState record.GenesisState) sdk.InitChainer {

	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {

		mapp.InitChainer(ctx, req)

		if genState.IsEmpty() {
			record.InitGenesis(ctx, keeper, record.DefaultGenesisState())
		} else {
			record.InitGenesis(ctx, keeper, genState)
		}
		return abci.ResponseInitChain{}
	}
}

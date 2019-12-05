package tests

import (
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"

	"github.com/xar-network/xar-network/x/record"
	"github.com/xar-network/xar-network/x/record/internal/keeper"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

var (
	ReceiverCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("receiverCoins")))
	SenderAccAddr        sdk.AccAddress

	RecordParams = types.RecordParams{
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

	RecordQueryParams = types.RecordQueryParams{
		Limit:  30,
		Sender: SenderAccAddr}
)

// initialize the mock application for this module
func getMockApp(t *testing.T, genState record.GenesisState, genAccs []exported.Account) (
	mapp *mock.App, k keeper.Keeper, addrs []sdk.AccAddress,
	pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {
	mapp = mock.NewApp()
	types.RegisterCodec(mapp.Cdc)
	keyRecord := sdk.NewKVStoreKey(types.StoreKey)

	pk := mapp.ParamsKeeper

	k = record.NewKeeper(mapp.Cdc, keyRecord, pk.Subspace("testrecord"), types.DefaultCodespace)

	mapp.Router().AddRoute(types.RouterKey, record.NewHandler(k))
	mapp.QueryRouter().AddRoute(types.QuerierRoute, keeper.NewQuerier(k))
	//mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, k, genState))

	require.NoError(t, mapp.CompleteSetup(keyRecord))
	valTokens := sdk.TokensFromConsensusPower(1000000000000)
	if len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(2,
			sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens)))
	}
	SenderAccAddr = genAccs[0].GetAddress()

	RecordInfo.Sender = SenderAccAddr

	mock.SetGenesis(mapp, genAccs)

	return mapp, k, addrs, pubKeys, privKeys
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

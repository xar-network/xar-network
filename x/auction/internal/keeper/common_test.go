package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/crypto"
	"github.com/xar-network/xar-network/x/auction"
	"github.com/xar-network/xar-network/x/auction/internal/keeper"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

func setUpMockApp() (*mock.App, keeper.Keeper, []sdk.AccAddress, []crypto.PrivKey) {
	// Create uninitialized mock app
	mapp := mock.NewApp()

	// Register codecs
	types.RegisterCodec(mapp.Cdc)

	// Create keepers
	keyAuction := sdk.NewKVStoreKey(auction.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	blacklistedAddrs := make(map[string]bool)
	bankKeeper := bank.NewBaseKeeper(mapp.AccountKeeper, mapp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)
	maccPerms := map[string][]string{}
	supplyKeeper := supply.NewKeeper(mapp.Cdc, keySupply, mapp.AccountKeeper, bankKeeper, maccPerms)
	auctionKeeper := keeper.NewKeeper(mapp.Cdc, supplyKeeper, keyAuction, mapp.ParamsKeeper.Subspace(types.DefaultParamspace))

	// Mount and load the stores
	err := mapp.CompleteSetup(keyAuction)
	if err != nil {
		panic("mock app setup failed")
	}

	// Create a bunch (ie 10) of pre-funded accounts to use for tests
	genAccs, addrs, _, privKeys := mock.CreateGenAccounts(10, sdk.NewCoins(sdk.NewInt64Coin("token1", 100), sdk.NewInt64Coin("token2", 100)))
	mock.SetGenesis(mapp, genAccs)

	return mapp, auctionKeeper, addrs, privKeys
}

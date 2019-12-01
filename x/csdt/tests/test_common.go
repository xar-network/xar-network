package csdt

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
	"github.com/xar-network/xar-network/x/oracle"
)

// Mock app is an ABCI app with an in memory database.
// This function creates an app, setting up the keepers, routes, begin and end blockers.
// But leaves it to the tests to call InitChain (done by calling mock.SetGenesis)
// The app works by submitting ABCI messages.
//  - InitChain sets up the app db from genesis.
//  - BeginBlock starts the delivery of a new block
//  - DeliverTx delivers a tx
//  - EndBlock signals the end of a block
//  - Commit ?
func setUpMockAppWithoutGenesis() (*mock.App, keeper.Keeper) {
	// Create uninitialized mock app
	mapp := mock.NewApp()

	// Register codecs
	types.RegisterCodec(mapp.Cdc)

	// Create keepers
	keyCSDT := sdk.NewKVStoreKey("csdt")
	keyPriceFeed := sdk.NewKVStoreKey(oracle.StoreKey)
	keySupply := sdk.NewKVStoreKey("supply")
	keyAccount := sdk.NewKVStoreKey("account")
	oracleKeeper := oracle.NewKeeper(keyPriceFeed, mapp.Cdc, oracle.DefaultCodespace)
	bankKeeper := bank.NewBaseKeeper(mapp.AccountKeeper, mapp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})
	accountKeeper := auth.NewAccountKeeper(mapp.Cdc, keyAccount, mapp.ParamsKeeper.Subspace("accountSubspace"), auth.ProtoBaseAccount)

	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter, supply.Burner},
	}

	supplyKeeper := supply.NewKeeper(mapp.Cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	csdtKeeper := keeper.NewKeeper(mapp.Cdc, keyCSDT, mapp.ParamsKeeper.Subspace("csdtSubspace"), oracleKeeper, bankKeeper, supplyKeeper)

	// Register routes
	mapp.Router().AddRoute("csdt", csdt.NewHandler(csdtKeeper))

	mapp.SetInitChainer(
		func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
			res := mapp.InitChainer(ctx, req)
			csdt.InitGenesis(ctx, csdtKeeper, csdt.DefaultGenesisState()) // Create a default genesis state, then set the keeper store to it
			return res
		},
	)

	// Mount and load the stores
	err := mapp.CompleteSetup(keyPriceFeed, keyCSDT)
	if err != nil {
		panic("mock app setup failed")
	}

	return mapp, csdtKeeper
}

// Avoid cluttering test cases with long function name
func i(in int64) sdk.Int                    { return sdk.NewInt(in) }
func d(str string) sdk.Dec                  { return sdk.MustNewDecFromStr(str) }
func c(denom string, amount int64) sdk.Coin { return sdk.NewInt64Coin(denom, amount) }
func cs(coins ...sdk.Coin) sdk.Coins        { return sdk.NewCoins(coins...) }

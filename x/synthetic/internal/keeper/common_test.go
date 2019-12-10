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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/crypto"
	"github.com/xar-network/xar-network/x/oracle"
	"github.com/xar-network/xar-network/x/synthetic/internal/keeper"
	"github.com/xar-network/xar-network/x/synthetic/internal/types"
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
func setUpMockAppWithoutGenesis() (*mock.App, keeper.Keeper, []sdk.AccAddress, []crypto.PrivKey) {
	// Create uninitialized mock app
	mapp := mock.NewApp()

	// Register codecs
	types.RegisterCodec(mapp.Cdc)
	supply.RegisterCodec(mapp.Cdc)

	// Create keepers
	keyCSDT := sdk.NewKVStoreKey(types.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter, supply.Burner},
	}

	oracleKeeper := oracle.NewKeeper(keyOracle, mapp.Cdc, mapp.ParamsKeeper.Subspace(oracle.DefaultParamspace), oracle.DefaultCodespace)
	bankKeeper := bank.NewBaseKeeper(mapp.AccountKeeper, mapp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})
	supplyKeeper := supply.NewKeeper(mapp.Cdc, keySupply, mapp.AccountKeeper, bankKeeper, maccPerms)
	csdtKeeper := keeper.NewKeeper(mapp.Cdc, keyCSDT, mapp.ParamsKeeper.Subspace(types.DefaultParamspace), oracleKeeper, bankKeeper, supplyKeeper)

	// Mount and load the stores
	err := mapp.CompleteSetup(keyOracle, keyCSDT, keySupply)
	if err != nil {
		panic("mock app setup failed")
	}

	// Create a bunch (ie 10) of pre-funded accounts to use for tests
	genAccs, addrs, _, privKeys := mock.CreateGenAccounts(10, sdk.NewCoins(sdk.NewInt64Coin("token1", 100), sdk.NewInt64Coin("token2", 100)))
	mock.SetGenesis(mapp, genAccs)

	return mapp, csdtKeeper, addrs, privKeys
}

// Avoid cluttering test cases with long function name
func i(in int64) sdk.Int                    { return sdk.NewInt(in) }
func d(str string) sdk.Dec                  { return sdk.MustNewDecFromStr(str) }
func c(denom string, amount int64) sdk.Coin { return sdk.NewInt64Coin(denom, amount) }
func cs(coins ...sdk.Coin) sdk.Coins        { return sdk.NewCoins(coins...) }

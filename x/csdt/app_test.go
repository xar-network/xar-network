/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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

package csdt_test

import (
	"log"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
	"github.com/xar-network/xar-network/x/oracle"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
)

const chainID = ""

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) auth.StdTx {
	// Make the transaction free
	fee := auth.StdFee{
		Amount: sdk.NewCoins(sdk.NewInt64Coin("foocoin", 0)),
		Gas:    200000,
	}

	sigs := make([]auth.StdSignature, len(priv))
	memo := "testmemotestmemo"

	for i, p := range priv {
		sig, err := p.Sign(auth.StdSignBytes(chainID, accnums[i], seq[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:    p.PubKey(),
			Signature: sig,
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, memo)
}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, cdc *codec.Codec, app *baseapp.BaseApp, header abci.Header, msgs []sdk.Msg,
	accNums, seq []uint64, expSimPass, expPass bool, priv ...crypto.PrivKey,
) sdk.Result {

	tx := GenTx(msgs, accNums, seq, priv...)

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	require.Nil(t, err)

	// Must simulate now as CheckTx doesn't run Msgs anymore
	res := app.Simulate(txBytes, tx)

	if expSimPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	res = app.Deliver(tx)

	if expPass {
		require.Equal(t, sdk.CodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
	}

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	return res
}

func TestApp_CreateModifyDeleteCSDT(t *testing.T) {
	// Setup
	mapp, keeper := setUpMockAppWithoutGenesis()
	genAccs, addrs, _, privKeys := mock.CreateGenAccounts(1, cs(c("uftm", 100)))
	testAddr := addrs[0]
	testPrivKey := privKeys[0]
	mock.SetGenesis(mapp, genAccs)
	// setup oracle, TODO can this be shortened a bit?
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	keeper.SetParams(ctx, types.DefaultParams())
	keeper.GetSupply().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	oracleParams := oracle.DefaultParams()
	oracleParams.Assets = oracle.Assets{
		oracle.Asset{
			AssetCode:  "uftm",
			BaseAsset:  "uftm",
			QuoteAsset: csdt.StableDenom,
			Oracles: oracle.Oracles{
				oracle.Oracle{
					Address: addrs[0],
				},
			},
		},
	}
	oracleParams.Nominees = []string{addrs[0].String()}

	keeper.GetOracle().SetParams(ctx, oracleParams)
	_, _ = keeper.GetOracle().SetPrice(
		ctx, addrs[0], "uftm",
		sdk.MustNewDecFromStr("1.00"),
		time.Now().Add(time.Hour*1))
	_ = keeper.GetOracle().SetCurrentPrices(ctx)
	mapp.EndBlock(abci.RequestEndBlock{})
	mapp.Commit()

	currP := keeper.GetOracle().GetCurrentPrice(ctx, "uftm")
	log.Println(currP)

	// Create CSDT
	msgs := []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(10), "uftm", i(5))}
	SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{0}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(types.StableDenom, 5), c("uftm", 90)))

	// Modify CSDT
	msgs = []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(40), "uftm", i(5))}
	SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{1}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(types.StableDenom, 10), c("uftm", 50)))

	// Delete CSDT
	msgs = []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(-50), "uftm", i(-10))}
	SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{2}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c("uftm", 100)))

	// deposit
	msgs = []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(10), "uftm", i(5))}
	SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{3}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(types.StableDenom, 5), c("uftm", 90)))

	msgs = []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(0), "uftm", i(1))}
	SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{4}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(types.StableDenom, 6), c("uftm", 90)))

	addr := mapp.AccountKeeper.GetAccount(ctx, testAddr)
	log.Println(addr)
}

func TestApp_ParamExport(t *testing.T) {
	// Setup
	mapp, keeper := setUpMockAppWithoutGenesis()
	genAccs, _, _, _ := mock.CreateGenAccounts(1, cs(c("uftm", 100)))
	mock.SetGenesis(mapp, genAccs)
	// setup oracle, TODO can this be shortened a bit?
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	keeper.SetParams(ctx, types.DefaultParams())

	genState := csdt.ExportGenesis(ctx, keeper)
	require.Equal(t, 1, len(genState.Params.CollateralParams))
}

// Avoid cluttering test cases with long function name
func i(in int64) sdk.Int                    { return sdk.NewInt(in) }
func d(str string) sdk.Dec                  { return sdk.MustNewDecFromStr(str) }
func c(denom string, amount int64) sdk.Coin { return sdk.NewInt64Coin(denom, amount) }
func cs(coins ...sdk.Coin) sdk.Coins        { return sdk.NewCoins(coins...) }

func setUpMockAppWithoutGenesis() (*mock.App, csdt.Keeper) {
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
	csdtKeeper := csdt.NewKeeper(mapp.Cdc, keyCSDT, mapp.ParamsKeeper.Subspace(types.DefaultParamspace), oracleKeeper, bankKeeper, supplyKeeper, csdt.ModuleName)

	// Register routes
	mapp.Router().AddRoute("csdt", csdt.NewHandler(csdtKeeper))
	// Mount and load the stores
	err := mapp.CompleteSetup(keyOracle, keyCSDT, keySupply)
	if err != nil {
		panic("mock app setup failed")
	}

	return mapp, csdtKeeper
}

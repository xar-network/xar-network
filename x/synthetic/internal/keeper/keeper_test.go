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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/oracle"
	"github.com/xar-network/xar-network/x/synthetic/internal/types"
)

// How could one reduce the number of params in the test cases. Create a table driven test for each of the 4 add/withdraw collateral/debt?

func TestKeeper_ModifyCSDT(t *testing.T) {
	_, addrs := mock.GeneratePrivKeyAddressPairs(2)
	ownerAddr := addrs[0]

	// setup keeper
	mapp, keeper, _, _ := setUpMockAppWithoutGenesis()
	// initialize csdt owner account with coins
	genAcc := auth.BaseAccount{
		Address: ownerAddr,
		Coins:   sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(1000))),
	}

	mock.SetGenesis(mapp, []exported.Account{&genAcc})
	// create a new context
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	// setup store state
	oracleParams := oracle.DefaultParams()
	oracleParams.Assets = oracle.Assets{
		oracle.Asset{
			AssetCode:  "sbtc",
			BaseAsset:  "sbtc",
			QuoteAsset: types.StableDenom,
			Oracles: oracle.Oracles{
				oracle.Oracle{
					Address: addrs[1],
				},
			},
		},
	}
	oracleParams.Nominees = []string{addrs[1].String()}

	keeper.GetOracle().SetParams(ctx, oracleParams)
	_, _ = keeper.GetOracle().SetPrice(
		ctx, addrs[1], "sbtc",
		sdk.MustNewDecFromStr("1.0"),
		time.Now().Add(time.Hour*1))
	_ = keeper.GetOracle().SetCurrentPrices(ctx)

	keeper.GetSupply().SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(1000)))))

	// call func under test
	keeper.SetParams(ctx, types.DefaultParams())
	err := keeper.BuySynthetic(ctx, ownerAddr, sdk.NewCoin("sbtc", sdk.NewInt(1000)))
	mapp.EndBlock(abci.RequestEndBlock{})
	mapp.Commit()

	require.NoError(t, err)
	err = keeper.SellSynthetic(ctx, ownerAddr, sdk.NewCoin("sbtc", sdk.NewInt(1000)))

	require.NoError(t, err)
}

func TestMarketBalance(t *testing.T) {
	testRatioChange(t)
	testAddFee(t)
}

func testRatioChange(t *testing.T) {
	mb := types.EmptyMarketBalance("asd")

	mb.IncreaseLongVolume(sdk.NewInt(2000))
	require.Equal(t, mb.Imbalance.Ratio, float64(0))

	mb.IncreaseShortVolume(sdk.NewInt(1000))
	require.Equal(t, mb.Imbalance.Ratio, float64(1)) // 100% diff
	mb.Flash()

	mb.IncreaseLongVolume(sdk.NewInt(1500))
	mb.IncreaseShortVolume(sdk.NewInt(1000))
	require.Equal(t, mb.Imbalance.Ratio, float64(0.5)) // 50% diff
}

func testAddFee(t *testing.T) {
	mb := types.EmptyMarketBalance("asd")

	mb.IncreaseLongVolume(sdk.NewInt(2000))
	mb.IncreaseShortVolume(sdk.NewInt(1000))

	testAmt := sdk.NewInt(100)
	val := mb.AddFee(testAmt)
	assumedVal := sdk.NewInt(105) // 100% of an imbalance should lead to a 5% of a fee
	require.True(t, val.Equal(assumedVal))
	mb.Flash()
	testFlash(t, &mb)

	mb.IncreaseLongVolume(sdk.NewInt(1500))
	mb.IncreaseShortVolume(sdk.NewInt(1000))
	testAmt = sdk.NewInt(100)
	val = mb.AddFee(testAmt)
	assumedVal = sdk.NewInt(101) // 50% of an imbalance should lead to a 1% of a fee
	require.True(t, val.Equal(assumedVal))
}

func testFlash(t *testing.T, mb *types.MarketBalance) {
	require.True(t, mb.LongVolume.Equal(sdk.ZeroInt()))
	require.True(t, mb.ShortVolume.Equal(sdk.ZeroInt()))
	require.Equal(t, mb.Imbalance.Ratio, float64(0))
}

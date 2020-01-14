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
		Coins:   sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(1005))),
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

func startNewBlock(app *mock.App) {
	header := newBlockHeader(app)
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
}

func endBlock(app *mock.App) {
	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()
}

func newBlockHeader(app *mock.App) abci.Header {
	//ctx.BlockTime()
	return abci.Header{Height: app.LastBlockHeight() + 1, Time:time.Now()}
}

func TestMarketBalance(t *testing.T) {
	testRatioChange(t)
	testAddFee(t)
	testGetFee(t)
}

func testRatioChange(t *testing.T) {
	mb := types.EmptyMarketBalance("asd", 0)

	mb.IncreaseLongVolume(sdk.NewInt(2000))
	require.True(t, mb.Imbalance.Ratio.Equal(sdk.MustNewDecFromStr("0")))

	mb.IncreaseShortVolume(sdk.NewInt(1000))
	require.True(t, mb.Imbalance.Ratio.Equal(sdk.MustNewDecFromStr("1"))) // 100% diff
	mb.FlashVolumes()

	mb.IncreaseLongVolume(sdk.NewInt(1500))
	mb.IncreaseShortVolume(sdk.NewInt(1000))
	require.True(t, mb.Imbalance.Ratio.Equal(sdk.MustNewDecFromStr("0.5"))) // 50% diff
}

func testAddFee(t *testing.T) {
	mb := types.EmptyMarketBalance("asd", 0)

	mb.IncreaseLongVolume(sdk.NewInt(2000))
	mb.IncreaseShortVolume(sdk.NewInt(1000))

	testAmt := sdk.NewInt(100)
	val := mb.AddFee(testAmt)
	assumedVal := sdk.NewInt(105) // 100% of an imbalance should lead to a 5% of a fee
	require.True(t, val.Equal(assumedVal))
	mb.FlashVolumes()
	testFlash(t, &mb)

	mb.IncreaseLongVolume(sdk.NewInt(1500))
	mb.IncreaseShortVolume(sdk.NewInt(1000))
	testAmt = sdk.NewInt(100)
	val = mb.AddFee(testAmt)
	assumedVal = sdk.NewInt(101) // 50% of an imbalance should lead to a 1% of a fee
	require.True(t, val.Equal(assumedVal))
}

func testGetFee(t *testing.T) {
	mb := types.EmptyMarketBalance("asd", 0)

	mb.IncreaseLongVolume(sdk.NewInt(2000))
	mb.IncreaseShortVolume(sdk.NewInt(1000))

	testAmt := sdk.NewInt(100)
	val := mb.GetFeeForAmount(testAmt)
	assumedVal := sdk.NewInt(5) // 100% of an imbalance should lead to a 5% of a fee
	require.True(t, val.Equal(assumedVal))
	mb.FlashVolumes()
	testFlash(t, &mb)

	mb.IncreaseLongVolume(sdk.NewInt(1500))
	mb.IncreaseShortVolume(sdk.NewInt(1000))
	testAmt = sdk.NewInt(100)
	val = mb.GetFeeForAmount(testAmt)
	assumedVal = sdk.NewInt(1) // 50% of an imbalance should lead to a 1% of a fee
	require.True(t, val.Equal(assumedVal))
}

func testFlash(t *testing.T, mb *types.MarketBalance) {
	require.True(t, mb.LongVolume.Equal(sdk.ZeroInt()))
	require.True(t, mb.ShortVolume.Equal(sdk.ZeroInt()))
	require.True(t, mb.Imbalance.Ratio.Equal(sdk.ZeroDec()))
}

func TestSnapshots(t *testing.T) {
	testBasicSnapshotFuncs(t)
	testFeesWithSnapshot(t)
}

func TestEndBlock(t *testing.T) {
	testOnEndBlockCounter(t)
}

func testOnEndBlockCounter(t *testing.T) {
	snaps := types.NewVolumeSnapshots(2, nil)
	mb := types.NewMarketBalance("asd", snaps, 1, time.Duration(0))

	newLong := sdk.NewInt(2000)
	newShort := sdk.NewInt(1000)
	mb.IncreaseLongVolume(newLong)
	mb.IncreaseShortVolume(newShort)
	mb.OnEndBlock(abci.Header{})
	mb.OnEndBlock(abci.Header{})

	require.True(t, len(mb.VolumeSnapshots.Snapshots) == 2)
	require.True(t, mb.VolumeSnapshots.Snapshots[0].LongVolume.Equal(newLong))
	require.True(t, mb.VolumeSnapshots.Snapshots[0].ShortVolume.Equal(newShort))
	require.True(t, mb.VolumeSnapshots.Snapshots[1].LongVolume.Equal(sdk.ZeroInt()))
	require.True(t, mb.VolumeSnapshots.Snapshots[1].ShortVolume.Equal(sdk.ZeroInt()))
}

func TestOnEndBlockTimer(t *testing.T) {

	snaps := types.NewVolumeSnapshots(2, nil)
	tn := time.Now()
	mb := types.NewMarketBalance("asd", snaps, 1, time.Hour)

	newLong := sdk.NewInt(2000)
	newShort := sdk.NewInt(1000)
	mb.IncreaseLongVolume(newLong)
	mb.IncreaseShortVolume(newShort)
	mb.OnEndBlock(abci.Header{Time: tn.Add(time.Hour)})
	mb.OnEndBlock(abci.Header{Time: tn.Add(time.Hour).Add(time.Hour)})

	require.True(t, len(mb.VolumeSnapshots.Snapshots) == 2)
	require.True(t, mb.VolumeSnapshots.Snapshots[0].LongVolume.Equal(newLong))
	require.True(t, mb.VolumeSnapshots.Snapshots[0].ShortVolume.Equal(newShort))
	require.True(t, mb.VolumeSnapshots.Snapshots[1].LongVolume.Equal(sdk.ZeroInt()))
	require.True(t, mb.VolumeSnapshots.Snapshots[1].ShortVolume.Equal(sdk.ZeroInt()))
}

func testFeesWithSnapshot(t *testing.T) {
	mb := getMarketBalanceAndSnapshots()
	wsnap := mb.VolumeSnapshots.GetWeightedVolumes()
	assumedRatio := wsnap.LongVolume.Quo(wsnap.ShortVolume).Sub(sdk.OneInt())
	require.True(t, assumedRatio.Equal(mb.Imbalance.Ratio.TruncateInt()))

	testAmt := sdk.NewInt(100)
	val := mb.GetFeeForAmount(testAmt)
	assumedVal := sdk.NewInt(5) // 100% of an imbalance should lead to a 5% of a fee
	require.True(t, val.Equal(assumedVal))
}

func testBasicSnapshotFuncs(t *testing.T) {
	wsnap := getDefaultSnapshots().GetWeightedVolumes()
	assumedShortVal := sdk.NewInt(551) // 551 = 100 + 90 + 80 + 70 + 60 + 50 + 40 + 30 + 20 + 10 + 1
	assumedLongVal := sdk.NewInt(1102)
	require.True(t, wsnap.LongVolume.Equal(assumedLongVal))
	require.True(t, wsnap.ShortVolume.Equal(assumedShortVal))
}

func getMarketBalanceAndSnapshots() *types.MarketBalance {
	mb := types.EmptyMarketBalance("asd", 0)
	mb.VolumeSnapshots = *getDefaultSnapshots()
	mb.Recalculate()
	return &mb
}

func getDefaultSnapshots() *types.VolumeSnapshots {
	coeffs := []sdk.Int{
		sdk.NewInt(100),
		sdk.NewInt(90),
		sdk.NewInt(80),
		sdk.NewInt(70),
		sdk.NewInt(60),
		sdk.NewInt(50),
		sdk.NewInt(40),
		sdk.NewInt(30),
		sdk.NewInt(20),
		sdk.NewInt(10),
	}
	two := sdk.NewInt(2)
	one := sdk.OneInt()
	snap := types.NewVolumeSnapshots(len(coeffs)+1, coeffs)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	snap.AddSnapshotValues(two, one)
	return &snap
}

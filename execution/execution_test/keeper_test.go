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

package execution_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/testutil"
	"github.com/xar-network/xar-network/testutil/mockapp"
	"github.com/xar-network/xar-network/testutil/testflags"
	uexstore "github.com/xar-network/xar-network/types/store"
	"github.com/xar-network/xar-network/x/denominations"
	types2 "github.com/xar-network/xar-network/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

func TestKeeper_ExecuteAndCancelExpired(t *testing.T) {
	testflags.UnitTest(t)
	app := mockapp.New(t)
	nominee := testutil.RandAddr()
	buyer := testutil.RandAddr()
	seller := testutil.RandAddr()

	app.SupplyKeeper.SetSupply(app.Ctx, supply.NewSupply(sdk.Coins{}))
	marketParams := app.MarketKeeper.GetParams(app.Ctx)
	marketParams.Nominees = []string{nominee.String()}
	app.MarketKeeper.SetParams(app.Ctx, marketParams)

	err := app.SupplyKeeper.MintCoins(app.Ctx, denominations.ModuleName, sdk.NewCoins(sdk.NewCoin("tst1", sdk.NewInt(1000000000000)), sdk.NewCoin("tst2", sdk.NewInt(1000000000000))))
	require.NoError(t, err)
	require.NoError(t, app.SupplyKeeper.SendCoinsFromModuleToAccount(app.Ctx, denominations.ModuleName, buyer, sdk.NewCoins(sdk.NewCoin("tst1", sdk.NewInt(10000000000)))))
	require.NoError(t, app.SupplyKeeper.SendCoinsFromModuleToAccount(app.Ctx, denominations.ModuleName, buyer, sdk.NewCoins(sdk.NewCoin("tst2", sdk.NewInt(10000000000)))))
	require.NoError(t, app.SupplyKeeper.SendCoinsFromModuleToAccount(app.Ctx, denominations.ModuleName, seller, sdk.NewCoins(sdk.NewCoin("tst1", sdk.NewInt(10000000000)))))
	require.NoError(t, app.SupplyKeeper.SendCoinsFromModuleToAccount(app.Ctx, denominations.ModuleName, seller, sdk.NewCoins(sdk.NewCoin("tst2", sdk.NewInt(10000000000)))))
	market := types2.NewMsgCreateMarket(nominee, "tst1", "tst2")
	mkt, err := app.MarketKeeper.CreateMarket(app.Ctx, market.Nominee.String(), market.BaseAsset, market.QuoteAsset)
	require.NoError(t, err)

	_, err = app.OrderKeeper.Post(app.Ctx, buyer, mkt.ID, matcheng.Bid, testutil.ToBaseUnits(1), testutil.ToBaseUnits(1), 100)
	require.NoError(t, err)

	ctx := app.Ctx.WithBlockHeight(602)
	bids := [][2]uint64{
		{1, 10},
		{2, 10},
		{3, 10},
	}
	asks := [][2]uint64{
		{2, 10},
		{3, 10},
		{4, 10},
	}
	for _, bid := range bids {
		_, err = app.OrderKeeper.Post(ctx, buyer, mkt.ID, matcheng.Bid, testutil.ToBaseUnits(bid[0]), testutil.ToBaseUnits(bid[1]), 100)
		require.NoError(t, err)
	}
	for _, ask := range asks {
		_, err = app.OrderKeeper.Post(ctx, seller, mkt.ID, matcheng.Ask, testutil.ToBaseUnits(ask[0]), testutil.ToBaseUnits(ask[1]), 100)
		require.NoError(t, err)
	}
	require.NoError(t, app.ExecutionKeeper.ExecuteAndCancelExpired(ctx))
	t.Run("should expire orders out of TIF", func(t *testing.T) {
		assert.False(t, app.OrderKeeper.Has(ctx, uexstore.NewEntityID(1)))
	})
	t.Run("should update quantities of partially filled orders", func(t *testing.T) {
		ord3, err := app.OrderKeeper.Get(ctx, uexstore.NewEntityID(3))
		require.NoError(t, err)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(5), ord3.Quantity)
		ord4, err := app.OrderKeeper.Get(ctx, uexstore.NewEntityID(4))
		require.NoError(t, err)
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(5), ord4.Quantity)
	})

	// perform next round of cancellation after since orders are
	// deleted on cancellation
	ctx = app.Ctx.WithBlockHeight(703)
	require.NoError(t, app.ExecutionKeeper.ExecuteAndCancelExpired(ctx))

	t.Run("should delete completely filled orders", func(t *testing.T) {
		assert.False(t, app.OrderKeeper.Has(ctx, uexstore.NewEntityID(5)))
	})
	t.Run("all executed orders should exchange coins", func(t *testing.T) {
		// seller should have 9990 asset 1, because two orders were
		// partially executed (for 5 each), then expired.

		sellerAsset1Bal := sdk.NewUintFromBigInt(app.BankKeeper.GetCoins(ctx, seller).AmountOf("tst1").BigInt())
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(90), sellerAsset1Bal)
		// 10020 because two orders executed at clearing price 2:
		// 10000 + 5 * 2 + 5 * 2 = 10020
		sellerAsset2Bal := sdk.NewUintFromBigInt(app.BankKeeper.GetCoins(ctx, seller).AmountOf("tst2").BigInt())
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(120), sellerAsset2Bal)

		buyerAsset1Bal := sdk.NewUintFromBigInt(app.BankKeeper.GetCoins(ctx, buyer).AmountOf("tst1").BigInt())
		// the orders with prices 1 and 2 receives partial fills of 5
		// the other orders expired.
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(110), buyerAsset1Bal)

		buyerAsset2Bal := sdk.NewUintFromBigInt(app.BankKeeper.GetCoins(ctx, buyer).AmountOf("tst2").BigInt())
		// clearing of 2. two of buyer's orders were rationed for a total of 10
		// asset 2 credited.
		testutil.AssertEqualUints(t, testutil.ToBaseUnits(80), buyerAsset2Bal)
	})
}

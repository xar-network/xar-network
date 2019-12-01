package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"

	"github.com/xar-network/xar-network/x/csdt"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
	"github.com/xar-network/xar-network/x/oracle"
)

func TestKeeper_SeizeAndStartCollateralAuction(t *testing.T) {
	// Setup
	ctx, k := setupTestKeepers()

	_, addrs := mock.GeneratePrivKeyAddressPairs(1)

	oracle.InitGenesis(ctx, k.oracleKeeper, oracleGenesis(addrs[0]))
	_, err := k.oracleKeeper.SetPrice(ctx, addrs[0], "btc", sdk.MustNewDecFromStr("8000.00"), time.Now().Add(time.Hour*1))
	require.NoError(t, err)

	k.oracleKeeper.SetCurrentPrices(ctx)
	csdt.InitGenesis(ctx, k.csdtKeeper, csdtDefaultGenesis())

	dp := defaultParams()
	k.liquidatorKeeper.SetParams(ctx, dp)
	_, err = k.bankKeeper.AddCoins(ctx, addrs[0], cs(c("btc", 100)))
	require.NoError(t, err)

	err = k.csdtKeeper.ModifyCSDT(ctx, addrs[0], "btc", i(3), i(16000))
	require.NoError(t, err)

	_, err = k.oracleKeeper.SetPrice(ctx, addrs[0], "btc", sdk.MustNewDecFromStr("7999.99"), time.Now().Add(time.Hour*1))
	require.NoError(t, err)
	k.oracleKeeper.SetCurrentPrices(ctx)

	addr, perms := k.supplyKeeper.GetModuleAddressAndPermissions(types.ModuleName)
	require.Equal(t, "cosmos1eu2ta269haf6j6z3lsj79a8rq3hsmnhuxj34g9", addr.String())
	require.Equal(t, 1, len(perms))

	// Run test function
	csdt, found := k.csdtKeeper.GetCSDT(ctx, addrs[0], "btc")
	require.NoError(t, err)
	auctionID, err := k.liquidatorKeeper.SeizeAndStartCollateralAuction(ctx, addrs[0], "btc")

	// Check CDP
	require.NoError(t, err)
	csdt, found = k.csdtKeeper.GetCSDT(ctx, addrs[0], "btc")

	require.True(t, found)
	require.Equal(t, csdt.CollateralAmount, cs(c("btc", 2)))                 // original amount - params.CollateralAuctionSize
	require.Equal(t, csdt.Debt, cs(c(k.csdtKeeper.GetStableDenom(), 10667))) // original debt scaled by amount of collateral removed
	// Check auction exists
	_, found = k.auctionKeeper.GetAuction(ctx, auctionID)
	require.True(t, found)
	// TODO check auction values are correct?
}

func TestKeeper_StartDebtAuction(t *testing.T) {
	// Setup
	ctx, k := setupTestKeepers()
	k.liquidatorKeeper.SetParams(ctx, defaultParams())
	initSDebt := types.SeizedDebt{i(2000), i(0)}
	k.liquidatorKeeper.SetSeizedDebt(ctx, initSDebt)

	// Execute
	auctionID, err := k.liquidatorKeeper.StartDebtAuction(ctx)

	// Check
	require.NoError(t, err)
	require.Equal(t,
		types.SeizedDebt{
			initSDebt.Total,
			initSDebt.SentToAuction.Add(k.liquidatorKeeper.GetParams(ctx).DebtAuctionSize),
		},
		k.liquidatorKeeper.GetSeizedDebt(ctx),
	)
	_, found := k.auctionKeeper.GetAuction(ctx, auctionID)
	require.True(t, found)
	// TODO check auction values are correct?
}

// func TestKeeper_StartSurplusAuction(t *testing.T) {
// 	// Setup
// 	ctx, k := setupTestKeepers()
// 	initSurplus := i(2000)
// 	k.liquidatorKeeper.bankKeeper.AddCoins(ctx, k.csdtKeeper.GetLiquidatorAccountAddress(), cs(sdk.NewCoin(k.csdtKeeper.GetStableDenom(), initSurplus)))
// 	k.liquidatorKeeper.setSeizedDebt(ctx, i(0))

// 	// Execute
// 	auctionID, err := k.liquidatorKeeper.StartSurplusAuction(ctx)

// 	// Check
// 	require.NoError(t, err)
// 	require.Equal(t,
// 		initSurplus.Sub(SurplusAuctionSize),
// 		k.liquidatorKeeper.bankKeeper.GetCoins(ctx,
// 			k.csdtKeeper.GetLiquidatorAccountAddress(),
// 		).AmountOf(k.csdtKeeper.GetStableDenom()),
// 	)
// 	_, found := k.auctionKeeper.GetAuction(ctx, auctionID)
// 	require.True(t, found)
// }

func TestKeeper_partialSeizeCSDT(t *testing.T) {
	// Setup
	ctx, k := setupTestKeepers()

	_, addrs := mock.GeneratePrivKeyAddressPairs(1)

	oracle.InitGenesis(ctx, k.oracleKeeper, oracleGenesis(addrs[0]))

	k.oracleKeeper.SetPrice(ctx, addrs[0], "btc", sdk.MustNewDecFromStr("8000.00"), time.Now().Add(time.Hour*1))
	k.oracleKeeper.SetCurrentPrices(ctx)
	k.bankKeeper.AddCoins(ctx, addrs[0], cs(c("btc", 100)))
	csdt.InitGenesis(ctx, k.csdtKeeper, csdtDefaultGenesis())
	k.liquidatorKeeper.SetParams(ctx, defaultParams())

	k.csdtKeeper.ModifyCSDT(ctx, addrs[0], "btc", i(3), i(16000))

	k.oracleKeeper.SetPrice(ctx, addrs[0], "btc", sdk.MustNewDecFromStr("7999.99"), time.Now().Add(time.Hour*1))
	k.oracleKeeper.SetCurrentPrices(ctx)

	// Run test function
	err := k.liquidatorKeeper.PartialSeizeCSDT(ctx, addrs[0], "btc", i(2), i(10000))

	// Check
	require.NoError(t, err)
	csdt, found := k.csdtKeeper.GetCSDT(ctx, addrs[0], "btc")
	require.True(t, found)
	require.Equal(t, cs(c(csdt.CollateralDenom, 1)), csdt.CollateralAmount)
	require.Equal(t, cs(c(k.csdtKeeper.GetStableDenom(), 6000)), csdt.Debt)
}

func TestKeeper_GetSetSeizedDebt(t *testing.T) {
	// Setup
	ctx, k := setupTestKeepers()
	debt := types.SeizedDebt{i(234247645), i(2343)}

	// Run test function
	k.liquidatorKeeper.SetSeizedDebt(ctx, debt)
	readDebt := k.liquidatorKeeper.GetSeizedDebt(ctx)

	// Check
	require.Equal(t, debt, readDebt)
}

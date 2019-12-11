package uniswap

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
	"log"
	"testing"
	"time"
)

const nonNativeDenomTest = "asd"

func TestSwap(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewInt(0), 1)
	tm, err := time.Parse("2006-01-02T15:04:05.000Z", "2022-04-23T18:25:43.511Z")
	if err != nil {
		require.NoError(t, err)
	}

	oneCoin := sdk.NewInt(1)
	testCoinAmt := sdk.NewInt(14)
	nativeDenom := keeper.GetNativeDenom(ctx)
	testDenom := nonNativeDenomTest
	nativeCoinAmt := sdk.NewInt(10040)
	nonNativeCoinAmt := sdk.NewInt(151000)

	err = addLiquidityForTest(ctx, keeper, accs, nativeCoinAmt, nonNativeCoinAmt, testDenom)
	if err != nil {
		require.NoError(t, err)
	}

	TestCoin1 := sdk.NewCoin(nativeDenom, oneCoin)
	TestCoin2 := sdk.NewCoin(testDenom, testCoinAmt)

	userCoin1 := sdk.NewCoin(nativeDenom, nativeCoinAmt)
	userCoin2 := sdk.NewCoin(testDenom, nonNativeCoinAmt)

	msg := MsgSwapOrder{
		TestCoin1,
		TestCoin2,
		tm,
		accs[0].GetAddress(),
		accs[0].GetAddress(),
		false,
	}

	err = keeper.RecieveCoins(ctx, accs[0].GetAddress(), userCoin2, userCoin1)
	if err != nil {
		require.NoError(t, err)
	}

	moduleName := keeper.MustGetModuleName(nativeDenom, testDenom)

	rp1, found := keeper.GetReservePool(ctx, moduleName)
	if !found {
		require.NoError(t, err)
	}

	HandleMsgSwapOrder(ctx, msg, keeper)

	rp2, found := keeper.GetReservePool(ctx, moduleName)
	if !found {
		require.NoError(t, err)
	}

	expectedNativeDenomAmt := rp1.AmountOf(nativeDenom).Add(oneCoin)
	require.True(t, expectedNativeDenomAmt.Equal(rp2.AmountOf(nativeDenom)))

	expectedNonNativeDenomAmt := rp1.AmountOf(testDenom).Sub(testCoinAmt)
	require.True(t, expectedNonNativeDenomAmt.Equal(rp2.AmountOf(testDenom)))

	log.Println()
}

func TestDoubleSwap(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewInt(0), 1)
	tm, err := time.Parse("2006-01-02T15:04:05.000Z", "2022-04-23T18:25:43.511Z")
	if err != nil {
		require.NoError(t, err)
	}

	oneCoin := sdk.NewInt(1)
	testCoinAmt1 := sdk.NewInt(16)
	testCoinAmt2 := sdk.NewInt(10)
	nativeDenom := keeper.GetNativeDenom(ctx)
	//oneHundredCoin := sdk.NewInt(100)
	testDenom := nonNativeDenomTest
	testDenom2 := nonNativeDenomTest + "2"

	//outputTestCoin1 := sdk.NewCoin(testDenom, oneHundredCoin)
	//outputTestCoin2 := sdk.NewCoin(testDenom2, oneHundredCoin)

	nativeCoinAmt := sdk.NewInt(10040)
	nonNativeCoinAmt := sdk.NewInt(151000)

	err = addLiquidityForTest(ctx, keeper, accs, nativeCoinAmt, nonNativeCoinAmt, testDenom)
	if err != nil {
		require.NoError(t, err)
	}

	err = addLiquidityForTest(ctx, keeper, accs, nativeCoinAmt, nonNativeCoinAmt, testDenom2)
	if err != nil {
		require.NoError(t, err)
	}

	TestCoin1 := sdk.NewCoin(testDenom, testCoinAmt1)
	TestCoin2 := sdk.NewCoin(testDenom2, testCoinAmt2)

	userCoin1 := sdk.NewCoin(nativeDenom, nativeCoinAmt.Mul(sdk.NewInt(2)))
	userCoin2 := sdk.NewCoin(testDenom, nonNativeCoinAmt)
	userCoin3 := sdk.NewCoin(testDenom2, nonNativeCoinAmt)

	msg := MsgSwapOrder{
		TestCoin1,
		TestCoin2,
		tm,
		accs[0].GetAddress(),
		accs[0].GetAddress(),
		false,
	}

	coins := sdk.Coins{userCoin2, userCoin1, userCoin3}.Sort()
	err = keeper.RecieveCoins(ctx, accs[0].GetAddress(), coins...)
	if err != nil {
		require.NoError(t, err)
	}

	moduleName1 := keeper.MustGetModuleName(nativeDenom, testDenom)
	moduleName2 := keeper.MustGetModuleName(nativeDenom, testDenom2)

	rpA1, found := keeper.GetReservePool(ctx, moduleName1)
	if !found {
		require.NoError(t, err)
	}

	rpB1, found := keeper.GetReservePool(ctx, moduleName2)
	if !found {
		require.NoError(t, err)
	}

	HandleMsgSwapOrder(ctx, msg, keeper)

	rpA2, found := keeper.GetReservePool(ctx, moduleName1)
	if !found {
		require.NoError(t, err)
	}

	rpB2, found := keeper.GetReservePool(ctx, moduleName2)
	if !found {
		require.NoError(t, err)
	}

	expectedNativeDenomAmt := rpA1.AmountOf(nativeDenom).Sub(oneCoin)
	expectedNonNativeDenomAmt := rpA1.AmountOf(testDenom).Add(testCoinAmt1)
	require.True(t, expectedNativeDenomAmt.Equal(rpA2.AmountOf(nativeDenom)))
	require.True(t, expectedNonNativeDenomAmt.Equal(rpA2.AmountOf(testDenom)))

	expectedNativeDenomAmt = rpB1.AmountOf(nativeDenom).Add(oneCoin)
	expectedNonNativeDenomAmt = rpB1.AmountOf(testDenom2).Sub(testCoinAmt1)
	require.True(t, expectedNativeDenomAmt.Equal(rpB2.AmountOf(nativeDenom)))
	require.True(t, expectedNonNativeDenomAmt.Equal(rpB2.AmountOf(testDenom2)))
}

func addLiquidityForTest(ctx sdk.Context, keeper Keeper, accs []exported.Account, nativeAmt, nonNativeAmt sdk.Int, denom string) error {
	if len(accs) == 0 {
		return fmt.Errorf("len is not enough")
	}

	var nonNativeDenomAmt = nonNativeAmt
	var nativeDenomAmt = nativeAmt
	var minReward = sdk.NewInt(1)

	t, err := time.Parse("2006-01-02T15:04:05.000Z", "2022-04-23T18:25:43.511Z")
	if err != nil {
		return err
	}

	nonNativeDeposit := sdk.Coin{Denom: denom, Amount: nonNativeDenomAmt}
	nativeDeposit := sdk.Coin{Denom: keeper.GetNativeDenom(ctx), Amount: nativeDenomAmt}

	err = keeper.RecieveCoins(ctx, accs[0].GetAddress(), sdk.Coins{nonNativeDeposit}...)
	if err != nil {
		return err
	}

	err = keeper.RecieveCoins(ctx, accs[0].GetAddress(), sdk.Coins{nativeDeposit}...)
	if err != nil {
		return err
	}

	msg := types.MsgAddLiquidity{
		Deposit:       nonNativeDeposit,
		DepositAmount: nativeDenomAmt,
		MinReward:     minReward,
		Deadline:      t,
		Sender:        accs[0].GetAddress(),
	}
	keeper.CreateReservePool(ctx, keeper.MustGetModuleName(keeper.GetNativeDenom(ctx), denom))
	keeper.AddInitialLiquidity(ctx, &msg)
	return nil
}

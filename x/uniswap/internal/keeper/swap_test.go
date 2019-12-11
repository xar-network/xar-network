package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/xar-network/xar-network/x/uniswap/internal/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	native = sdk.DefaultBondDenom
)

func TestIsDoubleSwap(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 0)

	cases := []struct {
		name         string
		denom1       string
		denom2       string
		isDoubleSwap bool
	}{
		{"denom1 is native", native, "btc", false},
		{"denom2 is native", "btc", native, false},
		{"neither denom is native", "eth", "btc", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doubleSwap := keeper.IsDoubleSwap(ctx, tc.denom1, tc.denom2)
			require.Equal(t, tc.isDoubleSwap, doubleSwap)
		})
	}
}

func TestGetAmount(t *testing.T) {
	defaultTestcase(t)
	defaultTestcaseErrors(t)

	testGetDoubleswapAmount(t)
	testDoubleSwapAmountErrors(t)
}

func defaultTestcaseErrors(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 1)

	oneCoin := sdk.NewInt(1)
	//oneHundredCoin := sdk.NewInt(100)
	nativeDenom := keeper.GetNativeDenom(ctx)
	testDenom := nonNativeDenomTest

	panicAssert := func() {
		keeper.GetInputAmount(ctx, oneCoin, nativeDenom, nativeDenom)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.GetInputAmount(ctx, oneCoin, testDenom, testDenom)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.GetOutputAmount(ctx, oneCoin, nativeDenom, nativeDenom)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.GetOutputAmount(ctx, oneCoin, testDenom, testDenom)
	}

	oCoin := sdk.NewCoin(testDenom, sdk.NewInt(1))
	iCoin := sdk.NewCoin(nativeDenom, sdk.NewInt(1))
	panicAssert = func() {
		keeper.InputAmount(ctx, oCoin, oCoin.Denom)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.InputAmount(ctx, iCoin, iCoin.Denom)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.OutputAmount(ctx, oCoin, oCoin.Denom)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.OutputAmount(ctx, iCoin, iCoin.Denom)
	}
	require.Panics(t, panicAssert)

	b := keeper.DenomIsNative(ctx, "asdd")
	require.False(t, b)

	b = keeper.DenomIsNative(ctx, keeper.GetNativeDenom(ctx))
	require.True(t, b)
}

func defaultTestcase(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewInt(0), 1)

	oneCoin := sdk.NewInt(1)
	//oneHundredCoin := sdk.NewInt(100)
	nativeDenom := keeper.GetNativeDenom(ctx)
	testDenom := nonNativeDenomTest
	nativeCoinAmt := sdk.NewInt(10040)
	nonNativeCoinAmt := sdk.NewInt(151000)

	err := addLiquidityForTest(ctx, keeper, accs, nativeCoinAmt, nonNativeCoinAmt, testDenom)
	if err != nil {
		panic(err)
	}

	amt := keeper.GetInputAmount(ctx, oneCoin, nativeDenom, testDenom)
	if !amt.Equal(sdk.NewInt(1)) {
		t.Error("incorrect amount calculation")
	}

	amt = keeper.GetOutputAmount(ctx, oneCoin, nativeDenom, testDenom)
	if !amt.Equal(sdk.NewInt(14)) {
		t.Error("incorrect amount calculation")
	}

	amt = keeper.InputAmount(ctx, sdk.NewCoin(testDenom, oneCoin), nativeDenom)
	if !amt.Equal(sdk.NewInt(1)) {
		t.Error("incorrect amount calculation")
	}

	amt = keeper.OutputAmount(ctx, sdk.NewCoin(nativeDenom, oneCoin), testDenom)
	if !amt.Equal(sdk.NewInt(14)) {
		t.Error("incorrect amount calculation")
	}
}

const nonNativeDenomTest = "asd"

func testGetDoubleswapAmount(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewInt(0), 1)

	oneHundredCoin := sdk.NewInt(100)
	testDenom := nonNativeDenomTest
	testDenom2 := nonNativeDenomTest + "2"

	outputTestCoin1 := sdk.NewCoin(testDenom, oneHundredCoin)
	outputTestCoin2 := sdk.NewCoin(testDenom2, oneHundredCoin)

	nativeCoinAmt := sdk.NewInt(10040)
	nonNativeCoinAmt := sdk.NewInt(151000)

	err := addLiquidityForTest(ctx, keeper, accs, nativeCoinAmt, nonNativeCoinAmt, testDenom)
	if err != nil {
		require.NoError(t, err)
	}

	err = addLiquidityForTest(ctx, keeper, accs, nativeCoinAmt, nonNativeCoinAmt, testDenom2)
	if err != nil {
		require.NoError(t, err)
	}

	nativeAmt, nonNativeAmt := keeper.DoubleSwapOutputAmount(ctx, outputTestCoin1, outputTestCoin2)
	require.Equal(t, nativeAmt, sdk.NewInt(6))
	require.Equal(t, nonNativeAmt, sdk.NewInt(89))

	nonNativeAmt, nativeAmt = keeper.DoubleSwapInputAmount(ctx, outputTestCoin1, outputTestCoin2)
	require.Equal(t, nativeAmt, sdk.NewInt(7))
	require.Equal(t, nonNativeAmt, sdk.NewInt(106))
}

func testDoubleSwapAmountErrors(t *testing.T) {
	equalDenomPanic(t)
	missingModuleNamePanic(t)
}

func equalDenomPanic(t *testing.T) {
	ctx, keeper, _ := createTestInput(t, sdk.NewInt(0), 1)

	oneHundredCoin := sdk.NewInt(100)
	testDenom := nonNativeDenomTest
	//testDenom2 := nonNativeDenomTest + "2"

	outputTestCoin1 := sdk.NewCoin(testDenom, oneHundredCoin)
	outputTestCoin2 := sdk.NewCoin(testDenom, oneHundredCoin)

	//nativeAmt, nonNativeAmt :=
	panicAssert := func() {
		keeper.DoubleSwapOutputAmount(ctx, outputTestCoin1, outputTestCoin2)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.DoubleSwapInputAmount(ctx, outputTestCoin1, outputTestCoin2)
	}
	require.Panics(t, panicAssert)
}

func missingModuleNamePanic(t *testing.T) {
	ctx, keeper, accs := createTestInput(t, sdk.NewInt(0), 1)

	oneHundredCoin := sdk.NewInt(100)
	testDenom := nonNativeDenomTest
	testDenom2 := nonNativeDenomTest + "2"

	outputTestCoin1 := sdk.NewCoin(testDenom, oneHundredCoin)
	outputTestCoin2 := sdk.NewCoin(testDenom2, oneHundredCoin)

	//nativeAmt, nonNativeAmt :=
	panicAssert := func() {
		keeper.DoubleSwapOutputAmount(ctx, outputTestCoin1, outputTestCoin2)
	}
	require.Panics(t, panicAssert)

	panicAssert = func() {
		keeper.DoubleSwapInputAmount(ctx, outputTestCoin1, outputTestCoin2)
	}
	require.Panics(t, panicAssert)

	tm, err := time.Parse("2006-01-02T15:04:05.000Z", "2022-04-23T18:25:43.511Z")
	if err != nil {
		require.NoError(t, err)
	}

	msg := types.MsgAddLiquidity{
		Deposit:       outputTestCoin1,
		DepositAmount: sdk.NewInt(10),
		MinReward:     sdk.NewInt(1),
		Deadline:      tm,
		Sender:        accs[0].GetAddress(),
	}

	keeper.CreateReservePool(ctx, keeper.MustGetModuleName(keeper.GetNativeDenom(ctx), outputTestCoin1.Denom))
	keeper.AddInitialLiquidity(ctx, &msg)

	panicAssert = func() {
		keeper.DoubleSwapOutputAmount(ctx, outputTestCoin1, outputTestCoin2)
	}
	require.Panics(t, panicAssert)
	panicAssert = func() {
		keeper.DoubleSwapInputAmount(ctx, outputTestCoin1, outputTestCoin2)
	}
	require.Panics(t, panicAssert)
}

func addLiquidityForTest(ctx sdk.Context, keeper Keeper, accs []exported.Account,nativeAmt, nonNativeAmt sdk.Int, denom string) error {
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

	_, err = keeper.bk.AddCoins(ctx, accs[0].GetAddress(), sdk.Coins{nonNativeDeposit})
	if err != nil {
		return err
	}

	_, err = keeper.bk.AddCoins(ctx, accs[0].GetAddress(), sdk.Coins{nativeDeposit})
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


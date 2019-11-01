package csdt

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

// How could one reduce the number of params in the test cases. Create a table driven test for each of the 4 add/withdraw collateral/debt?

func TestKeeper_ModifyCSDT(t *testing.T) {
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	ownerAddr := addrs[0]

	type state struct { // TODO this allows invalid state to be set up, should it?
		CSDT             CSDT
		OwnerCoins      sdk.Coins
		GlobalDebt      sdk.Int
		CollateralState CollateralState
	}
	type args struct {
		owner              sdk.AccAddress
		collateralDenom    string
		changeInCollateral sdk.Int
		changeInDebt       sdk.Int
	}

	tests := []struct {
		name       string
		priorState state
		price      string
		// also missing CSDTModuleParams
		args          args
		expectPass    bool
		expectedState state
	}{
		{
			"addCollateralAndDecreaseDebt",
			state{CSDT{ownerAddr, "xrp", i(100), i(2)}, cs(c("xrp", 10), c(StableDenom, 2)), i(2), CollateralState{"xrp", i(2)}},
			"10.345",
			args{ownerAddr, "xrp", i(10), i(-1)},
			true,
			state{CSDT{ownerAddr, "xrp", i(110), i(1)}, cs( /*  0xrp  */ c(StableDenom, 1)), i(1), CollateralState{"xrp", i(1)}},
		},
		{
			"removeTooMuchCollateral",
			state{CSDT{ownerAddr, "xrp", i(1000), i(200)}, cs(c("xrp", 10), c(StableDenom, 10)), i(200), CollateralState{"xrp", i(200)}},
			"1.00",
			args{ownerAddr, "xrp", i(-601), i(0)},
			false,
			state{CSDT{ownerAddr, "xrp", i(1000), i(200)}, cs(c("xrp", 10), c(StableDenom, 10)), i(200), CollateralState{"xrp", i(200)}},
		},
		{
			"withdrawTooMuchStableCoin",
			state{CSDT{ownerAddr, "xrp", i(1000), i(200)}, cs(c("xrp", 10), c(StableDenom, 10)), i(200), CollateralState{"xrp", i(200)}},
			"1.00",
			args{ownerAddr, "xrp", i(0), i(301)},
			false,
			state{CSDT{ownerAddr, "xrp", i(1000), i(200)}, cs(c("xrp", 10), c(StableDenom, 10)), i(200), CollateralState{"xrp", i(200)}},
		},
		{
			"createCSDTAndWithdrawStable",
			state{CSDT{}, cs(c("xrp", 10), c(StableDenom, 10)), i(0), CollateralState{"xrp", i(0)}},
			"1.00",
			args{ownerAddr, "xrp", i(5), i(2)},
			true,
			state{CSDT{ownerAddr, "xrp", i(5), i(2)}, cs(c("xrp", 5), c(StableDenom, 12)), i(2), CollateralState{"xrp", i(2)}},
		},
		{
			"emptyCSDT",
			state{CSDT{ownerAddr, "xrp", i(1000), i(200)}, cs(c("xrp", 10), c(StableDenom, 201)), i(200), CollateralState{"xrp", i(200)}},
			"1.00",
			args{ownerAddr, "xrp", i(-1000), i(-200)},
			true,
			state{CSDT{}, cs(c("xrp", 1010), c(StableDenom, 1)), i(0), CollateralState{"xrp", i(0)}},
		},
		{
			"invalidCollateralType",
			state{CSDT{}, cs(c("shitcoin", 5000000)), i(0), CollateralState{}},
			"0.000001",
			args{ownerAddr, "shitcoin", i(5000000), i(1)}, // ratio of 5:1
			false,
			state{CSDT{}, cs(c("shitcoin", 5000000)), i(0), CollateralState{}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// setup keeper
			mapp, keeper := setUpMockAppWithoutGenesis()
			// initialize csdt owner account with coins
			genAcc := auth.BaseAccount{
				Address: ownerAddr,
				Coins:   tc.priorState.OwnerCoins,
			}
			mock.SetGenesis(mapp, []auth.Account{&genAcc})
			// create a new context
			header := abci.Header{Height: mapp.LastBlockHeight() + 1}
			mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
			ctx := mapp.BaseApp.NewContext(false, header)
			// setup store state
			keeper.pricefeed.AddAsset(ctx, "xrp", "xrp test")
			keeper.pricefeed.SetPrice(
				ctx, sdk.AccAddress{}, "xrp",
				sdk.MustNewDecFromStr(tc.price),
				i(10))
			keeper.pricefeed.SetCurrentPrices(ctx)
			if tc.priorState.CSDT.CollateralDenom != "" { // check if the prior CSDT should be created or not (see if an empty one was specified)
				keeper.setCSDT(ctx, tc.priorState.CSDT)
			}
			keeper.setGlobalDebt(ctx, tc.priorState.GlobalDebt)
			if tc.priorState.CollateralState.Denom != "" {
				keeper.setCollateralState(ctx, tc.priorState.CollateralState)
			}

			// call func under test
			err := keeper.ModifyCSDT(ctx, tc.args.owner, tc.args.collateralDenom, tc.args.changeInCollateral, tc.args.changeInDebt)
			mapp.EndBlock(abci.RequestEndBlock{})
			mapp.Commit()

			// check for err
			if tc.expectPass {
				require.NoError(t, err, fmt.Sprint(err))
			} else {
				require.Error(t, err)
			}
			// get new state for verification
			actualCSDT, found := keeper.GetCSDT(ctx, tc.args.owner, tc.args.collateralDenom)
			actualGDebt := keeper.GetGlobalDebt(ctx)
			actualCstate, _ := keeper.GetCollateralState(ctx, tc.args.collateralDenom)
			// check state
			require.Equal(t, tc.expectedState.CSDT, actualCSDT)
			if tc.expectedState.CSDT.CollateralDenom == "" { // if the expected CSDT is blank, then expect the CSDT to have been deleted (hence not found)
				require.False(t, found)
			} else {
				require.True(t, found)
			}
			require.Equal(t, tc.expectedState.GlobalDebt, actualGDebt)
			require.Equal(t, tc.expectedState.CollateralState, actualCstate)
			// check owner balance
			mock.CheckBalance(t, mapp, ownerAddr, tc.expectedState.OwnerCoins)
		})
	}
}

// TODO change to table driven test to test more test cases
func TestKeeper_PartialSeizeCSDT(t *testing.T) {
	// Setup
	const collateral = "xrp"
	mapp, keeper := setUpMockAppWithoutGenesis()
	genAccs, addrs, _, _ := mock.CreateGenAccounts(1, cs(c(collateral, 100)))
	testAddr := addrs[0]
	mock.SetGenesis(mapp, genAccs)
	// setup pricefeed
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	keeper.pricefeed.AddAsset(ctx, collateral, "test description")
	keeper.pricefeed.SetPrice(
		ctx, sdk.AccAddress{}, collateral,
		sdk.MustNewDecFromStr("1.00"),
		i(10))
	keeper.pricefeed.SetCurrentPrices(ctx)
	// Create CSDT
	err := keeper.ModifyCSDT(ctx, testAddr, collateral, i(10), i(5))
	require.NoError(t, err)
	// Reduce price
	keeper.pricefeed.SetPrice(
		ctx, sdk.AccAddress{}, collateral,
		sdk.MustNewDecFromStr("0.90"),
		i(10))
	keeper.pricefeed.SetCurrentPrices(ctx)

	// Seize entire CSDT
	err = keeper.PartialSeizeCSDT(ctx, testAddr, collateral, i(10), i(5))

	// Check
	require.NoError(t, err)
	_, found := keeper.GetCSDT(ctx, testAddr, collateral)
	require.False(t, found)
	collateralState, found := keeper.GetCollateralState(ctx, collateral)
	require.True(t, found)
	require.Equal(t, sdk.ZeroInt(), collateralState.TotalDebt)
}

func TestKeeper_GetCSDTs(t *testing.T) {
	// setup keeper
	mapp, keeper := setUpMockAppWithoutGenesis()
	mock.SetGenesis(mapp, []auth.Account(nil))
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	// setup CSDTs
	_, addrs := mock.GeneratePrivKeyAddressPairs(2)
	csdts := CSDTs{
		{addrs[0], "xrp", i(4000), i(5)},
		{addrs[1], "xrp", i(4000), i(2000)},
		{addrs[0], "btc", i(10), i(20)},
	}
	for _, csdt := range csdts {
		keeper.setCSDT(ctx, csdt)
	}

	// Check nil params returns all CSDTs
	returnedCsdts, err := keeper.GetCSDTs(ctx, "", sdk.Dec{})
	require.NoError(t, err)
	require.Equal(t,
		CSDTs{
			{addrs[0], "btc", i(10), i(20)},
			{addrs[1], "xrp", i(4000), i(2000)},
			{addrs[0], "xrp", i(4000), i(5)}},
		returnedCsdts,
	)
	// Check correct CSDTs filtered by collateral and sorted
	returnedCsdts, err = keeper.GetCSDTs(ctx, "xrp", d("0.00000001"))
	require.NoError(t, err)
	require.Equal(t,
		CSDTs{
			{addrs[1], "xrp", i(4000), i(2000)},
			{addrs[0], "xrp", i(4000), i(5)}},
		returnedCsdts,
	)
	returnedCsdts, err = keeper.GetCSDTs(ctx, "xrp", sdk.Dec{})
	require.NoError(t, err)
	require.Equal(t,
		CSDTs{
			{addrs[1], "xrp", i(4000), i(2000)},
			{addrs[0], "xrp", i(4000), i(5)}},
		returnedCsdts,
	)
	returnedCsdts, err = keeper.GetCSDTs(ctx, "xrp", d("0.9"))
	require.NoError(t, err)
	require.Equal(t,
		CSDTs{
			{addrs[1], "xrp", i(4000), i(2000)}},
		returnedCsdts,
	)
	// Check high price returns no CSDTs
	returnedCsdts, err = keeper.GetCSDTs(ctx, "xrp", d("999999999.99"))
	require.NoError(t, err)
	require.Equal(t,
		CSDTs(nil),
		returnedCsdts,
	)
	// Check unauthorized collateral denom returns error
	_, err = keeper.GetCSDTs(ctx, "a non existent coin", d("0.34023"))
	require.Error(t, err)
	// Check price without collateral returns error
	_, err = keeper.GetCSDTs(ctx, "", d("0.34023"))
	require.Error(t, err)
	// Check deleting a CSDT removes it
	keeper.deleteCSDT(ctx, csdts[0])
	returnedCsdts, err = keeper.GetCSDTs(ctx, "", sdk.Dec{})
	require.NoError(t, err)
	require.Equal(t,
		CSDTs{
			{addrs[0], "btc", i(10), i(20)},
			{addrs[1], "xrp", i(4000), i(2000)}},
		returnedCsdts,
	)
}
func TestKeeper_GetSetDeleteCSDT(t *testing.T) {
	// setup keeper, create CSDT
	mapp, keeper := setUpMockAppWithoutGenesis()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	csdt := CSDT{addrs[0], "xrp", i(412), i(56)}

	// write and read from store
	keeper.setCSDT(ctx, csdt)
	readCSDT, found := keeper.GetCSDT(ctx, csdt.Owner, csdt.CollateralDenom)

	// check before and after match
	require.True(t, found)
	require.Equal(t, csdt, readCSDT)

	// delete auction
	keeper.deleteCSDT(ctx, csdt)

	// check auction does not exist
	_, found = keeper.GetCSDT(ctx, csdt.Owner, csdt.CollateralDenom)
	require.False(t, found)
}
func TestKeeper_GetSetGDebt(t *testing.T) {
	// setup keeper, create GDebt
	mapp, keeper := setUpMockAppWithoutGenesis()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	gDebt := i(4120000)

	// write and read from store
	keeper.setGlobalDebt(ctx, gDebt)
	readGDebt := keeper.GetGlobalDebt(ctx)

	// check before and after match
	require.Equal(t, gDebt, readGDebt)
}

func TestKeeper_GetSetCollateralState(t *testing.T) {
	// setup keeper, create CState
	mapp, keeper := setUpMockAppWithoutGenesis()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	collateralState := CollateralState{"xrp", i(15400)}

	// write and read from store
	keeper.setCollateralState(ctx, collateralState)
	readCState, found := keeper.GetCollateralState(ctx, collateralState.Denom)

	// check before and after match
	require.Equal(t, collateralState, readCState)
	require.True(t, found)
}

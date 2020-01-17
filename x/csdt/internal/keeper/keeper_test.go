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

package keeper_test

import (
	"fmt"
	"github.com/xar-network/xar-network/types/fee"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
	"github.com/xar-network/xar-network/x/oracle"
)


func NewDefaultFee() fee.Fee {
	return fee.Fee{
		Numerator:            sdk.NewInt(0),
		Denominator:          sdk.NewInt(1),
		MinimumAdditionalFee: sdk.NewInt(0),
		MinimumSubFee:        sdk.NewInt(0),
	}
}

func DefaultTestParams() types.Params {
	return types.NewParams(
		types.DefaultGlobalDebt,
		types.DefaultCollateralParams,
		types.DefaultDebtParams,
		types.DefaultCircuitBreaker,
		[]string{},
		NewDefaultFee(),
	)
}

// How could one reduce the number of params in the test cases. Create a table driven test for each of the 4 add/withdraw collateral/debt?

func TestKeeper_ModifyCSDT(t *testing.T) {
	_, addrs := mock.GeneratePrivKeyAddressPairs(2)
	ownerAddr := addrs[0]

	type state struct { // TODO this allows invalid state to be set up, should it?
		CSDT            CSDT
		OwnerCoins      sdk.Coins
		GlobalDebt      sdk.Int
		CollateralState CollateralState
		ModuleCoins     sdk.Coins
	}
	type args struct {
		owner              sdk.AccAddress
		collateralDenom    string
		changeInCollateral sdk.Int
		debtDenom    	   string
		changeInDebt       sdk.Int
	}

	tests := []struct {
		name       string
		params	   		*types.Params
		poolSnapshot	*types.PoolSnapshot
		priorState state
		price      string
		// also missing CSDTModuleParams
		args          args
		expectPass    bool
		expectedState state
		foundCSDT	  bool
	}{
		{
			"addCollateralAndDecreaseDebt",
			nil,
			nil,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 100)),
				Debt:             cs(c(StableDenom, 2)),
			}, cs(c("uftm", 10), c(StableDenom, 2)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(2)}, cs(c("uftm", 100))},
			"10.345",
			args{ownerAddr, "uftm", i(10), StableDenom, i(-1)},
			true,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 110)),
				Debt:             cs(c(StableDenom, 1)),
			}, cs( c(StableDenom, 1)), i(1), CollateralState{Denom: "uftm", TotalDebt: i(12)}, cs(c("uftm", 110))},
			true,
		},
		{
			"removeTooMuchCollateral",
			nil,
			nil,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 1000)),
				Debt:             cs(c(StableDenom, 200)),
			}, cs(c("uftm", 10), c(StableDenom, 10)), i(200), CollateralState{Denom: "uftm", TotalDebt: i(200)}, cs(c("uftm", 1000))},
			"1.00",
			args{ownerAddr, "uftm", i(-801), "uftm", i(0)},
			false,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 1000)),
				Debt:             cs(c(StableDenom, 200)),
			}, cs(c("uftm", 10), c(StableDenom, 10)), i(200), CollateralState{Denom: "uftm", TotalDebt: i(200)}, cs(c("uftm", 1000))},
			true,
		},
		{
			"withdrawTooMuchStableCoin",
			nil,
			nil,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 300)),
				Debt:             cs(c(StableDenom, 200)),
			}, cs(c("uftm", 10), c(StableDenom, 10)), i(200), CollateralState{Denom: "uftm", TotalDebt: i(200)}, cs(c("uftm", 300))},
			"1.00",
			args{ownerAddr, "uftm", i(0), "uftm", i(500)},
			false,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 300)),
				Debt:             cs(c(StableDenom, 200)),
			}, cs(c("uftm", 10), c(StableDenom, 10)), i(200), CollateralState{Denom: "uftm", TotalDebt: i(200)}, cs(c("uftm", 300))},
			true,
		},
		{
			"createCSDTAndWithdrawStable",
			nil,
			nil,
			state{CSDT{}, cs(c("uftm", 10), c(StableDenom, 10)), i(0), CollateralState{Denom: "uftm", TotalDebt: i(0)}, cs(c("uftm", 0))},
			"1.00",
			args{ownerAddr, "uftm", i(5), StableDenom, i(2)},
			true,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 5)),
				Debt:             cs(c(StableDenom, 2)),
			}, cs(c("uftm", 5), c(StableDenom, 12)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(5)}, cs(c("uftm", 5))},
			true,
		},
		{
			"emptyCSDTUtfm",
			nil,
			nil,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 1000)),
				Debt:             cs(c(StableDenom, 0)),
			}, cs(c("uftm", 10), c(StableDenom, 201)), i(200), CollateralState{Denom: "uftm", TotalDebt: i(1000)}, cs(c("uftm", 1000))},
			"1.00",
			args{ownerAddr, "uftm", i(-1000), StableDenom, i(0)},
			true,
			state{CSDT{}, cs(c("uftm", 1010), c(StableDenom, 201)), i(0), CollateralState{Denom: "uftm", TotalDebt: i(0)}, cs(c("uftm", 0))},
			false,
		},
		{
			"emptyCSDTStable",
			nil,
			nil,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 0)),
				Debt:             cs(c(StableDenom, 200)),
			}, cs(c("uftm", 10), c(StableDenom, 201)), i(200), CollateralState{Denom: StableDenom, TotalDebt: i(200)}, cs(c("uftm", 1000))},
			"1.00",
			args{ownerAddr, "uftm", i(0), StableDenom, i(-200)},
			true,
			state{CSDT{}, cs(c("uftm", 10), c(StableDenom, 1)), i(0), CollateralState{Denom: "uftm", TotalDebt: i(0)}, cs(c("uftm", 0))},
			false,
		},
		{
			"invalidCollateralType",
			nil,
			nil,
			state{CSDT{}, cs(c("shitcoin", 5000000)), i(0), CollateralState{}, cs(c("uftm", 0))},
			"0.000001",
			args{ownerAddr, "shitcoin", i(5000000), StableDenom, i(1)}, // ratio of 5:1
			false,
			state{CSDT{}, cs(c("shitcoin", 5000000)), i(0), CollateralState{}, cs(c("uftm", 0))},
			false,
		},
		{
			"addCollateralAndDecreaseDebtNotStable",
			nil,
			nil,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 100)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 10)),
			}, cs(c("uftm", 20), c(StableDenom, 2)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(2)}, cs(c("uftm", 100))},
			"10.345",
			args{ownerAddr, "uftm", i(10), "uftm", i(-1)},
			true,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 110)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 9)),
			}, cs(c("uftm", 9), c(StableDenom, 2)), i(1), CollateralState{Denom: "uftm", TotalDebt: i(11)}, cs(c("uftm", 110))},
			true,
		},
		{
			"addCollateralAndIncreaseDebtNotStable",
			nil,
			nil,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 100)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 10)),
			}, cs(c("uftm", 20), c(StableDenom, 2)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(2)}, cs(c("uftm", 100))},
			"10.345",
			args{ownerAddr, "uftm", i(10), "uftm", i(10)},
			true,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 110)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 20)),
			}, cs(c("uftm", 20), c(StableDenom, 2)), i(1), CollateralState{Denom: "uftm", TotalDebt: i(22)}, cs(c("uftm", 110))},
			true,
		},
		{
			"returnCollateralWithLimitDetect",
			&types.Params{
				CollateralParams: types.CollateralParams{
					types.CollateralParam{
						Denom:            "uftm",
						LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
						DebtLimit:        sdk.NewCoins(sdk.NewCoin("uftm", sdk.NewInt(500000000000))),
						DecreaseLimits:   []types.PoolDecreaseLimitParam{
							{
								BorderTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
								Period:     time.Hour,
								MaxPercent: sdk.NewInt(1),
							},
						},
					},
				},
				DebtParams:       types.DefaultDebtParams,
				GlobalDebtLimit:  types.DefaultGlobalDebt,
				CircuitBreaker:   types.DefaultCircuitBreaker,
				Fee:              NewDefaultFee(),
				Nominees:         []string{},
			},
			&types.PoolSnapshot{
				ByLimits: []types.PoolSnapValue{
					{
						Limit: types.PoolDecreaseLimitParam{
							BorderTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							Period:     time.Hour,
							MaxPercent: sdk.NewInt(1),
						},
						Val:   sdk.NewCoin("uftm", sdk.NewInt(10000)),
					},
				},
			},
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 100)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 10)),
			}, cs(c("uftm", 20), c(StableDenom, 2)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(2)}, cs(c("uftm", 100))},
			"10.345",
			args{ownerAddr, "uftm", i(-10), "uftm", i(0)},
			false,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 100)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 10)),
			}, cs(c("uftm", 20), c(StableDenom, 2)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(2)}, cs(c("uftm", 100))},
			true,
		},
		{
			"returnCollateralWithOutLimitDetect",
			&types.Params{
				CollateralParams: types.CollateralParams{
					types.CollateralParam{
						Denom:            "uftm",
						LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
						DebtLimit:        sdk.NewCoins(sdk.NewCoin("uftm", sdk.NewInt(500000000000))),
						DecreaseLimits:   []types.PoolDecreaseLimitParam{
							{
								BorderTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
								Period:     time.Hour,
								MaxPercent: sdk.NewInt(90),
							},
						},
					},
				},
				DebtParams:       types.DefaultDebtParams,
				GlobalDebtLimit:  types.DefaultGlobalDebt,
				CircuitBreaker:   types.DefaultCircuitBreaker,
				Fee:              NewDefaultFee(),
				Nominees:         []string{},
			},
			&types.PoolSnapshot{
				ByLimits: []types.PoolSnapValue{
					{
						Limit: types.PoolDecreaseLimitParam{
							BorderTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							Period:     time.Hour,
							MaxPercent: sdk.NewInt(1),
						},
						Val:   sdk.NewCoin("uftm", sdk.NewInt(10000)),
					},
				},
			},
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 100)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 10)),
			}, cs(c("uftm", 20), c(StableDenom, 2)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(2)}, cs(c("uftm", 100))},
			"10.345",
			args{ownerAddr, "uftm", i(-1), "uftm", i(0)},
			true,
			state{CSDT{
				Owner:            ownerAddr,
				CollateralAmount: cs(c("uftm", 99)),
				Debt:             cs(c(StableDenom, 2), c("uftm", 10)),
			}, cs(c("uftm", 21), c(StableDenom, 2)), i(2), CollateralState{Denom: "uftm", TotalDebt: i(1)}, cs(c("uftm", 100))},
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// setup keeper
			mapp, keeper, _, _ := setUpMockAppWithoutGenesis()
			// initialize csdt owner account with coins
			genAcc := auth.BaseAccount{
				Address: ownerAddr,
				Coins:   tc.priorState.OwnerCoins,
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
					AssetCode:  "uftm",
					BaseAsset:  "uftm",
					QuoteAsset: StableDenom,
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
				ctx, addrs[1], "uftm",
				sdk.MustNewDecFromStr(tc.price),
				time.Now().Add(time.Hour*1))
			_ = keeper.GetOracle().SetCurrentPrices(ctx)

			keeper.SetCSDT(ctx, tc.priorState.CSDT)
			if tc.priorState.CollateralState.Denom != "" {
				keeper.SetCollateralState(ctx, tc.priorState.CollateralState)
			}

			keeper.GetSupply().SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewCoin(StableDenom, tc.priorState.GlobalDebt))))

			keeper.GetSupply().MintCoins(ctx, types.ModuleName, tc.priorState.ModuleCoins)

			// call func under test
			params := DefaultTestParams()
			if tc.params != nil {
				params  = *tc.params
			}
			keeper.SetParams(ctx, params)

			if tc.poolSnapshot != nil {
				keeper.SetPoolSnapshot(ctx, *tc.poolSnapshot)
			}

			err := keeper.ModifyCSDT(ctx, tc.args.owner, types.NewSignedCoin(tc.args.collateralDenom, tc.args.changeInCollateral), types.NewSignedCoin(tc.args.debtDenom, tc.args.changeInDebt))
			mapp.EndBlock(abci.RequestEndBlock{})
			mapp.Commit()

			// check for err
			if tc.expectPass {
				require.NoError(t, err, fmt.Sprint(err))
			} else {
				require.Error(t, err)
			}
			// get new state for verification
			actualCSDT, found := keeper.GetCSDT(ctx, tc.args.owner)
			actualCstate, _ := keeper.GetCollateralState(ctx, tc.args.collateralDenom)

			// check state
			require.Equal(t, tc.expectedState.CSDT, actualCSDT)
			require.True(t, found || !tc.foundCSDT)
			require.Equal(t, tc.expectedState.CollateralState, actualCstate)
			// check owner balance
			mock.CheckBalance(t, mapp, ownerAddr, tc.expectedState.OwnerCoins)
		})
	}
}

// TODO change to table driven test to test more test cases
func TestKeeper_PartialSeizeCSDT(t *testing.T) {
	// Setup
	const collateral = "uftm"
	mapp, keeper, _, _ := setUpMockAppWithoutGenesis()
	genAccs, addrs, _, _ := mock.CreateGenAccounts(2, cs(c(collateral, 100)))

	testAddr := addrs[0]
	mock.SetGenesis(mapp, genAccs)
	// setup oracle
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)

	// setup store state
	oracleParams := oracle.DefaultParams()
	oracleParams.Assets = oracle.Assets{
		oracle.Asset{
			AssetCode:  collateral,
			BaseAsset:  collateral,
			QuoteAsset: StableDenom,
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
		ctx, addrs[1], collateral,
		sdk.MustNewDecFromStr("1.00"),
		time.Now().Add(time.Hour*1))
	_ = keeper.GetOracle().SetCurrentPrices(ctx)

	// Create CSDT
	keeper.SetParams(ctx, DefaultTestParams())
	keeper.GetSupply().SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewCoin(collateral, sdk.NewInt(200)))))

	err := keeper.ModifyCSDT(ctx, testAddr, types.NewSignedCoin(collateral, i(10)), types.NewSignedCoin(StableDenom, i(5)))
	require.NoError(t, err)
	// Reduce price
	_, _ = keeper.GetOracle().SetPrice(
		ctx, addrs[1], collateral,
		sdk.MustNewDecFromStr("0.50"),
		time.Now().Add(time.Hour*1))
	_ = keeper.GetOracle().SetCurrentPrices(ctx)
	// Seize entire CSDT
	err = keeper.PartialSeizeCSDT(ctx, testAddr, collateral, i(10), StableDenom, i(5))

	// Check
	require.NoError(t, err)
	_, found := keeper.GetCSDT(ctx, testAddr)
	require.False(t, found)
	collateralState, found := keeper.GetCollateralState(ctx, collateral)
	require.True(t, found)

	require.Equal(t, sdk.ZeroInt(), collateralState.TotalDebt)
}

// TODO change to table driven test to test more test cases
func TestKeeper_CollateralParams(t *testing.T) {
	// Setup
	const collateral = "uftm"
	mapp, keeper, _, _ := setUpMockAppWithoutGenesis()
	genAccs, addrs, _, _ := mock.CreateGenAccounts(2, cs(c(collateral, 100)))

	//testAddr := addrs[0]
	mock.SetGenesis(mapp, genAccs)
	// setup oracle
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)

	params := DefaultTestParams()
	params.Nominees = []string{addrs[1].String()}

	// Create CSDT
	keeper.SetParams(ctx, params)
	keeper.GetSupply().SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewCoin(collateral, sdk.NewInt(200)))))

	// Try to add a denom that already exists and fail
	collateralParam := types.CollateralParam{
		Denom:            "uftm",
		LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
		DebtLimit:        sdk.NewCoins(sdk.NewCoin(StableDenom, sdk.NewInt(500000000000))),
	}
	err := keeper.AddCollateralParam(ctx, addrs[1].String(), collateralParam)
	require.Error(t, err)
	// Try to set an existing denom

	err = keeper.SetCollateralParam(ctx, addrs[1].String(), collateralParam)
	require.NoError(t, err)

	// Try to set with non authority
	err = keeper.SetCollateralParam(ctx, addrs[0].String(), collateralParam)
	require.Error(t, err)
	// Try to add when not a nominee
	err = keeper.AddCollateralParam(ctx, addrs[0].String(), collateralParam)
	require.Error(t, err)
	collateralParam = types.CollateralParam{
		Denom:            "uftm2",
		LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
		DebtLimit:        sdk.NewCoins(sdk.NewCoin(StableDenom, sdk.NewInt(500000000000))),
	}
	err = keeper.SetCollateralParam(ctx, addrs[1].String(), collateralParam)
	require.Error(t, err)

	// Add successfully
	err = keeper.AddCollateralParam(ctx, addrs[1].String(), collateralParam)
	require.NoError(t, err)
}

func TestKeeper_GetCSDTs(t *testing.T) {
	// setup keeper
	mapp, keeper, _, _ := setUpMockAppWithoutGenesis()
	mock.SetGenesis(mapp, []exported.Account(nil))
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	// setup CSDTs
	_, addrs := mock.GeneratePrivKeyAddressPairs(2)
	csdts := CSDTs{
		{Owner: addrs[0], CollateralAmount: cs(c("uftm", 10)), Debt: cs(c(StableDenom, 20))},
		{Owner: addrs[1], CollateralAmount: cs(c("uftm", 4000)), Debt: cs(c(StableDenom, 2000))},
	}
	for _, csdt := range csdts {
		keeper.SetCSDT(ctx, csdt)
	}

	// Check deleting a CSDT removes it
	keeper.DeleteCSDT(ctx, csdts[0])
	returnedCsdts, err := keeper.GetCSDTs(ctx)
	require.NoError(t, err)
	require.Equal(t,
		CSDTs{
			{Owner: addrs[1], CollateralAmount: cs(c("uftm", 4000)), Debt: cs(c(StableDenom, 2000))},
		},
		returnedCsdts,
	)
}
func TestKeeper_GetSetDeleteCSDT(t *testing.T) {
	// setup keeper, create CSDT
	mapp, keeper, _, _ := setUpMockAppWithoutGenesis()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	csdt := CSDT{Owner: addrs[0], CollateralAmount: cs(c("uftm", 412)), Debt: cs(c(StableDenom, 56))}

	// write and read from store
	keeper.SetCSDT(ctx, csdt)
	readCSDT, found := keeper.GetCSDT(ctx, csdt.Owner)

	// check before and after match
	require.True(t, found)
	require.Equal(t, csdt, readCSDT)

	// delete auction
	keeper.DeleteCSDT(ctx, csdt)

	// check auction does not exist
	_, found = keeper.GetCSDT(ctx, csdt.Owner)
	require.False(t, found)
}

func TestKeeper_GetSetCollateralState(t *testing.T) {
	// setup keeper, create CState
	mapp, keeper, _, _ := setUpMockAppWithoutGenesis()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	collateralState := CollateralState{Denom: "uftm", TotalDebt: i(15400)}

	// write and read from store
	keeper.SetCollateralState(ctx, collateralState)
	readCState, found := keeper.GetCollateralState(ctx, collateralState.Denom)

	// check before and after match
	require.Equal(t, collateralState, readCState)
	require.True(t, found)
}

// shorten for easier reading
type (
	CSDT            = types.CSDT
	CSDTs           = types.CSDTs
	CollateralState = types.CollateralState
)

const (
	StableDenom = types.StableDenom
)

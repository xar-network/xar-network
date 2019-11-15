package csdt

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// Test the bank functionality of the CSDT keeper
func TestKeeper_AddSubtractGetCoins(t *testing.T) {
	_, addrs := mock.GeneratePrivKeyAddressPairs(1)
	normalAddr := addrs[0]

	tests := []struct {
		name          string
		address       sdk.AccAddress
		shouldAdd     bool
		amount        sdk.Coins
		expectedCoins sdk.Coins
	}{
		{"addNormalAddress", normalAddr, true, cs(c(types.StableDenom, 53)), cs(c(types.StableDenom, 153), c(types.GovDenom, 100))},
		{"subNormalAddress", normalAddr, false, cs(c(types.StableDenom, 53)), cs(c(types.StableDenom, 47), c(types.GovDenom, 100))},
		{"addLiquidatorStable", keeper.LiquidatorAccountAddress, true, cs(c(types.StableDenom, 53)), cs(c(types.StableDenom, 153))},
		{"subLiquidatorStable", keeper.LiquidatorAccountAddress, false, cs(c(types.StableDenom, 53)), cs(c(types.StableDenom, 47))},
		{"addLiquidatorGov", keeper.LiquidatorAccountAddress, true, cs(c(types.GovDenom, 53)), cs(c(types.StableDenom, 100))},  // no change to balance
		{"subLiquidatorGov", keeper.LiquidatorAccountAddress, false, cs(c(types.GovDenom, 53)), cs(c(types.StableDenom, 100))}, // no change to balance
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// setup k
			mapp, k := setUpMockAppWithoutGenesis()
			// initialize an account with coins
			genAcc := auth.BaseAccount{
				Address: normalAddr,
				Coins:   cs(c(types.StableDenom, 100), c(types.GovDenom, 100)),
			}
			mock.SetGenesis(mapp, []exported.Account{&genAcc})

			// create a new context and setup the liquidator account
			header := abci.Header{Height: mapp.LastBlockHeight() + 1}
			mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
			ctx := mapp.BaseApp.NewContext(false, header)
			k.SetLiquidatorModuleAccount(ctx, keeper.LiquidatorModuleAccount{Coins: cs(c(types.StableDenom, 100))}) // set gov coin "balance" to zero

			// perform the test action
			var err sdk.Error
			if tc.shouldAdd {
				_, err = k.AddCoins(ctx, tc.address, tc.amount)
			} else {
				_, err = k.SubtractCoins(ctx, tc.address, tc.amount)
			}

			mapp.EndBlock(abci.RequestEndBlock{})
			mapp.Commit()

			// check balances are as expected
			require.NoError(t, err)
			require.Equal(t, tc.expectedCoins, k.GetCoins(ctx, tc.address))
		})
	}
}

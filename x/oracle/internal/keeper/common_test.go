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
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/xar-network/xar-network/x/oracle/internal/keeper"
	"github.com/xar-network/xar-network/x/oracle/internal/types"
)

type testHelper struct {
	mApp     *mock.App
	keeper   keeper.Keeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

func getMockApp(t *testing.T, numGenAccs int, genState types.GenesisState, genAccs []authexported.Account) testHelper {
	mApp := mock.NewApp()
	types.RegisterCodec(mApp.Cdc)
	keyPricefeed := sdk.NewKVStoreKey(types.StoreKey)

	pk := mApp.ParamsKeeper
	keeper := keeper.NewKeeper(keyPricefeed, mApp.Cdc, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace)

	require.NoError(t, mApp.CompleteSetup(keyPricefeed))

	valTokens := sdk.TokensFromConsensusPower(42)
	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs,
			sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens)))
	}

	mock.SetGenesis(mApp, genAccs)
	return testHelper{mApp, keeper, addrs, pubKeys, privKeys}
}

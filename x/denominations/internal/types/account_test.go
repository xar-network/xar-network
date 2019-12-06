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

package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

func TestFreeze(t *testing.T) {
	_, pub1, addr1 := KeyTestPubAddr()

	acc1 := auth.NewBaseAccount(addr1, nil, pub1, 1, 2)

	// Not enough starting coins to move
	account1 := types.NewFreezeAccount(acc1, nil)
	err := account1.FreezeCoins(NewTestCoins("abc", 1))
	require.NotNil(t, err)
	err = account1.FreezeCoins(NewTestCoins("abc", 0))
	require.NotNil(t, err)

	// Not enough coins chosen
	coinSymbol := "ab1"
	coins := NewTestCoins(coinSymbol, 100)
	frozenCoins := NewTestCoins(coinSymbol, 0)

	acc2 := auth.NewBaseAccount(addr1, coins, pub1, 1, 2)

	account2 := types.NewFreezeAccount(acc2, frozenCoins)

	err = account2.FreezeCoins(NewTestCoins(coinSymbol, 0))
	require.NotNil(t, err)

	err = account2.FreezeCoins(nil)
	require.NotNil(t, err)

	// Too many coins to freeze
	err = account2.FreezeCoins(NewTestCoins(coinSymbol, 101))
	require.NotNil(t, err)

	// Can freeze
	require.Equal(t, true, sdk.NewInt(100).Equal(account2.GetCoins().AmountOf(coinSymbol)))
	err = account2.FreezeCoins(NewTestCoins(coinSymbol, 2))
	require.Nil(t, err)
	require.Equal(t, true, sdk.NewInt(98).Equal(account2.GetCoins().AmountOf(coinSymbol)))
	require.Equal(t, true, sdk.NewInt(2).Equal(account2.FrozenCoins.AmountOf(coinSymbol)))
}

package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestFreeze(t *testing.T) {
	_, pub1, addr1 := KeyTestPubAddr()

	// Not enough starting coins to move
	account1 := NewCustomAccount(addr1, nil, nil, pub1, 1, 2)
	err := account1.FreezeCoins(NewTestCoins("abc", 1))
	require.NotNil(t, err)
	err = account1.FreezeCoins(NewTestCoins("abc", 0))
	require.NotNil(t, err)

	// Not enough coins chosen
	coinSymbol := "ab1"
	coins := NewTestCoins(coinSymbol, 100)
	frozenCoins := NewTestCoins(coinSymbol, 0)
	account2 := NewCustomAccount(addr1, coins, frozenCoins, pub1, 1, 2)

	err = account2.FreezeCoins(NewTestCoins(coinSymbol, 0))
	require.NotNil(t, err)

	err = account2.FreezeCoins(nil)
	require.NotNil(t, err)

	// Too many coins to freeze
	err = account2.FreezeCoins(NewTestCoins(coinSymbol, 101))
	require.NotNil(t, err)

	// Can freeze
	err = account2.FreezeCoins(NewTestCoins(coinSymbol, 1))
	require.Nil(t, err)
	require.Equal(t, true, types.NewInt(99).Equal(account2.Coins.AmountOf(coinSymbol)))
	require.Equal(t, true, types.NewInt(1).Equal(account2.FrozenCoins.AmountOf(coinSymbol)))
}

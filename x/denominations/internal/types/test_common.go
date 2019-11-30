// nolint noalias
package types

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// coins to more than cover the fee
func NewTestCoins(symbol string, amount int64) sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin(symbol, amount),
	}
}

func KeyTestPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

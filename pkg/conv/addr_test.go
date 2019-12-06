/*

Copyright 2019 All in Bits, Inc
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

package conv

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAccAddressFromECDSAPubKey(t *testing.T) {
	pub := "cosmospub1addwnpepq23xtezz0cm48uxmwlr7wkz8mh39hytc5pxzwgjahvnes8yvywh924qf238"
	decPub := sdk.MustGetAccPubKeyBech32(pub).(secp256k1.PubKeySecp256k1)
	btcPub, err := btcec.ParsePubKey(decPub[:], btcec.S256())
	require.NoError(t, err)
	addr := AccAddressFromECDSAPubKey(btcPub.ToECDSA())
	assert.Equal(t, "cosmos14hvaduk4ghcre8h44rrg4nmhx3j6wa24kj00kg", addr.String())
}

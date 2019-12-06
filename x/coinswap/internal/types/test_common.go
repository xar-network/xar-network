/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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

package types

import (
	"time"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: deadcode unused
var (
	amt = sdk.NewInt(100)

	senderPk    = ed25519.GenPrivKey().PubKey()
	recipientPk = ed25519.GenPrivKey().PubKey()
	sender      = sdk.AccAddress(senderPk.Address())
	recipient   = sdk.AccAddress(recipientPk.Address())

	denom0   = "atom"
	denom1   = "btc"
	unidenom = "unibtc"

	input             = sdk.NewCoin(denom0, sdk.NewInt(1000))
	output            = sdk.NewCoin(denom1, sdk.NewInt(500))
	withdrawLiquidity = sdk.NewCoin(unidenom, sdk.NewInt(500))
	deadline          = time.Now().Unix()

	emptyAddr sdk.AccAddress
	emptyTime int64
)

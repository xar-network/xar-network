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

package app

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/xar-network/xar-network/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// This will fail half the time with the second output being 173
// This is due to secp256k1 signatures not being constant size.
// nolint: vet
func ExampleTxSendSize() {
	cdc := app.MakeCodec()
	var gas uint64 = 1

	priv1 := secp256k1.GenPrivKeySecp256k1([]byte{0})
	addr1 := sdk.AccAddress(priv1.PubKey().Address())
	priv2 := secp256k1.GenPrivKeySecp256k1([]byte{1})
	addr2 := sdk.AccAddress(priv2.PubKey().Address())
	coins := sdk.Coins{sdk.NewCoin("denom", sdk.NewInt(10))}
	msg1 := bank.MsgMultiSend{
		Inputs:  []bank.Input{bank.NewInput(addr1, coins)},
		Outputs: []bank.Output{bank.NewOutput(addr2, coins)},
	}
	fee := auth.NewStdFee(gas, coins)
	signBytes := auth.StdSignBytes("example-chain-ID",
		1, 1, fee, []sdk.Msg{msg1}, "")
	sig, err := priv1.Sign(signBytes)
	if err != nil {
		return
	}
	sigs := []auth.StdSignature{{nil, sig}}
	tx := auth.NewStdTx([]sdk.Msg{msg1}, fee, sigs, "")
	fmt.Println(len(cdc.MustMarshalBinaryBare([]sdk.Msg{msg1})))
	fmt.Println(len(cdc.MustMarshalBinaryBare(tx)))
	// output: 80
	// 169
}

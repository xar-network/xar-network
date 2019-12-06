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

package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/x/oracle/internal/types"
)

func TestMsgSort(t *testing.T) {
	from := sdk.AccAddress([]byte("someName"))
	price, _ := sdk.NewDecFromStr("1")
	expiry := time.Now()

	msg := types.NewMsgPostPrice(from, "uftm", price, expiry)

	fee := auth.NewStdFee(200000, nil)
	stdTx := auth.NewStdTx([]sdk.Msg{msg}, fee, []auth.StdSignature{}, "")
	signBytes := auth.StdSignBytes("xar-chain-dora", 4, 1, stdTx.Fee, stdTx.Msgs, stdTx.Memo)

	t.Logf("%s", signBytes)
	signed := auth.StdSignBytes(
		"xar-chain-dora", 4, 1, auth.NewStdFee(200000, nil), []sdk.Msg{msg}, "",
	)
	t.Logf("%s", signed)
}

func TestMsgPlaceBid_ValidateBasic(t *testing.T) {
	addr := sdk.AccAddress([]byte("someName"))
	// oracles := []Oracle{Oracle{
	// 	OracleAddress: addr.String(),
	// }}
	price, _ := sdk.NewDecFromStr("0.3005")
	expiry := time.Now().Add(time.Hour * 2)
	negativeExpiry := time.Now()
	negativePrice, _ := sdk.NewDecFromStr("-3.05")

	tests := []struct {
		name       string
		msg        types.MsgPostPrice
		expectPass bool
	}{
		{"normal", types.MsgPostPrice{addr, "xrp", price, expiry}, true},
		{"emptyAddr", types.MsgPostPrice{sdk.AccAddress{}, "xrp", price, expiry}, false},
		{"emptyAsset", types.MsgPostPrice{addr, "", price, expiry}, false},
		{"negativePrice", types.MsgPostPrice{addr, "xrp", negativePrice, expiry}, false},
		{"negativeExpiry", types.MsgPostPrice{addr, "xrp", price, negativeExpiry}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectPass {
				require.Nil(t, tc.msg.ValidateBasic())
			} else {
				require.NotNil(t, tc.msg.ValidateBasic())
			}
		})
	}
}

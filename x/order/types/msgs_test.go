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

package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"
	"github.com/xar-network/xar-network/x/order/types"
)

func TestMsgSort(t *testing.T) {
	addr := sdk.AccAddress([]byte("someName"))
	price := sdk.NewUintFromString("3005")
	quantity := sdk.NewUintFromString("10")
	marketID := store.NewEntityID(1)

	msg := types.NewMsgPost(addr, marketID, matcheng.Bid, price, quantity, 600)

	require.Equal(t, `{"direction":"BID","market_id":"1","owner":"cosmos1wdhk6e2wv9kk2j88d92","price":"3005","quantity":"10","time_in_force":600}`, string(msg.GetSignBytes()))
	signed := auth.StdSignBytes(
		"xar-chain-zafx", 4, 1, auth.NewStdFee(200000, nil), []sdk.Msg{msg}, "",
	)
	require.Equal(t, `{"account_number":"4","chain_id":"xar-chain-zafx","fee":{"amount":[],"gas":"200000"},"memo":"","msgs":[{"direction":"BID","market_id":"1","owner":"cosmos1wdhk6e2wv9kk2j88d92","price":"3005","quantity":"10","time_in_force":600}],"sequence":"1"}`, string(signed))
}

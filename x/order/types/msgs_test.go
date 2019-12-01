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

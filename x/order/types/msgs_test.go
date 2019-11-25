package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"
)

func TestMsgSort(t *testing.T) {
	addr := sdk.AccAddress([]byte("someName"))
	price := sdk.NewUintFromString("3005")
	quantity := sdk.NewUintFromString("10")
	marketID := store.NewEntityID(1)

	msg := NewMsgPost(addr, marketID, matcheng.Bid, price, quantity, 600)

	t.Errorf("%s", msg.GetSignBytes())
	signed := auth.StdSignBytes(
		"xar-chain-zafx", 4, 1, auth.NewStdFee(200000, nil), []sdk.Msg{msg}, "",
	)
	t.Errorf("%s", signed)
}

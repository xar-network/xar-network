package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

func TestMsgPlaceBid_ValidateBasic(t *testing.T) {
	addr := sdk.AccAddress([]byte("someName"))
	tests := []struct {
		name       string
		msg        types.MsgPlaceBid
		expectPass bool
	}{
		{"normal", types.MsgPlaceBid{0, addr, sdk.NewInt64Coin("csdt", 10), sdk.NewInt64Coin("ftm", 20)}, true},
		{"emptyAddr", types.MsgPlaceBid{0, sdk.AccAddress{}, sdk.NewInt64Coin("csdt", 10), sdk.NewInt64Coin("ftm", 20)}, false},
		{"negativeBid", types.MsgPlaceBid{0, addr, sdk.Coin{"csdt", sdk.NewInt(-10)}, sdk.NewInt64Coin("ftm", 20)}, false},
		{"negativeLot", types.MsgPlaceBid{0, addr, sdk.NewInt64Coin("csdt", 10), sdk.Coin{"ftm", sdk.NewInt(-20)}}, false},
		{"zerocoins", types.MsgPlaceBid{0, addr, sdk.NewInt64Coin("csdt", 0), sdk.NewInt64Coin("ftm", 0)}, true},
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

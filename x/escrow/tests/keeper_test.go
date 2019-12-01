package tests

import (
	"testing"

	"github.com/hashgard/hashgard/x/box/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestGetLockBoxByAddress(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	handler := box.NewHandler(keeper)

	boxInfo := GetLockBoxInfo()
	cap := 10
	for i := 0; i < cap; i++ {
		keeper.GetBankKeeper().AddCoins(ctx, SenderAccAddr, sdk.NewCoins(boxInfo.TotalAmount.Token))
		msg := msgs.NewMsgLockBox(SenderAccAddr, boxInfo)
		res := handler(ctx, msg)
		require.True(t, res.IsOK())
	}
	issues := keeper.GetBoxByAddress(ctx, types.Lock, SenderAccAddr)

	require.Len(t, issues, cap)
}

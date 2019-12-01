package tests

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestLockBoxEndBlocker(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.GetBankKeeper().SetSendEnabled(ctx, true)
	handler := box.NewHandler(keeper)

	inactiveQueue := keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	boxParams := GetLockBoxInfo()

	keeper.GetBankKeeper().AddCoins(ctx, SenderAccAddr, sdk.NewCoins(boxParams.TotalAmount.Token))

	msg := msgs.NewMsgLockBox(SenderAccAddr, boxParams)

	res := handler(ctx, msg)
	require.True(t, res.IsOK())
	var id string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id)
	boxInfo := keeper.GetBox(ctx, id)

	coins := keeper.GetBankKeeper().GetCoins(ctx, boxInfo.Owner)
	require.Equal(t, coins.AmountOf(boxInfo.TotalAmount.Token.Denom), sdk.ZeroInt())

	coins = keeper.GetDepositedCoins(ctx, boxInfo.Id)
	require.True(t, coins.IsEqual(sdk.NewCoins(boxInfo.TotalAmount.Token)))

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newHeader = ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Lock.EndTime, 0)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	box.EndBlocker(ctx, keeper)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	coins = keeper.GetBankKeeper().GetCoins(ctx, boxInfo.Owner)
	require.Equal(t, coins.AmountOf(boxInfo.TotalAmount.Token.Denom), boxInfo.TotalAmount.Token.Amount)

	coins = keeper.GetDepositedCoins(ctx, boxInfo.Id)
	require.True(t, coins.IsZero())
}

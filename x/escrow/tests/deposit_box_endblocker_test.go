package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashgard/hashgard/x/box/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/hashgard/hashgard/x/box/types"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestDepositBoxEndBlocker(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.GetBankKeeper().SetSendEnabled(ctx, true)
	handler := box.NewHandler(keeper)

	inactiveQueue := keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	boxInfo := createDepositBox(t, ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.Deposit.Interest.Token))

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	msgBoxInterest := msgs.NewMsgBoxInterestInject(boxInfo.Id, boxInfo.Owner, sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
		boxInfo.Deposit.Interest.Token.Amount.Add(sdk.NewInt(1))))
	res := handler(ctx, msgBoxInterest)
	require.False(t, res.IsOK())

	msgBoxInterest = msgs.NewMsgBoxInterestInject(boxInfo.Id, boxInfo.Owner, boxInfo.Deposit.Interest.Token)
	res = handler(ctx, msgBoxInterest)
	require.True(t, res.IsOK())

	coins := keeper.GetDepositedCoins(ctx, boxInfo.Id)
	require.True(t, coins.IsEqual(sdk.NewCoins(boxInfo.Deposit.Interest.Token)))

	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Deposit.StartTime, 0)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	box.EndBlocker(ctx, keeper)
	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	depositBox := keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, depositBox.Status, types.BoxInjecting)

	keeper.GetBankKeeper().AddCoins(ctx, TransferAccAddr, sdk.Coins{boxInfo.TotalAmount.Token})

	inject := boxInfo.TotalAmount.Token.Amount.Quo(sdk.NewInt(2))

	msgBoxInject := msgs.NewMsgBoxInject(boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		boxInfo.TotalAmount.Token.Amount.Add(sdk.NewInt(1))))
	res = handler(ctx, msgBoxInject)
	require.False(t, res.IsOK())

	msgBoxInject = msgs.NewMsgBoxInject(boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		boxInfo.Deposit.Price.Add(sdk.NewInt(1))))
	res = handler(ctx, msgBoxInject)
	require.False(t, res.IsOK())

	msgBoxInject = msgs.NewMsgBoxInject(boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		inject))
	res = handler(ctx, msgBoxInject)
	require.True(t, res.IsOK())

	depositBox = keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, depositBox.Deposit.Share, inject.Quo(boxInfo.Deposit.Price))

	newHeader = ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Deposit.EstablishTime, 0)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	box.EndBlocker(ctx, keeper)
	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	depositBox = keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, depositBox.Status, types.DepositBoxInterest)
	coins = keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)
	require.Equal(t, coins.AmountOf(boxInfo.Id), inject.Quo(depositBox.Deposit.Price))

	newHeader = ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Deposit.MaturityTime, 0)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	box.EndBlocker(ctx, keeper)
	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	depositBox = keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, depositBox.Status, types.BoxFinished)

	coins = keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)
	totalInterest := sdk.ZeroInt()
	for _, coin := range coins {
		if utils.IsId(coin.Denom) {
			interest, _, err := keeper.ProcessBoxWithdraw(ctx, coin.Denom, TransferAccAddr)
			require.Nil(t, err)
			totalInterest = totalInterest.Add(interest)
		}
	}
	depositBox = keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, totalInterest, depositBox.Deposit.WithdrawalInterest)

	coins = keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)
	for _, coin := range coins {
		require.False(t, utils.IsId(coin.Denom))
	}

	require.Equal(t, coins.AmountOf(boxInfo.TotalAmount.Token.Denom), boxInfo.TotalAmount.Token.Amount)
	require.Equal(t, coins.AmountOf(boxInfo.Deposit.Interest.Token.Denom), totalInterest)
}

func TestDepositBoxNotEnoughIteratorEndBlocker(t *testing.T) {
	str := fmt.Sprintf("%saa%s%d", types.IDPreStr, strconv.FormatUint(types.BoxMaxId, 36), 999)
	fmt.Println(str)
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.GetBankKeeper().SetSendEnabled(ctx, true)
	handler := box.NewHandler(keeper)

	inactiveQueue := keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	boxInfo := createDepositBox(t, ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.Deposit.Interest.Token))

	injection := boxInfo.Deposit.Interest.Token.Amount.Quo(sdk.NewInt(2))

	msgBoxInterest := msgs.NewMsgBoxInterestInject(boxInfo.Id, boxInfo.Owner,
		sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
			injection))
	res := handler(ctx, msgBoxInterest)
	require.True(t, res.IsOK())

	coins := keeper.GetBankKeeper().GetCoins(ctx, boxInfo.Owner)
	require.Equal(t, coins.AmountOf(boxInfo.Deposit.Interest.Token.Denom),
		boxInfo.Deposit.Interest.Token.Amount.Sub(injection))

	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Deposit.StartTime, 0)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	box.EndBlocker(ctx, keeper)
	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	depositBox := keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, depositBox.Status, types.BoxClosed)

	coins = keeper.GetBankKeeper().GetCoins(ctx, boxInfo.Owner)
	require.Equal(t, coins.AmountOf(boxInfo.Deposit.Interest.Token.Denom), boxInfo.Deposit.Interest.Token.Amount)
}
func TestDepositBoxNotEnoughDepositEndBlocker(t *testing.T) {
	str := fmt.Sprintf("%saa%s%d", types.IDPreStr, strconv.FormatUint(types.BoxMaxId, 36), 999)
	fmt.Println(str)
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.GetBankKeeper().SetSendEnabled(ctx, true)
	handler := box.NewHandler(keeper)

	inactiveQueue := keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	boxInfo := createDepositBox(t, ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.Deposit.Interest.Token))

	msgBoxInterest := msgs.NewMsgBoxInterestInject(boxInfo.Id, boxInfo.Owner, boxInfo.Deposit.Interest.Token)
	res := handler(ctx, msgBoxInterest)
	require.True(t, res.IsOK())

	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Deposit.StartTime, 0)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	box.EndBlocker(ctx, keeper)
	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	keeper.GetBankKeeper().AddCoins(ctx, TransferAccAddr, sdk.NewCoins(sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		boxInfo.Deposit.BottomLine)))

	deposit := boxInfo.Deposit.BottomLine.Quo(sdk.NewInt(2))

	msgBoxInject := msgs.NewMsgBoxInject(boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		deposit))
	res = handler(ctx, msgBoxInject)
	require.True(t, res.IsOK())

	coins := keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)
	require.Equal(t, coins.AmountOf(boxInfo.TotalAmount.Token.Denom), boxInfo.Deposit.BottomLine.Sub(deposit))

	newHeader = ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Deposit.EstablishTime, 0)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	box.EndBlocker(ctx, keeper)
	inactiveQueue = keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	depositBox := keeper.GetBox(ctx, boxInfo.Id)
	require.NotNil(t, depositBox)
	require.Equal(t, depositBox.Status, types.BoxClosed)

	msgBoxCancel := msgs.NewMsgBoxInjectCancel(boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		deposit))
	res = handler(ctx, msgBoxCancel)
	require.True(t, res.IsOK())

	depositBox = keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, depositBox.Status, types.BoxClosed)

	coins = keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)

	require.Equal(t, coins.AmountOf(boxInfo.TotalAmount.Token.Denom), boxInfo.Deposit.BottomLine)
}

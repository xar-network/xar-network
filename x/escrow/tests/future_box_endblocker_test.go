package tests

import (
	"testing"
	"time"

	"github.com/hashgard/hashgard/x/box/utils"

	"github.com/hashgard/hashgard/x/box/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/box/msgs"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestFutureBoxEndBlocker(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.GetBankKeeper().SetSendEnabled(ctx, true)
	handler := box.NewHandler(keeper)

	boxInfo := createFutureBox(t, ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.TotalAmount.Token))

	msgDeposit := msgs.NewMsgBoxInject(boxInfo.Id, boxInfo.Owner, boxInfo.TotalAmount.Token)
	res := handler(ctx, msgDeposit)
	require.True(t, res.IsOK())

	newBoxInfo := keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, newBoxInfo.Status, types.BoxActived)

	var address sdk.AccAddress

	coins := keeper.GetDepositedCoins(ctx, boxInfo.Id)
	require.True(t, coins.IsEqual(sdk.NewCoins(boxInfo.TotalAmount.Token)))

	for _, v := range boxInfo.Future.Receivers {
		for j, rec := range v {
			if j == 0 {
				address, _ = sdk.AccAddressFromBech32(rec)
				coins = keeper.GetBankKeeper().GetCoins(ctx, address)
				continue
			}
			amount, _ := sdk.NewIntFromString(rec)
			boxDenom := utils.GetCoinDenomByFutureBoxSeq(boxInfo.Id, j)
			require.Equal(t, coins.AmountOf(boxDenom), amount)
		}
	}

	newHeader := ctx.BlockHeader()
	newHeader.Time = time.Unix(boxInfo.Future.TimeLine[len(boxInfo.Future.TimeLine)-1], 0)
	ctx = ctx.WithBlockHeader(newHeader)

	box.EndBlocker(ctx, keeper)
	inactiveQueue := keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	newBoxInfo = keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, newBoxInfo.Status, types.BoxFinished)

	for _, v := range boxInfo.Future.Receivers {
		address, _ = sdk.AccAddressFromBech32(v[0])
		coins = keeper.GetBankKeeper().GetCoins(ctx, address)
		totalAmount := sdk.ZeroInt()
		for i, coin := range coins {
			sleep := boxInfo.Future.TimeLine[i] - time.Now().Unix()
			if sleep > 0 {
				time.Sleep(time.Duration(sleep) * time.Second)
			}
			_, _, err := keeper.ProcessBoxWithdraw(ctx, coin.Denom, address)
			require.Nil(t, err)
			totalAmount = totalAmount.Add(coin.Amount)
		}
		coins1 := keeper.GetBankKeeper().GetCoins(ctx, address)
		require.Equal(t, coins1.AmountOf(boxInfo.TotalAmount.Token.Denom), totalAmount)
	}
	coins = keeper.GetDepositedCoins(ctx, boxInfo.Id)
	require.True(t, coins.IsZero())
}

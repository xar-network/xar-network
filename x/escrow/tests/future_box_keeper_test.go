package tests

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/hashgard/hashgard/x/box/types"
	issueutils "github.com/hashgard/hashgard/x/issue/utils"

	"github.com/hashgard/hashgard/x/box"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func createFutureBox(t *testing.T, ctx sdk.Context, keeper box.Keeper) *types.BoxInfo {
	boxInfo := GetFutureBoxInfo()

	handler := box.NewHandler(keeper)
	msg := msgs.NewMsgFutureBox(SenderAccAddr, boxInfo)
	res := handler(ctx, msg)
	require.True(t, res.IsOK())

	var id string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id)

	box := keeper.GetBox(ctx, id)
	require.Equal(t, box.Name, boxInfo.Name)

	return box
}

func TestFutureBoxAdd(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	boxInfo := createFutureBox(t, ctx, keeper)

	err := keeper.CreateBox(ctx, boxInfo)
	require.Nil(t, err)
	box := keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, boxInfo.Name, box.Name)
}

func TestFutureBoxCancelDeposit(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	boxInfo := createFutureBox(t, ctx, keeper)

	err := keeper.CreateBox(ctx, boxInfo)
	require.Nil(t, err)
	box := keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, boxInfo.Name, box.Name)

	keeper.GetBankKeeper().AddCoins(ctx, TransferAccAddr, sdk.NewCoins(boxInfo.TotalAmount.Token))

	inject := issueutils.MulDecimals(sdk.NewInt(1000), TestTokenDecimals)
	fetch := issueutils.MulDecimals(sdk.NewInt(500), TestTokenDecimals)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr,
		sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
			issueutils.MulDecimals(sdk.NewInt(10000), TestTokenDecimals)), types.Inject)
	require.Error(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr,
		sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, inject), types.Inject)
	require.Nil(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		issueutils.MulDecimals(sdk.NewInt(5000), TestTokenDecimals)), types.Cancel)
	require.Error(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, fetch), types.Cancel)
	require.Nil(t, err)

	newBoxInfo := keeper.GetBox(ctx, boxInfo.Id)
	require.Equal(t, newBoxInfo.Future.Injects[0].Amount, inject.Sub(fetch))

	coins := keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)
	require.Equal(t, coins.AmountOf(boxInfo.TotalAmount.Token.Denom), boxInfo.TotalAmount.Token.Amount.Sub(inject).Add(fetch))
}

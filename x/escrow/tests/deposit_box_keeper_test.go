package tests

import (
	"testing"

	"github.com/hashgard/hashgard/x/box/msgs"
	issueutils "github.com/hashgard/hashgard/x/issue/utils"

	"github.com/hashgard/hashgard/x/box/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func createDepositBox(t *testing.T, ctx sdk.Context, keeper box.Keeper) *types.BoxInfo {
	boxInfo := GetDepositBoxInfo()

	handler := box.NewHandler(keeper)
	msg := msgs.NewMsgDepositBox(SenderAccAddr, boxInfo)
	res := handler(ctx, msg)
	require.True(t, res.IsOK())

	var id string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id)

	box := keeper.GetBox(ctx, id)
	require.Equal(t, box.Name, boxInfo.Name)

	return box
}

func TestDepositBoxCancelInterest(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	boxInfo := createDepositBox(t, ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.Deposit.Interest.Token))
	keeper.GetBankKeeper().AddCoins(ctx, TransferAccAddr, sdk.NewCoins(boxInfo.Deposit.Interest.Token))

	injection := boxInfo.Deposit.Interest.Token.Amount.Quo(sdk.NewInt(2))

	_, err := keeper.InjectDepositBoxInterest(ctx, boxInfo.Id, boxInfo.Owner, sdk.NewCoin("error",
		issueutils.MulDecimals(sdk.NewInt(1000), TestTokenDecimals)))
	require.Error(t, err)
	_, err = keeper.InjectDepositBoxInterest(ctx, boxInfo.Id, boxInfo.Owner, sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
		boxInfo.Deposit.Interest.Token.Amount.Add(sdk.NewInt(1))))
	require.Error(t, err)
	_, err = keeper.InjectDepositBoxInterest(ctx, boxInfo.Id, boxInfo.Owner, sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
		injection))
	require.Nil(t, err)
	_, err = keeper.InjectDepositBoxInterest(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
		injection))
	require.Nil(t, err)

	_, err = keeper.CancelInterestFromDepositBox(ctx, boxInfo.Id, boxInfo.Owner, sdk.NewCoin("error",
		issueutils.MulDecimals(sdk.NewInt(1000), TestTokenDecimals)))
	require.Error(t, err)
	_, err = keeper.CancelInterestFromDepositBox(ctx, boxInfo.Id, boxInfo.Owner, sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
		injection.Mul(sdk.NewInt(10))))
	require.Error(t, err)

	_, err = keeper.CancelInterestFromDepositBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
		injection))
	require.Nil(t, err)

	coins := keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)
	require.Equal(t, coins.AmountOf(boxInfo.Deposit.Interest.Token.Denom), boxInfo.Deposit.Interest.Token.Amount)

	boxInfo = keeper.GetBox(ctx, boxInfo.Id)
	require.Len(t, boxInfo.Deposit.InterestInjects, 1)

	_, err = keeper.CancelInterestFromDepositBox(ctx, boxInfo.Id, boxInfo.Owner, sdk.NewCoin(boxInfo.Deposit.Interest.Token.Denom,
		injection))
	require.Nil(t, err)

	coins = keeper.GetBankKeeper().GetCoins(ctx, boxInfo.Owner)
	require.Equal(t, coins.AmountOf(boxInfo.Deposit.Interest.Token.Denom), boxInfo.Deposit.Interest.Token.Amount)

}
func TestDepositBoxCancelDeposit(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	boxInfo := createDepositBox(t, ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, boxInfo.Owner, sdk.NewCoins(boxInfo.Deposit.Interest.Token))

	_, err := keeper.InjectDepositBoxInterest(ctx, boxInfo.Id, boxInfo.Owner, boxInfo.Deposit.Interest.Token)
	require.Nil(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, sdk.NewInt(10000)), types.Inject)
	require.Error(t, err)

	boxInfo = keeper.GetBox(ctx, boxInfo.Id)
	err = keeper.ProcessDepositBoxByEndBlocker(ctx, boxInfo)
	require.Nil(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, sdk.NewInt(10000)), types.Inject)
	require.Error(t, err)

	keeper.GetBankKeeper().AddCoins(ctx, TransferAccAddr, sdk.NewCoins(boxInfo.TotalAmount.Token))

	inject := issueutils.MulDecimals(sdk.NewInt(1000), TestTokenDecimals)
	fetch := issueutils.MulDecimals(sdk.NewInt(500), TestTokenDecimals)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr,
		sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, issueutils.MulDecimals(sdk.NewInt(100000), TestTokenDecimals)), types.Inject)
	require.Error(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, inject), types.Inject)
	require.Nil(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom,
		issueutils.MulDecimals(sdk.NewInt(10000), TestTokenDecimals)), types.Cancel)
	require.Error(t, err)

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, fetch), types.Cancel)
	require.Nil(t, err)

	coins := keeper.GetBankKeeper().GetCoins(ctx, TransferAccAddr)
	require.Equal(t, coins.AmountOf(boxInfo.TotalAmount.Token.Denom), boxInfo.TotalAmount.Token.Amount.Sub(inject).Add(fetch))

	_, err = keeper.ProcessInjectBox(ctx, boxInfo.Id, TransferAccAddr, sdk.NewCoin(boxInfo.TotalAmount.Token.Denom, inject), types.Inject)
	require.Nil(t, err)

	boxInfo = keeper.GetBox(ctx, boxInfo.Id)

	require.Equal(t, boxInfo.Deposit.TotalInject, inject.Add(inject).Sub(fetch))
}

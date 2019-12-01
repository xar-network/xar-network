package tests

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box/types"

	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/box/msgs"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestLockBoxImportExportQueues(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	handler := box.NewHandler(keeper)

	boxInfo := GetLockBoxInfo()

	keeper.GetBankKeeper().AddCoins(ctx, SenderAccAddr, sdk.NewCoins(boxInfo.TotalAmount.Token))

	msg := msgs.NewMsgLockBox(SenderAccAddr, boxInfo)
	res := handler(ctx, msg)
	require.True(t, res.IsOK())
	var id1 string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id1)

	keeper.GetBankKeeper().AddCoins(ctx, SenderAccAddr, sdk.NewCoins(boxInfo.TotalAmount.Token))
	msg = msgs.NewMsgLockBox(SenderAccAddr, boxInfo)
	res = handler(ctx, msg)
	require.True(t, res.IsOK())
	var id2 string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id2)

	genAccs := mapp.AccountKeeper.GetAllAccounts(ctx)

	// Export the state and import it into a new Mock App
	genState := box.ExportGenesis(ctx, keeper)
	mapp2, keeper2, _, _, _, _ := getMockApp(t, genState, genAccs)

	header = abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp2.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx2 := mapp2.BaseApp.NewContext(false, abci.Header{})

	boxInfo1 := keeper2.GetBox(ctx2, id1)
	require.NotNil(t, boxInfo1)
	boxInfo2 := keeper2.GetBox(ctx2, id2)
	require.NotNil(t, boxInfo2)

	require.True(t, boxInfo1.Status == types.LockBoxLocked)
	require.True(t, boxInfo2.Status == types.LockBoxLocked)

	ctx2 = ctx2.WithBlockTime(time.Unix(boxInfo.Lock.EndTime, 0))

	box.EndBlocker(ctx2, keeper2)

	boxInfo1 = keeper2.GetBox(ctx2, id1)
	require.NotNil(t, boxInfo1)
	boxInfo2 = keeper2.GetBox(ctx2, id2)
	require.NotNil(t, boxInfo2)

	require.True(t, boxInfo1.Status == types.LockBoxUnlocked)
	require.True(t, boxInfo2.Status == types.LockBoxUnlocked)
}

func TestDepositBoxImportExportQueues(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	handler := box.NewHandler(keeper)

	boxInfo := GetDepositBoxInfo()

	msg := msgs.NewMsgDepositBox(SenderAccAddr, boxInfo)
	res := handler(ctx, msg)
	require.True(t, res.IsOK())
	var id1 string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id1)

	keeper.GetBankKeeper().AddCoins(ctx, SenderAccAddr, sdk.NewCoins(boxInfo.Deposit.Interest.Token))
	msgBoxInterest := msgs.NewMsgBoxInterestInject(id1, SenderAccAddr, boxInfo.Deposit.Interest.Token)
	res = handler(ctx, msgBoxInterest)
	require.True(t, res.IsOK())

	ctx = ctx.WithBlockTime(time.Unix(boxInfo.Deposit.StartTime, 0))
	box.EndBlocker(ctx, keeper)

	keeper.GetBankKeeper().AddCoins(ctx, TransferAccAddr, sdk.Coins{boxInfo.TotalAmount.Token})
	msgBoxInject := msgs.NewMsgBoxInject(id1, TransferAccAddr, boxInfo.TotalAmount.Token)
	res = handler(ctx, msgBoxInject)
	require.True(t, res.IsOK())

	ctx = ctx.WithBlockTime(time.Unix(boxInfo.Deposit.EstablishTime, 0))
	box.EndBlocker(ctx, keeper)

	boxInfo = GetDepositBoxInfo()

	msg = msgs.NewMsgDepositBox(SenderAccAddr, boxInfo)
	res = handler(ctx, msg)
	require.True(t, res.IsOK())
	var id2 string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id2)

	keeper.GetBankKeeper().AddCoins(ctx, SenderAccAddr, sdk.NewCoins(boxInfo.Deposit.Interest.Token))
	msg = msgs.NewMsgDepositBox(SenderAccAddr, boxInfo)
	res = handler(ctx, msg)
	require.True(t, res.IsOK())
	var id3 string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id3)

	msgBoxInterest = msgs.NewMsgBoxInterestInject(id3, SenderAccAddr, boxInfo.Deposit.Interest.Token)
	res = handler(ctx, msgBoxInterest)
	require.True(t, res.IsOK())

	genAccs := mapp.AccountKeeper.GetAllAccounts(ctx)

	// Export the state and import it into a new Mock App
	genState := box.ExportGenesis(ctx, keeper)
	mapp2, keeper2, _, _, _, _ := getMockApp(t, genState, genAccs)

	header = abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp2.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx2 := mapp2.BaseApp.NewContext(false, abci.Header{})

	boxInfo1 := keeper2.GetBox(ctx2, id1)
	require.NotNil(t, boxInfo1)
	boxInfo2 := keeper2.GetBox(ctx2, id2)
	require.NotNil(t, boxInfo2)
	boxInfo3 := keeper2.GetBox(ctx2, id3)
	require.NotNil(t, boxInfo3)

	require.True(t, boxInfo1.Status == types.DepositBoxInterest)
	require.True(t, boxInfo2.Status == types.BoxCreated)
	require.True(t, boxInfo3.Status == types.BoxCreated)

	ctx2 = ctx2.WithBlockTime(time.Unix(boxInfo.Deposit.MaturityTime, 0))
	box.EndBlocker(ctx2, keeper2)

	boxInfo1 = keeper2.GetBox(ctx2, id1)
	require.NotNil(t, boxInfo1)
	boxInfo2 = keeper2.GetBox(ctx2, id2)
	require.True(t, boxInfo2.Status == types.BoxClosed)
	boxInfo3 = keeper2.GetBox(ctx2, id3)
	require.NotNil(t, boxInfo3)

	require.True(t, boxInfo1.Status == types.BoxFinished)
	require.True(t, boxInfo3.Status == types.BoxInjecting)
}

func TestFutureBoxImportExportQueues(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, box.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	handler := box.NewHandler(keeper)

	boxInfo := GetFutureBoxInfo()

	msg := msgs.NewMsgFutureBox(SenderAccAddr, boxInfo)
	res := handler(ctx, msg)
	require.True(t, res.IsOK())
	var id1 string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id1)

	keeper.GetBankKeeper().AddCoins(ctx, SenderAccAddr, sdk.NewCoins(boxInfo.TotalAmount.Token))
	msgDeposit := msgs.NewMsgBoxInject(id1, SenderAccAddr, boxInfo.TotalAmount.Token)
	res = handler(ctx, msgDeposit)
	require.True(t, res.IsOK())

	boxInfo = GetFutureBoxInfo()

	msg = msgs.NewMsgFutureBox(SenderAccAddr, boxInfo)
	res = handler(ctx, msg)
	require.True(t, res.IsOK())
	var id2 string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &id2)

	genAccs := mapp.AccountKeeper.GetAllAccounts(ctx)

	// Export the state and import it into a new Mock App
	genState := box.ExportGenesis(ctx, keeper)
	mapp2, keeper2, _, _, _, _ := getMockApp(t, genState, genAccs)

	header = abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp2.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx2 := mapp2.BaseApp.NewContext(false, abci.Header{})

	boxInfo1 := keeper2.GetBox(ctx2, id1)
	require.NotNil(t, boxInfo1)
	boxInfo2 := keeper2.GetBox(ctx2, id2)
	require.NotNil(t, boxInfo2)

	require.True(t, boxInfo1.Status == types.BoxActived)
	require.True(t, boxInfo2.Status == types.BoxInjecting)

	ctx2 = ctx2.WithBlockTime(time.Unix(boxInfo.Future.TimeLine[len(boxInfo.Future.TimeLine)-1], 0))
	box.EndBlocker(ctx2, keeper2)

	boxInfo1 = keeper2.GetBox(ctx2, id1)
	require.NotNil(t, boxInfo1)
	boxInfo2 = keeper2.GetBox(ctx2, id2)
	require.True(t, boxInfo2.Status == types.BoxClosed)

	require.True(t, boxInfo1.Status == types.BoxFinished)

}

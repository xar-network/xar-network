package keeper_test

import (
	"testing"

	"github.com/xar-network/xar-network/x/issue/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/xar-network/xar-network/x/issue"
)

func TestCreateIssue(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)
	coinIssue := k.GetIssue(ctx, CoinIssueInfo.IssueId)
	require.Equal(t, coinIssue.TotalSupply, CoinIssueInfo.TotalSupply)
	coin := sdk.Coin{Denom: CoinIssueInfo.IssueId, Amount: sdk.NewInt(5000)}
	err = k.GetBankKeeper().SendCoins(ctx, SenderAccAddr, ReceiverCoinsAccAddr,
		sdk.NewCoins(coin))
	require.Nil(t, err)
	coinIssue = k.GetIssue(ctx, CoinIssueInfo.IssueId)
	require.True(t, coinIssue.TotalSupply.Equal(CoinIssueInfo.TotalSupply))
	acc := mapp.AccountKeeper.GetAccount(ctx, ReceiverCoinsAccAddr)
	amount := acc.GetCoins().AmountOf(CoinIssueInfo.IssueId)
	flag1 := amount.Equal(coin.Amount)
	require.True(t, flag1)
}

func TestGetIssues(t *testing.T) {
	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	cap := 10
	for i := 0; i < cap; i++ {
		_, err := k.CreateIssue(ctx, &CoinIssueInfo)
		require.Nil(t, err)
	}
	issues := k.GetIssues(ctx, CoinIssueInfo.Issuer.String())

	require.Len(t, issues, cap)
}

func TestMint(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)
	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)
	_, err = k.Mint(ctx, CoinIssueInfo.IssueId, sdk.NewInt(10000), SenderAccAddr, SenderAccAddr)
	require.Nil(t, err)
	coinIssue := k.GetIssue(ctx, CoinIssueInfo.IssueId)
	require.True(t, coinIssue.TotalSupply.Equal(sdk.NewInt(20000)))
}

func TestBurnOwner(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	_, err = k.BurnOwner(ctx, CoinIssueInfo.IssueId, sdk.NewInt(5000), SenderAccAddr)
	require.Nil(t, err)

	err = k.DisableFeature(ctx, CoinIssueInfo.Owner, CoinIssueInfo.IssueId, types.BurnOwner)
	require.Nil(t, err)

	_, err = k.BurnOwner(ctx, CoinIssueInfo.IssueId, sdk.NewInt(5000), SenderAccAddr)
	require.Error(t, err)

}

func TestBurnHolder(t *testing.T) {
	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.GetBankKeeper().SendCoins(ctx, SenderAccAddr, ReceiverCoinsAccAddr, sdk.NewCoins(sdk.NewCoin(CoinIssueInfo.IssueId, sdk.NewInt(10000))))
	require.Nil(t, err)

	_, err = k.BurnHolder(ctx, CoinIssueInfo.IssueId, sdk.NewInt(5000), ReceiverCoinsAccAddr)
	require.Nil(t, err)

	err = k.DisableFeature(ctx, CoinIssueInfo.Owner, CoinIssueInfo.IssueId, types.BurnHolder)
	require.Nil(t, err)

	_, err = k.BurnHolder(ctx, CoinIssueInfo.IssueId, sdk.NewInt(5000), ReceiverCoinsAccAddr)
	require.Error(t, err)

}

func TestBurnFrom(t *testing.T) {
	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.GetBankKeeper().SendCoins(ctx, SenderAccAddr, ReceiverCoinsAccAddr, sdk.NewCoins(sdk.NewCoin(CoinIssueInfo.IssueId, sdk.NewInt(10000))))
	require.Nil(t, err)

	_, err = k.BurnFrom(ctx, CoinIssueInfo.IssueId, sdk.NewInt(5000), SenderAccAddr, ReceiverCoinsAccAddr)
	require.Nil(t, err)

	err = k.DisableFeature(ctx, CoinIssueInfo.Owner, CoinIssueInfo.IssueId, types.BurnFrom)
	require.Nil(t, err)

	_, err = k.BurnFrom(ctx, CoinIssueInfo.IssueId, sdk.NewInt(5000), ReceiverCoinsAccAddr, ReceiverCoinsAccAddr)
	require.Error(t, err)
}

func TestApprove(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.Approve(ctx, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(5000))
	require.Nil(t, err)

	amount := k.Allowance(ctx, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId)

	require.Equal(t, amount, sdk.NewInt(5000))

}
func TestSendFrom(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.SendFrom(ctx, TransferAccAddr, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(1000))
	require.Error(t, err)

	err = k.Approve(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(5000))
	require.Nil(t, err)

	err = k.SendFrom(ctx, TransferAccAddr, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(6000))
	require.Error(t, err)

	err = k.SendFrom(ctx, TransferAccAddr, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(3000))
	require.Nil(t, err)

	amount := k.Allowance(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId)
	require.Equal(t, amount, sdk.NewInt(2000))

}

func TestSendFromByFreeze(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.Approve(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(5000))
	require.Nil(t, err)

	err = k.Freeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, ReceiverCoinsAccAddr, types.FreezeIn)
	require.Nil(t, err)

	err = k.SendFrom(ctx, TransferAccAddr, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(30000))
	require.Error(t, err)

	err = k.Freeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, SenderAccAddr, types.FreezeOut)
	require.Nil(t, err)

	err = k.SendFrom(ctx, TransferAccAddr, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(30000000))
	require.Error(t, err)

	err = k.UnFreeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, SenderAccAddr, types.FreezeInAndOut)
	require.Nil(t, err)

	err = k.UnFreeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, ReceiverCoinsAccAddr, types.FreezeInAndOut)
	require.Nil(t, err)

	err = k.SendFrom(ctx, TransferAccAddr, SenderAccAddr, ReceiverCoinsAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(3000))
	require.Nil(t, err)
}

func TestIncreaseApproval(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.Approve(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(5000))
	require.Nil(t, err)

	k.IncreaseApproval(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(1000))
	require.Nil(t, err)

	amount := k.Allowance(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId)

	require.Equal(t, amount, sdk.NewInt(6000))

}

func TestDecreaseApproval(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.Approve(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(5000))
	require.Nil(t, err)

	k.DecreaseApproval(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(6000))
	require.Nil(t, err)

	amount := k.Allowance(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId)

	require.Equal(t, amount, sdk.NewInt(0))

	err = k.Approve(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(5000))
	require.Nil(t, err)

	k.DecreaseApproval(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId, sdk.NewInt(4000))
	require.Nil(t, err)

	amount = k.Allowance(ctx, SenderAccAddr, TransferAccAddr, CoinIssueInfo.IssueId)

	require.Equal(t, amount, sdk.NewInt(1000))

}

func TestFreeze(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	err = k.Freeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, TransferAccAddr, types.FreezeIn)
	require.Nil(t, err)

	err = k.Freeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, TransferAccAddr, types.FreezeOut)
	require.Nil(t, err)

	freeze := k.GetFreeze(ctx, TransferAccAddr, CoinIssueInfo.IssueId)
	require.NotZero(t, freeze.String())
	require.NotZero(t, freeze.String())

	err = k.UnFreeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, TransferAccAddr, types.FreezeIn)
	require.Nil(t, err)

	err = k.UnFreeze(ctx, CoinIssueInfo.IssueId, SenderAccAddr, TransferAccAddr, types.FreezeOut)
	require.Nil(t, err)

	freeze = k.GetFreeze(ctx, TransferAccAddr, CoinIssueInfo.IssueId)
	require.Equal(t, false, freeze.Frozen)

}

func TestTransferOwnership(t *testing.T) {

	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	CoinIssueInfo.TotalSupply = sdk.NewInt(10000)

	_, err := k.CreateIssue(ctx, &CoinIssueInfo)
	require.Nil(t, err)

	owner := CoinIssueInfo.Owner

	err = k.TransferOwnership(ctx, CoinIssueInfo.IssueId, owner, TransferAccAddr)
	require.Nil(t, err)

	err = k.TransferOwnership(ctx, CoinIssueInfo.IssueId, owner, TransferAccAddr)
	require.Error(t, err)

	issueIDs := k.GetAddressIssues(ctx, owner.String())
	require.Len(t, issueIDs, 0)
	issueIDs = k.GetAddressIssues(ctx, TransferAccAddr.String())
	require.Len(t, issueIDs, 1)
	issueInfo := k.GetIssue(ctx, CoinIssueInfo.IssueId)
	require.Equal(t, owner, issueInfo.Issuer)
	require.Equal(t, TransferAccAddr, issueInfo.Owner)
}

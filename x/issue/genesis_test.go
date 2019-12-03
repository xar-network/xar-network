package issue_test

import (
	"testing"

	"github.com/xar-network/xar-network/x/issue"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/x/issue/internal/types"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestImportExportQueues(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, issue.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})
	keeper.GetSupplyKeeper().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	handler := issue.NewHandler(keeper)

	res := handler(ctx, types.NewMsgIssue(SenderAccAddr, &IssueParams))
	require.True(t, res.IsOK())

	var issueID1 string
	issueID1 = string(res.Data)
	require.NotNil(t, issueID1)

	res = handler(ctx, types.NewMsgIssue(SenderAccAddr, &IssueParams))
	require.True(t, res.IsOK())

	var issueID2 string
	issueID2 = string(res.Data)
	require.NotNil(t, issueID2)

	genAccs := mapp.AccountKeeper.GetAllAccounts(ctx)

	// Export the state and import it into a new Mock App
	genState := issue.ExportGenesis(ctx, keeper)
	mapp2, keeper2, _, _, _, _ := getMockApp(t, genState, genAccs)

	header = abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp2.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx2 := mapp2.BaseApp.NewContext(false, abci.Header{})

	issueInfo1 := keeper2.GetIssue(ctx2, issueID1)
	require.NotNil(t, issueInfo1)
	issueInfo2 := keeper2.GetIssue(ctx2, issueID2)
	require.NotNil(t, issueInfo2)
}

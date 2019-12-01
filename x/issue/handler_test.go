package issue_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/issue"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

func TestHandlerNewMsgIssue(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.NewContext(false, abci.Header{})
	mapp.InitChainer(ctx, abci.RequestInitChain{})

	handler := issue.NewHandler(keeper)

	res := handler(ctx, types.NewMsgIssue(SenderAccAddr, &IssueParams))
	require.True(t, res.IsOK())

	var issueID string
	issueID = string(res.Data)
	require.NotNil(t, issueID)
}

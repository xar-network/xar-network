package tests

import (
	"testing"

	"github.com/zar-network/zar-network/x/issue"
	"github.com/zar-network/zar-network/x/issue/msgs"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestHandlerNewMsgIssue(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.NewContext(false, abci.Header{})
	mapp.InitChainer(ctx, abci.RequestInitChain{})

	handler := issue.NewHandler(keeper)

	res := handler(ctx, msgs.NewMsgIssue(SenderAccAddr, &IssueParams))
	require.True(t, res.IsOK())

	var issueID string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &issueID)
	require.NotNil(t, issueID)
}

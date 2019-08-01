package tests

import (
	"fmt"
	"testing"

	"github.com/zar-network/zar-network/x/issue/utils"
	"github.com/tendermint/tendermint/crypto"

	"github.com/zar-network/zar-network/x/issue/params"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zar-network/zar-network/x/issue"
	queriers2 "github.com/zar-network/zar-network/x/issue/client/queriers"
	"github.com/zar-network/zar-network/x/issue/msgs"
	"github.com/zar-network/zar-network/x/issue/types"
)

func TestQueryIssue(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})

	querier := issue.NewQuerier(keeper)
	handler := issue.NewHandler(keeper)

	res := handler(ctx, msgs.NewMsgIssue(SenderAccAddr, &IssueParams))
	var issueID string
	keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &issueID)
	bz := getQueried(t, ctx, querier, queriers2.GetQueryIssuePath(issueID), types.QueryIssue, issueID)
	var issueInfo types.CoinIssueInfo
	keeper.Getcdc().MustUnmarshalJSON(bz, &issueInfo)

	require.Equal(t, issueInfo.GetIssueId(), issueID)
	require.Equal(t, issueInfo.GetName(), CoinIssueInfo.GetName())

}

func TestQueryIssues(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})
	//querier := issue.NewQuerier(keeper)
	handler := issue.NewHandler(keeper)
	cap := 10
	for i := 0; i < cap; i++ {
		handler(ctx, msgs.NewMsgIssue(SenderAccAddr, &IssueParams))
	}
	issues := keeper.List(ctx, params.IssueQueryParams{Limit: 10})
	require.Len(t, issues, cap)

}

func TestSearchIssues(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.NewContext(false, abci.Header{})
	querier := issue.NewQuerier(keeper)
	handler := issue.NewHandler(keeper)
	cap := 10
	for i := 0; i < cap; i++ {
		handler(ctx, msgs.NewMsgIssue(SenderAccAddr, &IssueParams))
	}
	bz := getQueried(t, ctx, querier, queriers2.GetQueryIssuePath("TES"), types.QuerySearch, "TES")
	var issues types.CoinIssues
	keeper.Getcdc().MustUnmarshalJSON(bz, &issues)
	require.Len(t, issues, cap)

}
func getQueried(t *testing.T, ctx sdk.Context, querier sdk.Querier, path string, querierRoute string, queryPathParam string) (res []byte) {
	query := abci.RequestQuery{
		Path: path,
		Data: nil,
	}
	bz, err := querier(ctx, []string{querierRoute, queryPathParam}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	return bz
}
func TestList(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})

	cap := 1000
	for i := 0; i < cap; i++ {
		CoinIssueInfo.SetIssuer(sdk.AccAddress(crypto.AddressHash([]byte(utils.GetRandomString(10)))))
		CoinIssueInfo.SetSymbol(utils.GetRandomString(6))
		_, err := keeper.CreateIssue(ctx, &CoinIssueInfo)
		if err != nil {
			fmt.Println(err.Error())
		}
		require.Nil(t, err)
	}

	issueId := ""
	for i := 0; i < 100; i++ {
		//fmt.Println("==================page:" + strconv.Itoa(i))
		issues := keeper.List(ctx, params.IssueQueryParams{StartIssueId: issueId, Owner: nil, Limit: 10})
		require.Len(t, issues, 10)
		for j, issue := range issues {
			if j > 0 {
				require.True(t, issues[j].IssueTime <= (issues[j-1].IssueTime))
			}
			//fmt.Println(issue.IssueId + "----" + issue.IssueTime.String())
			issueId = issue.IssueId
		}
	}
}

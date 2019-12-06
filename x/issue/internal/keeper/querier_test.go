/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package keeper_test

import (
	"fmt"
	"testing"

	"github.com/tendermint/tendermint/crypto"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/xar-network/xar-network/x/issue"
	"github.com/xar-network/xar-network/x/issue/internal/types"
	"github.com/xar-network/xar-network/x/issue/internal/keeper"
)

func TestQueryIssue(t *testing.T) {
	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})

	k.GetSupplyKeeper().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	querier := keeper.NewQuerier(k)
	handler := issue.NewHandler(k)

	res := handler(ctx, types.NewMsgIssue(SenderAccAddr, &IssueParams))
	var issueID string
	issueID = string(res.Data)
	bz := getQueried(t, ctx, querier, types.GetQueryIssuePath(issueID), types.QueryIssue, issueID)
	var issueInfo types.CoinIssueInfo
	mapp.Cdc.MustUnmarshalJSON(bz, &issueInfo)

	require.Equal(t, issueInfo.GetIssueId(), issueID)
	require.Equal(t, issueInfo.GetName(), CoinIssueInfo.GetName())

}

func TestQueryIssues(t *testing.T) {
	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})
	k.GetSupplyKeeper().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	//querier := issue.NewQuerier(k)
	handler := issue.NewHandler(k)
	cap := 10
	for i := 0; i < cap; i++ {
		handler(ctx, types.NewMsgIssue(SenderAccAddr, &IssueParams))
	}
	issues := k.List(ctx, types.IssueQueryParams{Limit: 10})
	require.Len(t, issues, cap)

}

func TestSearchIssues(t *testing.T) {
	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.NewContext(false, abci.Header{})
	k.GetSupplyKeeper().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	querier := keeper.NewQuerier(k)
	handler := issue.NewHandler(k)
	cap := 10
	for i := 0; i < cap; i++ {
		handler(ctx, types.NewMsgIssue(SenderAccAddr, &IssueParams))
	}
	bz := getQueried(t, ctx, querier, types.GetQueryIssuePath("TES"), types.QuerySearch, "TES")
	var issues types.CoinIssues
	mapp.Cdc.MustUnmarshalJSON(bz, &issues)
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
	mapp, k, _, _, _, _ := getMockApp(t, issue.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})
	k.GetSupplyKeeper().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

	cap := 1000
	for i := 0; i < cap; i++ {
		CoinIssueInfo.SetIssuer(sdk.AccAddress(crypto.AddressHash([]byte(types.GetRandomString(10)))))
		CoinIssueInfo.SetSymbol(types.GetRandomString(6))
		_, err := k.CreateIssue(ctx, &CoinIssueInfo)
		if err != nil {
			fmt.Println(err.Error())
		}
		require.Nil(t, err)
	}

	issueId := ""
	for i := 0; i < 100; i++ {
		//fmt.Println("==================page:" + strconv.Itoa(i))
		issues := k.List(ctx, types.IssueQueryParams{StartIssueId: issueId, Owner: nil, Limit: 10})
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

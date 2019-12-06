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

package issue_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
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
	keeper.GetSupplyKeeper().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	mapp.InitChainer(ctx, abci.RequestInitChain{})

	handler := issue.NewHandler(keeper)

	res := handler(ctx, types.NewMsgIssue(SenderAccAddr, &IssueParams))
	require.True(t, res.IsOK())

	var issueID string
	issueID = string(res.Data)
	require.NotNil(t, issueID)
}

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

package tests

import (
	"testing"

	"github.com/xar-network/xar-network/x/record"

	"github.com/xar-network/xar-network/x/record/internal/types"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestImportExportQueues(t *testing.T) {
	mapp, keeper, _, _, _ := getMockApp(t, record.DefaultGenesisState(), nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	handler := record.NewHandler(keeper)

	// record hash 1
	res := handler(ctx, types.NewMsgRecord(SenderAccAddr, &RecordParams))
	require.True(t, res.IsOK())

	// record hash 2
	RecordParams.Hash = "BC38CAEE32149BEF4CCFAEAB518EC9A5FBC85AE6AC8D5A9F6CD710FAF5E4A2B9"
	res2 := handler(ctx, types.NewMsgRecord(SenderAccAddr, &RecordParams))
	require.True(t, res2.IsOK())
}

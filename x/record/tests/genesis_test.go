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

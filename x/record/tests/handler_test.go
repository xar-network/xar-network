package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/record"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

func TestHandlerNewMsgRecord(t *testing.T) {
	mapp, keeper, _, _, _ := getMockApp(t, record.GenesisState{}, nil)
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.NewContext(false, abci.Header{})
	mapp.InitChainer(ctx, abci.RequestInitChain{})

	handler := record.NewHandler(keeper)

	msg := types.NewMsgRecord(SenderAccAddr, &RecordParams)
	err := msg.ValidateBasic()
	require.Nil(t, err)

	res := handler(ctx, msg)
	require.True(t, res.IsOK())
}

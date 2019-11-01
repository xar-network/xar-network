package tests

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/record/params"

	"github.com/xar-network/xar-network/x/record"
)

func TestCreateRecord(t *testing.T) {
	mapp, keeper, _, _, _ := getMockApp(t, record.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	err := keeper.CreateRecord(ctx, &RecordInfo)
	require.Nil(t, err)
	recordRes := keeper.GetRecord(ctx, RecordInfo.Hash)
	require.Equal(t, recordRes.Hash, RecordInfo.Hash)
	require.Equal(t, recordRes.Name, RecordInfo.Name)
	require.Equal(t, recordRes.RecordType, RecordInfo.RecordType)
	require.Equal(t, recordRes.Author, RecordInfo.Author)
	require.Equal(t, recordRes.RecordNo, RecordInfo.RecordNo)
}

func TestCreateRecordDuplicated(t *testing.T) {
	mapp, keeper, _, _, _ := getMockApp(t, record.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	err := keeper.CreateRecord(ctx, &RecordInfo)
	require.Nil(t, err)
	err2 := keeper.CreateRecord(ctx, &RecordInfo)
	require.NotNil(t, err2)
}

func TestGetRecords(t *testing.T) {
	mapp, keeper, _, _, _ := getMockApp(t, record.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	cap := 10
	for i := 0; i < cap; i++ {
		RecordInfo.Hash = RecordInfo.Hash[0:len(RecordInfo.Hash)-2] + strconv.Itoa(i)
		err := keeper.CreateRecord(ctx, &RecordInfo)
		require.Nil(t, err)
	}
	records := keeper.List(ctx, params.RecordQueryParams{
		Sender: RecordInfo.Sender,
		Limit:  99999999,
	})

	require.Len(t, records, cap)
}

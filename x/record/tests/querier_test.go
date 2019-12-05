package tests

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/record"
	"github.com/xar-network/xar-network/x/record/client/queriers"
	"github.com/xar-network/xar-network/x/record/internal/keeper"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

func TestQueryRecords(t *testing.T) {
	mapp, k, _, _, _ := getMockApp(t, record.GenesisState{}, nil)

	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})

	ctx := mapp.NewContext(false, abci.Header{})
	//querier := record.NewQuerier(k)
	handler := record.NewHandler(k)
	cap := 10
	for i := 0; i < cap; i++ {
		RecordParams.Hash = RecordParams.Hash[0:len(RecordParams.Hash)-1] + strconv.Itoa(i)
		handler(ctx, types.NewMsgRecord(SenderAccAddr, &RecordParams))
	}

	// query all
	queryParams := k.Getcdc().MustMarshalJSON(RecordQueryParams)
	bz := getQueried(t, ctx, keeper.NewQuerier(k), queriers.GetQueryRecordsPath(), types.QueryRecords, "", queryParams)
	var records []*types.RecordInfo
	k.Getcdc().MustUnmarshalJSON(bz, &records)
	require.Len(t, records, cap)
	require.Equal(t, "rec174876e800", records[len(records)-1].ID)

	// query by sender
	RecordQueryParams.Sender = SenderAccAddr
	queryParams2 := k.Getcdc().MustMarshalJSON(RecordQueryParams)
	bz2 := getQueried(t, ctx, keeper.NewQuerier(k), queriers.GetQueryRecordsPath(), types.QueryRecords, "", queryParams2)
	var records2 []*types.RecordInfo
	k.Getcdc().MustUnmarshalJSON(bz2, &records2)
	require.Len(t, records2, cap)
	require.Equal(t, "rec174876e800", records2[len(records)-1].ID)

	// query with start id and sender
	RecordQueryParams.StartRecordId = "rec174876e805"
	queryParams3 := k.Getcdc().MustMarshalJSON(RecordQueryParams)
	bz3 := getQueried(t, ctx, keeper.NewQuerier(k), queriers.GetQueryRecordsPath(), types.QueryRecords, "", queryParams3)
	var records3 []*types.RecordInfo
	k.Getcdc().MustUnmarshalJSON(bz3, &records3)
	require.Len(t, records3, 5)

	// query with start id
	RecordQueryParams.Sender = nil
	queryParams4 := k.Getcdc().MustMarshalJSON(RecordQueryParams)
	bz4 := getQueried(t, ctx, keeper.NewQuerier(k), queriers.GetQueryRecordsPath(), types.QueryRecords, "", queryParams4)
	var records4 []*types.RecordInfo
	k.Getcdc().MustUnmarshalJSON(bz4, &records4)
	require.Len(t, records4, 5)
}

func getQueried(t *testing.T, ctx sdk.Context, querier sdk.Querier, path string, querierRoute string, queryPathParam string, queryParam []byte) (res []byte) {
	query := abci.RequestQuery{
		Path: path,
		Data: queryParam,
	}
	bz, err := querier(ctx, []string{querierRoute, queryPathParam}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	return bz
}

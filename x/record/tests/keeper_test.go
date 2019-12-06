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
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/record/internal/types"

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
	records := keeper.List(ctx, types.RecordQueryParams{
		Sender: RecordInfo.Sender,
		Limit:  99999999,
	})

	require.Len(t, records, cap)
}

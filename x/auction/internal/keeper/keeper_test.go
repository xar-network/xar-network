/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
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
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/auction/internal/keeper"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

func TestKeeper_SetGetDeleteAuction(t *testing.T) {
	// setup k, create auction
	mapp, k, addresses, _ := setUpMockApp()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header}) // Without this it panics about "invalid memory address or nil pointer dereference"
	ctx := mapp.BaseApp.NewContext(false, header)
	auction, _ := types.NewForwardAuction(addresses[0], sdk.NewInt64Coin("csdt", 100), sdk.NewInt64Coin("ftm", 0), types.EndTime(1000))
	id := types.ID(5)
	auction.SetID(id)

	// write and read from store
	k.SetAuction(ctx, &auction)
	readAuction, found := k.GetAuction(ctx, id)

	// check before and after match
	require.True(t, found)
	require.Equal(t, &auction, readAuction)
	t.Log(auction)
	t.Log(readAuction.GetID())
	// check auction is in queue
	iter := k.GetQueueIterator(ctx, 100000)
	require.Equal(t, 1, len(convertIteratorToSlice(mapp.Cdc, k, iter)))
	iter.Close()

	// delete auction
	k.DeleteAuction(ctx, id)

	// check auction does not exist
	_, found = k.GetAuction(ctx, id)
	require.False(t, found)
	// check auction not in queue
	iter = k.GetQueueIterator(ctx, 100000)
	require.Equal(t, 0, len(convertIteratorToSlice(mapp.Cdc, k, iter)))
	iter.Close()

}

// TODO convert to table driven test with more test cases
func TestKeeper_ExpiredAuctionQueue(t *testing.T) {
	// setup k
	mapp, k, _, _ := setUpMockApp()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	// create an example queue
	type queue []struct {
		endTime   types.EndTime
		auctionID types.ID
	}
	q := queue{{1000, 0}, {1300, 2}, {5200, 1}}

	// write and read queue
	for _, v := range q {
		k.InsertIntoQueue(ctx, v.endTime, v.auctionID)
	}
	iter := k.GetQueueIterator(ctx, 1000)

	// check before and after match
	i := 0
	for ; iter.Valid(); iter.Next() {
		var auctionID types.ID
		mapp.Cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &auctionID)
		require.Equal(t, q[i].auctionID, auctionID)
		i++
	}

}

func convertIteratorToSlice(cdc *amino.Codec, k keeper.Keeper, iterator sdk.Iterator) []types.ID {
	var queue []types.ID
	for ; iterator.Valid(); iterator.Next() {
		var auctionID types.ID
		cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &auctionID)
		queue = append(queue, auctionID)
	}
	return queue
}

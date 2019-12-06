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

package auction_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/auction"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

func TestKeeper_EndBlocker(t *testing.T) {
	// setup keeper and auction
	mapp, keeper, addresses, _ := setUpMockApp()
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)

	params := types.DefaultAuctionParams()
	params.MaxBidDuration = 3 * 1
	params.MaxAuctionDuration = 2 * 24 * 1
	keeper.SetParams(ctx, params)

	seller := addresses[0]
	keeper.StartForwardAuction(ctx, seller, sdk.NewInt64Coin("token1", 20), sdk.NewInt64Coin("token2", 0))

	// run the endblocker, simulating a block height after auction expiry

	expiryBlock := ctx.BlockHeight() + int64(params.MaxAuctionDuration)
	auction.EndBlocker(ctx.WithBlockHeight(expiryBlock), keeper)

	// check auction has been closed
	_, found := keeper.GetAuction(ctx, 0)
	require.False(t, found)
}

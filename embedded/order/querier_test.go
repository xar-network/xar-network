/*

Copyright 2019 All in Bits, Inc
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

package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/xar-network/xar-network/testutil"
	"github.com/xar-network/xar-network/testutil/testflags"
	"github.com/xar-network/xar-network/types"
	"github.com/xar-network/xar-network/types/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestQuerier(t *testing.T) {
	testflags.UnitTest(t)
	cdc := codec.New()
	db := dbm.NewMemDB()
	k := NewKeeper(db, cdc)
	ctx := testutil.DummyContext()
	q := NewQuerier(k)
	doListQuery := func(req ListQueryRequest) (ListQueryResult, error) {
		reqB := serializeRequestQuery(cdc, req)
		var res ListQueryResult
		resB, err := q(ctx, []string{"list"}, reqB)
		if err != nil {
			return res, err
		}
		cdc.MustUnmarshalJSON(resB, &res)
		return res, nil
	}

	t.Run("should return no more than 50 orders in descending order", func(t *testing.T) {
		id := store.NewEntityID(0)

		for i := 0; i < 55; i++ {
			id = id.Inc()
			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: store.NewEntityID(2),
				ID:       id,
			}))
		}

		res, err := doListQuery(ListQueryRequest{})
		require.NoError(t, err)

		assert.Equal(t, 50, len(res.Orders))
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(55), res.Orders[0].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(6), res.Orders[49].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(5), res.NextID)
	})
	t.Run("should work with an offset", func(t *testing.T) {
		id := store.NewEntityID(0)

		for i := 0; i < 55; i++ {
			id = id.Inc()
			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: store.NewEntityID(2),
				ID:       id,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Start: store.NewEntityID(7),
		})
		require.NoError(t, err)

		assert.Equal(t, 7, len(res.Orders))
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(7), res.Orders[0].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(1), res.Orders[6].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(0), res.NextID)
	})
	t.Run("should support filter by address alongside offset", func(t *testing.T) {
		id := store.NewEntityID(0)
		genOwner := testutil.RandAddr()
		for i := 0; i < 110; i++ {
			id = id.Inc()
			var owner sdk.AccAddress
			if i%2 == 0 {
				owner = genOwner
			}

			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: store.NewEntityID(2),
				ID:       id,
				Owner:    owner,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Start: store.NewEntityID(104),
			Owner: genOwner,
		})
		require.NoError(t, err)

		assert.Equal(t, 50, len(res.Orders))
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(109), res.Orders[0].ID)
		testutil.AssertEqualEntityIDs(t, store.NewEntityID(11), res.Orders[49].ID)
	})
	t.Run("should support limit value", func(t *testing.T) {
		id := store.NewEntityID(0)
		genOwner := testutil.RandAddr()
		for i := 0; i < 110; i++ {
			id = id.Inc()
			var owner sdk.AccAddress
			if i%2 == 0 {
				owner = genOwner
			}

			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: store.NewEntityID(2),
				ID:       id,
				Owner:    owner,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Start: store.NewEntityID(104),
			Owner: genOwner,
			Limit: 23,
		})
		require.NoError(t, err)

		assert.Equal(t, 23, len(res.Orders))

		res, err = doListQuery(ListQueryRequest{
			Start: store.NewEntityID(104),
			Owner: genOwner,
			Limit: 7,
		})
		require.NoError(t, err)

		assert.Equal(t, 7, len(res.Orders))
	})
	t.Run("should support filter by market id", func(t *testing.T) {
		id := store.NewEntityID(0)
		genOwner := testutil.RandAddr()
		for i := 0; i < 110; i++ {
			id = id.Inc()
			var owner sdk.AccAddress
			if i%2 == 0 {
				owner = genOwner
			}

			var market store.EntityID
			if i%2 == 0 {
				market = store.NewEntityID(2)
			} else {
				market = store.NewEntityID(1)
			}

			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: market,
				ID:       id,
				Owner:    owner,
			}))
		}
		res, err := doListQuery(ListQueryRequest{
			MarketID: []store.EntityID{store.NewEntityID(1)},
			Limit: 1000,
		})
		require.NoError(t, err)

		assert.Greater(t, len(res.Orders), 0)
		for _, ord := range res.Orders {
			assert.Equal(t, ord.MarketID, store.NewEntityID(1))
		}

		res, err = doListQuery(ListQueryRequest{
			MarketID: []store.EntityID{store.NewEntityID(2)},
			Limit: 1000,
		})
		require.NoError(t, err)

		assert.Greater(t, len(res.Orders), 0)
		for _, ord := range res.Orders {
			assert.Equal(t, ord.MarketID, store.NewEntityID(2))
		}
	})
	t.Run("should support filter by status", func(t *testing.T) {
		id := store.NewEntityID(0)
		genOwner := testutil.RandAddr()
		for i := 0; i < 110; i++ {
			id = id.Inc()
			var owner sdk.AccAddress
			if i%2 == 0 {
				owner = genOwner
			}

			var market store.EntityID
			if i%2 == 0 {
				market = store.NewEntityID(2)
			} else {
				market = store.NewEntityID(1)
			}

			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: market,
				ID:       id,
				Owner:    owner,
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Status: []string{"OPEN"},
			Limit: 1000,
		})
		require.NoError(t, err)

		assert.Greater(t, len(res.Orders), 0)
		for _, ord := range res.Orders {
			assert.Equal(t, ord.Status, "OPEN")
		}

		res, err = doListQuery(ListQueryRequest{
			Status: []string{"CLOSE"},
			Limit: 1000,
		})
		require.NoError(t, err)

		assert.Equal(t, len(res.Orders), 0)
	})
	t.Run("should support filter by create time", func(t *testing.T) {
		id := store.NewEntityID(0)
		genOwner := testutil.RandAddr()
		for i := 0; i < 110; i++ {
			id = id.Inc()
			var owner sdk.AccAddress
			if i%2 == 0 {
				owner = genOwner
			}

			var market store.EntityID
			if i%2 == 0 {
				market = store.NewEntityID(2)
			} else {
				market = store.NewEntityID(1)
			}

			require.NoError(t, k.OnEvent(types.OrderCreated{
				MarketID: market,
				ID:       id,
				Owner:    owner,
				CreatedTime: int64(i),
			}))
		}

		res, err := doListQuery(ListQueryRequest{
			Limit: 1000,
			UnixTimeAfter: 0,
			UnixTimeBefore: 1000,
		})
		require.NoError(t, err)

		assert.Equal(t, len(res.Orders), 110)

		res, err = doListQuery(ListQueryRequest{
			Limit: 1000,
			UnixTimeAfter: 100,
			UnixTimeBefore: 100,
		})
		require.NoError(t, err)

		assert.Equal(t, len(res.Orders), 1)
		assert.Equal(t, res.Orders[0].CreatedTime, int64(100))

		res, err = doListQuery(ListQueryRequest{
			Limit: 1000,
			UnixTimeAfter: 51,
			UnixTimeBefore: 100,
		})
		require.NoError(t, err)

		assert.Equal(t, len(res.Orders), 50)
		for _, ord := range res.Orders {
			assert.GreaterOrEqual(t, ord.CreatedTime, int64(51))
			assert.LessOrEqual(t, ord.CreatedTime, int64(100))
		}
	})

	t.Run("should return an error if the request does not deserialize", func(t *testing.T) {
		_, err := q(ctx, []string{"list"}, abci.RequestQuery{Data: []byte("foo")})
		require.Error(t, err)
	})
}

func serializeRequestQuery(cdc *codec.Codec, req ListQueryRequest) abci.RequestQuery {
	data := cdc.MustMarshalBinaryBare(req)

	return abci.RequestQuery{
		Data: data,
	}
}

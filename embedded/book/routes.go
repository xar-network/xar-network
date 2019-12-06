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

package book

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/xar-network/xar-network/embedded"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/xar-network/xar-network/embedded/node"
	"github.com/xar-network/xar-network/embedded/order"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.Handle("/markets/{marketID}/book", bookHandler(ctx, cdc)).Methods("GET")
}

func bookHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mktId := vars["marketID"]

		resJSON, _, err := ctx.QueryWithData(fmt.Sprintf("custom/book/get/%s", mktId), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(resJSON) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no spread at this block")
			return
		}

		block, err := node.LatestBlock(ctx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch latest block")
			return
		}

		var orders []order.Order
		cdc.MustUnmarshalJSON(resJSON, &orders)

		qRes := QueryResult{
			MarketID:    store.NewEntityIDFromString(mktId),
			BlockNumber: block.Block.Height,
			Bids:        make([]QueryResultEntry, 0),
			Asks:        make([]QueryResultEntry, 0),
		}

		bidPrices := make(map[string]QueryResultEntry)
		askPrices := make(map[string]QueryResultEntry)

		for _, o := range orders {
			var m map[string]QueryResultEntry
			if o.Direction == matcheng.Bid {
				m = bidPrices
			} else {
				m = askPrices
			}

			entry, ok := m[o.Price.String()]
			if ok {
				entry.Quantity = entry.Quantity.Add(o.Quantity.Sub(o.QuantityFilled))
				m[o.Price.String()] = entry
			} else {
				entry = QueryResultEntry{
					Price:    o.Price,
					Quantity: o.Quantity.Sub(o.QuantityFilled),
				}
				m[o.Price.String()] = entry
			}
		}

		for _, entry := range bidPrices {
			qRes.Bids = append(qRes.Bids, entry)
		}
		for _, entry := range askPrices {
			qRes.Asks = append(qRes.Asks, entry)
		}

		sort.Slice(qRes.Bids, func(i, j int) bool {
			return qRes.Bids[i].Price.LT(qRes.Bids[j].Price)
		})
		sort.Slice(qRes.Asks, func(i, j int) bool {
			return qRes.Asks[i].Price.LT(qRes.Asks[j].Price)
		})

		embedded.PostProcessResponse(w, ctx, qRes)
	}
}

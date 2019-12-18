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
	"net/http"
	"net/url"
	"strconv"

	"github.com/xar-network/xar-network/embedded"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/xar-network/xar-network/embedded/auth"
	"github.com/xar-network/xar-network/types/store"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.Handle("/user/orders", auth.DefaultAuthMW(getOrdersHandler(ctx, cdc))).Methods("GET")
}

func getOrdersHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := auth.MustGetKBFromSession(r)
		q := r.URL.Query()

		unixTimeAfter, unixTimeBefore := getTimeLimitFromQuery(q)
		req := ListQueryRequest{
			Owner: owner.GetAddr(),
			Limit: getLimitFromQuery(q),
			MarketID: getMarketIDFromQuery(q),
			Status: getStatusFromQuery(q),
			UnixTimeAfter: unixTimeAfter,
			UnixTimeBefore: unixTimeBefore,
		}
		if start, ok := q["start"]; ok {
			req.Start = store.NewEntityIDFromString(start[0])
		}

		resB, _, err := ctx.QueryWithData("custom/embeddedorder/list", cdc.MustMarshalBinaryBare(req))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		embedded.PostProcessResponse(w, ctx, resB)
	}
}

func getLimitFromQuery(q url.Values) int {
	return int(getInt64FromQuery(q, "limit"))
}

func getMarketIDFromQuery(q url.Values) []store.EntityID {
	marketIDs := make([]store.EntityID, 0)

	if marketIdList, ok := q["market_ids"]; ok {
		for _, marketID := range marketIdList {
			marketIDs = append(marketIDs, store.NewEntityIDFromString(marketID))
		}
	}

	return marketIDs
}

func getStatusFromQuery(q url.Values) []string {
	statuses := make([]string, 0)

	if statusesList, ok := q["statuses"]; ok {
		for _, status := range statusesList {
			statuses = append(statuses, status)
		}
	}

	return statuses
}

func getTimeLimitFromQuery(q url.Values) (int64, int64) {
	after := getInt64FromQuery(q, "after")
	before := getInt64FromQuery(q, "before")

	if after != 0 && before != 0 && after > before {
		return 0, 0
	}

	return after, before
}

func getInt64FromQuery(q url.Values, name string) int64 {
	tm := int64(0)

	tmStrList, ok := q[name]
	if ok && len(tmStrList) != 0 && tmStrList[0] != "" {
		var err error
		tm, err = strconv.ParseInt(tmStrList[0], 10, 64)
		if err != nil {
			return 0
		}
	}

	return tm
}

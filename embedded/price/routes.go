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

package price

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xar-network/xar-network/embedded"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/xar-network/xar-network/pkg/conv"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.Handle("/markets/{marketID}/candles", candlesHandler(ctx, cdc)).Methods("GET")
	r.Handle("/markets/{marketID}/daily", dailyHandler(ctx, cdc)).Methods("GET")
}

func candlesHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mktID := vars["marketID"]
		now := time.Now()
		params := CandleQueryParams{
			From:     now.AddDate(0, 0, -1),
			To:       now,
			Interval: CandleInterval60M,
		}

		q := r.URL.Query()
		if start, ok := q["start"]; ok {
			startDate, err := conv.ParseISO8601(start[0])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid start date")
				return
			}
			params.From = startDate
		}
		if end, ok := q["end"]; ok {
			endDate, err := conv.ParseISO8601(end[0])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid end date")
				return
			}
			params.To = endDate
		}
		if granularity, ok := q["granularity"]; ok {
			cInterval, err := NewCandleIntervalFromString(granularity[0])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid granularity")
				return
			}
			params.Interval = cInterval
		}

		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}

		res, height, err := ctx.QueryWithData(fmt.Sprintf("custom/price/candles/%s", mktID), cdc.MustMarshalBinaryBare(params))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx = ctx.WithHeight(height)

		embedded.PostProcessResponse(w, ctx, res)
	}
}

func dailyHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		mktID := vars["marketID"]
		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}
		res, height, err := ctx.QueryWithData(fmt.Sprintf("custom/price/daily/%s", mktID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx = ctx.WithHeight(height)

		embedded.PostProcessResponse(w, ctx, res)
	}
}

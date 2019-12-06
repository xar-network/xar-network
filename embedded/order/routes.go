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

		req := ListQueryRequest{
			Owner: owner.GetAddr(),
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

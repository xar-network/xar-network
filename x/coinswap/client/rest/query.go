/*

Copyright 2016 All in Bits, Inc
Copyright 2017 IRIS Foundation Ltd.
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

package rest

import (
	"net/http"

	"github.com/cosmos-sdk/cosmos/client/context"
	"github.com/cosmos-sdk/cosmos/codec"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Query liquidity
	r.HandleFunc(
		"/coinswap/liquidities/{id}",
		queryLiquidityHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

// queryLiquidityHandlerFn performs liquidity information query
func queryLiquidityHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return queryLiquidity(cliCtx, cdc, "custom/coinswap/liquidities/{id}")
}

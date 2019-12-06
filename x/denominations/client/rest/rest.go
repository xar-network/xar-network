/*

Copyright 2016 All in Bits, Inc
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
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	restName = "token"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	// Queries
	r.HandleFunc(fmt.Sprintf("/%s/tokens", storeName), symbolsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/tokens/{%s}", storeName, restName), findTokenHandler(cliCtx, storeName)).Methods("GET")

	// Transactions
	r.HandleFunc(fmt.Sprintf("/%s/tokens", storeName), issueTokenHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/tokens/mint", storeName), mintHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/tokens/burn", storeName), burnHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/tokens/freeze", storeName), freezeHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/tokens/unfreeze", storeName), unfreezeHandler(cliCtx)).Methods("PUT")

}

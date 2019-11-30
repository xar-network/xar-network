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

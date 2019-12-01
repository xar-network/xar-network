package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	boxrest "github.com/hashgard/hashgard/x/box/client/rest"
	"github.com/hashgard/hashgard/x/box/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.Lock, types.QueryParams), boxrest.BoxQueryParamsHandlerFn(cdc, cliCtx, types.Lock)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.Lock, types.QueryBox, boxrest.ID), boxrest.BoxQueryHandlerFn(cdc, cliCtx, types.Lock)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.Lock, types.QuerySearch, boxrest.Name), boxrest.BoxSearchHandlerFn(cdc, cliCtx, types.Lock)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.Lock, types.QueryList), boxrest.BoxListHandlerFn(cdc, cliCtx, types.Lock)).Methods("GET")
}

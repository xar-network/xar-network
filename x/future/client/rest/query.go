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
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.Future, types.QueryParams), boxrest.BoxQueryParamsHandlerFn(cdc, cliCtx, types.Future)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.Future, types.QueryBox, boxrest.ID), boxrest.BoxQueryHandlerFn(cdc, cliCtx, types.Future)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.Future, types.QuerySearch, boxrest.Name), boxrest.BoxSearchHandlerFn(cdc, cliCtx, types.Future)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.Future, types.QueryList), boxrest.BoxListHandlerFn(cdc, cliCtx, types.Future)).Methods("GET")
}

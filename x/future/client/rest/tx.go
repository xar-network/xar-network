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
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/%s", types.Future), boxrest.PostFutureBoxCreateHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}", types.Future, types.Inject, boxrest.ID, boxrest.Amount), boxrest.PostInjectHandlerFn(cdc, cliCtx, types.Future, types.Inject)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}", types.Future, types.Cancel, boxrest.ID, boxrest.Amount), boxrest.PostInjectHandlerFn(cdc, cliCtx, types.Future, types.Cancel)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.Future, types.Withdraw, boxrest.ID), boxrest.PostWithdrawHandlerFn(cdc, cliCtx, types.Future)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}/{%s}", types.Deposit, types.Feature, types.Disable, boxrest.ID, boxrest.Feature), boxrest.PostDisableFeatureHandlerFn(cdc, cliCtx, types.Future)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/description/{%s}", types.Future, boxrest.ID), boxrest.PostDescribeHandlerFn(cdc, cliCtx, types.Future)).Methods("POST")
}

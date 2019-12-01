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
	r.HandleFunc(fmt.Sprintf("/%s", types.Deposit), boxrest.PostDepositBoxCreateHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}/{%s}", types.Deposit, types.DepositBoxInterest, types.Inject, boxrest.ID, boxrest.Amount), boxrest.PostInterestHandlerFn(cdc, cliCtx, types.Inject)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}/{%s}", types.Deposit, types.DepositBoxInterest, types.Cancel, boxrest.ID, boxrest.Amount), boxrest.PostInterestHandlerFn(cdc, cliCtx, types.Cancel)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}", types.Deposit, types.Inject, boxrest.ID, boxrest.Amount), boxrest.PostInjectHandlerFn(cdc, cliCtx, types.Deposit, types.Inject)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}", types.Deposit, types.Cancel, boxrest.ID, boxrest.Amount), boxrest.PostInjectHandlerFn(cdc, cliCtx, types.Deposit, types.Cancel)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.Deposit, types.Withdraw, boxrest.ID), boxrest.PostWithdrawHandlerFn(cdc, cliCtx, types.Deposit)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/%s/%s/{%s}/{%s}", types.Deposit, types.Feature, types.Disable, boxrest.ID, boxrest.Feature), boxrest.PostDisableFeatureHandlerFn(cdc, cliCtx, types.Deposit)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/description/{%s}", types.Deposit, boxrest.ID), boxrest.PostDescribeHandlerFn(cdc, cliCtx, types.Deposit)).Methods("POST")
}

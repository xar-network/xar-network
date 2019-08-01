package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/zar-network/zar-network/x/issue/params"
	"github.com/zar-network/zar-network/x/issue/types"

	"github.com/zar-network/zar-network/x/issue/client/queriers"
	issueutils "github.com/zar-network/zar-network/x/issue/utils"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.QuerierRoute, types.QueryParams), queryParamsHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QueryIssue, IssueID), queryIssueHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.QuerierRoute, types.QueryIssues), queryIssuesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QuerySearch, Symbol), queryIssueSearchHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}", types.QuerierRoute, types.QueryFreeze, IssueID, restAddress), queryIssueFreezeHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QueryFreezes, IssueID), queryIssueFreezesHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}/{%s}", types.QuerierRoute, types.QueryAllowance, IssueID, restAddress, spenderAddress), queryIssueAllowanceHandlerFn(cdc, cliCtx)).Methods("GET")

}
func queryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queriers.QueryParams(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func queryIssueHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := queriers.QueryIssueByID(issueID, cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func queryIssueSearchHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars[Symbol]

		res, err := queriers.QueryIssueBySymbol(symbol, cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func queryIssuesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		address, err := sdk.AccAddressFromBech32(r.URL.Query().Get(restAddress))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		issueQueryParams := params.IssueQueryParams{
			StartIssueId: r.URL.Query().Get(restStartIssueId),
			Owner:        address,
			Limit:        30,
		}
		strNumLimit := r.URL.Query().Get(restLimit)
		if len(strNumLimit) > 0 {
			limit, err := strconv.Atoi(r.URL.Query().Get(restLimit))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			issueQueryParams.Limit = limit
		}

		res, err := queriers.QueryIssuesList(issueQueryParams, cdc, cliCtx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func queryIssueFreezeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		address, err := sdk.AccAddressFromBech32(vars[restAddress])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, err := queriers.QueryIssueFreeze(issueID, address, cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func queryIssueFreezesHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := queriers.QueryIssueFreezes(issueID, cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func queryIssueAllowanceHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		address, err := sdk.AccAddressFromBech32(vars[restAddress])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		spenderAddress, err := sdk.AccAddressFromBech32(vars[spenderAddress])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, err := queriers.QueryIssueAllowance(issueID, address, spenderAddress, cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

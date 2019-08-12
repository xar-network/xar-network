package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"
	"github.com/zar-network/zar-network/x/issue/types"
)

const (
	IssueID    = "issue-id"
	Feature    = "feature"
	AccAddress = "accAddress"
	From       = "from"
	FreezeType = "freeze-type"
	EndTime    = "end-time"
	Symbol     = "symbol"
	Amount     = "amount"
	To         = "to"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type PostIssueReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	types.IssueParams `json:"issue"`
}
type PostDescriptionReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Description string       `json:"description"`
}
type PostIssueBaseReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/issue", postIssueHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postIssueApproveHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/increase/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postIssueIncreaseApproval(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/decrease/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postIssueDecreaseApproval(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/burn/{%s}/{%s}", IssueID, Amount), postBurnHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/burn-from/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postBurnFromHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/freeze/{%s}/{%s}/{%s}/{%s}", FreezeType, IssueID, AccAddress, EndTime), postIssueFreezeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/unfreeze/{%s}/{%s}/{%s}", FreezeType, IssueID, AccAddress), postIssueUnFreezeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/send-from/{%s}/{%s}/{%s}/{%s}", IssueID, From, To, Amount), postIssueSendFrom(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/mint/{%s}/{%s}/{%s}", IssueID, Amount, To), postMintHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/ownership/transfer/{%s}/{%s}", IssueID, To), postTransferOwnershipHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/description/{%s}", IssueID), postDescribeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/feature/disable/{%s}/{%s}", IssueID, Feature), postDisableFeatureHandlerFn(cliCtx)).Methods("POST")
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.QuerierRoute, types.QueryParams), queryParamsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QueryIssue, IssueID), queryIssueHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.QuerierRoute, types.QueryIssues), queryIssuesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QuerySearch, Symbol), queryIssueSearchHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}", types.QuerierRoute, types.QueryFreeze, IssueID, restAddress), queryIssueFreezeHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QueryFreezes, IssueID), queryIssueFreezesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}/{%s}", types.QuerierRoute, types.QueryAllowance, IssueID, restAddress, spenderAddress), queryIssueAllowanceHandlerFn(cliCtx)).Methods("GET")
}

func postIssueHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(req.Description) > 0 && !json.Valid([]byte(req.Description)) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrCoinDescriptionNotValid().Error())
			return
		}
		// create the message
		msg := types.NewMsgIssue(fromAddr, &req.IssueParams)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postMintHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]

		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		num, err := strconv.ParseInt(vars[Amount], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount := sdk.NewInt(num)
		toAddr, err := sdk.AccAddressFromBech32(vars[To])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		account, err := cliCtx.GetAccount(fromAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		issueInfo, err := clientutils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIssueMint(issueID, fromAddress, amount, issueInfo.GetDecimals(), to)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postDisableFeatureHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		feature := vars[Feature]
		_, ok := types.Features[feature]
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrUnknownFeatures().Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		account, err := cliCtx.GetAccount(fromAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = clientutils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgIssueDisableFeature(issueID, fromAddress, feature)

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}

}
func postDescribeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req PostDescriptionReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		if len(req.Description) <= 0 || !json.Valid([]byte(req.Description)) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrCoinDescriptionNotValid().Error())
			return
		}
		msg := types.NewMsgIssueDescription(issueID, fromAddress, []byte(req.Description))
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		account, err := cliCtx.GetAccount(fromAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = clientutils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
func postTransferOwnershipHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		to, err := sdk.AccAddressFromBech32(vars[To])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIssueTransferOwnership(issueID, fromAddress, to)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		account, err := cliCtx.GetAccount(fromAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = clientutils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postIssueFreezeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return issueFreezeHandlerFn(cdc, cliCtx, true)
}
func postIssueUnFreezeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return issueFreezeHandlerFn(cdc, cliCtx, false)
}
func issueFreezeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, freeze bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		account, err := cliCtx.GetAccount(fromAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		vars := mux.Vars(r)

		msg, err := clientutils.GetIssueFreezeMsg(cdc, cliCtx, account, vars[FreezeType], vars[IssueID], vars[AccAddress], vars[EndTime], freeze)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postBurnHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cdc, cliCtx, types.BurnHolder)
}
func postBurnFromHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cdc, cliCtx, types.BurnFrom)
}
func postIssueSendFrom(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		from, err := sdk.AccAddressFromBech32(vars[From])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		to, err := sdk.AccAddressFromBech32(vars[To])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		num, err := strconv.ParseInt(vars[Amount], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount := sdk.NewInt(num)

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		account, err := cliCtx.GetAccount(sender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err = clientutils.GetIssueByID(cdc, cliCtx, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := clientutils.CheckAllowance(cdc, cliCtx, issueID, from, account.GetAddress(), amount); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err = clientutils.CheckFreeze(cdc, cliCtx, issueID, from, to); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIssueSendFrom(issueID, sender, from, to, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postIssueApproveHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cdc, cliCtx, types.Approve)
}
func postIssueIncreaseApproval(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cdc, cliCtx, types.IncreaseApproval)
}
func postIssueDecreaseApproval(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cdc, cliCtx, types.DecreaseApproval)
}
func issueApproveHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, approveType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)

		issueID := vars[IssueID]
		accAddress, err := sdk.AccAddressFromBech32(vars[AccAddress])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount, ok := sdk.NewIntFromString(vars[Amount])
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Amount not a valid int")
			return
		}

		account, err := cliCtx.GetAccount(fromAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, err := clientutils.GetIssueApproveMsg(cdc, cliCtx, issueID, account, accAddress, approveType, amount, false)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

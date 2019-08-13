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
	cutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	"github.com/zar-network/zar-network/x/issue/client/utils"
	"github.com/zar-network/zar-network/x/issue/internal/types"
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
	r.HandleFunc(fmt.Sprintf("/issue/approve/{%s}/{%s}/{%s}", "issue-id", "accAddress", "amount"), postIssueApproveHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/increase/{%s}/{%s}/{%s}", "issue-id", "accAddress", "amount"), postIssueIncreaseApproval(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/decrease/{%s}/{%s}/{%s}", "issue-id", "accAddress", "amount"), postIssueDecreaseApproval(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/burn/{%s}/{%s}", "issue-id", "amount"), postBurnHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/burn-from/{%s}/{%s}/{%s}", "issue-id", "accAddress", "amount"), postBurnFromHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/freeze/{%s}/{%s}/{%s}", "freeze-type", "issue-id", "accAddress"), postIssueFreezeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/unfreeze/{%s}/{%s}/{%s}", "freeze-type", "issue-id", "accAddress"), postIssueUnFreezeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/send-from/{%s}/{%s}/{%s}/{%s}", "issue-id", "from", "to", "amount"), postIssueSendFrom(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/mint/{%s}/{%s}/{%s}", "issue-id", "amount", "to"), postMintHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/ownership/transfer/{%s}/{%s}", "issue-id", "to"), postTransferOwnershipHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/description/{%s}", "issue-id"), postDescribeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/feature/disable/{%s}/{%s}", "issue-id", "feature"), postDisableFeatureHandlerFn(cliCtx)).Methods("POST")
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.QuerierRoute, types.QueryParams), queryParamsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QueryIssue, "issue-id"), queryIssueHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s", types.QuerierRoute, types.QueryIssues), queryIssuesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QuerySearch, "symbol"), queryIssueSearchHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}", types.QuerierRoute, types.QueryFreeze, "issue-id", restAddress), queryIssueFreezeHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}", types.QuerierRoute, types.QueryFreezes, "issue-id"), queryIssueFreezesHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}/{%s}", types.QuerierRoute, types.QueryAllowance, "issue-id", restAddress, spenderAddress), queryIssueAllowanceHandlerFn(cliCtx)).Methods("GET")
}

func postIssueHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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

		cutils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postMintHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]

		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		num, err := strconv.ParseInt(vars["amount"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount := sdk.NewInt(num)
		toAddr, err := sdk.AccAddressFromBech32(vars["to"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}cliCtx.
		account, err := cliCtx.GetAccount(fromAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		issueInfo, err := utils.IssueOwnerCheck(cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIssueMint(issueID, fromAddr, toAddr, amount, issueInfo.GetDecimals())
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cutils.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postDisableFeatureHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]

		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		feature := vars["feature"]
		_, ok := types.Features[feature]
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrUnknownFeatures().Error())
			return
		}
		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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
		_, err = utils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
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
		issueID := vars["issue-id"]
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
		_, err = utils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
func postTransferOwnershipHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]
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
		to, err := sdk.AccAddressFromBech32(vars["to"])
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
		_, err = utils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postIssueFreezeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return issueFreezeHandlerFn(cliCtx, true)
}
func postIssueUnFreezeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return issueFreezeHandlerFn(cliCtx, false)
}
func issueFreezeHandlerFn(cliCtx context.CLIContext, freeze bool) http.HandlerFunc {
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

		msg, err := utils.GetIssueFreezeMsg(cdc, cliCtx, account, vars["freeze-type"], vars["issue-id"], vars["accAddress"], vars[EndTime], freeze)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		types.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postBurnHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cliCtx, types.BurnHolder)
}
func postBurnFromHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cliCtx, types.BurnFrom)
}
func postIssueSendFrom(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars["issue-id"]
		if err := types.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		from, err := sdk.AccAddressFromBech32(vars["from"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		to, err := sdk.AccAddressFromBech32(vars["to"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		num, err := strconv.ParseInt(vars["amount"], 10, 64)
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
		_, err = utils.GetIssueByID(cliCtx, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := utils.CheckAllowance(cliCtx, issueID, from, account.GetAddress(), amount); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err = utils.CheckFreeze(cliCtx, issueID, from, to); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgIssueSendFrom(issueID, sender, from, to, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cutils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postIssueApproveHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cliCtx, types.Approve)
}
func postIssueIncreaseApproval(cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cliCtx, types.IncreaseApproval)
}
func postIssueDecreaseApproval(cliCtx context.CLIContext) http.HandlerFunc {
	return issueApproveHandlerFn(cliCtx, types.DecreaseApproval)
}
func issueApproveHandlerFn(cliCtx context.CLIContext, approveType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req PostIssueBaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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

		vars := mux.Vars(r)

		issueID := vars["issue-id"]
		accAddr, err := sdk.AccAddressFromBech32(vars["accAddress"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount, ok := sdk.NewIntFromString(vars["amount"])
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Amount not a valid int")
			return
		}

		account, err := cliCtx.GetAccount(fromAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg, err := utils.GetIssueApproveMsg(cliCtx, issueID, account, accAddr, approveType, amount, false)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cutils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

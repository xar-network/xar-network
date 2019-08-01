package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zar-network/zar-network/x/issue/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"
	"github.com/zar-network/zar-network/x/issue/msgs"
	"github.com/zar-network/zar-network/x/issue/params"
	"github.com/zar-network/zar-network/x/issue/types"
	issueutils "github.com/zar-network/zar-network/x/issue/utils"
)

type PostIssueReq struct {
	BaseReq            rest.BaseReq `json:"base_req"`
	params.IssueParams `json:"issue"`
}
type PostDescriptionReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Description string       `json:"description"`
}
type PostIssueBaseReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/issue", postIssueHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postIssueApproveHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/increase/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postIssueIncreaseApproval(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/approve/decrease/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postIssueDecreaseApproval(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/burn/{%s}/{%s}", IssueID, Amount), postBurnHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/burn-from/{%s}/{%s}/{%s}", IssueID, AccAddress, Amount), postBurnFromHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/freeze/{%s}/{%s}/{%s}/{%s}", FreezeType, IssueID, AccAddress, EndTime), postIssueFreezeHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/unfreeze/{%s}/{%s}/{%s}", FreezeType, IssueID, AccAddress), postIssueUnFreezeHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/send-from/{%s}/{%s}/{%s}/{%s}", IssueID, From, To, Amount), postIssueSendFrom(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/mint/{%s}/{%s}/{%s}", IssueID, Amount, To), postMintHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/ownership/transfer/{%s}/{%s}", IssueID, To), postTransferOwnershipHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/description/{%s}", IssueID), postDescribeHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/issue/feature/disable/{%s}/{%s}", IssueID, Feature), postDisableFeatureHandlerFn(cdc, cliCtx)).Methods("POST")

}
func postIssueHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostIssueReq
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

		if len(req.Description) > 0 && !json.Valid([]byte(req.Description)) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrCoinDescriptionNotValid().Error())
			return
		}
		// create the message
		msg := msgs.NewMsgIssue(fromAddress, &req.IssueParams)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postMintHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		num, err := strconv.ParseInt(vars[Amount], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		amount := sdk.NewInt(num)
		to, err := sdk.AccAddressFromBech32(vars[To])
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
		fromAddress, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			return
		}
		account, err := cliCtx.GetAccount(fromAddress)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		issueInfo, err := clientutils.IssueOwnerCheck(cdc, cliCtx, account, issueID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := msgs.NewMsgIssueMint(issueID, fromAddress, amount, issueInfo.GetDecimals(), to)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postDisableFeatureHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		feature := vars[Feature]
		_, ok := types.Features[feature]
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrUnknownFeatures().Error())
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
		msg := msgs.NewMsgIssueDisableFeature(issueID, fromAddress, feature)

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}

}
func postDescribeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
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
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrCoinDescriptionNotValid().Error())
			return
		}
		msg := msgs.NewMsgIssueDescription(issueID, fromAddress, []byte(req.Description))
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
		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
func postTransferOwnershipHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
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

		msg := msgs.NewMsgIssueTransferOwnership(issueID, fromAddress, to)
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

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

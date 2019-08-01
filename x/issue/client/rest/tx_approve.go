package rest

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"

	"github.com/zar-network/zar-network/x/issue/msgs"
	"github.com/zar-network/zar-network/x/issue/types"
	issueutils "github.com/zar-network/zar-network/x/issue/utils"
)

func postIssueSendFrom(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		issueID := vars[IssueID]
		if err := issueutils.CheckIssueId(issueID); err != nil {
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

		msg := msgs.NewMsgIssueSendFrom(issueID, sender, from, to, amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
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

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

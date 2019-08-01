package rest

import (
	"net/http"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"

	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

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

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

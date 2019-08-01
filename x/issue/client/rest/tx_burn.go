package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	clientutils "github.com/zar-network/zar-network/x/issue/client/utils"

	"github.com/zar-network/zar-network/x/issue/types"
	issueutils "github.com/zar-network/zar-network/x/issue/utils"
)

func postBurnHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cdc, cliCtx, types.BurnHolder)
}
func postBurnFromHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return postBurnFromAddressHandlerFn(cdc, cliCtx, types.BurnFrom)
}

func postBurnFromAddressHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, burnFromType string) http.HandlerFunc {
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
		if err := issueutils.CheckIssueId(issueID); err != nil {
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

		//burn sender
		accAddress := fromAddress

		if types.BurnFrom == burnFromType {
			//burn from holder address
			accAddress, err = sdk.AccAddressFromBech32(vars[AccAddress])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		msg, err := clientutils.GetBurnMsg(cdc, cliCtx, account, accAddress, issueID, amount, burnFromType, false)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

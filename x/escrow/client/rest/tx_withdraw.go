package rest

import (
	"net/http"

	"github.com/hashgard/hashgard/x/box/errors"
	boxutils "github.com/hashgard/hashgard/x/box/utils"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	clientutils "github.com/hashgard/hashgard/x/box/client/utils"
)

func PostWithdrawHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, boxType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[ID]
		if boxutils.GetBoxTypeByValue(id) != boxType {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrUnknownBox(id).Error())
			return
		}
		var req PostBoxBaseReq
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
		msg, err := clientutils.GetWithdrawMsg(cdc, cliCtx, account, id)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

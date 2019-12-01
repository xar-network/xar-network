package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/hashgard/hashgard/x/box/msgs"
	"github.com/hashgard/hashgard/x/box/params"
)

type PostLockBoxReq struct {
	BaseReq              rest.BaseReq `json:"base_req"`
	params.BoxLockParams `json:"box"`
}
type PostDepositBoxReq struct {
	BaseReq                 rest.BaseReq `json:"base_req"`
	params.BoxDepositParams `json:"box"`
}
type PostFutureBoxReq struct {
	BaseReq                rest.BaseReq `json:"base_req"`
	params.BoxFutureParams `json:"box"`
}

func PostLockBoxCreateHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostLockBoxReq
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

		params := params.BoxLockParams{
			Name:        req.Name,
			TotalAmount: req.TotalAmount,
			Description: req.Description,
			Lock:        req.Lock,
		}
		// create the message
		msg := msgs.NewMsgLockBox(fromAddress, &params)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
func PostDepositBoxCreateHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostDepositBoxReq
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

		params := params.BoxDepositParams{
			Name:        req.Name,
			TotalAmount: req.TotalAmount,
			Description: req.Description,
			Deposit:     req.Deposit,
		}
		// create the message
		msg := msgs.NewMsgDepositBox(fromAddress, &params)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
func PostFutureBoxCreateHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostFutureBoxReq
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

		params := params.BoxFutureParams{
			Name:        req.Name,
			TotalAmount: req.TotalAmount,
			Description: req.Description,
			Future:      req.Future,
		}
		// create the message
		msg := msgs.NewMsgFutureBox(fromAddress, &params)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

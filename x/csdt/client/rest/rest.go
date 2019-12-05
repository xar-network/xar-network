package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/xar-network/xar-network/x/csdt/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

/*
API Design:

Currently CSDTs do not have IDs so standard REST uri conventions (ie GET /csdts/{csdt-id}) don't work too well.

Get one or more csdts
	GET /csdts?collateralDenom={denom}&owner={address}&underCollateralizedAt={price}
Modify a CSDT (idempotent). Create is not separated out because conceptually all CSDTs already exist (just with zero collateral and debt). // TODO is making this idempotent actually useful?
	PUT /csdts
Get the module params, including authorized collateral denoms.
	GET /params
*/

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/csdts", getCsdtsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/csdts", modifyCsdtHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc("/csdts/params", getParamsHandlerFn(cliCtx)).Methods("GET")
}

const (
	RestOwner                 = "owner"
	RestCollateralDenom       = "collateralDenom"
	RestUnderCollateralizedAt = "underCollateralizedAt"
)

func getCsdtsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get parameters from the URL
		ownerBech32 := r.URL.Query().Get(RestOwner)
		collateralDenom := r.URL.Query().Get(RestCollateralDenom)
		priceString := r.URL.Query().Get(RestUnderCollateralizedAt)

		// Construct querier params
		querierParams := types.QueryCsdtsParams{}

		if len(ownerBech32) != 0 {
			owner, err := sdk.AccAddressFromBech32(ownerBech32)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			querierParams.Owner = owner
		}

		if len(collateralDenom) != 0 {
			// TODO validate denom
			querierParams.CollateralDenom = collateralDenom
		}

		if len(priceString) != 0 {
			price, err := sdk.NewDecFromStr(priceString)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			querierParams.UnderCollateralizedAt = price
		}

		querierParamsBz, err := cliCtx.Codec.MarshalJSON(querierParams)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// Get the CSDTs
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/csdt/%s", types.QueryGetCsdts), querierParamsBz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)

		// Return the CSDTs
		rest.PostProcessResponse(w, cliCtx, res)

	}
}

type ModifyCsdtRequestBody struct {
	BaseReq rest.BaseReq                `json:"base_req"`
	Csdt    types.MsgCreateOrModifyCSDT `json:"csdt"`
}

func modifyCsdtHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode PUT request body
		var requestBody ModifyCsdtRequestBody
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &requestBody) {
			return
		}
		requestBody.BaseReq = requestBody.BaseReq.Sanitize()
		if !requestBody.BaseReq.ValidateBasic(w) {
			return
		}

		// Get the stored CSDT
		querierParams := types.QueryCsdtsParams{
			Owner:           requestBody.Csdt.Sender,
			CollateralDenom: requestBody.Csdt.CollateralDenom,
		}
		querierParamsBz, err := cliCtx.Codec.MarshalJSON(querierParams)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/csdt/%s", types.QueryGetCsdts), querierParamsBz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)

		var csdts types.CSDTs
		err = cliCtx.Codec.UnmarshalJSON(res, &csdts)
		if len(csdts) != 1 || err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Create and return msg
		msg := types.NewMsgCreateOrModifyCSDT(
			requestBody.Csdt.Sender,
			requestBody.Csdt.CollateralDenom,
			requestBody.Csdt.CollateralChange,
			requestBody.Csdt.DebtChange,
		)
		utils.WriteGenerateStdTxResponse(w, cliCtx, requestBody.BaseReq, []sdk.Msg{msg})
	}
}

func getParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the params
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/csdt/%s", types.QueryGetParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		cliCtx = cliCtx.WithHeight(height)

		// Return the params
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

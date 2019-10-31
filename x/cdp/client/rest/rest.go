package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/xar-network/xar-network/x/cdp/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

/*
API Design:

Currently CDPs do not have IDs so standard REST uri conventions (ie GET /cdps/{cdp-id}) don't work too well.

Get one or more cdps
	GET /cdps?collateralDenom={denom}&owner={address}&underCollateralizedAt={price}
Modify a CDP (idempotent). Create is not separated out because conceptually all CDPs already exist (just with zero collateral and debt). // TODO is making this idempotent actually useful?
	PUT /cdps
Get the module params, including authorized collateral denoms.
	GET /params
*/

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/cdps", getCdpsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/cdps", modifyCdpHandlerFn(cliCtx)).Methods("PUT")
	r.HandleFunc("/cdps/params", getParamsHandlerFn(cliCtx)).Methods("GET")
}

const (
	RestOwner                 = "owner"
	RestCollateralDenom       = "collateralDenom"
	RestUnderCollateralizedAt = "underCollateralizedAt"
)

func getCdpsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get parameters from the URL
		ownerBech32 := r.URL.Query().Get(RestOwner)
		collateralDenom := r.URL.Query().Get(RestCollateralDenom)
		priceString := r.URL.Query().Get(RestUnderCollateralizedAt)

		// Construct querier params
		querierParams := types.QueryCdpsParams{}

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

		// Get the CDPs
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cdp/%s", types.QueryGetCdps), querierParamsBz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		// Return the CDPs
		rest.PostProcessResponse(w, cliCtx, res)

	}
}

type ModifyCdpRequestBody struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Cdp     types.CDP    `json:"cdp"`
}

func modifyCdpHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode PUT request body
		var requestBody ModifyCdpRequestBody
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &requestBody) {
			return
		}
		requestBody.BaseReq = requestBody.BaseReq.Sanitize()
		if !requestBody.BaseReq.ValidateBasic(w) {
			return
		}

		// Get the stored CDP
		querierParams := types.QueryCdpsParams{
			Owner:           requestBody.Cdp.Owner,
			CollateralDenom: requestBody.Cdp.CollateralDenom,
		}
		querierParamsBz, err := cliCtx.Codec.MarshalJSON(querierParams)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cdp/%s", types.QueryGetCdps), querierParamsBz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var cdps types.CDPs
		err = cliCtx.Codec.UnmarshalJSON(res, &cdps)
		if len(cdps) != 1 || err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Calculate CDP updates
		collateralDelta := requestBody.Cdp.CollateralAmount.Sub(cdps[0].CollateralAmount)
		debtDelta := requestBody.Cdp.Debt.Sub(cdps[0].Debt)

		// Create and return msg
		msg := types.NewMsgCreateOrModifyCDP(
			requestBody.Cdp.Owner,
			requestBody.Cdp.CollateralDenom,
			collateralDelta,
			debtDelta,
		)
		utils.WriteGenerateStdTxResponse(w, cliCtx, requestBody.BaseReq, []sdk.Msg{msg})
	}
}

func getParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the params
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cdp/%s", types.QueryGetParams), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Return the params
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

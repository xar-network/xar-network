package rest

import (
	"net/http"
	"strconv"

	"github.com/hashgard/hashgard/x/box/types"

	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/box/client/utils"

	"github.com/hashgard/hashgard/x/box/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/hashgard/hashgard/x/box/params"

	"github.com/hashgard/hashgard/x/box/client/queriers"
	boxutils "github.com/hashgard/hashgard/x/box/utils"
)

func BoxQueryParamsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, boxType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := queriers.QueryBoxParams(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var params box.Params
		cdc.MustUnmarshalJSON(res, &params)
		rest.PostProcessResponse(w, cdc, utils.GetBoxParams(params, boxType), cliCtx.Indent)
	}
}
func BoxQueryHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, boxType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars[ID]
		if boxutils.GetBoxTypeByValue(id) != boxType {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrUnknownBox(id).Error())
			return
		}
		if err := boxutils.CheckId(id); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := queriers.QueryBoxByID(id, cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var box types.BoxInfo
		cdc.MustUnmarshalJSON(res, &box)
		rest.PostProcessResponse(w, cdc, utils.GetBoxInfo(box), cliCtx.Indent)
	}
}
func BoxSearchHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, boxType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, err := queriers.QueryBoxByName(boxType, vars[Name], cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
func BoxListHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext, boxType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address, err := sdk.AccAddressFromBech32(r.URL.Query().Get(RestAddress))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		boxQueryParams := params.BoxQueryParams{
			StartId: r.URL.Query().Get(RestStartId),
			BoxType: boxType,
			Owner:   address,
			Limit:   30,
		}
		strNumLimit := r.URL.Query().Get(RestLimit)
		if len(strNumLimit) > 0 {
			limit, err := strconv.Atoi(r.URL.Query().Get(RestLimit))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			boxQueryParams.Limit = limit
		}

		res, err := queriers.QueryBoxsList(boxQueryParams, cdc, cliCtx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

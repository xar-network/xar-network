package market

import (
	"fmt"
	"net/http"

	"github.com/xar-network/xar-network/embedded"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.Handle("/markets", getMarkets(ctx, cdc)).Methods("GET")
}

func getMarkets(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := ctx.QueryWithData(fmt.Sprintf("custom/market/list"), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if res == nil {
			w.WriteHeader(404)
			return
		}
		ctx = ctx.WithHeight(height)

		embedded.PostProcessResponse(w, ctx, res)
	}
}

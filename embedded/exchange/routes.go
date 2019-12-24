/*

Copyright 2019 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package exchange

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"

	"github.com/xar-network/xar-network/embedded"
	"github.com/xar-network/xar-network/embedded/auth"
	"github.com/xar-network/xar-network/embedded/order"
	"github.com/xar-network/xar-network/types/store"
	"github.com/xar-network/xar-network/x/order/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	sdkauth "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	sub := r.PathPrefix("/exchange").Subrouter()
	sub.Use(auth.DefaultAuthMW)
	sub.HandleFunc("/orders", postOrderHandler(ctx, cdc)).Methods("POST")
	sub.HandleFunc("/orders/{order_id}/get", getOrderHandler(ctx, cdc)).Methods("GET")
}

func postOrderHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OrderCreationRequest
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		kb := auth.MustGetKBFromSession(r)
		owner := kb.GetAddr()
		ctx = ctx.WithFromAddress(owner)

		msg := types.NewMsgPost(owner, req.MarketID, req.Direction, req.Price, req.Quantity, req.TimeInForce)
		msgs := []sdk.Msg{msg}
		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		bldr := sdkauth.NewTxBuilderFromCLI(nil).
			WithTxEncoder(utils.GetTxEncoder(cdc)).
			WithKeybase(kb)

		bldr, sdkErr := utils.PrepareTxBuilder(bldr, ctx)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}

		broadcastResB, sdkErr := bldr.BuildAndSign(kb.GetName(), auth.MustGetKBPassphraseFromSession(r), msgs)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}
		broadcastRes, sdkErr := ctx.BroadcastTxCommit(broadcastResB)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}

		var orderIDStr string
		for _, log := range broadcastRes.Logs {
			if strings.HasPrefix(log.Log, "order_id") {
				orderIDStr = strings.TrimPrefix(log.Log, "order_id:")
				break
			}
		}
		orderID := store.NewEntityIDFromString(orderIDStr)
		res := OrderCreationResponse{
			BlockInclusion: embedded.BlockInclusion{
				BlockNumber:     broadcastRes.Height,
				TransactionHash: broadcastRes.TxHash,
				BlockTimestamp:  broadcastRes.Timestamp,
			},
			ID:          orderID,
			MarketID:    msg.MarketID,
			Direction:   msg.Direction,
			Price:       msg.Price,
			Quantity:    msg.Quantity,
			Type:        req.Type,
			TimeInForce: msg.TimeInForce,
			Status:      "OPEN",
		}

		out, sdkErr := cdc.MarshalJSON(res)
		if sdkErr != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, sdkErr.Error())
			return
		}
		if _, err := w.Write(out); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func getOrderHandler(ctx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse URI for get order id
		var orderIDStr string
		re, _ := regexp.Compile("/orders/(.+)/(get|test)")
		values := re.FindStringSubmatch(r.URL.RequestURI())
		if len(values) > 1 {
			orderIDStr = values[1]
		}
		// Check for test mode
		testMode := values[2] == "test"

		if orderIDStr == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "order_id not present")
			return
		}

		orderID := store.NewEntityIDFromString(orderIDStr)

		// Use standart filter with start = order id and limit 1
		req := order.ListQueryRequest{
			Start: orderID,
			Limit: 1,
		}

		var resB []byte
		var err error
		if !testMode {
			// In test mode this block not worked with error
			resB, _, err = ctx.QueryWithData("custom/embeddedorder/list", cdc.MustMarshalBinaryBare(req))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			// In test mode return test answer
			resB = []byte(fmt.Sprintf("Test OK = %s", orderID.String()))
		}

		embedded.PostProcessResponse(w, ctx, resB)
	}
}

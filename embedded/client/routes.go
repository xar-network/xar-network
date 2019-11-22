package client

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/xar-network/xar-network/embedded/auth"
	"github.com/xar-network/xar-network/embedded/balance"
	"github.com/xar-network/xar-network/embedded/batch"
	"github.com/xar-network/xar-network/embedded/book"
	"github.com/xar-network/xar-network/embedded/exchange"
	"github.com/xar-network/xar-network/embedded/fill"
	"github.com/xar-network/xar-network/embedded/market"
	"github.com/xar-network/xar-network/embedded/order"
	"github.com/xar-network/xar-network/embedded/price"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router, cdc *codec.Codec, enableFaucet bool) {
	r.Use(auth.HandleCORSMW)
	r.Use(auth.ProtectCSRFMW([]string{
		"/api/v1/faucet/transfer",
	}))
	sub := r.PathPrefix("/api/v1").Subrouter()
	auth.RegisterRoutes(ctx, sub, cdc)
	exchange.RegisterRoutes(ctx, sub, cdc)
	fill.RegisterRoutes(ctx, sub, cdc)
	market.RegisterRoutes(ctx, sub, cdc)
	order.RegisterRoutes(ctx, sub, cdc)
	balance.RegisterRoutes(ctx, sub, cdc, enableFaucet)
	price.RegisterRoutes(ctx, sub, cdc)
	book.RegisterRoutes(ctx, sub, cdc)
	batch.RegisterRoutes(ctx, sub, cdc)
	//ui.RegisterRoutes(ctx, r, cdc)
}

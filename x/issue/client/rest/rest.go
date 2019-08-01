package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
)

const (
	IssueID    = "issue-id"
	Feature    = "feature"
	AccAddress = "accAddress"
	From       = "from"
	FreezeType = "freeze-type"
	EndTime    = "end-time"
	Symbol     = "symbol"
	Amount     = "amount"
	To         = "to"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc)
}

package fill

import (
	"github.com/xar-network/xar-network/embedded"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Fill struct {
	OrderID     store.EntityID     `json:"order_id"`
	Owner       sdk.AccAddress     `json:"owner"`
	Pair        string             `json:"pair"`
	Direction   matcheng.Direction `json:"direction"`
	QtyFilled   sdk.Uint           `json:"qty_filled"`
	QtyUnfilled sdk.Uint           `json:"qty_unfilled"`
	BlockNumber int64              `json:"block_number"`
	Price       sdk.Uint           `json:"price"`
}

type QueryRequest struct {
	Owner      sdk.AccAddress
	StartBlock int64
	EndBlock   int64
}

type QueryResult struct {
	Fills []Fill
}

type RESTQueryResult struct {
	Fills []RESTFill `json:"fills"`
}

type RESTFill struct {
	BlockInclusion   embedded.BlockInclusion `json:"block_inclusion"`
	QuantityFilled   sdk.Uint                `json:"quantity_filled"`
	QuantityUnfilled sdk.Uint                `json:"quantity_unfilled"`
	Direction        matcheng.Direction      `json:"direction"`
	OrderID          store.EntityID          `json:"order_id"`
	Pair             string                  `json:"pair"`
	Price            sdk.Uint                `json:"price"`
	Owner            sdk.AccAddress          `json:"owner"`
}

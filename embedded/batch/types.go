package batch

import (
	"time"

	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Batch struct {
	BlockNumber   int64                     `json:"block_number"`
	BlockTime     time.Time                 `json:"block_time"`
	MarketID      store.EntityID            `json:"market_id"`
	ClearingPrice sdk.Uint                  `json:"clearing_price"`
	Bids          []matcheng.AggregatePrice `json:"bids"`
	Asks          []matcheng.AggregatePrice `json:"asks"`
}

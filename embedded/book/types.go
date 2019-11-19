package book

import (
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Book struct {
	Bids []matcheng.AggregatePrice `json:"bids"`
	Asks []matcheng.AggregatePrice `json:"asks"`
}

type QueryResultEntry struct {
	Price    sdk.Uint `json:"price"`
	Quantity sdk.Uint `json:"quantity"`
}

type QueryResult struct {
	MarketID    store.EntityID     `json:"market_id"`
	BlockNumber int64              `json:"block_number"`
	Bids        []QueryResultEntry `json:"bids"`
	Asks        []QueryResultEntry `json:"asks"`
}

package exchange

import (
	"github.com/xar-network/xar-network/embedded"
	"github.com/xar-network/xar-network/pkg/matcheng"
	"github.com/xar-network/xar-network/types/store"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OrderCreationRequest struct {
	MarketID    store.EntityID     `json:"market_id"`
	Direction   matcheng.Direction `json:"direction"`
	Price       sdk.Uint           `json:"price"`
	Quantity    sdk.Uint           `json:"quantity"`
	Type        string             `json:"type"`
	TimeInForce uint16             `json:"time_in_force"`
}

type OrderCreationResponse struct {
	BlockInclusion embedded.BlockInclusion `json:"block_inclusion"`
	ID             store.EntityID          `json:"id"`
	MarketID       store.EntityID          `json:"market_id"`
	Direction      matcheng.Direction      `json:"direction"`
	Price          sdk.Uint                `json:"price"`
	Quantity       sdk.Uint                `json:"quantity"`
	Type           string                  `json:"type"`
	TimeInForce    uint16                  `json:"time_in_force"`
	Status         string                  `json:"status"`
}

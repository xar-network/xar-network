package types

import (
	"github.com/xar-network/xar-network/types/store"
)

const (
	ModuleName = "market"
	RouterKey  = ModuleName
	StoreKey   = RouterKey
)

type Market struct {
	ID              store.EntityID
	BaseAssetDenom  string
	QuoteAssetDenom string
}

func NewMarket(
	id store.EntityID,
	baseAsset string,
	quoteAsset string,
) Market {
	return Market{
		ID:              id,
		BaseAssetDenom:  baseAsset,
		QuoteAssetDenom: quoteAsset,
	}
}

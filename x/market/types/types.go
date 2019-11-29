package types

import (
	"fmt"

	"github.com/xar-network/xar-network/types/store"
)

type Markets []Market

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

// implement fmt.Stringer
func (m Market) String() string {
	return fmt.Sprintf(`Market:
	ID: %s
	Base Asset: %s
	Quote Asset: %s`,
		m.ID.String(), m.BaseAssetDenom, m.QuoteAssetDenom)
}

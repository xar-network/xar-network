package types

import (
	"fmt"
	"strings"
)

// implement fmt.Stringer
func (a PendingPriceAsset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AssetCode: %s`, a.AssetCode))
}

// PendingPriceAsset struct that contains the info about the asset which price is still to be determined
type PendingPriceAsset struct {
	AssetCode string `json:"asset_code"`
}

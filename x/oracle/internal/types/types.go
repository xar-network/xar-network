package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// QueryCurrentPrice command for current price queries
	QueryCurrentPrice = "price"
	// QueryRawPrices command for raw price queries
	QueryRawPrices = "rawprices"
	// QueryAssets command for assets query
	QueryAssets = "assets"
)

// implement fmt.Stringer
func (cp CurrentPrice) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AssetCode: %s
Price: %s
Expiry: %s`, cp.AssetCode, cp.Price, cp.Expiry))
}

// implement fmt.Stringer
func (pp PostedPrice) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AssetCode: %s
OracleAddress: %s
Price: %s
Expiry: %s`, pp.AssetCode, pp.OracleAddress, pp.Price, pp.Expiry))
}

// implement fmt.Stringer
func (a Asset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AssetCode: %s
Description: %s`, a.AssetCode, a.Description))
}

// implement fmt.Stringer
func (a PendingPriceAsset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`AssetCode: %s`, a.AssetCode))
}

// QueryRawPricesResp response to a rawprice query
type QueryRawPricesResp []string

// implement fmt.Stringer
func (n QueryRawPricesResp) String() string {
	return strings.Join(n[:], "\n")
}

// QueryAssetsResp response to a assets query
type QueryAssetsResp []string

// implement fmt.Stringer
func (n QueryAssetsResp) String() string {
	return strings.Join(n[:], "\n")
}

// Asset struct that represents an asset in the oracle
type Asset struct {
	AssetCode   string `json:"asset_code"` // The nft id or the asset denom
	Description string `json:"description"`
}

// Oracle struct that documents which address an oracle is using
type Oracle struct {
	OracleAddress string `json:"oracle_address"`
}

// PendingPriceAsset struct that contains the info about the asset which price is still to be determined
type PendingPriceAsset struct {
	AssetCode string `json:"asset_code"`
}

// CurrentPrice struct that contains the metadata of a current price for a particular asset in the oracle module.
type CurrentPrice struct {
	AssetCode string  `json:"asset_code"`
	Price     sdk.Dec `json:"price"`
	Expiry    sdk.Int `json:"expiry"`
}

// PostedPrice struct represented a price for an asset posted by a specific oracle
type PostedPrice struct {
	AssetCode     string  `json:"asset_code"`
	OracleAddress string  `json:"oracle_address"`
	Price         sdk.Dec `json:"price"`
	Expiry        sdk.Int `json:"expiry"`
}

// SortDecs provides the interface needed to sort sdk.Dec slices
type SortDecs []sdk.Dec

func (a SortDecs) Len() int           { return len(a) }
func (a SortDecs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortDecs) Less(i, j int) bool { return a[i].LT(a[j]) }

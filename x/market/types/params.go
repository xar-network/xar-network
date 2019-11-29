package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Parameter keys
var (
	// ParamStoreKeyAuctionParams Param store key for auction params
	KeyMarkets = []byte(ModuleName)
	KeyAuth    = []byte(ModuleAuth)
)

type Params struct {
	Markets []Market
	POA     string
}

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		subspace.NewParamSetPair(KeyMarkets, &p.Markets),
		subspace.NewParamSetPair(KeyAuth, &p.POA),
	}
}

// NewParams creates a new Params object
func NewParams(markets []Market, poa string) Params {
	return Params{
		Markets: markets,
		POA:     poa,
	}
}

// DefaultParams default params
func DefaultParams() Params {
	return NewParams(Markets{}, "")
}

// String implements fmt.stringer
func (p Params) String() string {
	out := "Params:\n"
	out += "PoA: " + p.POA + "\n"
	for _, a := range p.Markets {
		out += a.String()
	}
	return strings.TrimSpace(out)
}

// ParamSubspace defines the expected Subspace interface for parameters
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

// Validate ensure that params have valid values
func (p Params) Validate() error {
	// iterate over assets and verify them
	for _, market := range p.Markets {
		if !market.ID.IsDefined() {
			return fmt.Errorf("invalid id: %s. missing id", market.String())
		}
	}
	return nil
}

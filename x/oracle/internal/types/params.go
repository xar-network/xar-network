/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	// KeyAssets store key for assets
	KeyAssets   = []byte("oracleassets")
	KeyNominees = []byte("oraclenominees")
)

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params params for oracle. Can be altered via governance
type Params struct {
	Assets   []Asset  `json:"assets" yaml:"assets"` //  Array containing the assets supported by the oracle
	Nominees []string `json:"nominees" yaml:"nominees"`
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of oracle module's parameters.
func (p Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyAssets, Value: &p.Assets},
		{Key: KeyNominees, Value: &p.Nominees},
	}
}

// NewParams creates a new AssetParams object
func NewParams(assets []Asset, nominees []string) Params {
	return Params{
		Assets:   assets,
		Nominees: nominees,
	}
}

// DefaultParams default params for oracle
func DefaultParams() Params {
	return NewParams(Assets{}, []string{})
}

// String implements fmt.stringer
func (p Params) String() string {
	out := "Params:\n"
	for _, a := range p.Assets {
		out += a.String()
	}
	for _, a := range p.Nominees {
		out += a
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
	for _, asset := range p.Assets {
		if asset.AssetCode == "" {
			return fmt.Errorf("invalid asset: %s. missing asset code", asset.String())
		}
	}
	return nil
}

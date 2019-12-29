/*

Copyright 2016 All in Bits, Inc
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

package coinswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/coinswap/internal/types"
)

// TODO: ...

// GenesisState - coinswap genesis state
type GenesisState struct {
	Params types.Params `json:"params"`
}

// NewGenesisState is the constructor function for GenesisState
func NewGenesisState(params types.Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultParams())
}

// InitGenesis new coinswap genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	if err := types.ValidateParams(data.Params); err != nil {
		panic(err)
	}

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(types.DefaultParams())
}

// ValidateGenesis - placeholder function
func ValidateGenesis(data GenesisState) error {
	if err := types.ValidateParams(data.Params); err != nil {
		return err
	}
	return nil
}

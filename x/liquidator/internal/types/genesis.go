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

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	Params LiquidatorParams `json:"liquidator_params" yaml:"liquidator_params"`
}

// DefaultGenesisState returns a default genesis state
// TODO pick better values
func DefaultGenesisState() GenesisState {
	return GenesisState{
		DefaultParams(),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	return nil
}

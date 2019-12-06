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
	"bytes"
)

// GenesisAuctions type for an array of auctions
type GenesisAuctions []Auction

// GenesisState - auction state that must be provided at genesis
type GenesisState struct {
	AuctionParams AuctionParams   `json:"auction_params" yaml:"auction_params"`
	Auctions      GenesisAuctions `json:"genesis_auctions" yaml:"genesis_auctions"`
}

// NewGenesisState returns a new genesis state object for auctions module
func NewGenesisState(ap AuctionParams, ga GenesisAuctions) GenesisState {
	return GenesisState{
		AuctionParams: ap,
		Auctions:      ga,
	}
}

// DefaultGenesisState defines default genesis state for auction module
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultAuctionParams(), GenesisAuctions{})
}

// Equal checks whether two GenesisState structs are equivalent
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis validates genesis inputs. Returns error if validation of any input fails.
func ValidateGenesis(data GenesisState) error {
	if err := data.AuctionParams.Validate(); err != nil {
		return err
	}
	return nil
}

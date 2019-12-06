/*

Copyright 2019 All in Bits, Inc
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

package market

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/types/store"
	"github.com/xar-network/xar-network/x/market/types"
)

type GenesisState struct {
	Markets  types.Markets `json:"markets" yaml:"markets"`
	Nominees []string      `json:"nominees" yaml:"nominees"`
}

func NewGenesisState(markets types.Markets, nominees []string) GenesisState {
	return GenesisState{Markets: markets, Nominees: nominees}
}

func ValidateGenesis(data GenesisState) error {
	currentId := store.ZeroEntityID
	for _, market := range data.Markets {
		currentId = currentId.Inc()
		if !currentId.Equals(market.ID) {
			return errors.New("Invalid Market: ID must monotonically increase.")
		}
		if market.BaseAssetDenom == "" {
			return errors.New("Invalid Market: Must specify a non-zero base asset denom.")
		}
		if market.QuoteAssetDenom == "" {
			return errors.New("Invalid Market: Must specify a non-zero quote asset denom.")
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Markets: types.Markets{
			{
				ID:              store.NewEntityID(1),
				BaseAssetDenom:  "uftm",
				QuoteAssetDenom: "uzar",
			},
			{
				ID:              store.NewEntityID(2),
				BaseAssetDenom:  "ueur",
				QuoteAssetDenom: "uzar",
			},
			{
				ID:              store.NewEntityID(3),
				BaseAssetDenom:  "uusd",
				QuoteAssetDenom: "uzar",
			},
			{
				ID:              store.NewEntityID(4),
				BaseAssetDenom:  "ubtc",
				QuoteAssetDenom: "uzar",
			},
		},
		Nominees: []string{},
	}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetParams(ctx, types.NewParams(data.Markets, []string{}))
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	return GenesisState{Markets: params.Markets, Nominees: params.Nominees}
}

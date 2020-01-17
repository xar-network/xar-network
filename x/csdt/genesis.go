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

package csdt

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/types/fee"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	Params     types.Params `json:"params"`
	CSDTs      types.CSDTs  `json:"csdts" yaml:"csdts"`
	// don't need to setup CollateralStates as they are created as needed
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		types.Params{
			GlobalDebtLimit: sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(25000000000000))),
			CollateralParams: types.CollateralParams{
				{
					Denom:            "ubtc",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
				},
				{
					Denom:            "ubnb",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
				},
				{
					Denom:            "ueth",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
				},
				{
					Denom:            "uftm",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
				},
				{
					Denom:            "uzar",
					LiquidationRatio: sdk.MustNewDecFromStr("1.3"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
				},
			},
			Fee: fee.FromPercentString("0"),
		},
		types.CSDTs{},
	}
}

func NewGenesisState(params types.Params, globalDebt sdk.Int) GenesisState {
	return GenesisState{
		Params:     params,
		CSDTs:      types.CSDTs{},
	}
}

// InitGenesis sets the genesis state in the keeper.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	k.SetParams(ctx, data.Params)

	for _, csdt := range data.CSDTs {
		k.SetCSDT(ctx, csdt)
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	// TODO implement this
	// validate denoms
	// check collateral debt limits sum to global limit?
	// check limits are > 0
	// check ratios are > 1
	// check no repeated denoms

	// check global debt is zero - force the chain to always start with zero stable coin, otherwise collateralStatus's will need to be set up as well.
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	params := k.GetParams(ctx)
	csdts := types.CSDTs{}

	l, err := k.GetCSDTs(ctx)
	if err != nil {
		panic(err)
	} else {
		csdts = append(csdts, l...)
	}

	return GenesisState{
		Params:     params,
		CSDTs:      csdts,
	}
}

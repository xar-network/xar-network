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

package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/types/fee"
	"time"
)

const (
	defaultBlocksPerSnapshot = 0
	defaultSnapshotLimit     = 100
	defaultTimerInterval     = time.Hour * 24
)

// Parameter keys
var (
	KeySyntheticParams     = []byte("SyntheticParams")
	KeyNominees            = []byte("Nominees")
	KeyFee                 = []byte("Fee")
	KeyMarketBalanceParam  = []byte("MarketBalanceParam")
	DefaultSyntheticParams = SyntheticParams{SyntheticParam{
		Denom: "sbtc",
	}}
	DefaultMarketBalanceParam = MarketBalanceParam{
		defaultBlocksPerSnapshot,
		defaultTimerInterval,
		defaultSnapshotLimit,
		nil,
	}
)

// Params governance parameters for synthetic module
type Params struct {
	SyntheticParams    SyntheticParams    `json:"synthetic_params" yaml:"synthetic_params"`
	Nominees           []string           `json:"nominees" yaml:"nominees"`
	Fee                fee.Fee            `json:"fee" yaml:"fee"`
	MarketBalanceParam MarketBalanceParam `json:"market_balance_param" yaml:"market_balance_param"`
}

func (p Params) IsSyntheticPresent(denom string) bool {
	// search for matching denom, return
	for _, sp := range p.SyntheticParams {
		if sp.Denom == denom {
			return true
		}
	}
	return false
}

func (p Params) GetSyntheticParam(denom string) SyntheticParam {
	// search for matching denom, return
	for _, sp := range p.SyntheticParams {
		if sp.Denom == denom {
			return sp
		}
	}
	// panic if not found, to be safe
	panic("synthetic params not found in module params")
}

// String implements fmt.Stringer
func (p Params) String() string {
	return fmt.Sprintf(`Params:
	Synthetic Params: %s
	Nominees: %s`,
		p.SyntheticParams,
		p.Nominees,
	)
}

// NewParams returns a new params object
func NewParams(
	syntheticParams SyntheticParams,
	nominees []string,
	fee fee.Fee,
	mbParam MarketBalanceParam,
) Params {
	return Params{
		SyntheticParams:    syntheticParams,
		Nominees:           nominees,
		Fee:                fee,
		MarketBalanceParam: mbParam,
	}
}

// DefaultParams returns default params for synthetic module
func DefaultParams() Params {
	return NewParams(
		DefaultSyntheticParams,
		[]string{},
		fee.FromPercentString("0.005"),
		DefaultMarketBalanceParam,
	)
}

type SyntheticParam struct {
	Denom string `json:"denom" yaml:"denom"`
}

// String implements fmt.Stringer
func (sp SyntheticParam) String() string {
	return fmt.Sprintf(`Synthetic:
	MarketDenom: %s`, sp.Denom)
}

// SyntheticParams array of SyntheticParam
type SyntheticParams []SyntheticParam

// String implements fmt.Stringer
func (sps SyntheticParams) String() string {
	out := "Synthetic Params\n"
	for _, sp := range sps {
		out += fmt.Sprintf("%s\n", sp)
	}
	return out
}

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeySyntheticParams, Value: &p.SyntheticParams},
		{Key: KeyNominees, Value: &p.Nominees},
		{Key: KeyFee, Value: &p.Fee},
		{Key: KeyMarketBalanceParam, Value: &p.MarketBalanceParam},
	}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	syntheticDupMap := make(map[string]int)
	for _, sp := range p.SyntheticParams {
		_, found := syntheticDupMap[sp.Denom]
		if found {
			return fmt.Errorf("duplicate collateral denom: %s", sp.Denom)
		}
		syntheticDupMap[sp.Denom] = 1
	}
	return nil
}

type MarketBalanceParam struct {
	BlocksPerSnapshot int           `json:"blocks_per_snapshot" yaml:"blocks_per_snapshot"`
	TimerInterval     time.Duration `json:"timer_interval" yaml:"timer_interval"`
	SnapshotLimit     int           `json:"snapshot_limit" yaml:"snapshot_limit"`
	Coefficients      []sdk.Int     `json:"coefficients" yaml:"coefficients"`
}

func (m MarketBalanceParam) String() string {
	return fmt.Sprintf(`BlocksPerSnapshot: %v
						SnapshotLimit: %v
						TimerInterval: %v`, m.BlocksPerSnapshot, m.SnapshotLimit, m.TimerInterval)
}

package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyMintDenom = []byte("MintDenom")
	KeyParams    = []byte("MintParameters")
)

// TODO Divide into two? One "base class" holding Denom and inflation for use in Genesis and a "subclass" with current state
type InterestAsset struct {
	Denom    string  `json:"denom" yaml:"denom"`
	Interest sdk.Dec `json:"interest" yaml:"interest"`
	Accum    sdk.Dec `json:"accum" yaml:"accum"`
}

type InterestAssets = []InterestAsset

type InterestState struct {
	LastAppliedTime   time.Time      `json:"last_applied" yaml:"last_applied"`
	LastAppliedHeight sdk.Int        `json:"last_applied_height" yaml:"last_applied_height"`
	InterestAssets    InterestAssets `json:"assets" yaml:"assets"`
}

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterType(KeyParams, InterestState{})
}

func NewInterestState(assets ...string) InterestState {
	if len(assets)%2 != 0 {
		panic("Unable to parse asset parameters")
	}

	result := make(InterestAssets, 0)
	for i := 0; i < len(assets); i += 2 {
		interest, err := sdk.NewDecFromStr(assets[i+1])
		if err != nil {
			panic(err)
		}

		result = append(result, InterestAsset{
			Denom:    assets[i],
			Interest: interest,
			Accum:    sdk.NewDec(0),
		})
	}

	return InterestState{
		InterestAssets:    result,
		LastAppliedTime:   time.Now().UTC(),
		LastAppliedHeight: sdk.ZeroInt(),
	}
}

func DefaultInterestState() InterestState {
	return NewInterestState()
}

// validate params
func ValidateInterestState(is InterestState) error {
	// Check for duplicates
	{
		duplicateDenoms := make(map[string]interface{})
		for _, asset := range is.InterestAssets {
			duplicateDenoms[strings.ToLower(asset.Denom)] = true
		}

		if len(duplicateDenoms) != len(is.InterestAssets) {
			return fmt.Errorf("interest parameters contain duplicate denominations")
		}
	}

	// Check for negative interest
	{
		for _, asset := range is.InterestAssets {
			if asset.Interest.IsNegative() {
				return fmt.Errorf("interest parameters contain an asset with negative interest: %v", asset.Denom)
			}
		}
	}

	return nil
}

func (is InterestState) String() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Last inflation: %v\n", is.LastAppliedTime))
	result.WriteString("Interest state:\n")
	for _, asset := range is.InterestAssets {
		result.WriteString(fmt.Sprintf("\tDenom: %v\t\t\tInterest: %v\t\tAccum: %v\n", asset.Denom, asset.Interest, asset.Accum))
	}

	return result.String()
}

func (is *InterestState) FindByDenom(denom string) *InterestAsset {
	for i, a := range is.InterestAssets {
		if a.Denom == denom {
			return &is.InterestAssets[i]
		}
	}
	return nil
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/pkg/matcheng"
)

type Imbalance struct {
	Direction matcheng.Direction `json:"direction" yaml:"direction"`
	Ratio     sdk.Dec            `json:"ratio" yaml:"ratio"`
}

func (i Imbalance) IsImbalanced() bool {
	return i.Ratio.Equal(sdk.ZeroDec())
}

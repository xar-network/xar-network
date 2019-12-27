package types

import (
	"github.com/xar-network/xar-network/pkg/matcheng"
)

type Imbalance struct {
	Direction matcheng.Direction `json:"direction" yaml:"direction"`
	Ratio     float64            `json:"ratio" yaml:"ratio"`
}

func (i Imbalance) IsImbalanced() bool {
	return i.Ratio == 0
}

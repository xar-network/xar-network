package types

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const DefaultParamspace = "swap"

var KeyCoinDenom []byte = []byte("CoinDenom")

type PoolConfig struct {
	CoinDenom string
}

func (config *PoolConfig) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyCoinDenom, Value: &config.CoinDenom},
	}
}

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&PoolConfig{})
}

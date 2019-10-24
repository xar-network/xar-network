package uniswap

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSwap{}, "cosmos-sdk/MsgSwap", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}

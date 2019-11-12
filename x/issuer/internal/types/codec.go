package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIncreaseCredit{}, "issuer/MsgIncreaseCredit", nil)
	cdc.RegisterConcrete(MsgDecreaseCredit{}, "issuer/MsgDecreaseCredit", nil)
	cdc.RegisterConcrete(MsgRevokeLiquidityProvider{}, "issuer/MsgRevokeLiquidityProvider", nil)
	cdc.RegisterConcrete(MsgSetInterest{}, "issuer/MsgSetInterest", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

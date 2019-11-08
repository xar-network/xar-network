package types

import "github.com/cosmos/cosmos-sdk/codec"

var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateIssuer{}, "authority/MsgCreateIssuer", nil)
	cdc.RegisterConcrete(MsgDestroyIssuer{}, "authority/MsgDestroyIssuer", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

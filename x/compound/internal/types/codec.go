package types

import "github.com/cosmos/cosmos-sdk/codec"

// generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}

// RegisterCodec registers concrete types on the codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateOrModifyCSDT{}, "csdt/MsgCreateOrModifyCSDT", nil)
	cdc.RegisterConcrete(MsgTransferCSDT{}, "csdt/MsgTransferCSDT", nil)
}

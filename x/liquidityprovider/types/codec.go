package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&LiquidityProviderAccount{}, "liquidityprovider/LiquidityProviderAccount", nil)
	cdc.RegisterConcrete(MsgMintTokens{}, "liquidityprovider/MsgMintTokens", nil)
	cdc.RegisterConcrete(MsgBurnTokens{}, "liquidityprovider/MsgBurnTokens", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

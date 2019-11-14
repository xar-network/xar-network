package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgMintTokens{}, "liquidityprovider/MsgMintTokens", nil)
	cdc.RegisterConcrete(MsgBurnTokens{}, "liquidityprovider/MsgBurnTokens", nil)

	cdc.RegisterConcrete(&LiquidityProviderAccount{}, "liquidityprovider/LiquidityProviderAccount", nil)
}

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "denominations/MsgIssueToken", nil)
	cdc.RegisterConcrete(MsgMintCoins{}, "denominations/MsgMintCoins", nil)
	cdc.RegisterConcrete(MsgBurnCoins{}, "denominations/MsgBurnCoins", nil)
	cdc.RegisterConcrete(MsgFreezeCoins{}, "denominations/MsgFreezeCoins", nil)
	cdc.RegisterConcrete(MsgUnfreezeCoins{}, "denominations/MsgUnfreezeCoins", nil)

	cdc.RegisterConcrete(&FreezeAccount{}, "denominations/FreezeAccount", nil)
}

func init() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}

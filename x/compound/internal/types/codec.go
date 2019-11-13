package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateCompound{}, "compound/CreateCompound", nil)
	cdc.RegisterConcrete(MsgSupplyMarket{}, "moneymarket/SupplyCompound", nil)
	cdc.RegisterConcrete(MsgBorrowFromMarket{}, "moneymarket/BorrowFromMarket", nil)
	cdc.RegisterConcrete(MsgRedeemFromMarket{}, "moneymarket/RedeemFromMarket", nil)
	cdc.RegisterConcrete(MsgRepayToMarket{}, "moneymarket/RepayToMarket", nil)
}

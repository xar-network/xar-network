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
	cdc.RegisterInterface((*InterestRateModel)(nil), nil)
	cdc.RegisterConcrete(CsdtInterest{}, "csdt/CsdtInterest", nil)
	cdc.RegisterConcrete(MsgCreateOrModifyCSDT{}, "csdt/MsgCreateOrModifyCSDT", nil)
	cdc.RegisterConcrete(MsgDepositCollateral{}, "csdt/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawCollateral{}, "csdt/MsgWithdrawCollateral", nil)
	cdc.RegisterConcrete(MsgSettleDebt{}, "csdt/MsgSettleDebt", nil)
	cdc.RegisterConcrete(MsgWithdrawDebt{}, "csdt/MsgWithdrawDebt", nil)
	cdc.RegisterConcrete(MsgTransferCSDT{}, "csdt/MsgTransferCSDT", nil)
	cdc.RegisterConcrete(MsgAddCollateralParam{}, "csdt/MsgAddCollateralParam", nil)
	cdc.RegisterConcrete(MsgSetCollateralParam{}, "csdt/MsgSetCollateralParam", nil)

}

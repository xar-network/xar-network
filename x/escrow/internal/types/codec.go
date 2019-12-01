package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var MsgCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgLockBox{}, "escrow/MsgLockBox", nil)
	cdc.RegisterConcrete(MsgDepositBox{}, "escrow/MsgDepositBox", nil)
	cdc.RegisterConcrete(MsgFutureBox{}, "escrow/MsgFutureBox", nil)
	cdc.RegisterConcrete(MsgBoxInterestInject{}, "escrow/MsgBoxInterestInject", nil)
	cdc.RegisterConcrete(MsgBoxInterestCancel{}, "escrow/MsgBoxInterestCancel", nil)
	cdc.RegisterConcrete(MsgBoxInject{}, "escrow/MsgBoxInject", nil)
	cdc.RegisterConcrete(MsgBoxInjectCancel{}, "escrow/MsgBoxInjectCancel", nil)
	cdc.RegisterConcrete(MsgBoxWithdraw{}, "escrow/MsgBoxWithdraw", nil)
	cdc.RegisterConcrete(MsgBoxDescription{}, "escrow/MsgBoxDescription", nil)
	cdc.RegisterConcrete(MsgBoxDisableFeature{}, "escrow/MsgBoxDisableFeature", nil)

	cdc.RegisterConcrete(LockBoxInfo{}, "escrow/LockBoxInfo", nil)
	cdc.RegisterConcrete(DepositBoxInfo{}, "escrow/DepositBoxInfo", nil)
	cdc.RegisterConcrete(FutureBoxInfo{}, "escrow/FutureBoxInfo", nil)

	cdc.RegisterInterface((*Box)(nil), nil)
	cdc.RegisterConcrete(&BoxInfo{}, "escrow/BoxInfo", nil)
}

//nolint
func init() {
	RegisterCodec(MsgCdc)
}

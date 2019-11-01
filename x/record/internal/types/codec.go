package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRecord{}, "record/MsgRecord", nil)

	cdc.RegisterInterface((*Record)(nil), nil)
	cdc.RegisterConcrete(&RecordInfo{}, "record/RecordInfo", nil)
}

//nolint
func init() {
	RegisterCodec(ModuleCdc)
}

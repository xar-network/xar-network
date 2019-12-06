/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

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
	cdc.RegisterConcrete(MsgDepositCollateral{}, "csdt/MsgDepositCollateral", nil)
	cdc.RegisterConcrete(MsgWithdrawCollateral{}, "csdt/MsgWithdrawCollateral", nil)
	cdc.RegisterConcrete(MsgSettleDebt{}, "csdt/MsgSettleDebt", nil)
	cdc.RegisterConcrete(MsgWithdrawDebt{}, "csdt/MsgWithdrawDebt", nil)
	cdc.RegisterConcrete(MsgTransferCSDT{}, "csdt/MsgTransferCSDT", nil)
	cdc.RegisterConcrete(MsgAddCollateralParam{}, "csdt/MsgAddCollateralParam", nil)
	cdc.RegisterConcrete(MsgSetCollateralParam{}, "csdt/MsgSetCollateralParam", nil)
}

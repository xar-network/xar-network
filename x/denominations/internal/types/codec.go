/*

Copyright 2016 All in Bits, Inc
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

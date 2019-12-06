/*

Copyright 2016 All in Bits, Inc
Copyright 2018 public-chain
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

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssue{}, "issue/MsgIssue", nil)
	cdc.RegisterConcrete(MsgIssueTransferOwnership{}, "issue/MsgIssueTransferOwnership", nil)
	cdc.RegisterConcrete(MsgIssueDescription{}, "issue/MsgIssueDescription", nil)
	cdc.RegisterConcrete(MsgIssueMint{}, "issue/MsgIssueMint", nil)
	cdc.RegisterConcrete(MsgIssueBurnOwner{}, "issue/MsgIssueBurnOwner", nil)
	cdc.RegisterConcrete(MsgIssueBurnHolder{}, "issue/MsgIssueBurnHolder", nil)
	cdc.RegisterConcrete(MsgIssueBurnFrom{}, "issue/MsgIssueBurnFrom", nil)
	cdc.RegisterConcrete(MsgIssueDisableFeature{}, "issue/MsgIssueDisableFeature", nil)
	cdc.RegisterConcrete(MsgIssueApprove{}, "issue/MsgIssueApprove", nil)
	cdc.RegisterConcrete(MsgIssueSendFrom{}, "issue/MsgIssueSendFrom", nil)
	cdc.RegisterConcrete(MsgIssueIncreaseApproval{}, "issue/MsgIssueIncreaseApproval", nil)
	cdc.RegisterConcrete(MsgIssueDecreaseApproval{}, "issue/MsgIssueDecreaseApproval", nil)
	cdc.RegisterConcrete(MsgIssueFreeze{}, "issue/MsgIssueFreeze", nil)
	cdc.RegisterConcrete(MsgIssueUnFreeze{}, "issue/MsgIssueUnFreeze", nil)

	cdc.RegisterInterface((*Issue)(nil), nil)
	cdc.RegisterConcrete(&CoinIssueInfo{}, "issue/CoinIssueInfo", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

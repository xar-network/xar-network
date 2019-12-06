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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryParams    = "params"
	QueryIssues    = "list"
	QueryIssue     = "query"
	QueryAllowance = "allowance"
	QueryFreeze    = "freeze"
	QueryFreezes   = "freezes"
	QuerySearch    = "search"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "issue"
	// StoreKey is the store key string for issue
	StoreKey = ModuleName
	// RouterKey is the message route for issue
	RouterKey = ModuleName
	// QuerierRoute is the querier route for issue
	QuerierRoute = ModuleName
	// Parameter store default namestore
	DefaultParamspace = ModuleName
)

var (
	CoinMaxTotalSupply, _        = sdk.NewIntFromString("1000000000000000000000000000000000000")
	CoinIssueMaxId        uint64 = 999999999999
	CoinIssueMinId        uint64 = 100000000000
)

const (
	IDPreStr = "xar"
	Custom   = "custom"
)
const (
	TypeMsgIssue                  = "issue"
	TypeMsgIssueMint              = "issue_mint"
	TypeMsgIssueBurnOwner         = "issue_burn_owner"
	TypeMsgIssueBurnHolder        = "issue_burn_holder"
	TypeMsgIssueBurnFrom          = "issue_burn_from"
	TypeMsgIssueDisableFeature    = "issue_disable_feature"
	TypeMsgIssueDescription       = "issue_description"
	TypeMsgIssueTransferOwnership = "issue_transfer_ownership"
	TypeMsgIssueApprove           = "issue_approve"
	TypeMsgIssueSendFrom          = "issue_send_from"
	TypeMsgIssueIncreaseApproval  = "issue_increase_approval"
	TypeMsgIssueDecreaseApproval  = "issue_decrease_approval"
	TypeMsgIssueFreeze            = "issue_freeze"
	TypeMsgIssueUnFreeze          = "issue_unfreeze"
)
const (
	CodeInvalidGenesis       sdk.CodeType = 102
	CoinNameMinLength                     = 3
	CoinNameMaxLength                     = 32
	CoinSymbolMinLength                   = 2
	CoinSymbolMaxLength                   = 8
	CoinDescriptionMaxLength              = 1024
)

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
	// ModuleKey is the name of the module
	ModuleName = "record"
	// StoreKey is the store key string for record
	StoreKey = ModuleName
	// RouterKey is the message route for record
	RouterKey = ModuleName
	// QuerierRoute is the querier route for record
	QuerierRoute = ModuleName
	// Parameter store default namestore
	DefaultParamspace = ModuleName
)
const (
	DefaultCodespace sdk.CodespaceType = ModuleName
)

var (
	RecordMaxId uint64 = 999999999999
	RecordMinId uint64 = 100000000000
)

const (
	IDPreStr = "rec"
	Custom   = "custom"
)
const (
	QueryParams  = "params"
	QueryRecords = "list"
	QueryRecord  = "query"
	QuerySearch  = "search"
)

const (
	TypeMsgRecord = "record"
)
const (
	CodeInvalidGenesis   sdk.CodeType = 102
	HashLength                        = 64
	NameMinLength                     = 3
	NameMaxLength                     = 32
	AuthorMaxLength                   = 64
	RecordTypeMaxLength               = 32
	RecordNoMaxLength                 = 32
	DescriptionMaxLength              = 1024
)

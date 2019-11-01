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

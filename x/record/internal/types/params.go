package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Param for record
type RecordParams struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Hash        string `json:"hash"`
	RecordNo    string `json:"record_number"`
	RecordType  string `json:"record_type"`
	Description string `json:"description"`
}

// Record tags
var (
	TxCategory = "record"

	Action     = sdk.AttributeKeyAction
	Category   = "category"
	Sender     = sdk.AttributeKeySender
	ID         = "id"
	Name       = "name"
	Hash       = "hash"
	Author     = "author"
	RecordType = "record-type"
	RecordNo   = "record-number"
)

// Record interface
type Record interface {
	GetID() string
	SetID(string)

	GetHash() string
	SetHash(string)

	GetRecordNo() string
	SetRecordNo(string)

	GetSender() sdk.AccAddress
	SetSender(sdk.AccAddress)

	GetRecordTime() int64
	SetRecordTime(int64)

	GetName() string
	SetName(string)

	GetAuthor() string
	SetAuthor(string)

	GetRecordType() string
	SetRecordType(string)

	GetDescription() string
	SetDescription(string)

	String() string
}

// Records is an array of Record
type Records []RecordInfo

// Record Info
type RecordInfo struct {
	ID          string         `json:"id"`
	Hash        string         `json:"hash"`
	RecordNo    string         `json:"record_number"`
	Sender      sdk.AccAddress `json:"sender"`
	RecordTime  int64          `json:"record_time"`
	Name        string         `json:"name"`
	Author      string         `json:"author"`
	RecordType  string         `json:"record_type"`
	Description string         `json:"description"`
}

// Implements Record Interface
var _ Record = (*RecordInfo)(nil)

//nolint
func (ci RecordInfo) GetID() string {
	return ci.ID
}
func (ci *RecordInfo) SetID(id string) {
	ci.ID = id
}
func (ci RecordInfo) GetHash() string {
	return ci.Hash
}
func (ci *RecordInfo) SetHash(hash string) {
	ci.Hash = hash
}
func (ci RecordInfo) GetRecordNo() string {
	return ci.RecordNo
}
func (ci *RecordInfo) SetRecordNo(number string) {
	ci.RecordNo = number
}
func (ci RecordInfo) GetSender() sdk.AccAddress {
	return ci.Sender
}
func (ci *RecordInfo) SetSender(operator sdk.AccAddress) {
	ci.Sender = operator
}
func (ci RecordInfo) GetRecordTime() int64 {
	return ci.RecordTime
}
func (ci *RecordInfo) SetRecordTime(time int64) {
	ci.RecordTime = time
}
func (ci RecordInfo) GetName() string {
	return ci.Name
}
func (ci *RecordInfo) SetName(name string) {
	ci.Name = name
}
func (ci RecordInfo) GetAuthor() string {
	return ci.Author
}
func (ci *RecordInfo) SetAuthor(author string) {
	ci.Author = author
}
func (ci RecordInfo) GetRecordType() string {
	return ci.RecordType
}
func (ci *RecordInfo) SetRecordType(recordType string) {
	ci.RecordType = recordType
}
func (ci RecordInfo) GetDescription() string {
	return ci.Description
}
func (ci *RecordInfo) SetDescription(description string) {
	ci.Description = description
}

//nolint
func (ci RecordInfo) String() string {
	return fmt.Sprintf(`Record:
  ID:		          		    %s
  Hash:       		   		    %s
  RecordNo:          		    %s
  Sender:           			%s
  Name:             			%s
  Author:             	 		%s
  Description:           	 	%s
  RecordType:    	    		%s`,
		ci.ID, ci.Hash, ci.RecordNo, ci.Sender.String(), ci.Name, ci.Author, ci.Description, ci.RecordType)
}

//nolint
func (records Records) String() string {
	out := fmt.Sprintf("%-17s|%-17s|%-10s|%-44s|%-10s|%-6s|%-6s|%s\n",
		"Id", "Hash", "RecordNo", "Sender", "Name", "Author", "RecordType", "RecordTime")
	for _, record := range records {
		out += fmt.Sprintf("%-17s|%-17s|%-10s|%-44s|%-10s|%-6s|%-6s|%d\n",
			record.ID, record.Hash, record.RecordNo, record.GetSender().String(), record.Name, record.Author, record.RecordType, record.RecordTime)
	}
	return strings.TrimSpace(out)
}

const (
	CodeRecordExist          sdk.CodeType = 1
	CodeRecordHashNotValid   sdk.CodeType = 2
	CodeRecordIDNotValid     sdk.CodeType = 3
	CodeRecordNumberNotValid sdk.CodeType = 4
	CodeRecordNameNotValid   sdk.CodeType = 5
	CodeRecordAuthorNotValid sdk.CodeType = 6
	CodeRecordTypeNotValid   sdk.CodeType = 7
	CodeDescriptionNotValid  sdk.CodeType = 8
	CodeUnknownRecord        sdk.CodeType = 9
	CodeUnknownAuthor        sdk.CodeType = 10
)

//convert sdk.Error to error
func Errorf(err sdk.Error) error {
	return fmt.Errorf(err.Stacktrace().Error())
}

// Error constructors
func ErrRecordExist(record string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRecordExist, fmt.Sprintf("Record %s exist", record))
}
func ErrRecordHashNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRecordHashNotValid, fmt.Sprintf("Record hash length must be %d character", HashLength))
}
func ErrRecordNumberNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRecordNumberNotValid, fmt.Sprintf("RecordNo. max length is %d character", RecordNoMaxLength))
}
func ErrRecordAuthorNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRecordAuthorNotValid, fmt.Sprintf("Author max length is %d character", AuthorMaxLength))
}
func ErrRecordNameNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRecordNameNotValid, fmt.Sprintf("The length of the name is between %d and %d", NameMinLength, NameMaxLength))
}
func ErrRecordTypeNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRecordTypeNotValid, fmt.Sprintf("The max length of record type is %d", RecordTypeMaxLength))
}
func ErrDescriptionNotValid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDescriptionNotValid, fmt.Sprintf("The max length of description is %d", DescriptionMaxLength))
}
func ErrRecordIDNotValid(recordID string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRecordIDNotValid, fmt.Sprintf("Record-id %s is not a valid recordId", recordID))
}
func ErrUnknownRecord(record string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownRecord, fmt.Sprintf("Unknown record %s", record))
}
func ErrUnknownAuthor(author string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeUnknownAuthor, fmt.Sprintf("Unknown record %s", author))
}

// Param for query record
type RecordQueryParams struct {
	StartRecordId string         `json:"start_record_id"`
	Sender        sdk.AccAddress `json:"sender"`
	Limit         int            `json:"limit"`
}

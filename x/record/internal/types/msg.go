package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgFlag interface {
	sdk.Msg

	GetRecordId() string
	SetRecordId(string)

	GetSender() sdk.AccAddress
	SetSender(sdk.AccAddress)
}

// MsgRecord to allow a registered recordr
// to record new coins.
type MsgRecord struct {
	Sender        sdk.AccAddress `json:"sender" yaml:"sender"`
	*RecordParams `json:"params" yaml:"params"`
}

//New MsgRecord Instance
func NewMsgRecord(sender sdk.AccAddress, params *RecordParams) MsgRecord {
	return MsgRecord{sender, params}
}

// Route Implements Msg.
func (msg MsgRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRecord) Type() string { return TypeMsgRecord }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgRecord) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("Sender address cannot be empty")
	}
	// record hash must be 64 characters length
	res, err := regexp.Match("^[0-9a-zA-Z]{64}$", []byte(msg.RecordParams.Hash))
	if err != nil || !res {
		return ErrRecordHashNotValid()
	}
	if len(msg.Name) < NameMinLength || len(msg.Name) > NameMaxLength {
		return ErrRecordNameNotValid()
	}
	if len(msg.Author) > AuthorMaxLength {
		return ErrRecordAuthorNotValid()
	}
	if len(msg.RecordType) > RecordTypeMaxLength {
		return ErrRecordTypeNotValid()
	}
	if len(msg.RecordNo) > RecordNoMaxLength {
		return ErrRecordNumberNotValid()
	}
	if len(msg.Description) > DescriptionMaxLength {
		return ErrDescriptionNotValid()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgRecord) String() string {
	return fmt.Sprintf("MsgRecord{%s - %s}", "", msg.Sender.String())
}

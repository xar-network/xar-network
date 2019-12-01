package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgLockBox to allow a registered boxr
// to box new coins.
type MsgLockBox struct {
	Sender         sdk.AccAddress `json:"sender"`
	*BoxLockParams `json:"params"`
}

func NewMsgLockBox(sender sdk.AccAddress, params *BoxLockParams) MsgLockBox {
	return MsgLockBox{sender, params}
}

// Route Implements Msg.
func (msg MsgLockBox) Route() string { return RouterKey }

// Type Implements Msg.789
func (msg MsgLockBox) Type() string { return TypeMsgBoxCreate }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgLockBox) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("Sender address cannot be empty")
	}
	if msg.TotalAmount.Token.IsZero() || msg.TotalAmount.Token.Amount.IsNegative() {
		return ErrAmountNotValid("Token amount")
	}
	if len(msg.Name) > BoxNameMaxLength {
		return ErrBoxNameNotValid()
	}
	if len(msg.Description) > BoxDescriptionMaxLength {
		return ErrBoxDescriptionMaxLengthNotValid()
	}
	return nil
}
func (msg MsgLockBox) ValidateService() sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if msg.Lock.EndTime <= time.Now().Unix() {
		return ErrTimeNotValid("EndTime")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgLockBox) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgLockBox) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgLockBox) String() string {
	return fmt.Sprintf("MsgLockBox{%s - %s}", "", msg.Sender.String())
}

// MsgFutureBox to allow a registered boxr
// to box new coins.
type MsgFutureBox struct {
	Sender           sdk.AccAddress `json:"sender"`
	*BoxFutureParams `json:"params"`
}

func NewMsgFutureBox(sender sdk.AccAddress, params *BoxFutureParams) MsgFutureBox {
	return MsgFutureBox{sender, params}
}

// Route Implements Msg.
func (msg MsgFutureBox) Route() string { return RouterKey }

// Type Implements Msg.789
func (msg MsgFutureBox) Type() string { return TypeMsgBoxCreate }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgFutureBox) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("Sender address cannot be empty")
	}
	if msg.TotalAmount.Token.IsZero() || msg.TotalAmount.Token.Amount.IsNegative() {
		return ErrAmountNotValid("Token amount")
	}
	if len(msg.Name) > BoxNameMaxLength {
		return ErrBoxNameNotValid()
	}
	if len(msg.Description) > BoxDescriptionMaxLength {
		return ErrBoxDescriptionMaxLengthNotValid()
	}
	return nil
}

func (msg MsgFutureBox) ValidateService() sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if msg.Future.TimeLine == nil || msg.Future.Receivers == nil ||
		len(msg.Future.TimeLine) == 0 || len(msg.Future.Receivers) == 0 {
		return ErrNotSupportOperation()
	}
	if len(msg.Future.TimeLine) > BoxMaxInstalment {
		return ErrNotEnoughAmount()
	}
	for i, v := range msg.Future.TimeLine {
		if i == 0 {
			if v <= time.Now().Unix() {
				return ErrTimelineNotValid(msg.Future.TimeLine)
			}
			continue
		}
		if v <= msg.Future.TimeLine[i-1] {
			return ErrTimelineNotValid(msg.Future.TimeLine)
		}
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgFutureBox) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgFutureBox) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgFutureBox) String() string {
	return fmt.Sprintf("MsgFutureBox{%s - %s}", "", msg.Sender.String())
}

// MsgDepositBox to allow a registered boxr
// to box new coins.
type MsgDepositBox struct {
	Sender            sdk.AccAddress `json:"sender"`
	*BoxDepositParams `json:"params"`
}

func NewMsgDepositBox(sender sdk.AccAddress, params *BoxDepositParams) MsgDepositBox {
	return MsgDepositBox{sender, params}
}

// Route Implements Msg.
func (msg MsgDepositBox) Route() string { return RouterKey }

// Type Implements Msg.789
func (msg MsgDepositBox) Type() string { return TypeMsgBoxCreate }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgDepositBox) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("Sender address cannot be empty")
	}
	if msg.TotalAmount.Token.IsZero() || msg.TotalAmount.Token.Amount.IsNegative() {
		return ErrAmountNotValid("Token amount")
	}
	if len(msg.Name) > BoxNameMaxLength {
		return ErrBoxNameNotValid()
	}
	if len(msg.Description) > BoxDescriptionMaxLength {
		return ErrBoxDescriptionMaxLengthNotValid()
	}
	return nil
}
func (msg MsgDepositBox) ValidateService() sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	zero := sdk.ZeroInt()
	if msg.Deposit.StartTime <= time.Now().Unix() {
		return ErrTimeNotValid("StartTime")
	}
	if msg.Deposit.EstablishTime <= msg.Deposit.StartTime {
		return ErrTimeNotValid("EstablishTime")
	}
	if msg.Deposit.MaturityTime <= msg.Deposit.EstablishTime {
		return ErrTimeNotValid("MaturityTime")
	}
	if msg.Deposit.BottomLine.LT(zero) || msg.Deposit.BottomLine.GT(msg.TotalAmount.Token.Amount) {
		return ErrAmountNotValid("BottomLine")
	}
	if msg.Deposit.Interest.Token.Amount.LT(zero) {
		return ErrAmountNotValid("Interest")
	}
	if msg.Deposit.Price.LTE(zero) || !msg.TotalAmount.Token.Amount.Mod(msg.Deposit.Price).IsZero() {
		return ErrAmountNotValid("Price")
	}
	if !msg.Deposit.PerCoupon.Equal(utils.CalcInterestRate(msg.TotalAmount.Token.Amount, msg.Deposit.Price,
		msg.Deposit.Interest.Token, msg.Deposit.Interest.Decimals)) {
		return ErrAmountNotValid("PerCoupon")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDepositBox) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgDepositBox) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgDepositBox) String() string {
	return fmt.Sprintf("MsgDepositBox{%s - %s}", "", msg.Sender.String())
}

// MsgBoxWithdraw
type MsgBoxWithdraw struct {
	Id     string         `json:"id"`
	Sender sdk.AccAddress `json:"sender"`
}

//New MsgBoxWithdraw Instance
func NewMsgBoxWithdraw(boxId string, sender sdk.AccAddress) MsgBoxWithdraw {
	return MsgBoxWithdraw{boxId, sender}
}

// Route Implements Msg.
func (msg MsgBoxWithdraw) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBoxWithdraw) Type() string { return TypeMsgBoxWithdraw }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgBoxWithdraw) ValidateBasic() sdk.Error {
	if len(msg.Id) == 0 {
		return ErrUnknownBox("")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBoxWithdraw) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBoxWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBoxWithdraw) String() string {
	return fmt.Sprintf("MsgBoxWithdraw{%s}", msg.Id)
}

//New MsgBoxInterestInject Instance
func NewMsgBoxInterestInject(boxId string, sender sdk.AccAddress, interest sdk.Coin) MsgBoxInterestInject {
	return MsgBoxInterestInject{boxId, sender, interest}
}

// Route Implements Msg.
func (msg MsgBoxInterestInject) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBoxInterestInject) Type() string { return TypeMsgBoxInterestInject }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgBoxInterestInject) ValidateBasic() sdk.Error {
	if len(msg.Id) == 0 {
		return ErrUnknownBox("")
	}
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return ErrAmountNotValid(msg.Amount.Denom)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBoxInterestInject) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBoxInterestInject) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBoxInterestInject) String() string {
	return fmt.Sprintf("MsgBoxInterestInject{%s}", msg.Id)
}

// MsgBoxInterestInject
type MsgBoxInterestInject struct {
	Id     string         `json:"id"`
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Coin       `json:"amount"`
}

// MsgBoxInterestCancel
type MsgBoxInterestCancel struct {
	Id     string         `json:"id"`
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Coin       `json:"amount"`
}

//New MsgBoxInterestCancel Instance
func NewMsgBoxInterestCancel(boxId string, sender sdk.AccAddress, interest sdk.Coin) MsgBoxInterestCancel {
	return MsgBoxInterestCancel{boxId, sender, interest}
}

// Route Implements Msg.
func (msg MsgBoxInterestCancel) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBoxInterestCancel) Type() string { return TypeMsgBoxInterestCancel }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgBoxInterestCancel) ValidateBasic() sdk.Error {
	if len(msg.Id) == 0 {
		return ErrUnknownBox("")
	}
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return ErrAmountNotValid(msg.Amount.Denom)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBoxInterestCancel) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBoxInterestCancel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBoxInterestCancel) String() string {
	return fmt.Sprintf("MsgBoxInterestCancel{%s}", msg.Id)
}

// MsgBoxInject
type MsgBoxInject struct {
	Id     string         `json:"id"`
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Coin       `json:"amount"`
}

//New MsgBoxInject Instance
func NewMsgBoxInject(boxId string, sender sdk.AccAddress, amount sdk.Coin) MsgBoxInject {
	return MsgBoxInject{boxId, sender, amount}
}

// Route Implements Msg.
func (msg MsgBoxInject) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBoxInject) Type() string { return TypeMsgBoxInject }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgBoxInject) ValidateBasic() sdk.Error {
	if len(msg.Id) == 0 {
		return ErrUnknownBox("")
	}
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return ErrAmountNotValid(msg.Amount.Denom)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBoxInject) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBoxInject) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBoxInject) String() string {
	return fmt.Sprintf("MsgBoxInject{%s}", msg.Id)
}

// MsgBoxDisableFeature to allow a registered owner
type MsgBoxDisableFeature struct {
	Id      string         `json:"id"`
	Sender  sdk.AccAddress `json:"sender"`
	Feature string         `json:"feature"`
}

//New MsgBoxDisableFeature Instance
func NewMsgBoxDisableFeature(boxId string, sender sdk.AccAddress, feature string) MsgBoxDisableFeature {
	return MsgBoxDisableFeature{boxId, sender, feature}
}

//nolint
func (ci MsgBoxDisableFeature) GetId() string {
	return ci.Id
}
func (ci MsgBoxDisableFeature) SetId(boxId string) {
	ci.Id = boxId
}
func (ci MsgBoxDisableFeature) GetSender() sdk.AccAddress {
	return ci.Sender
}
func (ci MsgBoxDisableFeature) SetSender(sender sdk.AccAddress) {
	ci.Sender = sender
}
func (ci MsgBoxDisableFeature) GetFeature() string {
	return ci.Feature
}
func (ci MsgBoxDisableFeature) SetFeature(feature string) {
	ci.Feature = feature
}

// Route Implements Msg.
func (msg MsgBoxDisableFeature) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBoxDisableFeature) Type() string { return TypeMsgBoxDisableFeature }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgBoxDisableFeature) ValidateBasic() sdk.Error {
	if len(msg.Id) == 0 {
		return sdk.ErrInvalidAddress("Id cannot be empty")
	}
	_, ok := Features[msg.Feature]
	if !ok {
		return ErrUnknownFeatures()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBoxDisableFeature) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBoxDisableFeature) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBoxDisableFeature) String() string {
	return fmt.Sprintf("MsgBoxDisableFeature{%s}", msg.Id)
}

// MsgBoxDescription to allow a registered owner
// to box new coins.
type MsgBoxDescription struct {
	Id          string         `json:"id"`
	Sender      sdk.AccAddress `json:"sender"`
	Description []byte         `json:"description"`
}

//New MsgBoxDescription Instance
func NewMsgBoxDescription(boxId string, sender sdk.AccAddress, description []byte) MsgBoxDescription {
	return MsgBoxDescription{boxId, sender, description}
}

// Route Implements Msg.
func (msg MsgBoxDescription) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBoxDescription) Type() string { return TypeMsgBoxDescription }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgBoxDescription) ValidateBasic() sdk.Error {
	if len(msg.Id) == 0 {
		return ErrUnknownBox("")
	}
	if len(msg.Description) > BoxDescriptionMaxLength {
		return ErrBoxDescriptionMaxLengthNotValid()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBoxDescription) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBoxDescription) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBoxDescription) String() string {
	return fmt.Sprintf("MsgBoxDescription{%s}", msg.Id)
}

// MsgBoxInjectCancel
type MsgBoxInjectCancel struct {
	Id     string         `json:"id"`
	Sender sdk.AccAddress `json:"sender"`
	Amount sdk.Coin       `json:"amount"`
}

//New MsgBoxInjectCancel Instance
func NewMsgBoxInjectCancel(boxId string, sender sdk.AccAddress, amount sdk.Coin) MsgBoxInjectCancel {
	return MsgBoxInjectCancel{boxId, sender, amount}
}

// Route Implements Msg.
func (msg MsgBoxInjectCancel) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBoxInjectCancel) Type() string { return TypeMsgBoxCancel }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgBoxInjectCancel) ValidateBasic() sdk.Error {
	if len(msg.Id) == 0 {
		return ErrUnknownBox("")
	}
	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return ErrAmountNotValid(msg.Amount.Denom)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBoxInjectCancel) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBoxInjectCancel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgBoxInjectCancel) String() string {
	return fmt.Sprintf("MsgBoxInjectCancel{%s}", msg.Id)
}

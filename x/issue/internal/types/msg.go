package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgIssue to allow a registered issuer
// to issue new coins.
type MsgIssue struct {
	FromAddress  sdk.AccAddress `json:"from_address" yaml:"from_address"`
	*IssueParams `json:"params" yaml:"params"`
}

var _ sdk.Msg = MsgIssue{}

//NewMsgIssue Instance
func NewMsgIssue(fromAddr sdk.AccAddress, params *IssueParams) MsgIssue {
	return MsgIssue{fromAddr, params}
}

// Route Implements Msg.
func (msg MsgIssue) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssue) Type() string { return "issue" }

// ValidateBasic ensures addresses are valid and Coin is positive
func (msg MsgIssue) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	// Cannot issue zero or negative coins
	if msg.TotalSupply.IsZero() || !msg.TotalSupply.IsPositive() {
		return sdk.ErrInvalidCoins("cannot issue 0 or less supply")
	}
	if QuoDecimals(msg.TotalSupply, msg.Decimals).GT(CoinMaxTotalSupply) {
		return ErrCoinTotalSupplyMaxValueNotValid()
	}
	if len(msg.Name) < CoinNameMinLength || len(msg.Name) > CoinNameMaxLength {
		return ErrCoinNamelNotValid()
	}
	if len(msg.Symbol) < CoinSymbolMinLength || len(msg.Symbol) > CoinSymbolMaxLength {
		return ErrCoinSymbolNotValid()
	}
	if msg.Decimals > CoinDecimalsMaxValue {
		return ErrCoinDecimalsMaxValueNotValid()
	}
	if msg.Decimals%CoinDecimalsMultiple != 0 {
		return ErrCoinDecimalsMultipleNotValid()
	}
	if len(msg.Description) > CoinDescriptionMaxLength {
		return ErrCoinDescriptionMaxLengthNotValid()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssue) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssue) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueUnFreeze to allow a registered owner
type MsgIssueUnFreeze struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	FreezeType  string         `json:"freeze_type" yaml:"freeze_type"`
}

var _ sdk.Msg = MsgIssueUnFreeze{}

//New MsgIssueUnFreeze Instance
func NewMsgIssueUnFreeze(issueId string, fromAddr, toAddr sdk.AccAddress, freezeType string) MsgIssueUnFreeze {
	return MsgIssueUnFreeze{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		FreezeType:  freezeType,
	}
}

// Route Implements Msg.
func (msg MsgIssueUnFreeze) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueUnFreeze) Type() string { return "issue_unfreeze" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueUnFreeze) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	_, ok := FreezeTypes[msg.FreezeType]
	if !ok {
		return ErrUnknownFreezeType()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueUnFreeze) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueUnFreeze) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueTransferOwnership to allow a registered owner
// to issue new coins.
type MsgIssueTransferOwnership struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
}

var _ sdk.Msg = MsgIssueTransferOwnership{}

//New MsgIssueTransferOwnership Instance
func NewMsgIssueTransferOwnership(issueId string, fromAddr, toAddr sdk.AccAddress) MsgIssueTransferOwnership {
	return MsgIssueTransferOwnership{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
	}
}

// Route Implements Msg.
func (msg MsgIssueTransferOwnership) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueTransferOwnership) Type() string { return "issue_transfer_ownership" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueTransferOwnership) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("IssueId cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueTransferOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueTransferOwnership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueSendFrom to allow a registered owner
type MsgIssueSendFrom struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	From        sdk.AccAddress `json:"from" yaml:"from"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgIssueSendFrom{}

//New MsgIssueSendFrom Instance
func NewMsgIssueSendFrom(issueId string, fromAddr, from, toAddr sdk.AccAddress, amount sdk.Int) MsgIssueSendFrom {
	return MsgIssueSendFrom{
		IssueId:     issueId,
		FromAddress: fromAddr,
		From:        from,
		ToAddress:   toAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgIssueSendFrom) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueSendFrom) Type() string { return "issue_send_from" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueSendFrom) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if msg.Amount.IsNegative() {
		return sdk.ErrInvalidCoins("can't send negative amount")
	}
	if msg.From.Equals(msg.ToAddress) {
		return sdk.ErrInvalidCoins("can't send yourself")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueSendFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueSendFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueMint to allow a registered issuer
// to issue new coins.
type MsgIssueMint struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
	Decimals    uint           `json:"decimals" yaml:"decimals"`
}

var _ sdk.Msg = MsgIssueMint{}

//NewMsgIssueMint Instance
func NewMsgIssueMint(
	issueId string,
	fromAddr,
	toAddr sdk.AccAddress,
	amount sdk.Int,
	decimals uint,
) MsgIssueMint {
	return MsgIssueMint{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
		Decimals:    decimals,
	}
}

// Route Implements Msg.
func (msg MsgIssueMint) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueMint) Type() string { return "issue_mint" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueMint) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("cannot mint 0 or negative coin amounts")
	}
	if QuoDecimals(msg.Amount, msg.Decimals).GT(CoinMaxTotalSupply) {
		return ErrCoinTotalSupplyMaxValueNotValid()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueIncreaseApproval to allow a registered owner
type MsgIssueIncreaseApproval struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgIssueIncreaseApproval{}

//New MsgIssueIncreaseApproval Instance
func NewMsgIssueIncreaseApproval(
	issueId string,
	fromAddr,
	toAddr sdk.AccAddress,
	amount sdk.Int,
) MsgIssueIncreaseApproval {
	return MsgIssueIncreaseApproval{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgIssueIncreaseApproval) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueIncreaseApproval) Type() string { return "issue_increase_approval" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueIncreaseApproval) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if msg.Amount.IsNegative() {
		return sdk.ErrInvalidCoins("can't approve negative coin amount")
	}
	if msg.FromAddress.Equals(msg.ToAddress) {
		return sdk.ErrInvalidCoins("can't approve yourself")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueIncreaseApproval) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueIncreaseApproval) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueFreeze to allow a registered owner
type MsgIssueFreeze struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	FreezeType  string         `json:"freeze_type" yaml:"freeze_type"`
}

var _ sdk.Msg = MsgIssueFreeze{}

//New MsgIssueFreeze Instance
func NewMsgIssueFreeze(
	issueId string,
	fromAddr,
	toAddr sdk.AccAddress,
	freezeType string,
) MsgIssueFreeze {
	return MsgIssueFreeze{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		FreezeType:  freezeType,
	}
}

// Route Implements Msg.
func (msg MsgIssueFreeze) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueFreeze) Type() string { return "issue_freeze" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueFreeze) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	_, ok := FreezeTypes[msg.FreezeType]
	if !ok {
		return ErrUnknownFreezeType()
	}
	return nil
}
func (msg MsgIssueFreeze) ValidateService() sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueFreeze) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueFreeze) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueDisableFeature to allow a registered owner
type MsgIssueDisableFeature struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	Feature     string         `json:"feature" yaml:"feature"`
}

var _ sdk.Msg = MsgIssueDisableFeature{}

//New MsgIssueDisableFeature Instance
func NewMsgIssueDisableFeature(
	issueId string,
	fromAddr sdk.AccAddress,
	feature string,
) MsgIssueDisableFeature {
	return MsgIssueDisableFeature{
		IssueId:     issueId,
		FromAddress: fromAddr,
		Feature:     feature,
	}
}

// Route Implements Msg.
func (msg MsgIssueDisableFeature) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueDisableFeature) Type() string { return "issue_disable_feature" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueDisableFeature) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	_, ok := Features[msg.Feature]
	if !ok {
		return ErrUnknownFeatures()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueDisableFeature) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueDisableFeature) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueDescription to allow a registered owner
// to issue new coins.
type MsgIssueDescription struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	Description []byte         `json:"description" yaml:"description"`
}

var _ sdk.Msg = MsgIssueDescription{}

//New MsgIssueDescription Instance
func NewMsgIssueDescription(
	issueId string,
	fromAddr sdk.AccAddress,
	description []byte,
) MsgIssueDescription {
	return MsgIssueDescription{
		IssueId:     issueId,
		FromAddress: fromAddr,
		Description: description,
	}
}

// Route Implements Msg.
func (msg MsgIssueDescription) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueDescription) Type() string { return "issue_description" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueDescription) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	if len(msg.Description) > CoinDescriptionMaxLength {
		return ErrCoinDescriptionMaxLengthNotValid()
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueDescription) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueDescription) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueDecreaseApproval to allow a registered owner
type MsgIssueDecreaseApproval struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgIssueDecreaseApproval{}

//New MsgIssueDecreaseApproval Instance
func NewMsgIssueDecreaseApproval(
	issueId string,
	fromAddr, toAddr sdk.AccAddress,
	amount sdk.Int,
) MsgIssueDecreaseApproval {
	return MsgIssueDecreaseApproval{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgIssueDecreaseApproval) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueDecreaseApproval) Type() string { return "issue_decrease_approval" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueDecreaseApproval) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if msg.Amount.IsNegative() {
		return sdk.ErrInvalidCoins("can't approve negative coin amount")
	}
	if msg.FromAddress.Equals(msg.ToAddress) {
		return sdk.ErrInvalidCoins("can't approve yourself")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueDecreaseApproval) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueDecreaseApproval) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueBurnOwner to allow a registered owner
type MsgIssueBurnOwner struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgIssueBurnOwner{}

//New MsgIssueBurnOwner Instance
func NewMsgIssueBurnOwner(
	issueId string,
	fromAddr sdk.AccAddress,
	amount sdk.Int,
) MsgIssueBurnOwner {
	return MsgIssueBurnOwner{
		IssueId:     issueId,
		FromAddress: fromAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgIssueBurnOwner) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueBurnOwner) Type() string { return "issue_burn_owner" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueBurnOwner) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("cannot Burn 0 or negative coin amounts")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueBurnOwner) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueBurnOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueBurnHolder to allow a registered owner
type MsgIssueBurnHolder struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgIssueBurnHolder{}

//New NewMsgIssueBurnHolder Instance
func NewMsgIssueBurnHolder(
	issueId string,
	fromAddr sdk.AccAddress,
	amount sdk.Int,
) MsgIssueBurnHolder {
	return MsgIssueBurnHolder{
		IssueId:     issueId,
		FromAddress: fromAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgIssueBurnHolder) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueBurnHolder) Type() string { return "issue_burn_holder" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueBurnHolder) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("cannot Burn 0 or negative coin amounts")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueBurnHolder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueBurnHolder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueBurnFrom to allow a registered owner
type MsgIssueBurnFrom struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgIssueBurnFrom{}

//New NewMsgIssueBurnFrom Instance
func NewMsgIssueBurnFrom(issueId string,
	fromAddr, toAddr sdk.AccAddress,
	amount sdk.Int,
) MsgIssueBurnFrom {
	return MsgIssueBurnFrom{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgIssueBurnFrom) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueBurnFrom) Type() string { return "issue_burn_from" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueBurnFrom) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("cannot Burn 0 or negative coin amounts")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueBurnFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueBurnFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

// MsgIssueApprove to allow a registered owner
type MsgIssueApprove struct {
	IssueId     string         `json:"issue_id" yaml:"issue_id"`
	FromAddress sdk.AccAddress `json:"from_address" yaml:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address" yaml:"to_address"`
	Amount      sdk.Int        `json:"amount" yaml:"amount"`
}

var _ sdk.Msg = MsgIssueApprove{}

//New MsgIssueApprove Instance
func NewMsgIssueApprove(
	issueId string,
	fromAddr, toAddr sdk.AccAddress,
	amount sdk.Int,
) MsgIssueApprove {
	return MsgIssueApprove{
		IssueId:     issueId,
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgIssueApprove) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgIssueApprove) Type() string { return "issue_approve" }

// Implements Msg. Ensures addresses are valid and Coin is positive
func (msg MsgIssueApprove) ValidateBasic() sdk.Error {
	if len(msg.IssueId) == 0 {
		return sdk.ErrInvalidAddress("issueId cannot be empty")
	}
	// Cannot issue zero or negative coins
	if msg.Amount.IsNegative() {
		return sdk.ErrInvalidCoins("can't approve negative coin amount")
	}
	if msg.FromAddress.Equals(msg.ToAddress) {
		return sdk.ErrInvalidCoins("can't approve yourself")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueApprove) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}

type MsgFlag interface {
	sdk.Msg

	GetIssueId() string
	SetIssueId(string)

	GetSender() sdk.AccAddress
	SetSender(sdk.AccAddress)
}

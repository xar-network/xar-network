package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgIssueToken defines a IssueToken message
type MsgIssueToken struct {
	SourceAddress  sdk.AccAddress `json:"source_address"`
	Name           string         `json:"name"`
	Symbol         string         `json:"symbol"`
	OriginalSymbol string         `json:"original_symbol"`
	TotalSupply    int64          `json:"total_supply"`
	Mintable       bool           `json:"mintable"`
}

// NewMsgIssueToken is a constructor function for MsgIssueToken
func NewMsgIssueToken(sourceAddress sdk.AccAddress, name, symbol string, originalSymbol string,
	totalSupply int64, mintable bool) MsgIssueToken {
	return MsgIssueToken{
		SourceAddress:  sourceAddress,
		Name:           name,
		Symbol:         symbol,
		OriginalSymbol: originalSymbol,
		TotalSupply:    totalSupply,
		Mintable:       mintable,
	}
}

// Route should return the name of the module
func (msg MsgIssueToken) Route() string { return RouterKey }

// Type should return the action
func (msg MsgIssueToken) Type() string { return "issue_token" }

// ValidateBasic runs stateless checks on the message
func (msg MsgIssueToken) ValidateBasic() sdk.Error {
	if msg.SourceAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.SourceAddress.String())
	}
	if len(msg.Name) == 0 || len(msg.Symbol) == 0 || len(msg.OriginalSymbol) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Symbols cannot be empty")
	}
	if msg.TotalSupply < 1 {
		return sdk.ErrUnknownRequest("TotalSupply cannot be less than 1")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgIssueToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SourceAddress}
}

// MsgMintCoins defines the MintCoins message
type MsgMintCoins struct {
	Amount int64          `json:"amount"`
	Symbol string         `json:"symbol"`
	Owner  sdk.AccAddress `json:"owner"`
}

// NewMsgMintCoins is the constructor function for MsgMintCoins
func NewMsgMintCoins(amount int64, symbol string, owner sdk.AccAddress) MsgMintCoins {
	return MsgMintCoins{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgMintCoins) Route() string { return RouterKey }

// Type should return the action
func (msg MsgMintCoins) Type() string { return "mint_coins" }

// ValidateBasic runs stateless checks on the message
func (msg MsgMintCoins) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	if msg.Amount < 1 {
		return sdk.ErrUnknownRequest("Amount cannot be less than 1")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgMintCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgMintCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgBurnCoins defines the BurnCoins message
type MsgBurnCoins struct {
	Amount int64          `json:"amount"`
	Symbol string         `json:"symbol"`
	Owner  sdk.AccAddress `json:"owner"`
}

// NewMsgBurnCoins is the constructor function for MsgBurnCoins
func NewMsgBurnCoins(amount int64, symbol string, owner sdk.AccAddress) MsgBurnCoins {
	return MsgBurnCoins{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgBurnCoins) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBurnCoins) Type() string { return "burn_coins" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBurnCoins) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	if msg.Amount < 1 {
		return sdk.ErrUnknownRequest("Amount cannot be less than 1")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBurnCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBurnCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgFreezeCoins defines the FreezeCoins message
type MsgFreezeCoins struct {
	Amount int64          `json:"amount"`
	Symbol string         `json:"symbol"`
	Owner  sdk.AccAddress `json:"owner"`
}

// NewMsgFreezeCoins is the constructor function for MsgFreezeCoins
func NewMsgFreezeCoins(amount int64, symbol string, owner sdk.AccAddress) MsgFreezeCoins {
	return MsgFreezeCoins{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgFreezeCoins) Route() string { return RouterKey }

// Type should return the action
func (msg MsgFreezeCoins) Type() string { return "freeze_coins" }

// ValidateBasic runs stateless checks on the message
func (msg MsgFreezeCoins) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	if msg.Amount < 1 {
		return sdk.ErrUnknownRequest("Amount cannot be less than 1")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgFreezeCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgFreezeCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgUnfreezeCoins defines the UnfreezeCoins message
type MsgUnfreezeCoins struct {
	Amount int64          `json:"amount"`
	Symbol string         `json:"symbol"`
	Owner  sdk.AccAddress `json:"owner"`
}

// NewMsgUnfreezeCoins is the constructor function for MsgUnfreezeCoins
func NewMsgUnfreezeCoins(amount int64, symbol string, owner sdk.AccAddress) MsgUnfreezeCoins {
	return MsgUnfreezeCoins{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgUnfreezeCoins) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUnfreezeCoins) Type() string { return "unfreeze_coins" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnfreezeCoins) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Symbol cannot be empty")
	}
	if msg.Amount < 1 {
		return sdk.ErrUnknownRequest("Amount cannot be less than 1")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUnfreezeCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnfreezeCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

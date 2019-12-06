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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgIssueToken defines a IssueToken message
type MsgIssueToken struct {
	SourceAddress  sdk.AccAddress `json:"source_address" yaml:"source_address"`
	Owner          sdk.AccAddress `json:"owner" yaml:"owner"`
	Name           string         `json:"name" yaml:"name"`
	Symbol         string         `json:"symbol" yaml:"symbol"`
	OriginalSymbol string         `json:"original_symbol" yaml:"original_symbol"`
	MaxSupply      sdk.Int        `json:"max_supply" yaml:"max_supply"`
	Mintable       bool           `json:"mintable" yaml:"mintable"`
}

// NewMsgIssueToken is a constructor function for MsgIssueToken
func NewMsgIssueToken(sourceAddress, owner sdk.AccAddress, name, symbol string, originalSymbol string, maxSupply sdk.Int, mintable bool) MsgIssueToken {
	return MsgIssueToken{
		SourceAddress:  sourceAddress,
		Owner:          owner,
		Name:           name,
		Symbol:         symbol,
		OriginalSymbol: originalSymbol,
		MaxSupply:      maxSupply,
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
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Symbol) == 0 || len(msg.OriginalSymbol) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Symbols cannot be empty")
	}
	if msg.MaxSupply.LT(sdk.NewInt(1)) {
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
	Amount sdk.Int        `json:"amount" yaml:"amount"`
	Symbol string         `json:"symbol" yaml:"symbol"`
	Owner  sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewMsgMintCoins is the constructor function for MsgMintCoins
func NewMsgMintCoins(amount sdk.Int, symbol string, owner sdk.AccAddress) MsgMintCoins {
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
	if msg.Amount.LT(sdk.NewInt(1)) {
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
	Amount sdk.Int        `json:"amount" yaml:"amount"`
	Symbol string         `json:"symbol" yaml:"symbol"`
	Owner  sdk.AccAddress `json:"owner" yaml:"owner"`
}

// NewMsgBurnCoins is the constructor function for MsgBurnCoins
func NewMsgBurnCoins(amount sdk.Int, symbol string, owner sdk.AccAddress) MsgBurnCoins {
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
	if msg.Amount.LT(sdk.NewInt(1)) {
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
	Amount  sdk.Int        `json:"amount" yaml:"amount"`
	Symbol  string         `json:"symbol" yaml:"symbol"`
	Owner   sdk.AccAddress `json:"owner" yaml:"owner"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgFreezeCoins is the constructor function for MsgFreezeCoins
func NewMsgFreezeCoins(amount sdk.Int, symbol string, owner sdk.AccAddress, address sdk.AccAddress) MsgFreezeCoins {
	return MsgFreezeCoins{
		Amount:  amount,
		Symbol:  symbol,
		Owner:   owner,
		Address: address,
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
	if msg.Amount.LT(sdk.NewInt(1)) {
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
	Amount  sdk.Int        `json:"amount" yaml:"amount"`
	Symbol  string         `json:"symbol" yaml:"symbol"`
	Owner   sdk.AccAddress `json:"owner" yaml:"owner"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgUnfreezeCoins is the constructor function for MsgUnfreezeCoins
func NewMsgUnfreezeCoins(amount sdk.Int, symbol string, owner sdk.AccAddress, address sdk.AccAddress) MsgUnfreezeCoins {
	return MsgUnfreezeCoins{
		Amount:  amount,
		Symbol:  symbol,
		Owner:   owner,
		Address: address,
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
	if msg.Amount.LT(sdk.NewInt(1)) {
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

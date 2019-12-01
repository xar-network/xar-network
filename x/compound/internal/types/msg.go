package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateMarket defines the a new compound struct
type MsgCreateCompound struct {
	Name            string         `json:"name"`
	Denom           string         `json:"denom"`
	InterestRate    sdk.Coins      `json:"interest_tate"`
	Buyer           sdk.AccAddress `json:"buyer"`
	TokenName       string         `json:"token_name"`
	CollateralToken string         `json:"collateral_token"`
}

// NewMsgCreateCompound is the constructor function for MsgCreateCompound
func NewMsgCreateCompound(
	name string,
	denom string,
	interestRate sdk.Coins,
	buyer sdk.AccAddress,
	tokenName string,
	collateralToken string,
) MsgCreateCompound {

	return MsgCreateCompound{
		Name:            name,
		Denom:           denom,
		InterestRate:    interestRate,
		Buyer:           buyer,
		TokenName:       tokenName,
		CollateralToken: collateralToken,
	}
}

// Route should return the name of the module
func (msg MsgCreateCompound) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateCompound) Type() string { return "create_market" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateCompound) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if len(msg.Denom) == 0 {
		return sdk.ErrInsufficientCoins("Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateCompound) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateCompound) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

type MsgSupplyMarket struct {
	Market     string         `json:"compound"`
	LendTokens sdk.Coins      `json:"lend_tokens"`
	Supplier   sdk.AccAddress `json:"supplier"`
}

// NewMsgCreateMarket is the constructor function for MsgBuyName
func NewMsgSupplyMarket(market string, coins sdk.Coins, supplier sdk.AccAddress) MsgSupplyMarket {
	return MsgSupplyMarket{
		Market:     market,
		LendTokens: coins,
		Supplier:   supplier,
	}
}

// Route should return the name of the module
func (msg MsgSupplyMarket) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSupplyMarket) Type() string { return "supply_market" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSupplyMarket) ValidateBasic() sdk.Error {
	if msg.Supplier.Empty() {
		return sdk.ErrInvalidAddress(msg.Supplier.String())
	}
	if len(msg.Market) == 0 {
		return sdk.ErrUnknownRequest("Market cannot be empty")
	}
	if len(msg.LendTokens) == 0 {
		return sdk.ErrInsufficientCoins("You must supply at least one token")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSupplyMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSupplyMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Supplier}
}

//Borrow from Market

type MsgBorrowFromMarket struct {
	Market           string         `json:"market"`
	BorrowTokens     sdk.Coins      `json:"borrowTokens"`
	CollateralTokens sdk.Coins      `json:"collateralTokens"`
	Supplier         sdk.AccAddress `json:"supplier"`
}

// NewMsgCreateMarket is the constructor function for MsgBuyName
func NewMsgBorrowFromMarket(market string, coins sdk.Coins, collateralcoins sdk.Coins, supplier sdk.AccAddress) MsgBorrowFromMarket {
	return MsgBorrowFromMarket{
		Market:           market,
		BorrowTokens:     coins,
		CollateralTokens: collateralcoins,
		Supplier:         supplier,
	}
}

// Route should return the name of the module
func (msg MsgBorrowFromMarket) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBorrowFromMarket) Type() string { return "borrow_from_market" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBorrowFromMarket) ValidateBasic() sdk.Error {
	if msg.Supplier.Empty() {
		return sdk.ErrInvalidAddress(msg.Supplier.String())
	}
	if len(msg.Market) == 0 {
		return sdk.ErrUnknownRequest("Market cannot be empty")
	}
	if len(msg.BorrowTokens) == 0 {
		return sdk.ErrInsufficientCoins("You must borrow at least one token")
	}
	if len(msg.CollateralTokens) == 0 {
		return sdk.ErrInsufficientCoins("You must supply at least one collateral token")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBorrowFromMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBorrowFromMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Supplier}
}

//Redeem from Market

type MsgRedeemFromMarket struct {
	Market       string         `json:"market"`
	RedeemTokens sdk.Coins      `json:"redeemTokens"`
	Supplier     sdk.AccAddress `json:"supplier"`
}

// NewMsgCreateMarket is the constructor function for MsgBuyName
func NewMsgRedeemFromMarket(market string, coins sdk.Coins, supplier sdk.AccAddress) MsgRedeemFromMarket {
	return MsgRedeemFromMarket{
		Market:       market,
		RedeemTokens: coins,
		Supplier:     supplier,
	}
}

// Route should return the name of the module
func (msg MsgRedeemFromMarket) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRedeemFromMarket) Type() string { return "redeem_from_market" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRedeemFromMarket) ValidateBasic() sdk.Error {
	if msg.Supplier.Empty() {
		return sdk.ErrInvalidAddress(msg.Supplier.String())
	}
	if len(msg.Market) == 0 {
		return sdk.ErrUnknownRequest("Market cannot be empty")
	}
	if len(msg.RedeemTokens) == 0 {
		return sdk.ErrInsufficientCoins("You must supply at least one token")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRedeemFromMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRedeemFromMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Supplier}
}

//Repay to Market

type MsgRepayToMarket struct {
	Market      string         `json:"market"`
	RepayTokens sdk.Coins      `json:"repayTokens"`
	Borrower    sdk.AccAddress `json:"supplier"`
}

// NewMsgCreateMarket is the constructor function for MsgBuyName
func NewMsgRepayToMarket(market string, coins sdk.Coins, supplier sdk.AccAddress) MsgRepayToMarket {
	return MsgRepayToMarket{
		Market:      market,
		RepayTokens: coins,
		Borrower:    supplier,
	}
}

// Route should return the name of the module
func (msg MsgRepayToMarket) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRepayToMarket) Type() string { return "repay_to_market" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRepayToMarket) ValidateBasic() sdk.Error {
	if msg.Borrower.Empty() {
		return sdk.ErrInvalidAddress(msg.Borrower.String())
	}
	if len(msg.Market) == 0 {
		return sdk.ErrUnknownRequest("Market cannot be empty")
	}
	if len(msg.RepayTokens) == 0 {
		return sdk.ErrInsufficientCoins("You must supply at least one token")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRepayToMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRepayToMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Borrower}
}

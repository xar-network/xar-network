package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Compound struct {
	Name             string         `json:"name"`
	Denom            string         `json:"denom"`
	Owner            sdk.AccAddress `json:"owner"`
	InterestRate     sdk.Int        `json:"interest_rate"`
	TokenSupply      sdk.Coins      `json:"token_supply"`
	TokenBorrowed    sdk.Coins      `json:"token_borrowed"`
	BorrowCollateral sdk.Coins      `json:"borrow_collateral"`
	TokenName        string         `json:"token_name"`
	CollateralToken  string         `json:"collateral_token"`
}

// Returns a new compound with default rate
func NewCompound() Compound {
	return Compound{
		InterestRate:     sdk.NewInt(0),
		TokenSupply:      sdk.Coins{sdk.NewInt64Coin("uftm", 0)},
		TokenBorrowed:    sdk.Coins{sdk.NewInt64Coin("uftm", 0)},
		BorrowCollateral: sdk.Coins{sdk.NewInt64Coin("ubtc", 0)},
	}
}

// implement fmt.Stringer
func (w Compound) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
		Name: %s
		Symbol: %s
		InterestRate: %s`, w.Owner, w.Name, w.Denom, w.InterestRate))
}

type CompoundPosition struct {
	Owner            sdk.AccAddress `json:"owner"`
	Market           string         `json:"market"`
	LendTokens       sdk.Coins      `json:"lend_tokens"`
	BorrowTokens     sdk.Coins      `json:"borrow_tokens"`
	BorrowCollateral sdk.Coins      `json:"Borrow_collateral"`
}

// Returns a new CompoundPosition with default rate
func NewCompoundPosition() CompoundPosition {
	return CompoundPosition{
		LendTokens:       sdk.Coins{sdk.NewInt64Coin("uftm", 0)},
		BorrowTokens:     sdk.Coins{sdk.NewInt64Coin("uftm", 0)},
		BorrowCollateral: sdk.Coins{sdk.NewInt64Coin("ubtc", 0)},
	}
}

// implement fmt.Stringer
func (w CompoundPosition) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
		Market: %s
		LendTokens: %s
		BorrowTokens: %s`, w.Owner, w.Market, w.LendTokens, w.BorrowTokens))
}

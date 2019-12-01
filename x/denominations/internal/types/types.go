package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// Freezer allows setting and getting frozen coins
type Freezer interface {
	// Get just the frozen coins
	GetFrozenCoins() sdk.Coins
	SetFrozenCoins(sdk.Coins) error

	// Freeze coins by a certain amount. It will reduce amount of coins available from GetCoins()
	FreezeCoins(sdk.Coins) error
	// Unfreeze coins by a certain amount. It will increase the amount of coins available from GetCoins()
	UnfreezeCoins(sdk.Coins) error
}

// CustomCoinAccount extends the built in account interface with extra abilities such as frozen coins
type CustomCoinAccount interface {
	exported.Account
	Freezer
}

// Token is a struct that contains all the metadata of the asset
type Token struct {
	Owner          sdk.AccAddress `json:"owner"`
	Name           string         `json:"name"`            // token name eg Fantom Chain Token
	Symbol         string         `json:"symbol"`          // unique token trade symbol eg FTM-000
	OriginalSymbol string         `json:"original_symbol"` // token symbol eg FTM
	TotalSupply    sdk.Coins      `json:"total_supply"`    // Total token supply
	MaxSupply      sdk.Int        `json:"max_supply"`      // Maximum mintable supply
	Mintable       bool           `json:"mintable"`
}

// NewToken returns a new token
func NewToken(name, symbol, originalSymbol string, totalSupply sdk.Int, maxSupply sdk.Int, owner sdk.AccAddress, mintable bool) *Token {
	return &Token{
		Name:           name,
		Symbol:         symbol,
		OriginalSymbol: originalSymbol,
		TotalSupply:    sdk.Coins{sdk.NewCoin(symbol, totalSupply)},
		MaxSupply:      maxSupply,
		Owner:          owner,
		Mintable:       mintable,
	}
}

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (t Token) ValidateBasic() sdk.Error {
	if t.Owner.Empty() {
		return sdk.ErrInternal(fmt.Sprintf("invalid TokenRecord: Value: %s. Error: Missing Owner", t.Owner))
	}
	if t.Symbol == "" {
		return sdk.ErrInternal(fmt.Sprintf("invalid TokenRecord: Owner: %s. Error: Missing Symbol", t.Symbol))
	}
	if t.TotalSupply == nil || t.TotalSupply.Len() == 0 {
		return sdk.ErrInternal(fmt.Sprintf("invalid TokenRecord: Symbol: %s. Error: Missing TotalSupply", t.TotalSupply))
	}
	if t.MaxSupply.IsZero() {
		return sdk.ErrInternal(fmt.Sprintf("invalid TokenRecord: Symbol: %s. Error: Missing TotalSupply", t.TotalSupply))
	}
	if t.Name == "" {
		return sdk.ErrInternal(fmt.Sprintf("invalid TokenRecord: Symbol: %s. Error: Missing Name", t.Name))
	}
	if t.OriginalSymbol == "" {
		return sdk.ErrInternal(fmt.Sprintf("invalid TokenRecord: Symbol: %s. Error: Missing OriginalSymbol", t.OriginalSymbol))
	}
	return nil
}

// String implements fmt.Stringer
func (t Token) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Name: %s
Symbol: %s
Original Symbol: %s
Total Supply %s
Max Supply %s
Mintable: %v`, t.Owner, t.Name, t.Symbol, t.OriginalSymbol, t.TotalSupply, t.MaxSupply, t.Mintable))
}

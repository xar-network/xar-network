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
	Mintable       bool           `json:"mintable"`
}

// NewToken returns a new token
func NewToken(name, symbol, originalSymbol string, totalSupply int64, owner sdk.AccAddress, mintable bool) *Token {
	return &Token{
		Name:           name,
		Symbol:         symbol,
		OriginalSymbol: originalSymbol,
		TotalSupply:    sdk.Coins{sdk.NewInt64Coin(symbol, totalSupply)},
		Owner:          owner,
		Mintable:       mintable,
	}
}

// String implements fmt.Stringer
func (t Token) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Name: %s
Symbol: %s
Original Symbol: %s
Total Supply %s
Mintable: %v`, t.Owner, t.Name, t.Symbol, t.OriginalSymbol, t.TotalSupply, t.Mintable))
}

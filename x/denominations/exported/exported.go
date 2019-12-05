package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// FreezeAccountI defines an account interface for freeze accounts that hold tokens in an escrow
type FreezeAccountI interface {
	exported.Account

	GetFrozenCoins() sdk.Coins
}

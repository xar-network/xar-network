package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zar-network/zar-network/x/auction"
	"github.com/zar-network/zar-network/x/cdp"
)

type cdpKeeper interface {
	GetCDP(sdk.Context, sdk.AccAddress, string) (cdp.CDP, bool)
	PartialSeizeCDP(sdk.Context, sdk.AccAddress, string, sdk.Int, sdk.Int) sdk.Error
	ReduceGlobalDebt(sdk.Context, sdk.Int) sdk.Error
	GetStableDenom() string // TODO can this be removed somehow?
	GetGovDenom() string
	GetLiquidatorAccountAddress() sdk.AccAddress // This won't need to exist once the module account is defined in this module (instead of in the cdp module)
}

type bankKeeper interface {
	GetCoins(sdk.Context, sdk.AccAddress) sdk.Coins
	AddCoins(sdk.Context, sdk.AccAddress, sdk.Coins) (sdk.Coins, sdk.Error)
	SubtractCoins(sdk.Context, sdk.AccAddress, sdk.Coins) (sdk.Coins, sdk.Error)
}

type auctionKeeper interface {
	StartForwardAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin) (auction.ID, sdk.Error)
	StartReverseAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin) (auction.ID, sdk.Error)
	StartForwardReverseAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin, sdk.AccAddress) (auction.ID, sdk.Error)
}

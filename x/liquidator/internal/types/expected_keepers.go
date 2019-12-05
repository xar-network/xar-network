package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/xar-network/xar-network/x/auction"
	"github.com/xar-network/xar-network/x/csdt"
)

type CsdtKeeper interface {
	GetCSDT(sdk.Context, sdk.AccAddress, string) (csdt.CSDT, bool)
	PartialSeizeCSDT(sdk.Context, sdk.AccAddress, string, sdk.Int, sdk.Int) sdk.Error
	ReduceGlobalDebt(sdk.Context, sdk.Int) sdk.Error
	GetStableDenom() string // TODO can this be removed somehow?
	GetGovDenom() string
}

type BankKeeper interface {
	GetCoins(sdk.Context, sdk.AccAddress) sdk.Coins
}

type AuctionKeeper interface {
	StartForwardAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin) (auction.ID, sdk.Error)
	StartReverseAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin) (auction.ID, sdk.Error)
	StartForwardReverseAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin, sdk.AccAddress) (auction.ID, sdk.Error)
}

type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply exported.SupplyI)
	SetSupply(ctx sdk.Context, supply exported.SupplyI)
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	GetModuleAddress(moduleName string) sdk.AccAddress
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error

	//testing
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)
}

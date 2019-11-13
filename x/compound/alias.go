package compound

import (
	"github.com/xar-network/xar-network/x/compound/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgBuyName      = types.NewMsgBuyName
	NewMsgSetName      = types.NewMsgSetName
	NewMsgCreateMarket = types.NewMsgCreateMarket
	NewMarket          = types.NewMarket
	NewMarketPosition  = types.NewMarketPosition
	NewWhois           = types.NewWhois
	ModuleCdc          = types.ModuleCdc
	RegisterCodec      = types.RegisterCodec
)

type (
	MsgSetName          = types.MsgSetName
	MsgBuyName          = types.MsgBuyName
	MsgCreateMarket     = types.MsgCreateMarket
	MoneyMarket         = types.MoneyMarket
	MarketPosition      = types.MarketPosition
	MsgSupplyMarket     = types.MsgSupplyMarket
	MsgBorrowFromMarket = types.MsgBorrowFromMarket
	MsgRedeemFromMarket = types.MsgRedeemFromMarket
	MsgRepayToMarket    = types.MsgRepayToMarket
	QueryResResolve     = types.QueryResResolve
	QueryResNames       = types.QueryResNames
	Whois               = types.Whois
)

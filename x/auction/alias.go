/**

Baseline from Kava Cosmos Module

**/

package auction

import (
	"github.com/xar-network/xar-network/x/auction/internal/keeper"
	"github.com/xar-network/xar-network/x/auction/internal/types"
)

type (
	Keeper = keeper.Keeper
	ID     = types.ID
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
)

var (
	// functions aliases
	NewIDFromString          = types.NewIDFromString
	NewBaseAuction           = types.NewBaseAuction
	NewForwardAuction        = types.NewForwardAuction
	NewReverseAuction        = types.NewReverseAuction
	NewForwardReverseAuction = types.NewForwardReverseAuction
	RegisterCodec            = types.RegisterCodec
	NewGenesisState          = types.NewGenesisState
	DefaultGenesisState      = types.DefaultGenesisState
	ValidateGenesis          = types.ValidateGenesis
	NewMsgPlaceBid           = types.NewMsgPlaceBid
	NewAuctionParams         = types.NewAuctionParams
	DefaultAuctionParams     = types.DefaultAuctionParams
	ParamKeyTable            = types.ParamKeyTable
	NewKeeper                = keeper.NewKeeper
	NewQuerier               = keeper.NewQuerier

	ModuleCdc = types.ModuleCdc
)

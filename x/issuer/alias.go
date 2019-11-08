package issuer

import (
	"github.com/xar-network/xar-network/x/issuer/internal/keeper"
	"github.com/xar-network/xar-network/x/issuer/internal/types"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	ModuleCdc = types.ModuleCdc
	NewKeeper = keeper.NewKeeper
	NewIssuer = types.NewIssuer
)

type (
	Keeper = keeper.Keeper
	Issuer = types.Issuer
)

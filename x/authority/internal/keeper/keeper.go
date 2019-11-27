package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/x/authority/internal/types"
	"github.com/xar-network/xar-network/x/issuer"
	"github.com/xar-network/xar-network/x/market"
	"github.com/xar-network/xar-network/x/oracle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	keyAuthorityAccAddress = "AuthorityAccountAddress"
)

type Keeper struct {
	storeKey sdk.StoreKey
	ik       issuer.Keeper
	ok       oracle.Keeper
	mk       market.Keeper
	sk       supply.Keeper
}

func NewKeeper(storeKey sdk.StoreKey, issuerKeeper issuer.Keeper, oracleKeeper oracle.Keeper, marketKeeper market.Keeper, supplyKeeper supply.Keeper) Keeper {
	return Keeper{
		ik:       issuerKeeper,
		ok:       oracleKeeper,
		mk:       marketKeeper,
		sk:       supplyKeeper,
		storeKey: storeKey,
	}
}

func (k Keeper) SetAuthority(ctx sdk.Context, authority sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	if store.Get([]byte(keyAuthorityAccAddress)) != nil {
		panic("Authority was already specified")
	}

	bz := types.ModuleCdc.MustMarshalBinaryBare(authority)
	store.Set([]byte(keyAuthorityAccAddress), bz)
}

func (k Keeper) GetAuthority(ctx sdk.Context) (authority sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(keyAuthorityAccAddress))
	types.ModuleCdc.MustUnmarshalBinaryBare(bz, &authority)
	return
}

func (k Keeper) CreateIssuer(ctx sdk.Context, authority sdk.AccAddress, issuerAddress sdk.AccAddress, denoms []string) sdk.Result {
	k.MustBeAuthority(ctx, authority)

	for _, denom := range denoms {
		if !types.ValidateDenom(denom) {
			return types.ErrInvalidDenom(denom).Result()
		}
	}

	i := issuer.NewIssuer(issuerAddress, denoms...)
	return k.ik.AddIssuer(ctx, i)
}

func (k Keeper) DestroyIssuer(ctx sdk.Context, authority sdk.AccAddress, issuerAddress sdk.AccAddress) sdk.Result {
	k.MustBeAuthority(ctx, authority)

	return k.ik.RemoveIssuer(ctx, issuerAddress)
}

func (k Keeper) CreateOracle(ctx sdk.Context, authority sdk.AccAddress, oracle sdk.AccAddress) sdk.Result {
	k.MustBeAuthority(ctx, authority)

	k.ok.AddOracle(ctx, oracle.String())
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) CreateMarket(ctx sdk.Context, authority sdk.AccAddress, baseAsset, quoteAsset string) sdk.Result {
	k.MustBeAuthority(ctx, authority)

	k.mk.Create(ctx, baseAsset, quoteAsset)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) AddSupply(ctx sdk.Context, authority sdk.AccAddress, amount sdk.Coins) sdk.Result {
	k.MustBeAuthority(ctx, authority)

	supply := k.sk.GetSupply(ctx)
	supply = supply.Inflate(amount)
	k.sk.SetSupply(ctx, supply)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) MustBeAuthority(ctx sdk.Context, address sdk.AccAddress) {
	authority := k.GetAuthority(ctx)
	if authority == nil {
		panic(types.ErrNoAuthorityConfigured())
	}

	if authority.Equals(address) {
		return
	}

	panic(types.ErrNotAuthority(address.String()))
}

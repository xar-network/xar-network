package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter
// methods for the various parts of the state machine
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	ak            auth.AccountKeeper
	sk            supply.Keeper
	paramSubspace subspace.Subspace
	codespace     sdk.CodespaceType
}

// NewKeeper creates new instances of the assetmanagement Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc *codec.Codec,
	ak auth.AccountKeeper,
	sk supply.Keeper,
	paramstore subspace.Subspace,
	codespace sdk.CodespaceType,
) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		ak:            ak,
		sk:            sk,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
		codespace:     codespace,
	}
}

// GetToken gets the entire Token metadata struct by symbol. False if not found, true otherwise
func (k Keeper) GetToken(ctx sdk.Context, symbol string) (*types.Token, error) {
	store := ctx.KVStore(k.storeKey)
	if !k.IsSymbolPresent(ctx, symbol) {
		return nil, fmt.Errorf("could not find Token for symbol '%s'", symbol)
	}
	bz := store.Get([]byte(symbol))
	var token types.Token
	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return &token, nil
}

// SetToken sets the entire Token metadata struct by symbol. Owner must be set. Returns success
func (k Keeper) SetToken(ctx sdk.Context, symbol string, token *types.Token) error {
	if token == nil {
		return errors.New("unable to store nil/empty token")
	}
	if token.Owner.Empty() {
		return fmt.Errorf("unable to store token because owner for symbol '%s' is empty", symbol)
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(symbol), k.cdc.MustMarshalBinaryBare(*token))
	return nil
}

// ResolveName - returns the name string that the symbol resolves to
func (k Keeper) ResolveName(ctx sdk.Context, symbol string) (string, error) {
	found, err := k.GetToken(ctx, symbol)
	if err == nil {
		return found.Name, nil
	}
	return "", fmt.Errorf("couldn't resolve name for symbol '%s' because: %s", symbol, err)
}

// SetName - sets the name string that a symbol resolves to
func (k Keeper) SetName(ctx sdk.Context, symbol string, name string) error {
	token, err := k.GetToken(ctx, symbol)
	if err == nil {
		token.Name = name
		return k.SetToken(ctx, symbol, token)
	}
	return fmt.Errorf("failed to set token name for symbol '%s' because: %s", symbol, err)
}

// HasOwner - returns whether or not the symbol already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, symbol string) (bool, error) {
	token, err := k.GetToken(ctx, symbol)
	if err == nil {
		return !token.Owner.Empty(), nil
	}
	return false, fmt.Errorf("unable to check owner for symbol '%s' because: %s", symbol, err)
}

// GetOwner - get the current owner of a symbol
func (k Keeper) GetOwner(ctx sdk.Context, symbol string) (sdk.AccAddress, error) {
	token, err := k.GetToken(ctx, symbol)
	if err == nil {
		return token.Owner, nil
	}
	return nil, fmt.Errorf("unable to get owner for symbol '%s' because: %s", symbol, err)
}

// SetOwner - sets the current owner of a symbol
func (k Keeper) SetOwner(ctx sdk.Context, symbol string, owner sdk.AccAddress) error {
	token, err := k.GetToken(ctx, symbol)
	if err == nil {
		token.Owner = owner
		return k.SetToken(ctx, symbol, token)
	}
	return fmt.Errorf("unable to set owner for symbol '%s' because: %s", symbol, err)
}

// GetTotalSupply - gets the current total supply of a symbol
func (k Keeper) GetTotalSupply(ctx sdk.Context, symbol string) (sdk.Coins, error) {
	token, err := k.GetToken(ctx, symbol)
	if err == nil {
		return token.TotalSupply, nil
	}
	return nil, fmt.Errorf("failed to get total supply for symbol '%s' because: %s", symbol, err)
}

// SetTotalSupply - sets the current total supply of a symbol
func (k Keeper) SetTotalSupply(ctx sdk.Context, symbol string, totalSupply sdk.Coins) error {
	token, err := k.GetToken(ctx, symbol)
	if err == nil {
		token.TotalSupply = totalSupply
		return k.SetToken(ctx, symbol, token)
	}
	return fmt.Errorf("failed to set total supply for symbol '%s' because: %s", symbol, err)
}

// GetTokensIterator - Get an iterator over all symbols in which the keys are the symbols and the values are the token
func (k Keeper) GetTokensIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

// IsSymbolPresent - Check if the symbol is present in the store or not
func (k Keeper) IsSymbolPresent(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(symbol))
}

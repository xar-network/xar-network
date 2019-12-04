package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
func (k Keeper) SetToken(ctx sdk.Context, owner sdk.Address, symbol string, token *types.Token) error {
	if token == nil {
		return errors.New("unable to store nil/empty token")
	}
	if token.Owner.Empty() {
		return fmt.Errorf("unable to store token because owner for symbol '%s' is empty", symbol)
	}
	err := token.ValidateBasic()
	if err != nil {
		return err
	}
	tkn, _ := k.GetToken(ctx, symbol)
	if tkn != nil {
		if !tkn.Owner.Equals(owner) {
			return errors.New("only owner can update the token")
		}
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
func (k Keeper) SetName(ctx sdk.Context, owner sdk.Address, symbol string, name string) error {
	token, err := k.GetToken(ctx, symbol)
	if err != nil {
		return err
	}
	token.Name = name
	return k.SetToken(ctx, owner, symbol, token)
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
func (k Keeper) SetOwner(ctx sdk.Context, owner sdk.Address, symbol string, newOwner sdk.AccAddress) error {
	token, err := k.GetToken(ctx, symbol)
	if err != nil {
		return err
	}
	token.Owner = newOwner
	return k.SetToken(ctx, owner, symbol, token)
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

func (k Keeper) IssueToken(ctx sdk.Context, nominee, owner sdk.AccAddress, token types.Token) sdk.Result {
	if !k.IsNominee(ctx, nominee.String()) {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account is not a nominee %s", nominee.String())).Result()
	}

	if k.IsSymbolPresent(ctx, token.Symbol) {
		return sdk.ErrInvalidCoins(token.Symbol).Result()
	}

	acc := k.ak.GetAccount(ctx, owner)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", owner.String())).Result()
	}

	er := k.SetToken(ctx, owner, token.Symbol, &token)
	if er != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to store new token: '%s'", er)).Result()
	}

	newSymbolLog := fmt.Sprintf("new_symbol=%s", token.Symbol)
	ctx.Logger().Info(newSymbolLog)
	return sdk.Result{
		Log: newSymbolLog,
	}
}

func (k Keeper) MintCoins(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int, denom string) sdk.Result {
	// Check if denom exists
	if !k.IsSymbolPresent(ctx, denom) {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	// Check if owner is assigned
	owner, err := k.GetOwner(ctx, denom)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", denom)).Result()
	}

	// Make sure minter is owner
	if !from.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Are the coins valid
	if !coins.IsValid() {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	// Does the owner exist
	acc := k.ak.GetAccount(ctx, from)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", from.String())).Result()
	}

	token, terr := k.GetToken(ctx, denom)
	if terr != nil {
		return sdk.ErrInvalidCoins(terr.Error()).Result()
	}

	// Is the token mintable
	if !token.Mintable {
		return sdk.ErrInvalidCoins("not mintable").Result()
	}

	// Any overflow

	supply := k.sk.GetSupply(ctx).GetTotal()

	newCoins := supply.Add(coins)
	if newCoins.IsAnyNegative() {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("supply overflow; %s < %s", newCoins.AmountOf(token.Symbol), amount)).Result()
	}

	// Max supply reached
	if token.MaxSupply.LT(newCoins.AmountOf(denom)) {
		return sdk.ErrInvalidCoins("max supply reached").Result()
	}

	er := k.sk.MintCoins(ctx, types.ModuleName, coins)
	if er != nil {
		return er.Result()
	}

	er = k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, from, coins)
	if er != nil {
		return er.Result()
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) BurnCoins(ctx sdk.Context, from sdk.AccAddress, amount sdk.Int, denom string) sdk.Result {
	// Check if denom exists
	if !k.IsSymbolPresent(ctx, denom) {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	// Check if owner is assigned
	owner, err := k.GetOwner(ctx, denom)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", denom)).Result()
	}

	// Make sure minter is owner
	if !from.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Are the coins valid
	if !coins.IsValid() {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	// Does the owner exist
	acc := k.ak.GetAccount(ctx, from)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", from.String())).Result()
	}

	token, terr := k.GetToken(ctx, denom)
	if terr != nil {
		return sdk.ErrInvalidCoins(terr.Error()).Result()
	}

	// Any overflow
	supply := k.sk.GetSupply(ctx).GetTotal()
	newCoins := supply.Sub(coins)
	if newCoins.IsAnyNegative() {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("supply overflow; %s < %s", newCoins.AmountOf(token.Symbol), amount)).Result()
	}

	// Max supply reached
	if token.MaxSupply.LT(newCoins.AmountOf(denom)) {
		return sdk.ErrInvalidCoins("max supply reached").Result()
	}

	_, hasNeg := acc.GetCoins().SafeSub(coins)
	if hasNeg {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient account funds; %s < %s", acc.GetCoins(), coins)).Result()
	}

	er := k.sk.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, coins)
	if er != nil {
		return er.Result()
	}

	er = k.sk.BurnCoins(ctx, types.ModuleName, coins)
	if er != nil {
		return er.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) FreezeCoins(ctx sdk.Context, from sdk.AccAddress, address sdk.AccAddress, amount sdk.Int, denom string) sdk.Result {
	// Check if denom exists
	if !k.IsSymbolPresent(ctx, denom) {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	// Check if owner is assigned
	owner, err := k.GetOwner(ctx, denom)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", denom)).Result()
	}

	// Make sure minter is owner
	if !from.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Are the coins valid
	if !coins.IsValid() {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	acc := k.ak.GetAccount(ctx, address)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", from.String())).Result()
	}

	_, terr := k.GetToken(ctx, denom)
	if terr != nil {
		return sdk.ErrInvalidCoins(terr.Error()).Result()
	}

	_, hasNeg := acc.GetCoins().SafeSub(coins)
	if hasNeg {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("insufficient account funds; %s < %s", acc.GetCoins(), coins)).Result()
	}

	// Todo: Validate you are allowed access to account?
	baseAcc := authtypes.NewBaseAccount(acc.GetAddress(), acc.GetCoins(), acc.GetPubKey(), acc.GetAccountNumber(), acc.GetSequence())
	var freezeAccount, ok = acc.(*types.FreezeAccount)
	if !ok {
		freezeAccount = types.NewFreezeAccount(baseAcc, nil)
	}
	er := freezeAccount.FreezeCoins(sdk.NewCoins(sdk.NewCoin(denom, amount)))
	if er != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to freeze coins: '%s'", err)).Result()
	}

	// Save changes to account
	k.ak.SetAccount(ctx, freezeAccount)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) UnfreezeCoins(ctx sdk.Context, from sdk.AccAddress, address sdk.AccAddress, amount sdk.Int, denom string) sdk.Result {
	// Check if denom exists
	if !k.IsSymbolPresent(ctx, denom) {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	// Check if owner is assigned
	owner, err := k.GetOwner(ctx, denom)
	if err != nil {
		return sdk.ErrUnknownAddress(
			fmt.Sprintf("Could not find the owner for the symbol '%s'", denom)).Result()
	}

	// Make sure minter is owner
	if !from.Equals(owner) { // Checks if the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Are the coins valid
	if !coins.IsValid() {
		return sdk.ErrInvalidCoins(denom).Result()
	}

	// Does the owner exist
	acc := k.ak.GetAccount(ctx, address)
	if acc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", from.String())).Result()
	}

	_, terr := k.GetToken(ctx, denom)
	if terr != nil {
		return sdk.ErrInvalidCoins(terr.Error()).Result()
	}

	newCoins := acc.GetCoins().Add(coins)
	if newCoins.IsAnyNegative() {
		return sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", newCoins, amount),
		).Result()
	}

	// Todo: Validate you are allowed access to account?
	var customAccount, ok = k.ak.GetAccount(ctx, address).(*types.FreezeAccount)
	if !ok {
		return sdk.ErrInternal("failed to get correct account type to unfreeze coins").Result()
	}
	er := customAccount.UnfreezeCoins(sdk.NewCoins(sdk.NewCoin(denom, amount)))
	if er != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to unfreeze coins: '%s'", err)).Result()
	}

	// Save changes to account
	k.ak.SetAccount(ctx, customAccount)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) IsNominee(ctx sdk.Context, nominee string) bool {
	params := k.GetParams(ctx)
	nominees := params.Nominees
	for _, v := range nominees {
		if v == nominee {
			return true
		}
	}
	return false
}

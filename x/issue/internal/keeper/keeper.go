package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

var _ Keeper = (*BaseKeeper)(nil)

type Keeper interface {
	GetBankKeeper() BankKeeper
	AddIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo)
	CreateIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) (sdk.Coins, sdk.Error)
	Fee(ctx sdk.Context, sender sdk.AccAddress, fee sdk.Coin) sdk.Error
	GetIssue(ctx sdk.Context, issueID string) *types.CoinIssueInfo
	GetIssues(ctx sdk.Context, accAddress string) []*types.CoinIssueInfo
	SearchIssues(ctx sdk.Context, symbol string) []*types.CoinIssueInfo
	List(ctx sdk.Context, params types.IssueQueryParams) []*types.CoinIssueInfo
	Iterator(ctx sdk.Context, startIssueId string) sdk.Iterator
	ListAll(ctx sdk.Context) []types.CoinIssueInfo
	DisableFeature(ctx sdk.Context, sender sdk.AccAddress, issueID string, feature string) sdk.Error
	CanMint(ctx sdk.Context, issueID string) bool
	Mint(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress, to sdk.AccAddress) (sdk.Coins, sdk.Error)
	BurnOwner(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress) (sdk.Coins, sdk.Error)
	BurnHolder(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress) (sdk.Coins, sdk.Error)
	BurnFrom(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress, who sdk.AccAddress) (sdk.Coins, sdk.Error)
	GetFreeze(ctx sdk.Context, accAddress sdk.AccAddress, issueID string) types.IssueFreeze
	GetFreezes(ctx sdk.Context, issueID string) []types.IssueAddressFreeze
	Freeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string) sdk.Error
	UnFreeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string) sdk.Error
	SetIssueDescription(ctx sdk.Context, issueID string, sender sdk.AccAddress, description []byte) sdk.Error
	TransferOwnership(ctx sdk.Context, issueID string, sender sdk.AccAddress, to sdk.AccAddress) sdk.Error
	Approve(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error
	IncreaseApproval(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, addedValue sdk.Int) sdk.Error
	DecreaseApproval(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, subtractedValue sdk.Int) sdk.Error
	CheckFreeze(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, issueID string) sdk.Error
	SendFrom(ctx sdk.Context, sender sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error
	Allowance(ctx sdk.Context, owner sdk.AccAddress, spender sdk.AccAddress, issueID string) (amount sdk.Int)
	GetAddressIssues(ctx sdk.Context, accAddress string) (issueIDs []string)
	GetSymbolIssues(ctx sdk.Context, symbol string) (issueIDs []string)
	SetParams(ctx sdk.Context, params types.Params) sdk.Error
	GetParams(ctx sdk.Context) (params types.Params)
	SetInitialIssueStartingIssueId(ctx sdk.Context, issueID uint64) sdk.Error
	GetLastIssueID(ctx sdk.Context) (issueID uint64)
	PeekCurrentIssueID(ctx sdk.Context) (issueID uint64, err sdk.Error)
	SetInterestRate(ctx sdk.Context, sender sdk.AccAddress, interestRate sdk.Dec, issueID string) sdk.Result
}

// expected bank keeper
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}

// expected fee collection keeper
type FeeCollectionKeeper interface {
	AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins
}

// Issue Keeper
type BaseKeeper struct {
	// The reference to the Paramstore to get and set issue specific params
	paramSpace params.Subspace
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// The reference to the CoinKeeper to modify balances
	ck BankKeeper

	// Reserved codespace
	codespace sdk.CodespaceType
	ik        InterestKeeper
}

//New issue keeper Instance
func NewKeeper(key sdk.StoreKey,
	paramSpace params.Subspace, ck BankKeeper, codespace sdk.CodespaceType) BaseKeeper {
	return BaseKeeper{
		storeKey:   key,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		ck:         ck,
		codespace:  codespace,
	}
}

//Keys set
//Set issue
func (k BaseKeeper) setIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyIssuer(coinIssueInfo.IssueId), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(coinIssueInfo))
	return nil
}

func (k BaseKeeper) SetInterestRate(ctx sdk.Context, sender sdk.AccAddress, interestRate sdk.Dec, issueID string) sdk.Result {
	_, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}
	return k.ik.SetInterest(ctx, interestRate, issueID)
}

//Get box bankKeeper
func (keeper BaseKeeper) GetBankKeeper() BankKeeper {
	return keeper.ck
}

//Set address
func (k BaseKeeper) setAddressIssues(ctx sdk.Context, accAddress string, issueIDs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueIDs)
	store.Set(KeyAddressIssues(accAddress), bz)
}

func (k BaseKeeper) deleteAddressIssues(ctx sdk.Context, accAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(KeyAddressIssues(accAddress))
}

//Remove address
func (k BaseKeeper) removeAddressIssues(ctx sdk.Context, accAddress string, issueID string) {
	issueIDs := k.GetAddressIssues(ctx, accAddress)
	for index := 0; index < len(issueIDs); {
		if issueIDs[index] == issueID {
			issueIDs = append(issueIDs[:index], issueIDs[index+1:]...)
			break
		}
		index++
	}
	if len(issueIDs) == 0 {
		k.deleteAddressIssues(ctx, accAddress)
	} else {
		k.setAddressIssues(ctx, accAddress, issueIDs)
	}
}

//Add address
func (k BaseKeeper) addAddressIssues(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) {
	issueIDs := k.GetAddressIssues(ctx, coinIssueInfo.GetOwner().String())
	issueIDs = append(issueIDs, coinIssueInfo.IssueId)
	k.setAddressIssues(ctx, coinIssueInfo.GetOwner().String(), issueIDs)

}

//Set symbol
func (k BaseKeeper) setSymbolIssues(ctx sdk.Context, symbol string, issueIDs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueIDs)
	store.Set(KeySymbolIssues(symbol), bz)
}

//Set freeze
func (k BaseKeeper) setFreeze(ctx sdk.Context, issueID string, accAddress sdk.AccAddress, freeze types.IssueFreeze) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyFreeze(issueID, accAddress), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(freeze))
	return nil
}

//Set approve
func (k BaseKeeper) setApprove(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyAllowed(issueID, sender, spender), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(amount))
	return nil
}

//Keys add
//Add a issue
func (k BaseKeeper) AddIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) {
	k.addAddressIssues(ctx, coinIssueInfo)

	issueIDs := k.GetSymbolIssues(ctx, coinIssueInfo.Symbol)
	issueIDs = append(issueIDs, coinIssueInfo.IssueId)
	k.setSymbolIssues(ctx, coinIssueInfo.Symbol, issueIDs)

	store := ctx.KVStore(k.storeKey)
	bz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(coinIssueInfo)
	store.Set(KeyIssuer(coinIssueInfo.IssueId), bz)
}

//Create a issue
func (k BaseKeeper) CreateIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) (sdk.Coins, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	id, err := k.getNewIssueID(store)
	if err != nil {
		return nil, err
	}
	issueID := KeyIssueIdStr(id)
	coinIssueInfo.IssueTime = ctx.BlockHeader().Time.Unix()
	coinIssueInfo.IssueId = issueID

	k.AddIssue(ctx, coinIssueInfo)

	coin := sdk.Coin{Denom: coinIssueInfo.IssueId, Amount: coinIssueInfo.TotalSupply}
	coins, err := k.ck.AddCoins(ctx, coinIssueInfo.Owner, sdk.NewCoins(coin))

	return coins, err
}

func (k BaseKeeper) Fee(ctx sdk.Context, sender sdk.AccAddress, fee sdk.Coin) sdk.Error {
	if fee.IsZero() || fee.IsNegative() {
		return nil
	}
	_, err := k.GetBankKeeper().SubtractCoins(ctx, sender, sdk.NewCoins(fee))
	if err != nil {
		return types.ErrNotEnoughFee()
	}
	return nil
}

//Returns issue by issueID
func (k BaseKeeper) GetIssue(ctx sdk.Context, issueID string) *types.CoinIssueInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyIssuer(issueID))
	if len(bz) == 0 {
		return nil
	}
	var coinIssueInfo types.CoinIssueInfo
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &coinIssueInfo)
	return &coinIssueInfo
}

//Returns issues by accAddress
func (k BaseKeeper) GetIssues(ctx sdk.Context, accAddress string) []*types.CoinIssueInfo {

	issueIDs := k.GetAddressIssues(ctx, accAddress)
	length := len(issueIDs)
	if length == 0 {
		return []*types.CoinIssueInfo{}
	}
	issues := make([]*types.CoinIssueInfo, 0, length)
	for _, v := range issueIDs {
		issues = append(issues, k.GetIssue(ctx, v))
	}

	return issues
}
func (k BaseKeeper) SearchIssues(ctx sdk.Context, symbol string) []*types.CoinIssueInfo {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, KeySymbolIssues(symbol))
	defer iterator.Close()
	list := make([]*types.CoinIssueInfo, 0, 1)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		if len(bz) == 0 {
			continue
		}
		issueIDs := make([]string, 0, 1)
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueIDs)

		for _, v := range issueIDs {
			list = append(list, k.GetIssue(ctx, v))
		}
	}
	return list
}
func (k BaseKeeper) List(ctx sdk.Context, params types.IssueQueryParams) []*types.CoinIssueInfo {
	if params.Owner != nil && !params.Owner.Empty() {
		return k.GetIssues(ctx, params.Owner.String())
	}
	iterator := k.Iterator(ctx, params.StartIssueId)
	defer iterator.Close()
	list := make([]*types.CoinIssueInfo, 0, params.Limit)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		if len(bz) == 0 {
			continue
		}
		var info types.CoinIssueInfo
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
		list = append(list, &info)
		if len(list) >= params.Limit {
			break
		}
	}
	return list
}
func (k BaseKeeper) Iterator(ctx sdk.Context, startIssueId string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	endIssueId := startIssueId

	if len(startIssueId) == 0 {
		endIssueId = KeyIssueIdStr(types.CoinIssueMaxId)
		startIssueId = KeyIssueIdStr(types.CoinIssueMinId - 1)
	} else {
		startIssueId = KeyIssueIdStr(types.CoinIssueMinId - 1)
	}
	iterator := store.ReverseIterator(KeyIssuer(startIssueId), KeyIssuer(endIssueId))
	return iterator
}
func (k BaseKeeper) ListAll(ctx sdk.Context) []types.CoinIssueInfo {
	iterator := k.Iterator(ctx, "")
	defer iterator.Close()
	list := make([]types.CoinIssueInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		if len(bz) == 0 {
			continue
		}
		var info types.CoinIssueInfo
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
		list = append(list, info)
	}
	return list
}

func (k BaseKeeper) getIssueByOwner(ctx sdk.Context, sender sdk.AccAddress, issueID string) (*types.CoinIssueInfo, sdk.Error) {
	coinIssueInfo := k.GetIssue(ctx, issueID)
	if coinIssueInfo == nil {
		return nil, types.ErrUnknownIssue()
	}
	if !coinIssueInfo.Owner.Equals(sender) {
		return nil, types.ErrOwnerMismatch()
	}
	return coinIssueInfo, nil
}

func (k BaseKeeper) finishMinting(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsMintingFinished() {
		return nil
	}
	coinIssueInfo.MintingFinished = true
	return k.setIssue(ctx, coinIssueInfo)
}

func (k BaseKeeper) DisableFeature(ctx sdk.Context, sender sdk.AccAddress, issueID string, feature string) sdk.Error {
	switch feature {
	case types.BurnOwner:
		return k.disableBurnOwner(ctx, sender, issueID)
	case types.BurnHolder:
		return k.disableBurnHolder(ctx, sender, issueID)
	case types.BurnFrom:
		return k.disableBurnFrom(ctx, sender, issueID)
	case types.Freeze:
		return k.disableFreeze(ctx, sender, issueID)
	case types.Minting:
		return k.finishMinting(ctx, sender, issueID)
	default:
		return types.ErrUnknownFeatures()
	}
}

func (k BaseKeeper) disableBurnOwner(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnOwnerDisabled() {
		return nil
	}
	coinIssueInfo.BurnOwnerDisabled = true
	return k.setIssue(ctx, coinIssueInfo)
}

func (k BaseKeeper) disableBurnHolder(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnHolderDisabled() {
		return nil
	}
	coinIssueInfo.BurnHolderDisabled = true
	return k.setIssue(ctx, coinIssueInfo)
}

func (k BaseKeeper) disableFreeze(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnFromDisabled() {
		return nil
	}
	coinIssueInfo.FreezeDisabled = true
	return k.setIssue(ctx, coinIssueInfo)
}

func (k BaseKeeper) disableBurnFrom(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnFromDisabled() {
		return nil
	}
	coinIssueInfo.BurnFromDisabled = true
	return k.setIssue(ctx, coinIssueInfo)
}

//Can mint a coin
func (k BaseKeeper) CanMint(ctx sdk.Context, issueID string) bool {
	coinIssueInfo := k.GetIssue(ctx, issueID)
	return !coinIssueInfo.MintingFinished
}

//Mint a coin
func (k BaseKeeper) Mint(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress, to sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return nil, err
	}
	if coinIssueInfo.IsMintingFinished() {
		return nil, types.ErrCanNotMint()
	}
	if types.QuoDecimals(coinIssueInfo.TotalSupply.Add(amount), coinIssueInfo.Decimals).GT(types.CoinMaxTotalSupply) {
		return nil, types.ErrCoinTotalSupplyMaxValueNotValid()
	}
	coin := sdk.Coin{Denom: coinIssueInfo.IssueId, Amount: amount}
	coins, err := k.ck.AddCoins(ctx, to, sdk.NewCoins(coin))
	if err != nil {
		return coins, err
	}
	coinIssueInfo.TotalSupply = coinIssueInfo.TotalSupply.Add(amount)
	return coins, k.setIssue(ctx, coinIssueInfo)
}
func (k BaseKeeper) BurnOwner(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return nil, err
	}
	if coinIssueInfo.IsBurnOwnerDisabled() {
		return nil, types.ErrCanNotBurn()
	}
	return k.burn(ctx, coinIssueInfo, amount, sender)
}

//Burn a coin
func (k BaseKeeper) BurnHolder(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo := k.GetIssue(ctx, issueID)
	if coinIssueInfo == nil {
		return nil, types.ErrUnknownIssue()
	}
	if coinIssueInfo.IsBurnHolderDisabled() {
		return nil, types.ErrCanNotBurn()
	}
	return k.burn(ctx, coinIssueInfo, amount, sender)
}
func (k BaseKeeper) burn(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo, amount sdk.Int, who sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coin := sdk.Coin{Denom: coinIssueInfo.IssueId, Amount: amount}
	coins, err := k.ck.SubtractCoins(ctx, who, sdk.NewCoins(coin))
	if err != nil {
		return nil, err
	}
	coinIssueInfo.TotalSupply = coinIssueInfo.TotalSupply.Sub(amount)
	return coins, k.setIssue(ctx, coinIssueInfo)
}

func (k BaseKeeper) BurnFrom(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress, who sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return nil, err
	}
	if who.Equals(coinIssueInfo.GetOwner()) {
		if coinIssueInfo.IsBurnOwnerDisabled() {
			return nil, types.ErrCanNotBurn()
		}
	} else {
		if coinIssueInfo.IsBurnFromDisabled() {
			return nil, types.ErrCanNotBurn()
		}
	}
	return k.burn(ctx, coinIssueInfo, amount, who)
}
func (k BaseKeeper) GetFreeze(ctx sdk.Context, accAddress sdk.AccAddress, issueID string) types.IssueFreeze {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyFreeze(issueID, accAddress))
	if len(bz) == 0 {
		return types.IssueFreeze{false}
	}
	var freeze types.IssueFreeze
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &freeze)
	return freeze
}

func (k BaseKeeper) GetFreezes(ctx sdk.Context, issueID string) []types.IssueAddressFreeze {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, PrefixFreeze(issueID))
	defer iterator.Close()
	list := make([]types.IssueAddressFreeze, 0)
	for ; iterator.Valid(); iterator.Next() {
		var freeze types.IssueFreeze
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &freeze)
		keys := strings.Split(string(iterator.Key()), KeyDelimiter)
		address := keys[len(keys)-1]
		list = append(list, types.IssueAddressFreeze{
			Address: address})
	}
	return list
}
func (k BaseKeeper) freeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string) sdk.Error {
	switch freezeType {
	case types.FreezeIn:
		return k.freezeIn(ctx, issueID, accAddress)
	case types.FreezeOut:
		return k.freezeOut(ctx, issueID, accAddress)
	case types.FreezeInAndOut:
		return k.freezeInAndOut(ctx, issueID, accAddress)
	}
	return types.ErrUnknownFreezeType()
}
func (k BaseKeeper) Freeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string) sdk.Error {
	issueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if issueInfo.IsFreezeDisabled() {
		return types.ErrCanNotFreeze()
	}
	return k.freeze(ctx, issueID, sender, accAddress, freezeType)
}
func (k BaseKeeper) UnFreeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string) sdk.Error {
	_, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	return k.freeze(ctx, issueID, sender, accAddress, freezeType)
}

func (k BaseKeeper) freezeIn(ctx sdk.Context, issueID string, accAddress sdk.AccAddress) sdk.Error {
	freeze := k.GetFreeze(ctx, accAddress, issueID)
	return k.setFreeze(ctx, issueID, accAddress, freeze)
}

func (k BaseKeeper) freezeOut(ctx sdk.Context, issueID string, accAddress sdk.AccAddress) sdk.Error {
	freeze := k.GetFreeze(ctx, accAddress, issueID)
	return k.setFreeze(ctx, issueID, accAddress, freeze)
}

func (k BaseKeeper) freezeInAndOut(ctx sdk.Context, issueID string, accAddress sdk.AccAddress) sdk.Error {
	freeze := k.GetFreeze(ctx, accAddress, issueID)
	return k.setFreeze(ctx, issueID, accAddress, freeze)
}

func (k BaseKeeper) SetIssueDescription(ctx sdk.Context, issueID string, sender sdk.AccAddress, description []byte) sdk.Error {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	coinIssueInfo.Description = string(description)
	return k.setIssue(ctx, coinIssueInfo)
}

//TransferOwnership
func (k BaseKeeper) TransferOwnership(ctx sdk.Context, issueID string, sender sdk.AccAddress, to sdk.AccAddress) sdk.Error {
	coinIssueInfo, err := k.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	coinIssueInfo.Owner = to
	k.removeAddressIssues(ctx, sender.String(), issueID)
	k.addAddressIssues(ctx, coinIssueInfo)
	return k.setIssue(ctx, coinIssueInfo)
}

// Approve the passed address to spend the specified amount of tokens on behalf of sender
func (k BaseKeeper) Approve(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error {
	return k.setApprove(ctx, sender, spender, issueID, amount)
}

//Increase the amount of tokens that an owner allowed to a spender
func (k BaseKeeper) IncreaseApproval(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, addedValue sdk.Int) sdk.Error {
	allowance := k.Allowance(ctx, sender, spender, issueID)
	return k.setApprove(ctx, sender, spender, issueID, allowance.Add(addedValue))
}

//Decrease the amount of tokens that an owner allowed to a spender
func (k BaseKeeper) DecreaseApproval(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, subtractedValue sdk.Int) sdk.Error {
	allowance := k.Allowance(ctx, sender, spender, issueID)
	allowance = allowance.Sub(subtractedValue)
	if allowance.LT(sdk.ZeroInt()) {
		allowance = sdk.ZeroInt()
	}
	return k.setApprove(ctx, sender, spender, issueID, allowance)
}
func (k BaseKeeper) CheckFreeze(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, issueID string) sdk.Error {
	freeze := k.GetFreeze(ctx, from, issueID)
	if freeze.Frozen {
		return types.ErrCanNotTransferOut()
	}
	freeze = k.GetFreeze(ctx, to, issueID)
	if freeze.Frozen {
		return types.ErrCanNotTransferIn()
	}
	return nil
}

//Transfer tokens from one address to another
func (k BaseKeeper) SendFrom(ctx sdk.Context, sender sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error {
	allowance := k.Allowance(ctx, from, sender, issueID)
	if allowance.LT(amount) {
		return types.ErrNotEnoughAmountToTransfer()
	}
	if err := k.CheckFreeze(ctx, from, to, issueID); err != nil {
		return err
	}
	err := k.SendCoins(ctx, from, to, sdk.Coins{sdk.NewCoin(issueID, amount)})
	if err != nil {
		return err
	}
	return k.Approve(ctx, from, sender, issueID, allowance.Sub(amount))
}

//Send coins
func (k BaseKeeper) SendCoins(ctx sdk.Context,
	fromAddr sdk.AccAddress, toAddr sdk.AccAddress,
	amt sdk.Coins) sdk.Error {
	return k.ck.SendCoins(ctx, fromAddr, toAddr, amt)
}

//Get the amount of tokens that an owner allowed to a spender
func (k BaseKeeper) Allowance(ctx sdk.Context, owner sdk.AccAddress, spender sdk.AccAddress, issueID string) (amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyAllowed(issueID, owner, spender))
	if bz == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &amount)
	return amount
}

//Get address from a issue
func (k BaseKeeper) GetAddressIssues(ctx sdk.Context, accAddress string) (issueIDs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyAddressIssues(accAddress))
	if bz == nil {
		return []string{}
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueIDs)
	return issueIDs
}

//Get issueIDs from a issue
func (k BaseKeeper) GetSymbolIssues(ctx sdk.Context, symbol string) (issueIDs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeySymbolIssues(symbol))
	if bz == nil {
		return []string{}
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueIDs)
	return issueIDs
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (k BaseKeeper) SetParams(ctx sdk.Context, params types.Params) sdk.Error {
	if !params.IssueFee.IsValid() {
		return sdk.NewError(k.codespace, types.CodeInvalidGenesis, "invalid issue fee set")
	}
	if !params.MintFee.IsValid() {
		return sdk.NewError(k.codespace, types.CodeInvalidGenesis, "invalid mint fee set")
	}
	k.paramSpace.SetParamSet(ctx, &params)
	return nil
}

// GetParams gets the auth module's parameters.
func (k BaseKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

//Set the initial issueCount
func (k BaseKeeper) SetInitialIssueStartingIssueId(ctx sdk.Context, issueID uint64) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextIssueID)
	if bz != nil {
		return sdk.NewError(k.codespace, types.CodeInvalidGenesis, "Initial IssueId already set")
	}
	bz = types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueID)
	store.Set(KeyNextIssueID, bz)
	return nil
}

// Get the last used issueID
func (k BaseKeeper) GetLastIssueID(ctx sdk.Context) (issueID uint64) {
	issueID, err := k.PeekCurrentIssueID(ctx)
	if err != nil {
		return 0
	}
	issueID--
	return
}

// Gets the next available issueID and increments it
func (k BaseKeeper) getNewIssueID(store sdk.KVStore) (issueID uint64, err sdk.Error) {
	bz := store.Get(KeyNextIssueID)
	if bz == nil {
		return 0, sdk.NewError(k.codespace, types.CodeInvalidGenesis, "InitialIssueID never set")
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueID)
	bz = types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueID + 1)
	store.Set(KeyNextIssueID, bz)
	return issueID, nil
}

// Peeks the next available IssueID without incrementing it
func (k BaseKeeper) PeekCurrentIssueID(ctx sdk.Context) (issueID uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyNextIssueID)
	if bz == nil {
		return 0, sdk.NewError(k.codespace, types.CodeInvalidGenesis, "InitialIssueID never set")
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueID)
	return issueID, nil
}

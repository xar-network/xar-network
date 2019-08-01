package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

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
type Keeper struct {
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper
	// The reference to the Paramstore to get and set issue specific params
	paramSpace params.Subspace
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// The reference to the CoinKeeper to modify balances
	ck BankKeeper
	// The reference to the FeeCollectionKeeper to add fee
	feeCollectionKeeper FeeCollectionKeeper

	// Reserved codespace
	codespace sdk.CodespaceType
}

//Get box bankKeeper
func (keeper Keeper) GetBankKeeper() BankKeeper {
	return keeper.ck
}

//Get box feeCollectionKeeper
func (keeper Keeper) GetFeeCollectionKeeper() FeeCollectionKeeper {
	return keeper.feeCollectionKeeper
}

//New issue keeper Instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramsKeeper params.Keeper,
	paramSpace params.Subspace, ck BankKeeper, feeCollectionKeeper FeeCollectionKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:            key,
		paramsKeeper:        paramsKeeper,
		paramSpace:          paramSpace.WithKeyTable(types.ParamKeyTable()),
		ck:                  ck,
		feeCollectionKeeper: feeCollectionKeeper,
		codespace:           codespace,
	}
}

//Keys set
//Set issue
func (k Keeper) setIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Set(KeyIssuer(coinIssueInfo.IssueId), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(coinIssueInfo))
	return nil
}

//Set address
func (keeper Keeper) setAddressIssues(ctx sdk.Context, accAddress string, issueIDs []string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueIDs)
	store.Set(KeyAddressIssues(accAddress), bz)
}
func (keeper Keeper) deleteAddressIssues(ctx sdk.Context, accAddress string) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyAddressIssues(accAddress))
}

//Remove address
func (keeper Keeper) removeAddressIssues(ctx sdk.Context, accAddress string, issueID string) {
	issueIDs := keeper.GetAddressIssues(ctx, accAddress)
	for index := 0; index < len(issueIDs); {
		if issueIDs[index] == issueID {
			issueIDs = append(issueIDs[:index], issueIDs[index+1:]...)
			break
		}
		index++
	}
	if len(issueIDs) == 0 {
		keeper.deleteAddressIssues(ctx, accAddress)
	} else {
		keeper.setAddressIssues(ctx, accAddress, issueIDs)
	}
}

//Add address
func (keeper Keeper) addAddressIssues(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) {
	issueIDs := keeper.GetAddressIssues(ctx, coinIssueInfo.GetOwner().String())
	issueIDs = append(issueIDs, coinIssueInfo.IssueId)
	keeper.setAddressIssues(ctx, coinIssueInfo.GetOwner().String(), issueIDs)

}

//Set symbol
func (keeper Keeper) setSymbolIssues(ctx sdk.Context, symbol string, issueIDs []string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueIDs)
	store.Set(KeySymbolIssues(symbol), bz)
}

//Set freeze
func (keeper Keeper) setFreeze(ctx sdk.Context, issueID string, accAddress sdk.AccAddress, freeze types.IssueFreeze) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(KeyFreeze(issueID, accAddress), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(freeze))
	return nil
}

//Set approve
func (keeper Keeper) setApprove(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(KeyAllowed(issueID, sender, spender), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(amount))
	return nil
}

//Keys add
//Add a issue
func (keeper Keeper) AddIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) {
	keeper.addAddressIssues(ctx, coinIssueInfo)

	issueIDs := keeper.GetSymbolIssues(ctx, coinIssueInfo.Symbol)
	issueIDs = append(issueIDs, coinIssueInfo.IssueId)
	keeper.setSymbolIssues(ctx, coinIssueInfo.Symbol, issueIDs)

	store := ctx.KVStore(keeper.storeKey)
	bz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(coinIssueInfo)
	store.Set(KeyIssuer(coinIssueInfo.IssueId), bz)
}

//Create a issue
func (keeper Keeper) CreateIssue(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo) (sdk.Coins, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	id, err := keeper.getNewIssueID(store)
	if err != nil {
		return nil, err
	}
	issueID := KeyIssueIdStr(id)
	coinIssueInfo.IssueTime = ctx.BlockHeader().Time.Unix()
	coinIssueInfo.IssueId = issueID

	keeper.AddIssue(ctx, coinIssueInfo)

	coin := sdk.Coin{Denom: coinIssueInfo.IssueId, Amount: coinIssueInfo.TotalSupply}
	coins, err := keeper.ck.AddCoins(ctx, coinIssueInfo.Owner, sdk.NewCoins(coin))

	return coins, err
}

func (keeper Keeper) Fee(ctx sdk.Context, sender sdk.AccAddress, fee sdk.Coin) sdk.Error {
	if fee.IsZero() || fee.IsNegative() {
		return nil
	}
	_, err := keeper.GetBankKeeper().SubtractCoins(ctx, sender, sdk.NewCoins(fee))
	if err != nil {
		return types.ErrNotEnoughFee()
	}
	_ = keeper.GetFeeCollectionKeeper().AddCollectedFees(ctx, sdk.NewCoins(fee))
	return nil
}

//Returns issue by issueID
func (keeper Keeper) GetIssue(ctx sdk.Context, issueID string) *types.CoinIssueInfo {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyIssuer(issueID))
	if len(bz) == 0 {
		return nil
	}
	var coinIssueInfo types.CoinIssueInfo
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &coinIssueInfo)
	return &coinIssueInfo
}

//Returns issues by accAddress
func (keeper Keeper) GetIssues(ctx sdk.Context, accAddress string) []*types.CoinIssueInfo {

	issueIDs := keeper.GetAddressIssues(ctx, accAddress)
	length := len(issueIDs)
	if length == 0 {
		return []*types.CoinIssueInfo{}
	}
	issues := make([]*types.CoinIssueInfo, 0, length)
	for _, v := range issueIDs {
		issues = append(issues, keeper.GetIssue(ctx, v))
	}

	return issues
}
func (keeper Keeper) SearchIssues(ctx sdk.Context, symbol string) []*types.CoinIssueInfo {
	store := ctx.KVStore(keeper.storeKey)

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
			list = append(list, keeper.GetIssue(ctx, v))
		}
	}
	return list
}
func (keeper Keeper) List(ctx sdk.Context, params types.IssueQueryParams) []*types.CoinIssueInfo {
	if params.Owner != nil && !params.Owner.Empty() {
		return keeper.GetIssues(ctx, params.Owner.String())
	}
	iterator := keeper.Iterator(ctx, params.StartIssueId)
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
func (keeper Keeper) Iterator(ctx sdk.Context, startIssueId string) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
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
func (keeper Keeper) ListAll(ctx sdk.Context) []types.CoinIssueInfo {
	iterator := keeper.Iterator(ctx, "")
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

func (keeper Keeper) getIssueByOwner(ctx sdk.Context, sender sdk.AccAddress, issueID string) (*types.CoinIssueInfo, sdk.Error) {
	coinIssueInfo := keeper.GetIssue(ctx, issueID)
	if coinIssueInfo == nil {
		return nil, types.ErrUnknownIssue()
	}
	if !coinIssueInfo.Owner.Equals(sender) {
		return nil, types.ErrOwnerMismatch()
	}
	return coinIssueInfo, nil
}

func (keeper Keeper) finishMinting(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsMintingFinished() {
		return nil
	}
	coinIssueInfo.MintingFinished = true
	return keeper.setIssue(ctx, coinIssueInfo)
}

func (keeper Keeper) DisableFeature(ctx sdk.Context, sender sdk.AccAddress, issueID string, feature string) sdk.Error {
	switch feature {
	case types.BurnOwner:
		return keeper.disableBurnOwner(ctx, sender, issueID)
	case types.BurnHolder:
		return keeper.disableBurnHolder(ctx, sender, issueID)
	case types.BurnFrom:
		return keeper.disableBurnFrom(ctx, sender, issueID)
	case types.Freeze:
		return keeper.disableFreeze(ctx, sender, issueID)
	case types.Minting:
		return keeper.finishMinting(ctx, sender, issueID)
	default:
		return types.ErrUnknownFeatures()
	}
}

func (keeper Keeper) disableBurnOwner(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnOwnerDisabled() {
		return nil
	}
	coinIssueInfo.BurnOwnerDisabled = true
	return keeper.setIssue(ctx, coinIssueInfo)
}

func (keeper Keeper) disableBurnHolder(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnHolderDisabled() {
		return nil
	}
	coinIssueInfo.BurnHolderDisabled = true
	return keeper.setIssue(ctx, coinIssueInfo)
}

func (keeper Keeper) disableFreeze(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnFromDisabled() {
		return nil
	}
	coinIssueInfo.FreezeDisabled = true
	return keeper.setIssue(ctx, coinIssueInfo)
}

func (keeper Keeper) disableBurnFrom(ctx sdk.Context, sender sdk.AccAddress, issueID string) sdk.Error {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if coinIssueInfo.IsBurnFromDisabled() {
		return nil
	}
	coinIssueInfo.BurnFromDisabled = true
	return keeper.setIssue(ctx, coinIssueInfo)
}

//Can mint a coin
func (keeper Keeper) CanMint(ctx sdk.Context, issueID string) bool {
	coinIssueInfo := keeper.GetIssue(ctx, issueID)
	return !coinIssueInfo.MintingFinished
}

//Mint a coin
func (keeper Keeper) Mint(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress, to sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
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
	coins, err := keeper.ck.AddCoins(ctx, to, sdk.NewCoins(coin))
	if err != nil {
		return coins, err
	}
	coinIssueInfo.TotalSupply = coinIssueInfo.TotalSupply.Add(amount)
	return coins, keeper.setIssue(ctx, coinIssueInfo)
}
func (keeper Keeper) BurnOwner(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return nil, err
	}
	if coinIssueInfo.IsBurnOwnerDisabled() {
		return nil, types.ErrCanNotBurn()
	}
	return keeper.burn(ctx, coinIssueInfo, amount, sender)
}

//Burn a coin
func (keeper Keeper) BurnHolder(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo := keeper.GetIssue(ctx, issueID)
	if coinIssueInfo == nil {
		return nil, types.ErrUnknownIssue()
	}
	if coinIssueInfo.IsBurnHolderDisabled() {
		return nil, types.ErrCanNotBurn()
	}
	return keeper.burn(ctx, coinIssueInfo, amount, sender)
}
func (keeper Keeper) burn(ctx sdk.Context, coinIssueInfo *types.CoinIssueInfo, amount sdk.Int, who sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coin := sdk.Coin{Denom: coinIssueInfo.IssueId, Amount: amount}
	coins, err := keeper.ck.SubtractCoins(ctx, who, sdk.NewCoins(coin))
	if err != nil {
		return nil, err
	}
	coinIssueInfo.TotalSupply = coinIssueInfo.TotalSupply.Sub(amount)
	return coins, keeper.setIssue(ctx, coinIssueInfo)
}

func (keeper Keeper) BurnFrom(ctx sdk.Context, issueID string, amount sdk.Int, sender sdk.AccAddress, who sdk.AccAddress) (sdk.Coins, sdk.Error) {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
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
	return keeper.burn(ctx, coinIssueInfo, amount, who)
}
func (keeper Keeper) GetFreeze(ctx sdk.Context, accAddress sdk.AccAddress, issueID string) types.IssueFreeze {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyFreeze(issueID, accAddress))
	if len(bz) == 0 {
		return types.NewIssueFreeze(0, 0)
	}
	var freeze types.IssueFreeze
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &freeze)
	return freeze
}

func (keeper Keeper) GetFreezes(ctx sdk.Context, issueID string) []types.IssueAddressFreeze {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, PrefixFreeze(issueID))
	defer iterator.Close()
	list := make([]types.IssueAddressFreeze, 0)
	for ; iterator.Valid(); iterator.Next() {
		var freeze types.IssueFreeze
		types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &freeze)
		keys := strings.Split(string(iterator.Key()), KeyDelimiter)
		address := keys[len(keys)-1]
		list = append(list, types.IssueAddressFreeze{
			Address:    address,
			OutEndTime: freeze.OutEndTime,
			InEndTime:  freeze.InEndTime})
	}
	return list
}
func (keeper Keeper) freeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string, endTime int64) sdk.Error {
	switch freezeType {
	case types.FreezeIn:
		return keeper.freezeIn(ctx, issueID, accAddress, endTime)
	case types.FreezeOut:
		return keeper.freezeOut(ctx, issueID, accAddress, endTime)
	case types.FreezeInAndOut:
		return keeper.freezeInAndOut(ctx, issueID, accAddress, endTime)
	}
	return types.ErrUnknownFreezeType()
}
func (keeper Keeper) Freeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string, endTime int64) sdk.Error {
	issueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	if issueInfo.IsFreezeDisabled() {
		return types.ErrCanNotFreeze()
	}
	return keeper.freeze(ctx, issueID, sender, accAddress, freezeType, endTime)
}
func (keeper Keeper) UnFreeze(ctx sdk.Context, issueID string, sender sdk.AccAddress, accAddress sdk.AccAddress, freezeType string) sdk.Error {
	_, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	return keeper.freeze(ctx, issueID, sender, accAddress, freezeType, types.UnFreezeEndTime)
}

func (keeper Keeper) freezeIn(ctx sdk.Context, issueID string, accAddress sdk.AccAddress, endTime int64) sdk.Error {
	freeze := keeper.GetFreeze(ctx, accAddress, issueID)
	freeze.InEndTime = endTime
	return keeper.setFreeze(ctx, issueID, accAddress, freeze)
}

func (keeper Keeper) freezeOut(ctx sdk.Context, issueID string, accAddress sdk.AccAddress, endTime int64) sdk.Error {
	freeze := keeper.GetFreeze(ctx, accAddress, issueID)
	freeze.OutEndTime = endTime
	return keeper.setFreeze(ctx, issueID, accAddress, freeze)
}

func (keeper Keeper) freezeInAndOut(ctx sdk.Context, issueID string, accAddress sdk.AccAddress, endTime int64) sdk.Error {
	freeze := keeper.GetFreeze(ctx, accAddress, issueID)
	freeze.InEndTime = endTime
	freeze.OutEndTime = endTime
	return keeper.setFreeze(ctx, issueID, accAddress, freeze)
}

func (keeper Keeper) SetIssueDescription(ctx sdk.Context, issueID string, sender sdk.AccAddress, description []byte) sdk.Error {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	coinIssueInfo.Description = string(description)
	return keeper.setIssue(ctx, coinIssueInfo)
}

//TransferOwnership
func (keeper Keeper) TransferOwnership(ctx sdk.Context, issueID string, sender sdk.AccAddress, to sdk.AccAddress) sdk.Error {
	coinIssueInfo, err := keeper.getIssueByOwner(ctx, sender, issueID)
	if err != nil {
		return err
	}
	coinIssueInfo.Owner = to
	keeper.removeAddressIssues(ctx, sender.String(), issueID)
	keeper.addAddressIssues(ctx, coinIssueInfo)
	return keeper.setIssue(ctx, coinIssueInfo)
}

// Approve the passed address to spend the specified amount of tokens on behalf of sender
func (keeper Keeper) Approve(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error {
	return keeper.setApprove(ctx, sender, spender, issueID, amount)
}

//Increase the amount of tokens that an owner allowed to a spender
func (keeper Keeper) IncreaseApproval(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, addedValue sdk.Int) sdk.Error {
	allowance := keeper.Allowance(ctx, sender, spender, issueID)
	return keeper.setApprove(ctx, sender, spender, issueID, allowance.Add(addedValue))
}

//Decrease the amount of tokens that an owner allowed to a spender
func (keeper Keeper) DecreaseApproval(ctx sdk.Context, sender sdk.AccAddress, spender sdk.AccAddress, issueID string, subtractedValue sdk.Int) sdk.Error {
	allowance := keeper.Allowance(ctx, sender, spender, issueID)
	allowance = allowance.Sub(subtractedValue)
	if allowance.LT(sdk.ZeroInt()) {
		allowance = sdk.ZeroInt()
	}
	return keeper.setApprove(ctx, sender, spender, issueID, allowance)
}
func (keeper Keeper) CheckFreeze(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, issueID string) sdk.Error {
	freeze := keeper.GetFreeze(ctx, from, issueID)
	if freeze.OutEndTime > 0 && freeze.OutEndTime > ctx.BlockHeader().Time.Unix() {
		return types.ErrCanNotTransferOut()
	}
	freeze = keeper.GetFreeze(ctx, to, issueID)
	if freeze.InEndTime > 0 && freeze.InEndTime > ctx.BlockHeader().Time.Unix() {
		return types.ErrCanNotTransferIn()
	}
	return nil
}

//Transfer tokens from one address to another
func (keeper Keeper) SendFrom(ctx sdk.Context, sender sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, issueID string, amount sdk.Int) sdk.Error {
	allowance := keeper.Allowance(ctx, from, sender, issueID)
	if allowance.LT(amount) {
		return types.ErrNotEnoughAmountToTransfer()
	}
	if err := keeper.CheckFreeze(ctx, from, to, issueID); err != nil {
		return err
	}
	err := keeper.SendCoins(ctx, from, to, sdk.Coins{sdk.NewCoin(issueID, amount)})
	if err != nil {
		return err
	}
	return keeper.Approve(ctx, from, sender, issueID, allowance.Sub(amount))
}

//Send coins
func (keeper Keeper) SendCoins(ctx sdk.Context,
	fromAddr sdk.AccAddress, toAddr sdk.AccAddress,
	amt sdk.Coins) sdk.Error {
	return keeper.ck.SendCoins(ctx, fromAddr, toAddr, amt)
}

//Get the amount of tokens that an owner allowed to a spender
func (keeper Keeper) Allowance(ctx sdk.Context, owner sdk.AccAddress, spender sdk.AccAddress, issueID string) (amount sdk.Int) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyAllowed(issueID, owner, spender))
	if bz == nil {
		return sdk.ZeroInt()
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &amount)
	return amount
}

//Get address from a issue
func (keeper Keeper) GetAddressIssues(ctx sdk.Context, accAddress string) (issueIDs []string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyAddressIssues(accAddress))
	if bz == nil {
		return []string{}
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueIDs)
	return issueIDs
}

//Get issueIDs from a issue
func (keeper Keeper) GetSymbolIssues(ctx sdk.Context, symbol string) (issueIDs []string) {
	store := ctx.KVStore(keeper.storeKey)
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
func (ak Keeper) SetParams(ctx sdk.Context, params types.Params) {
	ak.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (ak Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	ak.paramSpace.GetParamSet(ctx, &params)
	return
}

//Set the initial issueCount
func (keeper Keeper) SetInitialIssueStartingIssueId(ctx sdk.Context, issueID uint64) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextIssueID)
	if bz != nil {
		return sdk.NewError(keeper.codespace, types.CodeInvalidGenesis, "Initial IssueId already set")
	}
	bz = types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueID)
	store.Set(KeyNextIssueID, bz)
	return nil
}

// Get the last used issueID
func (keeper Keeper) GetLastIssueID(ctx sdk.Context) (issueID uint64) {
	issueID, err := keeper.PeekCurrentIssueID(ctx)
	if err != nil {
		return 0
	}
	issueID--
	return
}

// Gets the next available issueID and increments it
func (keeper Keeper) getNewIssueID(store sdk.KVStore) (issueID uint64, err sdk.Error) {
	bz := store.Get(KeyNextIssueID)
	if bz == nil {
		return 0, sdk.NewError(keeper.codespace, types.CodeInvalidGenesis, "InitialIssueID never set")
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueID)
	bz = types.ModuleCdc.MustMarshalBinaryLengthPrefixed(issueID + 1)
	store.Set(KeyNextIssueID, bz)
	return issueID, nil
}

// Peeks the next available IssueID without incrementing it
func (keeper Keeper) PeekCurrentIssueID(ctx sdk.Context) (issueID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextIssueID)
	if bz == nil {
		return 0, sdk.NewError(keeper.codespace, types.CodeInvalidGenesis, "InitialIssueID never set")
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(bz, &issueID)
	return issueID, nil
}

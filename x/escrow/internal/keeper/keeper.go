package keeper

import (
	"fmt"

	"github.com/hashgard/hashgard/x/box/config"

	"github.com/hashgard/hashgard/x/box/utils"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/hashgard/hashgard/x/box/errors"
	boxparams "github.com/hashgard/hashgard/x/box/params"
	"github.com/hashgard/hashgard/x/box/types"
	issueerr "github.com/hashgard/hashgard/x/issue/errors"
)

// Box Keeper
type Keeper struct {
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper
	// The reference to the Paramstore to get and set box specific params
	paramSpace params.Subspace
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// The reference to the CoinKeeper to modify balances
	ck BankKeeper
	// The reference to the IssueKeeper to get issue info
	ik IssueKeeper
	// The reference to the FeeCollectionKeeper to add fee
	feeCollectionKeeper FeeCollectionKeeper
	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec
	// Reserved codespace
	codespace sdk.CodespaceType
}

//Get box codec
func (keeper Keeper) Getcdc() *codec.Codec {
	return keeper.cdc
}

//Get box bankKeeper
func (keeper Keeper) GetBankKeeper() BankKeeper {
	return keeper.ck
}

//Get box issueKeeper
func (keeper Keeper) GetIssueKeeper() IssueKeeper {
	return keeper.ik
}

//Get box feeCollectionKeeper
func (keeper Keeper) GetFeeCollectionKeeper() FeeCollectionKeeper {
	return keeper.feeCollectionKeeper
}

//New box keeper Instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramsKeeper params.Keeper,
	paramSpace params.Subspace, ck BankKeeper, ik IssueKeeper,
	feeCollectionKeeper FeeCollectionKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:            key,
		paramsKeeper:        paramsKeeper,
		paramSpace:          paramSpace.WithKeyTable(config.ParamKeyTable()),
		ck:                  ck,
		ik:                  ik,
		feeCollectionKeeper: feeCollectionKeeper,
		cdc:                 cdc,
		codespace:           codespace,
	}
}

func (keeper Keeper) getDepositedCoinsAddress(id string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(fmt.Sprintf("boxDepositedCoins:%s", id))))

}
func (keeper Keeper) SendDepositedCoin(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins, id string) sdk.Error {
	toAddr := keeper.getDepositedCoinsAddress(id)
	return keeper.GetBankKeeper().SendCoins(ctx, fromAddr, toAddr, amt)
}
func (keeper Keeper) CancelDepositedCoin(ctx sdk.Context, toAddr sdk.AccAddress, amt sdk.Coins, id string) sdk.Error {
	fromAddr := keeper.getDepositedCoinsAddress(id)
	return keeper.GetBankKeeper().SendCoins(ctx, fromAddr, toAddr, amt)
}
func (keeper Keeper) SubDepositedCoin(ctx sdk.Context, amt sdk.Coins, id string) sdk.Error {
	_, err := keeper.GetBankKeeper().SubtractCoins(ctx, keeper.getDepositedCoinsAddress(id), amt)
	return err
}
func (keeper Keeper) GetDepositedCoins(ctx sdk.Context, id string) sdk.Coins {
	return keeper.GetBankKeeper().GetCoins(ctx, keeper.getDepositedCoinsAddress(id))
}

//Keys set
//Set box
func (keeper Keeper) setBox(ctx sdk.Context, box *types.BoxInfo) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(KeyBox(box.Id), keeper.cdc.MustMarshalBinaryLengthPrefixed(box))
}

//Set address
func (keeper Keeper) setAddress(ctx sdk.Context, boxType string, accAddress sdk.AccAddress, ids []string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(ids)
	store.Set(KeyAddress(boxType, accAddress), bz)
}

//Set name
func (keeper Keeper) setName(ctx sdk.Context, boxType string, name string, ids []string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(ids)
	store.Set(KeyName(boxType, name), bz)
}

//Keys Add
//Add box
func (keeper Keeper) AddBox(ctx sdk.Context, box *types.BoxInfo) {
	ids := keeper.GetIdsByAddress(ctx, box.BoxType, box.Owner)
	ids = append(ids, box.Id)
	keeper.setAddress(ctx, box.BoxType, box.Owner, ids)

	ids = keeper.GetIdsByName(ctx, box.BoxType, box.Name)
	ids = append(ids, box.Id)
	keeper.setName(ctx, box.BoxType, box.Name, ids)

	keeper.setBox(ctx, box)
}

//Keys remove
//Remove box
//func (keeper Keeper) RemoveBox(ctx sdk.Context, box *types.BoxInfo) {
//	store := ctx.KVStore(keeper.storeKey)
//	store.Delete(KeyName(box.BoxType, box.Name))
//	store.Delete(KeyAddress(box.BoxType, box.Owner))
//	store.Delete(KeyBox(box.Id))
//}
func (keeper Keeper) GetCoinDecimals(ctx sdk.Context, coin sdk.Coin) (uint, sdk.Error) {
	if coin.Denom == types.Agard {
		return types.AgardDecimals, nil
	}
	coinIssueInfo := keeper.GetIssueKeeper().GetIssue(ctx, coin.Denom)
	if coinIssueInfo == nil {
		return 0, issueerr.ErrUnknownIssue(coin.Denom)
	}
	return coinIssueInfo.Decimals, nil
}

func (keeper Keeper) Fee(ctx sdk.Context, sender sdk.AccAddress, fee sdk.Coin) sdk.Error {
	if fee.IsZero() || fee.IsNegative() {
		return nil
	}
	_, err := keeper.GetBankKeeper().SubtractCoins(ctx, sender, sdk.NewCoins(fee))
	if err != nil {
		return errors.ErrNotEnoughFee()
	}
	_ = keeper.GetFeeCollectionKeeper().AddCollectedFees(ctx, sdk.NewCoins(fee))
	return nil
}

//Keys return
//Return box by id
func (keeper Keeper) GetBox(ctx sdk.Context, id string) *types.BoxInfo {
	id = utils.GetIdFromBoxSeqID(id)
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyBox(id))
	if len(bz) == 0 {
		return nil
	}
	var box types.BoxInfo
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &box)
	return &box
}

//Return box by id and and check owner
func (keeper Keeper) GetBoxByOwner(ctx sdk.Context, sender sdk.AccAddress, id string) (*types.BoxInfo, sdk.Error) {
	box := keeper.GetBox(ctx, id)
	if box == nil {
		return nil, errors.ErrUnknownBox(id)
	}
	if !box.Owner.Equals(sender) {
		return nil, errors.ErrOwnerMismatch(id)
	}
	return box, nil
}

//Returns box list by type and accAddress
func (keeper Keeper) GetBoxByAddress(ctx sdk.Context, boxType string, accAddress sdk.AccAddress) []*types.BoxInfo {
	ids := keeper.GetIdsByAddress(ctx, boxType, accAddress)
	length := len(ids)
	if length == 0 {
		return []*types.BoxInfo{}
	}
	boxs := make([]*types.BoxInfo, 0, length)

	for _, v := range ids {
		boxs = append(boxs, keeper.GetBox(ctx, v))
	}
	return boxs
}
func (keeper Keeper) CanTransfer(ctx sdk.Context, id string) sdk.Error {
	if !utils.IsId(id) {
		return nil
	}
	box := keeper.GetBox(ctx, id)
	if box == nil {
		return nil
	}
	if box.IsTransferDisabled() {
		return errors.ErrCanNotTransfer(id)
	}
	return nil
}

//Queries

//Search box by name
func (keeper Keeper) SearchBox(ctx sdk.Context, boxType string, name string) []*types.BoxInfo {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, KeyName(boxType, name))
	defer iterator.Close()
	list := make([]*types.BoxInfo, 0, 1)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		if len(bz) == 0 {
			continue
		}
		ids := make([]string, 0, 1)
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ids)

		for _, v := range ids {
			list = append(list, keeper.GetBox(ctx, v))
		}
	}
	return list
}

//List
func (keeper Keeper) List(ctx sdk.Context, params boxparams.BoxQueryParams) []*types.BoxInfo {
	if params.Owner != nil && !params.Owner.Empty() {
		return keeper.GetBoxByAddress(ctx, params.BoxType, params.Owner)
	}
	iterator := keeper.Iterator(ctx, params.BoxType, params.StartId)
	defer iterator.Close()
	list := make([]*types.BoxInfo, 0, params.Limit)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		if len(bz) == 0 {
			continue
		}
		var info types.BoxInfo
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
		list = append(list, &info)
		if len(list) >= params.Limit {
			break
		}
	}
	return list
}
func (keeper Keeper) ListAll(ctx sdk.Context, boxType string) []types.BoxInfo {
	iterator := keeper.Iterator(ctx, boxType, "")
	defer iterator.Close()
	list := make([]types.BoxInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		if len(bz) == 0 {
			continue
		}
		var info types.BoxInfo
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
		list = append(list, info)
	}
	return list
}
func (keeper Keeper) Iterator(ctx sdk.Context, boxType string, startId string) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	endId := startId
	if len(startId) == 0 {
		endId = KeyIdStr(boxType, types.BoxMaxId)
		startId = KeyIdStr(boxType, types.BoxMinId-1)
	} else {
		startId = KeyIdStr(boxType, types.BoxMinId-1)
	}
	iterator := store.ReverseIterator(KeyBox(startId), KeyBox(endId))
	return iterator
}

//Create a box
func (keeper Keeper) CreateBox(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	decimal, err := keeper.GetCoinDecimals(ctx, box.TotalAmount.Token)
	if err != nil {
		return err
	}
	if box.TotalAmount.Decimals != decimal {
		return errors.ErrDecimalsNotValid(box.TotalAmount.Decimals)
	}
	store := ctx.KVStore(keeper.storeKey)
	id, err := keeper.getNewBoxID(store, box.BoxType)
	if err != nil {
		return err
	}
	box.Id = KeyIdStr(box.BoxType, id)
	box.CreatedTime = ctx.BlockHeader().Time.Unix()

	switch box.BoxType {
	case types.Lock:
		err = keeper.ProcessLockBoxCreate(ctx, box)
	case types.Deposit:
		err = keeper.ProcessDepositBoxCreate(ctx, box)
	case types.Future:
		err = keeper.ProcessFutureBoxCreate(ctx, box)
	default:
		return errors.ErrUnknownBoxType()
	}
	if err != nil {
		return err
	}
	keeper.AddBox(ctx, box)
	return nil
}

func (keeper Keeper) ProcessInjectBox(ctx sdk.Context, id string, sender sdk.AccAddress, amount sdk.Coin, operation string) (*types.BoxInfo, sdk.Error) {
	box := keeper.GetBox(ctx, id)
	if box == nil {
		return nil, errors.ErrUnknownBox(id)
	}
	if types.BoxInjecting != box.Status && types.BoxClosed != box.Status {
		return nil, errors.ErrNotAllowedOperation(box.Status)
	}
	switch box.BoxType {
	case types.Deposit:
		return box, keeper.processDepositBoxInject(ctx, box, sender, amount, operation)
	case types.Future:
		return box, keeper.processFutureBoxInject(ctx, box, sender, amount, operation)
	}
	return nil, errors.ErrUnknownBoxType()
}
func (keeper Keeper) ProcessBoxWithdraw(ctx sdk.Context, id string, sender sdk.AccAddress) (sdk.Int, *types.BoxInfo, sdk.Error) {
	if keeper.GetBankKeeper().GetCoins(ctx, sender).AmountOf(id).IsZero() {
		return sdk.ZeroInt(), nil, errors.ErrNotEnoughAmount()
	}
	boxType := utils.GetBoxTypeByValue(id)
	switch boxType {
	case types.Deposit:
		return keeper.processDepositBoxWithdraw(ctx, id, sender)
	case types.Future:
		boxInfo, err := keeper.processFutureBoxWithdraw(ctx, id, sender)
		return sdk.ZeroInt(), boxInfo, err
	}
	return sdk.ZeroInt(), nil, errors.ErrUnknownBoxType()
}

func (keeper Keeper) SetBoxDescription(ctx sdk.Context, id string, sender sdk.AccAddress, description []byte) (*types.BoxInfo, sdk.Error) {
	box, err := keeper.GetBoxByOwner(ctx, sender, id)
	if err != nil {
		return box, err
	}
	box.Description = string(description)
	keeper.setBox(ctx, box)
	return box, nil
}
func (keeper Keeper) DisableFeature(ctx sdk.Context, sender sdk.AccAddress, id string, feature string) (*types.BoxInfo, sdk.Error) {
	boxInfo, err := keeper.GetBoxByOwner(ctx, sender, id)
	if err != nil {
		return nil, err
	}
	switch feature {
	case types.Transfer:
		return boxInfo, keeper.disableTransfer(ctx, sender, boxInfo)
	default:
		return nil, errors.ErrUnknownFeatures()
	}
}
func (keeper Keeper) disableTransfer(ctx sdk.Context, sender sdk.AccAddress, boxInfo *types.BoxInfo) sdk.Error {
	if boxInfo.GetBoxType() == types.Lock {
		return errors.ErrNotSupportOperation()
	}
	if !boxInfo.IsTransferDisabled() {
		return nil
	}
	boxInfo.TransferDisabled = false
	keeper.setBox(ctx, boxInfo)
	return nil
}

//Send coins
func (keeper Keeper) SendCoins(ctx sdk.Context,
	fromAddr sdk.AccAddress, toAddr sdk.AccAddress,
	amt sdk.Coins) sdk.Error {
	return keeper.ck.SendCoins(ctx, fromAddr, toAddr, amt)
}

//Get name from a box
func (keeper Keeper) GetIdsByName(ctx sdk.Context, boxType string, name string) (ids []string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyName(boxType, name))
	if bz == nil {
		return []string{}
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ids)
	return ids
}

//Get address from a box
func (keeper Keeper) GetIdsByAddress(ctx sdk.Context, boxType string, accAddress sdk.AccAddress) (ids []string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyAddress(boxType, accAddress))
	if bz == nil {
		return []string{}
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ids)
	return ids
}

// BoxQueues

// Returns an iterator for all the box in the Active Queue that expire by time
func (keeper Keeper) ActiveBoxQueueIterator(ctx sdk.Context, endTime int64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(PrefixActiveQueue, sdk.PrefixEndBytes(PrefixActiveBoxQueueTime(endTime)))
}

// Inserts a id into the active box queue at time
func (keeper Keeper) InsertActiveBoxQueue(ctx sdk.Context, endTime int64, boxIdStr string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(boxIdStr)
	store.Set(KeyActiveBoxQueue(endTime, boxIdStr), bz)
}

// removes a id from the Active box Queue
func (keeper Keeper) RemoveFromActiveBoxQueue(ctx sdk.Context, endTime int64, boxIdStr string) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyActiveBoxQueue(endTime, boxIdStr))
}
func (keeper Keeper) RemoveFromActiveBoxQueueByKey(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(key)
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (ak Keeper) SetParams(ctx sdk.Context, params config.Params) {
	ak.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (ak Keeper) GetParams(ctx sdk.Context) (params config.Params) {
	ak.paramSpace.GetParamSet(ctx, &params)
	return
}

//Set the initial boxCount
func (keeper Keeper) SetInitialBoxStartingId(ctx sdk.Context, boxType string, id uint64) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextBoxID(boxType))
	if bz != nil {
		return sdk.NewError(keeper.codespace, types.CodeInvalidGenesis, "Initial Id already set")
	}
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(id)
	store.Set(KeyNextBoxID(boxType), bz)
	return nil
}

// Get the last used id
func (keeper Keeper) GetLastBoxID(ctx sdk.Context, boxType string) (id uint64) {
	id, err := keeper.PeekCurrentBoxID(ctx, boxType)
	if err != nil {
		return 0
	}
	id--
	return
}

// Gets the next available id and increments it
func (keeper Keeper) getNewBoxID(store sdk.KVStore, boxType string) (id uint64, err sdk.Error) {
	bz := store.Get(KeyNextBoxID(boxType))
	if bz == nil {
		return 0, sdk.NewError(keeper.codespace, types.CodeInvalidGenesis, "InitialBoxID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(id + 1)
	store.Set(KeyNextBoxID(boxType), bz)
	return id, nil
}

// Peeks the next available BoxID without incrementing it
func (keeper Keeper) PeekCurrentBoxID(ctx sdk.Context, boxType string) (id uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextBoxID(boxType))
	if bz == nil {
		return 0, sdk.NewError(keeper.codespace, types.CodeInvalidGenesis, "InitialBoxID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	return id, nil
}

func (keeper Keeper) ProcessLockBoxCreate(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	if err := keeper.SendDepositedCoin(ctx, box.Owner, sdk.Coins{box.TotalAmount.Token}, box.Id); err != nil {
		return err
	}
	_, err := keeper.ck.AddCoins(ctx, box.Owner, sdk.Coins{sdk.NewCoin(box.Id, box.TotalAmount.Token.Amount)})
	if err != nil {
		return err
	}
	keeper.InsertActiveBoxQueue(ctx, box.Lock.EndTime, box.Id)
	box.Status = types.LockBoxLocked
	return nil
}
func (keeper Keeper) ProcessLockBoxByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	if box.Status == types.LockBoxUnlocked {
		return nil
	}
	_, err := keeper.ck.SubtractCoins(ctx, box.Owner, sdk.Coins{sdk.NewCoin(box.Id, box.TotalAmount.Token.Amount)})
	if err != nil {
		return err
	}
	if err := keeper.CancelDepositedCoin(ctx, box.Owner, sdk.Coins{box.TotalAmount.Token}, box.Id); err != nil {
		return err
	}
	keeper.RemoveFromActiveBoxQueue(ctx, box.Lock.EndTime, box.Id)
	box.Status = types.LockBoxUnlocked
	keeper.setBox(ctx, box)
	return nil
}

//Process Future box

func (keeper Keeper) ProcessFutureBoxCreate(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	box.Status = types.BoxInjecting
	box.Future.TotalWithdrawal = sdk.ZeroInt()
	//keeper.InsertActiveBoxQueue(ctx, box.Future.TimeLine[0], keeper.GetFutureBoxSeqString(box, 0))
	keeper.InsertActiveBoxQueue(ctx, box.Future.TimeLine[0], box.Id)
	return nil
}
func (keeper Keeper) processFutureBoxInject(ctx sdk.Context, box *types.BoxInfo, sender sdk.AccAddress, amount sdk.Coin, operation string) sdk.Error {
	switch operation {
	case types.Inject:
		return keeper.injectFutureBox(ctx, box, sender, amount)
	case types.Cancel:
		return keeper.cancelDepositFromFutureBox(ctx, box, sender, amount)
	default:
		return errors.ErrUnknownOperation()
	}
}
func (keeper Keeper) injectFutureBox(ctx sdk.Context, box *types.BoxInfo, sender sdk.AccAddress, amount sdk.Coin) sdk.Error {
	if box.Future.TimeLine[0] < ctx.BlockHeader().Time.Unix() {
		return errors.ErrNotSupportOperation()
	}
	if box.TotalAmount.Token.Denom != amount.Denom {
		return errors.ErrAmountNotValid(amount.Denom)
	}
	totalDeposit := sdk.ZeroInt()
	if box.Future.Injects == nil {
		box.Future.Injects = []types.AddressInject{{Address: sender, Amount: amount.Amount}}
	} else {
		exist := false
		for i, v := range box.Future.Injects {
			totalDeposit = totalDeposit.Add(v.Amount)
			if v.Address.Equals(sender) {
				box.Future.Injects[i].Amount = box.Future.Injects[i].Amount.Add(amount.Amount)
				exist = true
			}
		}
		if !exist {
			box.Future.Injects = append(box.Future.Injects, types.NewAddressInject(sender, amount.Amount))
		}
	}
	totalDeposit = totalDeposit.Add(amount.Amount)
	if totalDeposit.GT(box.TotalAmount.Token.Amount) {
		return errors.ErrNotEnoughAmount()
	}
	if err := keeper.SendDepositedCoin(ctx, sender, sdk.Coins{amount}, box.Id); err != nil {
		return err
	}
	if totalDeposit.Equal(box.TotalAmount.Token.Amount) {
		if err := keeper.processFutureBoxDistribute(ctx, box); err != nil {
			return err
		}
		box.Status = types.BoxActived
		//keeper.RemoveFromActiveBoxQueue(ctx, box.Future.TimeLine[0], keeper.GetFutureBoxSeqString(box, 0))
		keeper.RemoveFromActiveBoxQueue(ctx, box.Future.TimeLine[0], box.Id)
	}
	keeper.setBox(ctx, box)
	return nil
}
func (keeper Keeper) cancelDepositFromFutureBox(ctx sdk.Context, box *types.BoxInfo, sender sdk.AccAddress, amount sdk.Coin) sdk.Error {
	if box.Status == types.BoxActived {
		return errors.ErrNotAllowedOperation(box.Status)
	}
	if box.TotalAmount.Token.Denom != amount.Denom {
		return errors.ErrAmountNotValid(amount.Denom)
	}
	if box.Future.Injects == nil {
		return errors.ErrNotEnoughAmount()
	}
	exist := false
	for i, v := range box.Future.Injects {
		if v.Address.Equals(sender) {
			if box.Future.Injects[i].Amount.LT(amount.Amount) {
				return errors.ErrNotEnoughAmount()
			}
			box.Future.Injects[i].Amount = box.Future.Injects[i].Amount.Sub(amount.Amount)
			exist = true
			break
		}
	}
	if !exist {
		return errors.ErrNotEnoughAmount()
	}

	if err := keeper.CancelDepositedCoin(ctx, sender, sdk.NewCoins(amount), box.Id); err != nil {
		return err
	}
	keeper.setBox(ctx, box)
	return nil
}
func (keeper Keeper) processFutureBoxDistribute(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	var address sdk.AccAddress
	var total = sdk.ZeroInt()
	for _, items := range box.Future.Receivers {
		for j, rec := range items {
			if j == 0 {
				addr, err := sdk.AccAddressFromBech32(rec)
				if err != nil {
					return sdk.ErrInvalidAddress(rec)
				}
				address = addr
				continue
			}
			amount, ok := sdk.NewIntFromString(rec)
			if !ok {
				return errors.ErrAmountNotValid(rec)
			}
			boxDenom := utils.GetCoinDenomByFutureBoxSeq(box.Id, j)
			_, err := keeper.GetBankKeeper().AddCoins(ctx, address, sdk.NewCoins(sdk.NewCoin(boxDenom, amount)))
			if err != nil {
				return err
			}
			total = total.Add(amount)
		}
	}
	if !total.Equal(box.TotalAmount.Token.Amount) {
		return errors.ErrAmountNotValid("Receivers")
	}
	//times := len(box.Future.TimeLine)
	//keeper.InsertActiveBoxQueue(ctx, box.Future.TimeLine[times-1], keeper.GetFutureBoxSeqString(box, times))
	keeper.InsertActiveBoxQueue(ctx, box.Future.TimeLine[len(box.Future.TimeLine)-1], box.Id)
	return nil
}

//func (keeper Keeper) GetFutureBoxSeqString(box *types.BoxInfo, seq int) string {
//	return fmt.Sprintf("%s:%d", box.Id, seq)
//}

func (keeper Keeper) ProcessFutureBoxByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	switch box.Status {
	case types.BoxInjecting:
		return keeper.processFutureBoxInjectByEndBlocker(ctx, box)
	case types.BoxActived:
		return keeper.processFutureBoxActiveByEndBlocker(ctx, box)
	default:
		return errors.ErrNotAllowedOperation(box.Status)
	}
}
func (keeper Keeper) processFutureBoxInjectByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	if types.BoxInjecting != box.Status {
		return errors.ErrNotAllowedOperation(box.Status)
	}
	if box.Future.Injects != nil {
		for _, v := range box.Future.Injects {
			if err := keeper.CancelDepositedCoin(ctx, v.Address, sdk.NewCoins(sdk.NewCoin(box.TotalAmount.Token.Denom, v.Amount)), box.Id); err != nil {
				return err
			}
		}
	}
	box.Status = types.BoxClosed
	//keeper.RemoveFromActiveBoxQueue(ctx, box.Future.TimeLine[0], keeper.GetFutureBoxSeqString(box, seq))
	keeper.RemoveFromActiveBoxQueue(ctx, box.Future.TimeLine[0], box.Id)
	//keeper.RemoveBox(ctx, box)
	keeper.setBox(ctx, box)
	return nil
}

func (keeper Keeper) processFutureBoxWithdraw(ctx sdk.Context, idSeq string, sender sdk.AccAddress) (*types.BoxInfo, sdk.Error) {
	box := keeper.GetBox(ctx, idSeq)
	if box == nil {
		return nil, errors.ErrUnknownBox(idSeq)
	}
	if types.Future != box.BoxType {
		return nil, errors.ErrNotSupportOperation()
	}
	if types.BoxCreated == box.Status {
		return nil, errors.ErrNotAllowedOperation(box.Status)
	}
	seq := utils.GetSeqFromFutureBoxSeq(idSeq)
	if box.Future.TimeLine[seq-1] > ctx.BlockHeader().Time.Unix() {
		return nil, errors.ErrNotAllowedOperation(types.BoxUndue)
	}
	amount := keeper.GetBankKeeper().GetCoins(ctx, sender).AmountOf(idSeq)
	_, err := keeper.GetBankKeeper().SubtractCoins(ctx, sender, sdk.NewCoins(sdk.NewCoin(idSeq, amount)))
	if err != nil {
		return nil, err
	}
	if err := keeper.CancelDepositedCoin(ctx, sender, sdk.NewCoins(sdk.NewCoin(box.TotalAmount.Token.Denom, amount)), box.Id); err != nil {
		return nil, err
	}
	box.Future.TotalWithdrawal = amount.Add(box.Future.TotalWithdrawal)
	keeper.setBox(ctx, box)
	return box, nil
}

func (keeper Keeper) processFutureBoxActiveByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	if types.BoxActived != box.Status {
		return errors.ErrNotAllowedOperation(box.Status)
	}
	box.Status = types.BoxFinished
	//keeper.RemoveFromActiveBoxQueue(ctx, box.Future.TimeLine[seq-1], keeper.GetFutureBoxSeqString(box, seq))
	keeper.RemoveFromActiveBoxQueue(ctx, box.Future.TimeLine[len(box.Future.TimeLine)-1], box.Id)
	keeper.setBox(ctx, box)
	return nil
}

// expected bank keeper
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SetSendEnabled(ctx sdk.Context, enabled bool)
}

// expected issue keeper
type IssueKeeper interface {
	GetIssue(ctx sdk.Context, issueID string) *types.CoinIssueInfo
}

// expected fee collection keeper
type FeeCollectionKeeper interface {
	AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins
}

//Process deposit box

func (keeper Keeper) ProcessDepositBoxCreate(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	decimal, err := keeper.GetCoinDecimals(ctx, box.Deposit.Interest.Token)
	if err != nil {
		return err
	}
	if box.Deposit.Interest.Decimals != decimal {
		return errors.ErrDecimalsNotValid(box.Deposit.Interest.Decimals)
	}
	box.Status = types.BoxCreated
	box.Deposit.TotalInject = sdk.ZeroInt()
	box.Deposit.WithdrawalInterest = sdk.ZeroInt()
	box.Deposit.WithdrawalShare = sdk.ZeroInt()
	//box.Deposit.TotalInterestInject = sdk.ZeroInt()
	box.Deposit.Share = sdk.ZeroInt()
	keeper.InsertActiveBoxQueue(ctx, box.Deposit.StartTime, box.Id)
	return nil
}

func (keeper Keeper) InjectDepositBoxInterest(ctx sdk.Context, id string, sender sdk.AccAddress, interest sdk.Coin) (*types.BoxInfo, sdk.Error) {
	box := keeper.GetBox(ctx, id)
	if box == nil {
		return nil, errors.ErrUnknownBox(id)
	}
	if types.BoxCreated != box.Status {
		return nil, errors.ErrNotAllowedOperation(box.Status)
	}
	if box.Deposit.Interest.Token.Denom != interest.Denom {
		return nil, errors.ErrInterestInjectNotValid(interest)
	}
	if len(box.Deposit.InterestInjects) >= types.BoxMaxInjectInterest {
		return nil, errors.ErrInterestInjectNotValid(interest)
	}
	totalInterest := sdk.ZeroInt()
	if box.Deposit.InterestInjects == nil {
		box.Deposit.InterestInjects = []types.AddressInject{{Address: sender, Amount: interest.Amount}}
	} else {
		exist := false
		for i, v := range box.Deposit.InterestInjects {
			totalInterest = totalInterest.Add(v.Amount)
			if v.Address.Equals(sender) {
				box.Deposit.InterestInjects[i].Amount = box.Deposit.InterestInjects[i].Amount.Add(interest.Amount)
				exist = true
			}
		}
		if !exist {
			box.Deposit.InterestInjects = append(box.Deposit.InterestInjects, types.NewAddressInject(sender, interest.Amount))
		}
	}
	totalInterest = totalInterest.Add(interest.Amount)
	if totalInterest.GT(box.Deposit.Interest.Token.Amount) {
		return nil, errors.ErrInterestInjectNotValid(interest)
	}
	if err := keeper.SendDepositedCoin(ctx, sender, sdk.Coins{interest}, id); err != nil {
		return nil, err
	}
	keeper.setBox(ctx, box)
	return box, nil
}
func (keeper Keeper) CancelInterestFromDepositBox(ctx sdk.Context, id string, sender sdk.AccAddress, interest sdk.Coin) (*types.BoxInfo, sdk.Error) {
	box := keeper.GetBox(ctx, id)
	if box == nil {
		return nil, errors.ErrUnknownBox(id)
	}
	if types.BoxCreated != box.Status {
		return nil, errors.ErrNotAllowedOperation(box.Status)
	}
	if box.Deposit.Interest.Token.Denom != interest.Denom {
		return nil, errors.ErrInterestInjectNotValid(interest)
	}
	if box.Deposit.InterestInjects == nil {
		return nil, errors.ErrInterestCancelNotValid(interest)
	} else {
		remove := -1
		for i, v := range box.Deposit.InterestInjects {
			if v.Address.Equals(sender) {
				if box.Deposit.InterestInjects[i].Amount.LT(interest.Amount) {
					return nil, errors.ErrNotEnoughAmount()
				}
				box.Deposit.InterestInjects[i].Amount = box.Deposit.InterestInjects[i].Amount.Sub(interest.Amount)
				if box.Deposit.InterestInjects[i].Amount.IsZero() {
					remove = i
				}
			}
		}
		if remove != -1 {
			box.Deposit.InterestInjects = append(box.Deposit.InterestInjects[:remove], box.Deposit.InterestInjects[remove+1:]...)
		}
	}
	if err := keeper.CancelDepositedCoin(ctx, sender, sdk.Coins{interest}, id); err != nil {
		return nil, err
	}
	keeper.setBox(ctx, box)
	return box, nil
}
func (keeper Keeper) processDepositBoxInject(ctx sdk.Context, box *types.BoxInfo, sender sdk.AccAddress, amount sdk.Coin, operation string) sdk.Error {
	switch operation {
	case types.Inject:
		return keeper.injectDepositBox(ctx, box, sender, amount)
	case types.Cancel:
		return keeper.cancelDepositFromDepositBox(ctx, box, sender, amount)
	}
	return errors.ErrUnknownOperation()
}
func (keeper Keeper) injectDepositBox(ctx sdk.Context, box *types.BoxInfo, sender sdk.AccAddress, amount sdk.Coin) sdk.Error {
	if !amount.Amount.Mod(box.Deposit.Price).IsZero() {
		return errors.ErrAmountNotValid(amount.Denom)
	}
	if box.TotalAmount.Token.Denom != amount.Denom {
		return errors.ErrAmountNotValid(amount.Denom)
	}
	box.Deposit.TotalInject = box.Deposit.TotalInject.Add(amount.Amount)
	if box.Deposit.TotalInject.GT(box.TotalAmount.Token.Amount) {
		return errors.ErrAmountNotValid(amount.Denom)
	}
	if err := keeper.SendDepositedCoin(ctx, sender, sdk.Coins{amount}, box.Id); err != nil {
		return err
	}
	share := amount.Amount.Quo(box.Deposit.Price)
	_, err := keeper.ck.AddCoins(ctx, sender, sdk.NewCoins(sdk.NewCoin(box.Id, share)))
	if err != nil {
		return err
	}
	box.Deposit.Share = box.Deposit.Share.Add(share)
	keeper.setBox(ctx, box)
	return nil
}
func (keeper Keeper) cancelDepositFromDepositBox(ctx sdk.Context, box *types.BoxInfo, sender sdk.AccAddress, amount sdk.Coin) sdk.Error {
	if !amount.Amount.Mod(box.Deposit.Price).IsZero() {
		return errors.ErrAmountNotValid(amount.Amount.String())
	}
	if box.TotalAmount.Token.Denom != amount.Denom {
		return errors.ErrAmountNotValid(amount.Denom)
	}
	share := amount.Amount.Quo(box.Deposit.Price)
	_, err := keeper.GetBankKeeper().SubtractCoins(ctx, sender, sdk.NewCoins(sdk.NewCoin(box.Id, share)))
	if err != nil {
		return err
	}
	if err := keeper.CancelDepositedCoin(ctx, sender, sdk.NewCoins(amount), box.Id); err != nil {
		return err
	}
	box.Deposit.Share = box.Deposit.Share.Sub(share)
	box.Deposit.TotalInject = box.Deposit.TotalInject.Sub(amount.Amount)
	//if types.BoxClosed == box.Status && box.Deposit.Share.IsZero() {
	//	keeper.RemoveBox(ctx, box)
	//} else {
	//	keeper.setBox(ctx, box)
	//}
	keeper.setBox(ctx, box)
	return nil
}
func (keeper Keeper) ProcessDepositBoxByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	switch box.Status {
	case types.BoxCreated:
		return keeper.processBoxCreatedByEndBlocker(ctx, box)
	case types.BoxInjecting:
		return keeper.processDepositBoxInjectByEndBlocker(ctx, box)
	case types.DepositBoxInterest:
		return keeper.processDepositBoxInterestByEndBlocker(ctx, box)
	default:
		return errors.ErrNotAllowedOperation(box.Status)
	}
}
func (keeper Keeper) backBoxInterestInjects(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	for _, v := range box.Deposit.InterestInjects {
		if err := keeper.CancelDepositedCoin(ctx, v.Address,
			sdk.NewCoins(sdk.NewCoin(box.Deposit.Interest.Token.Denom, v.Amount)), box.Id); err != nil {
			return err
		}
	}
	return nil
}
func (keeper Keeper) backBoxUnUsedInterestInjects(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	totalCoupon := box.TotalAmount.Token.Amount.Quo(box.Deposit.Price)
	if totalCoupon.Equal(box.Deposit.Share) {
		return nil
	}
	unused := utils.CalcInterest(box.Deposit.PerCoupon, totalCoupon.Sub(box.Deposit.Share), box.Deposit.Interest)
	interestInjectsLen := len(box.Deposit.InterestInjects)
	if interestInjectsLen == 0 {
		if err := keeper.CancelDepositedCoin(ctx, box.Deposit.InterestInjects[0].Address,
			sdk.NewCoins(sdk.NewCoin(box.Deposit.Interest.Token.Denom, unused)), box.Id); err != nil {
			return err
		}
		box.Deposit.InterestInjects[0].Amount = box.Deposit.InterestInjects[0].Amount.Sub(unused)
	} else {
		total := sdk.ZeroInt()
		for i, v := range box.Deposit.InterestInjects {
			var amount sdk.Int
			if i == interestInjectsLen-1 {
				amount = unused.Sub(total)
			} else {
				amount = sdk.NewDecFromInt(v.Amount).QuoInt(box.Deposit.Interest.Token.Amount).MulInt(unused).TruncateInt()
			}
			if err := keeper.CancelDepositedCoin(ctx, v.Address,
				sdk.NewCoins(sdk.NewCoin(box.Deposit.Interest.Token.Denom, amount)), box.Id); err != nil {
				return err
			}
			box.Deposit.InterestInjects[i].Amount = box.Deposit.InterestInjects[i].Amount.Sub(amount)
			total = total.Add(amount)
		}
	}
	box.Deposit.Interest.Token.Amount = box.Deposit.Interest.Token.Amount.Sub(unused)
	return nil
}
func (keeper Keeper) backBoxAllDeposit(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	for _, v := range box.Deposit.InterestInjects {
		if err := keeper.CancelDepositedCoin(ctx, v.Address,
			sdk.NewCoins(sdk.NewCoin(box.Deposit.Interest.Token.Denom, v.Amount)), box.Id); err != nil {
			return err
		}
	}
	return nil
}
func (keeper Keeper) processBoxCreatedByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	if box.Status != types.BoxCreated {
		return nil
	}
	totalInterest := sdk.ZeroInt()
	if box.Deposit.InterestInjects != nil {
		for _, v := range box.Deposit.InterestInjects {
			totalInterest = totalInterest.Add(v.Amount)
		}
	}
	keeper.RemoveFromActiveBoxQueue(ctx, box.Deposit.StartTime, box.Id)
	if box.Deposit.Interest.Token.Amount.Equal(totalInterest) {
		box.Status = types.BoxInjecting
		keeper.InsertActiveBoxQueue(ctx, box.Deposit.EstablishTime, box.Id)
		keeper.setBox(ctx, box)
	} else {
		box.Status = types.BoxClosed
		if err := keeper.backBoxInterestInjects(ctx, box); err != nil {
			return err
		}
	}
	keeper.setBox(ctx, box)
	return nil
}
func (keeper Keeper) processDepositBoxInjectByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	if box.Status != types.BoxInjecting {
		return nil
	}

	keeper.RemoveFromActiveBoxQueue(ctx, box.Deposit.EstablishTime, box.Id)

	box.Status = types.BoxClosed

	if box.Deposit.TotalInject.IsZero() || box.Deposit.TotalInject.LT(box.Deposit.BottomLine) {
		if err := keeper.backBoxInterestInjects(ctx, box); err != nil {
			return err
		}
	} else {
		if err := keeper.backBoxUnUsedInterestInjects(ctx, box); err != nil {
			return err
		}
		box.Status = types.DepositBoxInterest
		keeper.InsertActiveBoxQueue(ctx, box.Deposit.MaturityTime, box.Id)
	}
	keeper.setBox(ctx, box)
	return nil
}
func (keeper Keeper) processDepositBoxWithdraw(ctx sdk.Context, id string, sender sdk.AccAddress) (sdk.Int, *types.BoxInfo, sdk.Error) {
	box := keeper.GetBox(ctx, id)
	if box == nil {
		return sdk.ZeroInt(), nil, errors.ErrUnknownBox(id)
	}
	if types.Deposit != box.BoxType {
		return sdk.ZeroInt(), nil, errors.ErrNotSupportOperation()
	}
	if types.BoxFinished != box.Status {
		return sdk.ZeroInt(), nil, errors.ErrNotAllowedOperation(box.Status)
	}
	share := keeper.GetBankKeeper().GetCoins(ctx, sender).AmountOf(id)
	box.Deposit.WithdrawalShare = share.Add(box.Deposit.WithdrawalShare)
	_, err := keeper.GetBankKeeper().SubtractCoins(ctx, sender, sdk.NewCoins(sdk.NewCoin(id, share)))
	if err != nil {
		return sdk.ZeroInt(), nil, err
	}

	amount := share.Mul(box.Deposit.Price)
	if err := keeper.CancelDepositedCoin(ctx, sender, sdk.NewCoins(sdk.NewCoin(box.TotalAmount.Token.Denom, amount)), box.Id); err != nil {
		return sdk.ZeroInt(), nil, err
	}

	interest := sdk.ZeroInt()
	if box.Deposit.WithdrawalShare == box.Deposit.Share {
		totalInterest := sdk.ZeroInt()
		for _, v := range box.Deposit.InterestInjects {
			totalInterest = totalInterest.Add(v.Amount)
		}
		interest = totalInterest.Sub(box.Deposit.WithdrawalInterest)
	} else {
		interest = utils.CalcInterest(box.Deposit.PerCoupon, share, box.Deposit.Interest)
	}

	if err = keeper.CancelDepositedCoin(ctx, sender,
		sdk.NewCoins(sdk.NewCoin(box.Deposit.Interest.Token.Denom, interest)), box.Id); err != nil {
		return sdk.ZeroInt(), nil, err
	}

	box.Deposit.WithdrawalInterest = interest.Add(box.Deposit.WithdrawalInterest)
	keeper.setBox(ctx, box)
	return interest, box, nil

}
func (keeper Keeper) processDepositBoxInterestByEndBlocker(ctx sdk.Context, box *types.BoxInfo) sdk.Error {
	if box.Status != types.DepositBoxInterest {
		return nil
	}
	keeper.RemoveFromActiveBoxQueue(ctx, box.Deposit.MaturityTime, box.Id)
	box.Status = types.BoxFinished
	keeper.setBox(ctx, box)
	return nil
}

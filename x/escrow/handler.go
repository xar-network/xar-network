package escrow

import (
	"fmt"

	"github.com/xar-network/xar-network/x/escrow/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/escrow/internal/keeper"
)

// Handle all "box" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgLockBox:
			return HandleMsgLockBox(ctx, keeper, msg)
		case types.MsgDepositBox:
			return HandleMsgDepositBox(ctx, keeper, msg)
		case types.MsgFutureBox:
			return HandleMsgFutureBox(ctx, keeper, msg)
		case types.MsgBoxInterestInject:
			return HandleMsgBoxInterestInject(ctx, keeper, msg)
		case types.MsgBoxInterestCancel:
			return HandleMsgBoxInterestCancel(ctx, keeper, msg)
		case types.MsgBoxInject:
			return HandleMsgBoxInject(ctx, keeper, msg)
		case types.MsgBoxInjectCancel:
			return HandleMsgBoxInjectCancel(ctx, keeper, msg)
		case types.MsgBoxWithdraw:
			return HandleMsgBoxWithdraw(ctx, keeper, msg)
		case types.MsgBoxDescription:
			return HandleMsgBoxDescription(ctx, keeper, msg)
		case types.MsgBoxDisableFeature:
			return HandleMsgBoxDisableFeature(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized gov msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

//Handle MsgBoxDescription
func HandleMsgBoxDescription(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBoxDescription) sdk.Result {
	fee := keeper.GetParams(ctx).DescribeFee
	if err := keeper.Fee(ctx, msg.Sender, fee); err != nil {
		return err.Result()
	}

	boxInfo, err := keeper.SetBoxDescription(ctx, msg.Id, msg.Sender, msg.Description)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(msg.Id),
		Tags: types.GetBoxTags(msg.Id, boxInfo.BoxType, msg.Sender),
	}
}

//Handle MsgBoxInject
func HandleMsgBoxInject(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBoxInject) sdk.Result {
	boxInfo, err := keeper.ProcessInjectBox(ctx, msg.Id, msg.Sender, msg.Amount, types.Inject)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(msg.Id),
		Tags: types.GetBoxTags(msg.Id, boxInfo.BoxType, msg.Sender),
	}
}

//Handle MsgBoxInject
func HandleMsgBoxInjectCancel(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBoxInjectCancel) sdk.Result {
	boxInfo, err := keeper.ProcessInjectBox(ctx, msg.Id, msg.Sender, msg.Amount, types.Cancel)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(msg.Id),
		Tags: types.GetBoxTags(msg.Id, boxInfo.BoxType, msg.Sender),
	}
}

//Handle MsgBoxDisableFeature
func HandleMsgBoxDisableFeature(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBoxDisableFeature) sdk.Result {
	fee := keeper.GetParams(ctx).DisableFeatureFee
	if err := keeper.Fee(ctx, msg.Sender, fee); err != nil {
		return err.Result()
	}

	boxInfo, err := keeper.DisableFeature(ctx, msg.Sender, msg.Id, msg.Feature)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(msg.Id),
		Tags: types.GetBoxTags(msg.Id, boxInfo.BoxType, msg.Sender).
			AppendTag(tags.Feature, msg.GetFeature()).AppendTag(tags.Fee, fee.String()),
	}
}

//Handle MsgBoxInterestInject
func HandleMsgBoxInterestInject(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBoxInterestInject) sdk.Result {
	boxInfo, err := keeper.InjectDepositBoxInterest(ctx, msg.Id, msg.Sender, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(msg.Id),
		Tags: types.GetBoxTags(msg.Id, boxInfo.BoxType, msg.Sender),
	}
}

//Handle MsgBoxInterestCancel
func HandleMsgBoxInterestCancel(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBoxInterestCancel) sdk.Result {
	boxInfo, err := keeper.CancelInterestFromDepositBox(ctx, msg.Id, msg.Sender, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(msg.Id),
		Tags: types.GetBoxTags(msg.Id, boxInfo.BoxType, msg.Sender),
	}
}

//Handle MsgBoxWithdraw
func HandleMsgBoxWithdraw(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgBoxWithdraw) sdk.Result {
	interest, boxInfo, err := keeper.ProcessBoxWithdraw(ctx, msg.Id, msg.Sender)
	if err != nil {
		return err.Result()
	}

	resTags := types.GetBoxTags(msg.Id, boxInfo.BoxType, msg.Sender)
	if types.Deposit == boxInfo.BoxType {
		resTags = resTags.AppendTag(tags.Interest, interest.String())
	}

	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(msg.Id),
		Tags: resTags,
	}
}

//Handle MsgLockBox
func HandleMsgLockBox(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgLockBox) sdk.Result {
	fee := keeper.GetParams(ctx).LockCreateFee
	if err := keeper.Fee(ctx, msg.Sender, fee); err != nil {
		return err.Result()
	}

	box := &types.BoxInfo{
		Owner:            msg.Sender,
		Name:             msg.Name,
		BoxType:          types.Lock,
		TotalAmount:      msg.TotalAmount,
		Description:      msg.Description,
		TransferDisabled: true,
		Lock:             msg.Lock,
	}
	return createBox(ctx, keeper, box, fee)
}

//Handle MsgDepositBox
func HandleMsgDepositBox(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgDepositBox) sdk.Result {
	fee := keeper.GetParams(ctx).DepositBoxCreateFee
	if err := keeper.Fee(ctx, msg.Sender, fee); err != nil {
		return err.Result()
	}

	box := &types.BoxInfo{
		Owner:            msg.Sender,
		Name:             msg.Name,
		BoxType:          types.Deposit,
		TotalAmount:      msg.TotalAmount,
		Description:      msg.Description,
		TransferDisabled: msg.TransferDisabled,
		Deposit:          msg.Deposit,
	}
	return createBox(ctx, keeper, box, fee)
}

//Handle MsgFutureBox
func HandleMsgFutureBox(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgFutureBox) sdk.Result {
	fee := keeper.GetParams(ctx).FutureBoxCreateFee
	if err := keeper.Fee(ctx, msg.Sender, fee); err != nil {
		return err.Result()
	}
	box := &types.BoxInfo{
		Owner:            msg.Sender,
		Name:             msg.Name,
		BoxType:          types.Future,
		TotalAmount:      msg.TotalAmount,
		Description:      msg.Description,
		TransferDisabled: msg.TransferDisabled,
		Future:           msg.Future,
	}
	return createBox(ctx, keeper, box, fee)
}
func createBox(ctx sdk.Context, keeper keeper.Keeper, box *types.BoxInfo, fee sdk.Coin) sdk.Result {
	err := keeper.CreateBox(ctx, box)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Data: keeper.Getcdc().MustMarshalBinaryLengthPrefixed(box.Id),
		Tags: sdk.NewTags(
			tags.Category, box.BoxType,
			tags.BoxID, box.Id,
			tags.Sender, box.Owner.String(),
			tags.Fee, fee.String(),
		),
	}
}

// Called every block, process inflation, update validator set
func EndBlocker(ctx sdk.Context, keeper Keeper) sdk.Tags {
	logger := ctx.Logger().With("module", "x/"+types.ModuleName)
	resTags := sdk.NewTags()
	// fetch active proposals whose voting periods have ended (are passed the block time)
	activeIterator := keeper.ActiveBoxQueueIterator(ctx, ctx.BlockHeader().Time.Unix())
	defer activeIterator.Close()
	count := 0
	for ; activeIterator.Valid(); activeIterator.Next() {
		if count == types.BoxMaxInstalment {
			break
		}
		var id string
		//var seq int
		keeper.Getcdc().MustUnmarshalBinaryLengthPrefixed(activeIterator.Value(), &id)
		//if strings.Contains(id, types.KeyDelimiterString) {
		//	ids := strings.Split(id, types.KeyDelimiterString)
		//	id = ids[0]
		//	seq, _ = strconv.Atoi(ids[1])
		//}
		boxInfo := keeper.GetBox(ctx, id)
		if boxInfo == nil {
			panic(fmt.Sprintf("box %s does not exist", id))
			//continue
		}
		switch boxInfo.BoxType {
		case types.Lock:
			if err := keeper.ProcessLockBoxByEndBlocker(ctx, boxInfo); err != nil {
				panic(err)
			}
			logger.Debug(fmt.Sprintf("lockbox %s (%s) unlocked", id, boxInfo.Name))
			resTags = resTags.AppendTag(tags.BoxID, id).AppendTag(tags.Category, boxInfo.GetBoxType()).AppendTag(tags.Status, types.LockBoxUnlocked)
		case types.Deposit:
			if err := keeper.ProcessDepositBoxByEndBlocker(ctx, boxInfo); err != nil {
				panic(err)
			}
			logger.Debug(fmt.Sprintf("depositbox %s (%s) status:%s", id, boxInfo.Name, boxInfo.Status))
			resTags = resTags.AppendTag(tags.BoxID, id).AppendTag(tags.Category, boxInfo.GetBoxType()).AppendTag(tags.Status, boxInfo.Status)
		case types.Future:
			if err := keeper.ProcessFutureBoxByEndBlocker(ctx, boxInfo); err != nil {
				panic(err)
			}
			logger.Debug(fmt.Sprintf("futurebox %s (%s) status:%s", id, boxInfo.Name, boxInfo.Status))
			resTags = resTags.AppendTag(tags.BoxID, id).AppendTag(tags.Category, boxInfo.GetBoxType()).AppendTag(tags.Status, boxInfo.Status)
		}
		count = count + 1
	}
	return resTags
}

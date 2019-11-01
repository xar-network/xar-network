package record

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/x/record/internal/keeper"
	"github.com/xar-network/xar-network/x/record/internal/types"
)

// NewHandler all "record" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgRecord:
			return HandleMsgRecord(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized record msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

//HandleMsgRecord MsgRecord
func HandleMsgRecord(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgRecord) sdk.Result {
	recordInfo := types.RecordInfo{
		Sender:      msg.Sender,
		Hash:        msg.Hash,
		Name:        msg.Name,
		Author:      msg.Author,
		RecordType:  msg.RecordType,
		RecordNo:    msg.RecordNo,
		Description: msg.Description,
	}

	err := keeper.CreateRecord(ctx, &recordInfo)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Events: types.GetRecordTags(&recordInfo),
	}
}

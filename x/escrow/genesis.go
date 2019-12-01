package escrow

import (
	"bytes"
	"fmt"

	"github.com/xar-network/xar-network/x/escrow/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all box state that must be provided at genesis
type GenesisState struct {
	StartingLockId    uint64    `json:"starting_lock_id"`
	StartingDepositId uint64    `json:"starting_deposit_id"`
	StartingFutureId  uint64    `json:"starting_future_id"`
	LockBoxs          []BoxInfo `json:"lock_boxs"`
	DepositBoxs       []BoxInfo `json:"deposit_boxs"`
	FutureBoxs        []BoxInfo `json:"future_boxs"`
	Params            Params    `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(startingLockId uint64, startingDepositId uint64, startingFutureId uint64) GenesisState {
	return GenesisState{
		StartingLockId:    startingLockId,
		StartingDepositId: startingDepositId,
		StartingFutureId:  startingFutureId,
		Params:            types.DefaultParams(sdk.DefaultBondDenom),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.BoxMinId, types.BoxMinId, types.BoxMinId)
}

// Returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// Checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := MsgCdc.MustMarshalBinaryBare(data)
	b2 := MsgCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	if err := keeper.SetInitialBoxStartingId(ctx, types.Lock, data.StartingLockId); err != nil {
		panic(err)
	}
	if err := keeper.SetInitialBoxStartingId(ctx, types.Deposit, data.StartingDepositId); err != nil {
		panic(err)
	}
	if err := keeper.SetInitialBoxStartingId(ctx, types.Future, data.StartingFutureId); err != nil {
		panic(err)
	}

	keeper.SetParams(ctx, data.Params)

	for _, box := range data.LockBoxs {
		keeper.AddBox(ctx, &box)
		if box.Status == types.LockBoxLocked {
			keeper.InsertActiveBoxQueue(ctx, box.Lock.EndTime, box.Id)
		}
	}

	for _, box := range data.DepositBoxs {
		keeper.AddBox(ctx, &box)
		switch box.Status {
		case types.BoxCreated:
			keeper.InsertActiveBoxQueue(ctx, box.Deposit.StartTime, box.Id)
		case types.BoxInjecting:
			keeper.InsertActiveBoxQueue(ctx, box.Deposit.EstablishTime, box.Id)
		case types.DepositBoxInterest:
			keeper.InsertActiveBoxQueue(ctx, box.Deposit.MaturityTime, box.Id)
		}
	}

	for _, box := range data.FutureBoxs {
		keeper.AddBox(ctx, &box)
		switch box.Status {
		case types.BoxInjecting:
			keeper.InsertActiveBoxQueue(ctx, box.Future.TimeLine[0], box.Id)
		case types.BoxActived:
			keeper.InsertActiveBoxQueue(ctx, box.Future.TimeLine[len(box.Future.TimeLine)-1], box.Id)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	genesisState := GenesisState{}
	var err sdk.Error

	genesisState.StartingLockId, err = keeper.PeekCurrentBoxID(ctx, types.Lock)
	if err != nil {
		panic(err)
	}
	genesisState.StartingDepositId, err = keeper.PeekCurrentBoxID(ctx, types.Deposit)
	if err != nil {
		panic(err)
	}
	genesisState.StartingFutureId, err = keeper.PeekCurrentBoxID(ctx, types.Future)
	if err != nil {
		panic(err)
	}
	genesisState.LockBoxs = keeper.ListAll(ctx, types.Lock)
	genesisState.DepositBoxs = keeper.ListAll(ctx, types.Deposit)
	genesisState.FutureBoxs = keeper.ListAll(ctx, types.Future)

	genesisState.Params = keeper.GetParams(ctx)

	return genesisState
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {

	if data.Params.LockCreateFee.IsNegative() {
		return fmt.Errorf("invalid lock create fee: %s", data.Params.LockCreateFee.String())
	}
	if data.Params.DepositBoxCreateFee.IsNegative() {
		return fmt.Errorf("invalid deposit box create fee: %s", data.Params.DepositBoxCreateFee.String())
	}
	if data.Params.FutureBoxCreateFee.IsNegative() {
		return fmt.Errorf("invalid future box create fee: %s", data.Params.FutureBoxCreateFee.String())
	}
	if data.Params.DisableFeatureFee.IsNegative() {
		return fmt.Errorf("invalid disable feature fee: %s", data.Params.DisableFeatureFee.String())
	}
	if data.Params.DescribeFee.IsNegative() {
		return fmt.Errorf("invalid describe fee: %s", data.Params.DescribeFee.String())
	}
	return nil
}

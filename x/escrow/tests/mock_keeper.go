package tests

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/issue/types"
)

type IssueKeeper struct {
}

//New issue keeper Instance
func NewIssueKeeper() IssueKeeper {
	return IssueKeeper{}
}

//Returns issue by issueID
func (keeper IssueKeeper) GetIssue(ctx sdk.Context, issueID string) *types.CoinIssueInfo {

	coinIssueInfo := types.CoinIssueInfo{
		Decimals: TestTokenDecimals,
	}
	return &coinIssueInfo
}


// Wrapper struct
type MockHooks struct {
	keeper box.Keeper
}

// Create new box hooks
func NewMockHooks(bk box.Keeper) MockHooks { return MockHooks{bk} }

func (hooks MockHooks) CanSend(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (bool, sdk.Error) {
	return hooks.keeper.Hooks().CanSend(ctx, fromAddr, toAddr, amt)
}

func (hooks MockHooks) CheckMustMemoAddress(ctx sdk.Context, toAddr sdk.AccAddress, memo string) (bool, sdk.Error) {
	return false, nil
}

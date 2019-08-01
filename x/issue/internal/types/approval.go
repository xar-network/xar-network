package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Approve          = "approve"
	IncreaseApproval = "increaseApproval"
	DecreaseApproval = "decreaseApproval"
)

type Approval struct {
	Amount sdk.Int `json:"amount"`
}

func NewApproval(amount sdk.Int) Approval {
	return Approval{amount}
}

func (ci Approval) String() string {
	return fmt.Sprintf(`Amount:%s`, ci.Amount)
}

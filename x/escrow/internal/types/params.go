package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type LockBoxParams struct {
	LockCreateFee sdk.Coin `json:"lock_create_fee"`
	DescribeFee   sdk.Coin `json:"describe_fee"`
}
type DepositBoxParams struct {
	DepositBoxCreateFee sdk.Coin `json:"deposit_box_create_fee"`
	DisableFeatureFee   sdk.Coin `json:"disable_feature_fee"`
	DescribeFee         sdk.Coin `json:"describe_fee"`
}
type FutureBoxParams struct {
	FutureBoxCreateFee sdk.Coin `json:"future_box_create_fee"`
	DisableFeatureFee  sdk.Coin `json:"disable_feature_fee"`
	DescribeFee        sdk.Coin `json:"describe_fee"`
}

//nolint
func (dp LockBoxParams) String() string {
	return fmt.Sprintf(`Params:
  LockCreateFee:			%s
  DescribeFee:			%s`,
		dp.LockCreateFee.String(),
		dp.DescribeFee.String())
}

//nolint
func (dp DepositBoxParams) String() string {
	return fmt.Sprintf(`Params:
  DepositBoxCreateFee:			%s
  DisableFeatureFee:			%s
  DescribeFee:			%s`,
		dp.DepositBoxCreateFee.String(),
		dp.DisableFeatureFee.String(),
		dp.DescribeFee.String())
}

//nolint
func (dp FutureBoxParams) String() string {
	return fmt.Sprintf(`Params:
  FutureBoxCreateFee:			%s
  DisableFeatureFee:			%s
  DescribeFee:			%s`,
		dp.FutureBoxCreateFee.String(),
		dp.DisableFeatureFee.String(),
		dp.DescribeFee.String())
}

var (
	// key for constant fee parameter
	ParamStoreKeyLockCreateFee        = []byte("LockCreateFee")
	ParamStoreKeyDepositBoxCreateFee  = []byte("DepositBoxCreateFee")
	ParamStoreKeyFutureBoxCreateFee   = []byte("FutureBoxCreateFee")
	ParamStoreKeyBoxDisableFeatureFee = []byte("DisableFeatureFee")
	ParamStoreKeyBoxDescribeFee       = []byte("DescribeFee")
)
var _ params.ParamSet = &Params{}

// Param Config issue for issue
type Params struct {
	LockCreateFee       sdk.Coin `json:"lock_create_fee"`
	DepositBoxCreateFee sdk.Coin `json:"deposit_box_create_fee"`
	FutureBoxCreateFee  sdk.Coin `json:"future_box_create_fee"`
	DisableFeatureFee   sdk.Coin `json:"disable_feature_fee"`
	DescribeFee         sdk.Coin `json:"describe_fee"`
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{ParamStoreKeyLockCreateFee, &p.LockCreateFee},
		{ParamStoreKeyDepositBoxCreateFee, &p.DepositBoxCreateFee},
		{ParamStoreKeyFutureBoxCreateFee, &p.FutureBoxCreateFee},
		{ParamStoreKeyBoxDisableFeatureFee, &p.DisableFeatureFee},
		{ParamStoreKeyBoxDescribeFee, &p.DescribeFee},
	}
}

// Checks equality of Params
func (dp Params) Equal(dp2 Params) bool {
	b1 := msgs.MsgCdc.MustMarshalBinaryBare(dp)
	b2 := msgs.MsgCdc.MustMarshalBinaryBare(dp2)
	return bytes.Equal(b1, b2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams(denom string) Params {
	return Params{
		LockCreateFee:       sdk.NewCoin(denom, sdk.NewIntWithDecimal(1000, 18)),
		DepositBoxCreateFee: sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 18)),
		FutureBoxCreateFee:  sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 18)),
		DisableFeatureFee:   sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 18)),
		DescribeFee:         sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 18)),
	}
}

func (dp Params) String() string {
	return fmt.Sprintf(`Params:
  LockCreateFee:			%s
  DepositBoxCreateFee:			%s
  FutureBoxCreateFee:			%s
  DisableFeatureFee:			%s
  DescribeFee:			%s`,
		dp.LockCreateFee.String(),
		dp.DepositBoxCreateFee.String(),
		dp.FutureBoxCreateFee.String(),
		dp.DisableFeatureFee.String(),
		dp.DescribeFee.String())
}

// Param query box
type BoxQueryParams struct {
	StartId string         `json:"start_id"`
	Owner   sdk.AccAddress `json:"owner"`
	BoxType string         `json:"type"`
	Limit   int            `json:"limit"`
}

// Param query deposit
type BoxQueryDepositListParams struct {
	Id    string         `json:"id"`
	Owner sdk.AccAddress `json:"owner"`
	//Limit int            `json:"limit"`
}

type BoxLockParams struct {
	Name        string   `json:"name"`
	TotalAmount BoxToken `json:"total_amount"`
	Description string   `json:"description"`
	Lock        LockBox  `json:"lock"`
}

type BoxFutureParams struct {
	Name             string    `json:"name"`
	TotalAmount      BoxToken  `json:"total_amount"`
	Description      string    `json:"description"`
	TransferDisabled bool      `json:"transfer_disabled"`
	Future           FutureBox `json:"future"`
}

type BoxDepositParams struct {
	Name             string     `json:"name"`
	TotalAmount      BoxToken   `json:"total_amount"`
	Description      string     `json:"description"`
	TransferDisabled bool       `json:"transfer_disabled"`
	Deposit          DepositBox `json:"deposit"`
}

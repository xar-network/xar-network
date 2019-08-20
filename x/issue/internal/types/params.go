package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	// key for constant fee parameter
	ParamStoreKeyIssueFee         = []byte("IssueFee")
	ParamStoreKeyMintFee          = []byte("MintFee")
	ParamStoreKeyFreezeFee        = []byte("FreezeFee")
	ParamStoreKeyUnFreezeFee      = []byte("UnfreezeFee")
	ParamStoreKeyBurnFee          = []byte("BurnFee")
	ParamStoreKeyBurnFromFee      = []byte("BurnFromFee")
	ParamStoreKeyTransferOwnerFee = []byte("TransferOwnerFee")
	ParamStoreKeyDescribeFee      = []byte("DescribeFee")
)

var _ params.ParamSet = &Params{}

// Param Config issue for issue
type Params struct {
	IssueFee         sdk.Coin `json:"issue_fee"`
	MintFee          sdk.Coin `json:"mint_fee"`
	FreezeFee        sdk.Coin `json:"freeze_fee"`
	UnFreezeFee      sdk.Coin `json:"unfreeze_fee"`
	BurnFee          sdk.Coin `json:"burn_fee"`
	BurnFromFee      sdk.Coin `json:"burn_from_fee"`
	TransferOwnerFee sdk.Coin `json:"transfer_owner_fee"`
	DescribeFee      sdk.Coin `json:"describe_fee"`
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
		{ParamStoreKeyIssueFee, &p.IssueFee},
		{ParamStoreKeyMintFee, &p.MintFee},
		{ParamStoreKeyFreezeFee, &p.FreezeFee},
		{ParamStoreKeyUnFreezeFee, &p.UnFreezeFee},
		{ParamStoreKeyBurnFee, &p.BurnFee},
		{ParamStoreKeyBurnFromFee, &p.BurnFromFee},
		{ParamStoreKeyTransferOwnerFee, &p.TransferOwnerFee},
		{ParamStoreKeyDescribeFee, &p.DescribeFee},
	}
}

// Checks equality of Params
func (dp Params) Equal(dp2 Params) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(dp)
	b2 := ModuleCdc.MustMarshalBinaryBare(dp2)
	return bytes.Equal(b1, b2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams(denom string) Params {
	return Params{
		IssueFee:         sdk.NewCoin(denom, sdk.NewIntWithDecimal(20000, 18)),
		MintFee:          sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 18)),
		FreezeFee:        sdk.NewCoin(denom, sdk.NewIntWithDecimal(20000, 18)),
		UnFreezeFee:      sdk.NewCoin(denom, sdk.NewIntWithDecimal(20000, 18)),
		BurnFee:          sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 18)),
		BurnFromFee:      sdk.NewCoin(denom, sdk.NewIntWithDecimal(10000, 18)),
		TransferOwnerFee: sdk.NewCoin(denom, sdk.NewIntWithDecimal(20000, 18)),
		DescribeFee:      sdk.NewCoin(denom, sdk.NewIntWithDecimal(4000, 18)),
	}
}

func (dp Params) String() string {
	return fmt.Sprintf(`Params:
  IssueFee:			%s
  MintFee:			%s
  FreezeFee:			%s
  UnFreezeFee:			%s
  BurnFee:			%s
  BurnFromFee:			%s
  TransferOwnerFee:		%s
  DescribeFee:			%s`,
		dp.IssueFee.String(),
		dp.MintFee.String(),
		dp.FreezeFee.String(),
		dp.UnFreezeFee.String(),
		dp.BurnFee.String(),
		dp.BurnFromFee.String(),
		dp.TransferOwnerFee.String(),
		dp.DescribeFee.String(),
	)
}

// Param issue for issue
type IssueParams struct {
	Name               string  `json:"name"`
	Symbol             string  `json:"symbol"`
	TotalSupply        sdk.Int `json:"total_supply"`
	Decimals           uint    `json:"decimals"`
	Description        string  `json:"description"`
	BurnOwnerDisabled  bool    `json:"burn_owner_disabled"`
	BurnHolderDisabled bool    `json:"burn_holder_disabled"`
	BurnFromDisabled   bool    `json:"burn_from_disabled"`
	MintingFinished    bool    `json:"minting_finished"`
	FreezeDisabled     bool    `json:"freeze_disabled"`
}

// Param query issue for issue
type IssueQueryParams struct {
	StartIssueId string         `json:"start_issue_id"`
	Owner        sdk.AccAddress `json:"owner"`
	Limit        int            `json:"limit"`
}

func GetQueryIssuePath(issueID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", Custom, QuerierRoute, QueryIssue, issueID)
}
func GetQueryParamsPath() string {
	return fmt.Sprintf("%s/%s/%s", Custom, QuerierRoute, QueryParams)
}
func GetQueryIssueAllowancePath(issueID string, owner sdk.AccAddress, spender sdk.AccAddress) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/%s", Custom, QuerierRoute, QueryAllowance, issueID, owner.String(), spender.String())
}
func GetQueryIssueFreezePath(issueID string, accAddress sdk.AccAddress) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", Custom, QuerierRoute, QueryFreeze, issueID, accAddress.String())
}
func GetQueryIssueFreezesPath(issueID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", Custom, QuerierRoute, QueryFreezes, issueID)
}
func GetQueryIssueSearchPath(symbol string) string {
	return fmt.Sprintf("%s/%s/%s/%s", Custom, QuerierRoute, QuerySearch, symbol)
}
func GetQueryIssuesPath() string {
	return fmt.Sprintf("%s/%s/%s", Custom, QuerierRoute, QueryIssues)
}

package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetInterest struct {
	Denom        string
	InterestRate sdk.Dec
	Issuer       sdk.AccAddress
}

func (msg MsgSetInterest) Route() string { return ModuleName }

func (msg MsgSetInterest) Type() string { return "setInterest" }

func (msg MsgSetInterest) ValidateBasic() sdk.Error {
	if msg.InterestRate.IsNegative() {
		return ErrNegativeInterest()
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	return nil
}

func (msg MsgSetInterest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetInterest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

//Issue interface
type Issue interface {
	GetIssueId() string
	SetIssueId(string)

	GetIssuer() sdk.AccAddress
	SetIssuer(sdk.AccAddress)

	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress)

	GetIssueTime() int64
	SetIssueTime(int64)

	GetName() string
	SetName(string)

	GetTotalSupply() sdk.Int
	SetTotalSupply(sdk.Int)

	GetDescription() string
	SetDescription(string)

	IsBurnOwnerDisabled() bool
	SetBurnOwnerDisabled(bool)

	IsBurnHolderDisabled() bool
	SetBurnHolderDisabled(bool)

	IsBurnFromDisabled() bool
	SetBurnFromDisabled(bool)

	IsFreezeDisabled() bool
	SetFreezeDisabled(bool)

	IsMintingFinished() bool
	SetMintingFinished(bool)

	GetSymbol() string
	SetSymbol(string)

	String() string
}

// CoinIssues is an array of Issue
type CoinIssues []CoinIssueInfo

//Coin Issue Info
type CoinIssueInfo struct {
	IssueId            string         `json:"issue_id"`
	Issuer             sdk.AccAddress `json:"issuer"`
	Owner              sdk.AccAddress `json:"owner"`
	IssueTime          int64          `json:"issue_time"`
	Name               string         `json:"name"`
	Symbol             string         `json:"symbol"`
	TotalSupply        sdk.Int        `json:"total_supply"`
	Description        string         `json:"description"`
	BurnOwnerDisabled  bool           `json:"burn_owner_disabled"`
	BurnHolderDisabled bool           `json:"burn_holder_disabled"`
	BurnFromDisabled   bool           `json:"burn_from_disabled"`
	FreezeDisabled     bool           `json:"freeze_disabled"`
	MintingFinished    bool           `json:"minting_finished"`
}

// Implements Issue Interface
var _ Issue = (*CoinIssueInfo)(nil)

//nolint
func (ci CoinIssueInfo) GetIssueId() string {
	return ci.IssueId
}
func (ci *CoinIssueInfo) SetIssueId(issueId string) {
	ci.IssueId = issueId
}
func (ci CoinIssueInfo) GetIssuer() sdk.AccAddress {
	return ci.Issuer
}
func (ci *CoinIssueInfo) SetIssuer(issuer sdk.AccAddress) {
	ci.Issuer = issuer
}
func (ci CoinIssueInfo) GetOwner() sdk.AccAddress {
	return ci.Owner
}
func (ci *CoinIssueInfo) SetOwner(owner sdk.AccAddress) {
	ci.Owner = owner
}
func (ci CoinIssueInfo) GetIssueTime() int64 {
	return ci.IssueTime
}
func (ci *CoinIssueInfo) SetIssueTime(issueTime int64) {
	ci.IssueTime = issueTime
}
func (ci CoinIssueInfo) GetName() string {
	return ci.Name
}
func (ci *CoinIssueInfo) SetName(name string) {
	ci.Name = name
}
func (ci CoinIssueInfo) GetTotalSupply() sdk.Int {
	return ci.TotalSupply
}
func (ci *CoinIssueInfo) SetTotalSupply(totalSupply sdk.Int) {
	ci.TotalSupply = totalSupply
}
func (ci CoinIssueInfo) GetDescription() string {
	return ci.Description
}
func (ci *CoinIssueInfo) SetDescription(description string) {
	ci.Description = description
}

func (ci CoinIssueInfo) GetSymbol() string {
	return ci.Symbol
}
func (ci *CoinIssueInfo) SetSymbol(symbol string) {
	ci.Symbol = symbol
}
func (ci CoinIssueInfo) IsBurnOwnerDisabled() bool {
	return ci.BurnOwnerDisabled
}

func (ci CoinIssueInfo) SetBurnOwnerDisabled(burnOwnerDisabled bool) {
	ci.BurnOwnerDisabled = burnOwnerDisabled
}

func (ci CoinIssueInfo) IsBurnHolderDisabled() bool {
	return ci.BurnHolderDisabled
}

func (ci CoinIssueInfo) SetBurnHolderDisabled(burnFromDisabled bool) {
	ci.BurnHolderDisabled = burnFromDisabled
}

func (ci CoinIssueInfo) IsBurnFromDisabled() bool {
	return ci.BurnFromDisabled
}

func (ci CoinIssueInfo) SetBurnFromDisabled(burnFromDisabled bool) {
	ci.BurnFromDisabled = burnFromDisabled
}
func (ci CoinIssueInfo) IsFreezeDisabled() bool {
	return ci.FreezeDisabled
}

func (ci CoinIssueInfo) SetFreezeDisabled(freezeDisabled bool) {
	ci.FreezeDisabled = freezeDisabled
}
func (ci CoinIssueInfo) IsMintingFinished() bool {
	return ci.MintingFinished
}

func (ci CoinIssueInfo) SetMintingFinished(mintingFinished bool) {
	ci.MintingFinished = mintingFinished
}

//nolint
func (ci CoinIssueInfo) String() string {
	return fmt.Sprintf(`Issue:
  IssueId:          			%s
  Issuer:           			%s
  Owner:           				%s
  Name:             			%s
  Symbol:    	    			%s
  TotalSupply:      			%s
  IssueTime:					%d
  Description:	    			%s
  BurnOwnerDisabled:  			%t
  BurnHolderDisabled:  			%t
  BurnFromDisabled:  			%t
  FreezeDisabled:  				%t
  MintingFinished:  			%t `,
		ci.IssueId, ci.Issuer.String(), ci.Owner.String(), ci.Name, ci.Symbol, ci.TotalSupply.String(),
		ci.IssueTime, ci.Description, ci.BurnOwnerDisabled, ci.BurnHolderDisabled,
		ci.BurnFromDisabled, ci.FreezeDisabled, ci.MintingFinished)
}

//nolint
func (coinIssues CoinIssues) String() string {
	out := fmt.Sprintf("%-17s|%-44s|%-10s|%-6s|%-18s|%s\n",
		"IssueID", "Owner", "Name", "Symbol", "TotalSupply", "IssueTime")
	for _, issue := range coinIssues {
		out += fmt.Sprintf("%-17s|%-44s|%-10s|%-6s|%-18s|%d\n",
			issue.IssueId, issue.GetOwner().String(), issue.Name, issue.Symbol, issue.TotalSupply.String(), issue.IssueTime)
	}
	return strings.TrimSpace(out)
}

const (
	FreezeIn       = "in"
	FreezeOut      = "out"
	FreezeInAndOut = "in-out"
)

var FreezeTypes = map[string]int{FreezeIn: 1, FreezeOut: 1, FreezeInAndOut: 1}

type IssueFreeze struct {
	Frozen bool `json:"frozen"`
}

func (ci IssueFreeze) String() string {
	return fmt.Sprintf(`Frozen:\n
	Frozen:			%t`,
		ci.Frozen)
}

type IssueAddressFreeze struct {
	Address string `json:"address"`
}

type IssueAddressFreezeList []IssueAddressFreeze

func (ci IssueAddressFreeze) String() string {
	return fmt.Sprintf(`FreezeList:\n
	Address:			%s`,
		ci.Address)
}

//nolint
func (ci IssueAddressFreezeList) String() string {
	out := fmt.Sprintf("%-44s\n",
		"Address")
	for _, v := range ci {
		out += fmt.Sprintf("%-44s\n",
			v.Address)
	}
	return strings.TrimSpace(out)
}

const (
	BurnOwner  = "burn-owner"
	BurnHolder = "burn-holder"
	BurnFrom   = "burn-from"
	Freeze     = "freeze"
	Minting    = "minting"
)

var Features = map[string]int{BurnOwner: 1, BurnHolder: 1, BurnFrom: 1, Freeze: 1, Minting: 1}

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

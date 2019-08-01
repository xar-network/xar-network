package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

	GetDecimals() uint
	SetDecimals(uint)

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
	Decimals           uint           `json:"decimals"`
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
func (ci CoinIssueInfo) GetDecimals() uint {
	return ci.Decimals
}
func (ci *CoinIssueInfo) SetDecimals(decimals uint) {
	ci.Decimals = decimals
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
  Decimals:         			%d
  IssueTime:					%d
  Description:	    			%s
  BurnOwnerDisabled:  			%t 
  BurnHolderDisabled:  			%t 
  BurnFromDisabled:  			%t 
  FreezeDisabled:  				%t 
  MintingFinished:  			%t `,
		ci.IssueId, ci.Issuer.String(), ci.Owner.String(), ci.Name, ci.Symbol, ci.TotalSupply.String(),
		ci.Decimals, ci.IssueTime, ci.Description, ci.BurnOwnerDisabled, ci.BurnHolderDisabled,
		ci.BurnFromDisabled, ci.FreezeDisabled, ci.MintingFinished)
}

//nolint
func (coinIssues CoinIssues) String() string {
	out := fmt.Sprintf("%-17s|%-44s|%-10s|%-6s|%-18s|%-8s|%s\n",
		"IssueID", "Owner", "Name", "Symbol", "TotalSupply", "Decimals", "IssueTime")
	for _, issue := range coinIssues {
		out += fmt.Sprintf("%-17s|%-44s|%-10s|%-6s|%-18s|%-8d|%d\n",
			issue.IssueId, issue.GetOwner().String(), issue.Name, issue.Symbol, issue.TotalSupply.String(), issue.Decimals, issue.IssueTime)
	}
	return strings.TrimSpace(out)
}

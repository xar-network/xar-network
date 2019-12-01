package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Box tags
var (
	TxCategory = "escrow"
	Feature    = "feature"
	Fee        = "fee"
	Owner      = "owner"
	Operation  = "operation"
	Interest   = "interest"
	BoxID      = "id"
	Status     = "status"
	Seq        = "seq"
	Transfer   = "transfer"
)

type LockBox struct {
	EndTime int64 `json:"end_time"`
}

//nolint
func (bi LockBox) String() string {
	return fmt.Sprintf(`LockInfo:
  EndTime:			%d`,
		bi.EndTime)
}

var InterestOperation = map[string]uint{Inject: 1, Cancel: 2}

func CheckInterestOperation(interestOperation string) sdk.Error {
	_, ok := InterestOperation[interestOperation]
	if !ok {
		return sdk.ErrInternal("unknown interest operation:" + interestOperation)
	}
	return nil
}

const ()

var Features = map[string]int{Transfer: 1}

var DepositOperation = map[string]uint{Inject: 1, Cancel: 2}

func CheckDepositOperation(depositOperation string) sdk.Error {
	_, ok := DepositOperation[depositOperation]
	if !ok {
		return sdk.ErrInternal("unknown deposit operation:" + depositOperation)
	}
	return nil
}

type DepositBox struct {
	StartTime          int64           `json:"start_time"`
	EstablishTime      int64           `json:"establish_time"`
	MaturityTime       int64           `json:"maturity_time"`
	BottomLine         sdk.Int         `json:"bottom_line"`
	Interest           BoxToken        `json:"interest"`
	Price              sdk.Int         `json:"price"`
	PerCoupon          sdk.Dec         `json:"per_coupon"`
	Share              sdk.Int         `json:"share"`
	TotalInject        sdk.Int         `json:"total_inject"`
	WithdrawalShare    sdk.Int         `json:"withdrawal_share"`
	WithdrawalInterest sdk.Int         `json:"withdrawal_interest"`
	InterestInjects    []AddressInject `json:"interest_injects"`
}

type DepositBoxInjectInterest struct {
	Address  sdk.AccAddress `json:"address"`
	Amount   sdk.Int        `json:"amount"`
	Interest sdk.Int        `json:"interest"`
}

type DepositBoxInjectInterestList []DepositBoxInjectInterest

//nolint
func (bi DepositBoxInjectInterest) String() string {
	return fmt.Sprintf(`
  Address:			%s
  Amount:			%s
  Interest:			%s`,
		bi.Address.String(), bi.Amount.String(), bi.Interest.String())
}

//nolint
func (bi DepositBox) String() string {
	return fmt.Sprintf(`DepositInfo:
  StartTime:			%d
  EstablishTime:		%d
  MaturityTime:			%d
  BottomLine:			%s
  Interest:			%s
  Price:			%s
  PerCoupon:			%s
  Share:			%s
  TotalInject:			%s
  WithdrawalShare:			%s,
  WithdrawalInterest:			%s,
  InterestInject:			%s`,
		bi.StartTime,
		bi.EstablishTime,
		bi.MaturityTime,
		bi.BottomLine.String(),
		bi.Interest.String(),
		bi.Price.String(),
		bi.PerCoupon.String(),
		bi.Share.String(),
		bi.TotalInject.String(),
		bi.WithdrawalShare.String(),
		bi.WithdrawalInterest.String(),
		bi.InterestInjects)
}

//nolint
func (bi DepositBoxInjectInterestList) String() string {
	out := fmt.Sprintf("%-44s|%-40s|%s\n",
		"Address", "Amount", "Interest")
	for _, box := range bi {
		out += fmt.Sprintf("%-44s|%-40s|%s\n",
			box.Address.String(), box.Amount.String(), box.Interest.String())
	}
	return strings.TrimSpace(out)
}

const (
	Lock    = "lock"
	Deposit = "deposit"
	Future  = "future"
)

var BoxType = map[string]string{Lock: "aa", Deposit: "ab", Future: "ac"}

func GetMustBoxTypeValue(boxType string) string {
	value, ok := BoxType[boxType]
	if !ok {
		panic("unknown type")
	}
	return value
}

func CheckBoxType(boxType string) sdk.Error {
	_, ok := BoxType[boxType]
	if !ok {
		return sdk.ErrInternal("unknown type:" + boxType)
	}
	return nil
}

func GetBoxTypeValue(boxType string) (string, error) {
	value, ok := BoxType[boxType]
	if !ok {
		return "", fmt.Errorf("unknown type:%s", boxType)
	}
	return value, nil
}

type BoxToken struct {
	Token    sdk.Coin `json:"token"`
	Decimals uint     `json:"decimals"`
}

//nolint
func (bi BoxToken) String() string {
	return fmt.Sprintf(`
  Token:			%s
  Decimals:			%d`,
		bi.Token.String(), bi.Decimals)
}

//Box interface
type Box interface {
	GetId() string
	SetId(string)

	GetBoxType() string
	SetBoxType(string)

	GetStatus() string
	SetStatus(string)

	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress)

	GetCreatedTime() int64
	SetCreatedTime(int64)

	GetName() string
	SetName(string)

	GetTotalAmount() BoxToken
	SetTotalAmount(BoxToken)

	GetDescription() string
	SetDescription(string)

	IsTransferDisabled() bool
	SetTransferDisabled(bool)

	GetLock() LockBox
	SetLock(LockBox)

	GetDeposit() DepositBox
	SetDeposit(DepositBox)

	GetFuture() FutureBox
	SetFuture(FutureBox)

	String() string
}

// BoxInfos is an array of BoxInfo
type BoxInfos []BoxInfo

//type BaseBoxInfo struct {
//}
type BoxInfo struct {
	Id               string         `json:"id"`
	Status           string         `json:"status"`
	Owner            sdk.AccAddress `json:"owner"`
	Name             string         `json:"name"`
	BoxType          string         `json:"type"`
	CreatedTime      int64          `json:"created_time"`
	TotalAmount      BoxToken       `json:"total_amount"`
	Description      string         `json:"description"`
	TransferDisabled bool           `json:"transfer_disabled"`
	Lock             LockBox        `json:"lock"`
	Deposit          DepositBox     `json:"deposit"`
	Future           FutureBox      `json:"future"`
}

// Implements Box Interface
var _ Box = (*BoxInfo)(nil)

func (bi BoxInfo) GetId() string {
	return bi.Id
}
func (bi *BoxInfo) SetId(boxId string) {
	bi.Id = boxId
}
func (bi BoxInfo) GetBoxType() string {
	return bi.BoxType
}
func (bi *BoxInfo) SetBoxType(boxType string) {
	bi.BoxType = boxType
}
func (bi BoxInfo) GetStatus() string {
	return bi.Status
}
func (bi *BoxInfo) SetStatus(boxStatus string) {
	bi.Status = boxStatus
}
func (bi BoxInfo) GetOwner() sdk.AccAddress {
	return bi.Owner
}
func (bi *BoxInfo) SetOwner(owner sdk.AccAddress) {
	bi.Owner = owner
}
func (bi BoxInfo) GetCreatedTime() int64 {
	return bi.CreatedTime
}
func (bi *BoxInfo) SetCreatedTime(createdTime int64) {
	bi.CreatedTime = createdTime
}
func (bi BoxInfo) GetName() string {
	return bi.Name
}
func (bi *BoxInfo) SetName(name string) {
	bi.Name = name
}
func (bi BoxInfo) GetTotalAmount() BoxToken {
	return bi.TotalAmount
}
func (bi *BoxInfo) SetTotalAmount(totalAmount BoxToken) {
	bi.TotalAmount = totalAmount
}
func (bi BoxInfo) GetDescription() string {
	return bi.Description
}
func (bi *BoxInfo) SetDescription(description string) {
	bi.Description = description
}

func (bi BoxInfo) IsTransferDisabled() bool {
	return bi.TransferDisabled
}

func (bi *BoxInfo) SetTransferDisabled(transferDisabled bool) {
	bi.TransferDisabled = transferDisabled
}

func (bi BoxInfo) GetLock() LockBox {
	return bi.Lock
}
func (bi *BoxInfo) SetLock(lock LockBox) {
	bi.Lock = lock
}

func (bi BoxInfo) GetDeposit() DepositBox {
	return bi.Deposit
}
func (bi *BoxInfo) SetDeposit(deposit DepositBox) {
	bi.Deposit = deposit
}

func (bi BoxInfo) GetFuture() FutureBox {
	return bi.Future
}
func (bi *BoxInfo) SetFuture(future FutureBox) {
	bi.Future = future
}

type AddressInject struct {
	Address sdk.AccAddress `json:"address"`
	Amount  sdk.Int        `json:"amount"`
}

func NewAddressInject(address sdk.AccAddress, amount sdk.Int) AddressInject {
	return AddressInject{address, amount}
}
func (bi AddressInject) String() string {
	return fmt.Sprintf(`
  Address:			%s
  Amount:			%s`,
		bi.Address.String(), bi.Amount.String())
}

//nolint
func (bi BoxInfo) String() string {
	return fmt.Sprintf(`Box:
  Id: 	         			%s
  Status:					%s
  Owner:           				%s
  Name:             			%s
  TotalAmount:      			%s
  CreatedTime:					%d
  Description:	    			%s
  TransferDisabled:			%t`,
		bi.Id, bi.Status, bi.Owner.String(), bi.Name, bi.TotalAmount.String(),
		bi.CreatedTime, bi.Description, bi.TransferDisabled)
}

//nolint
func (bi BoxInfos) String() string {
	out := fmt.Sprintf("%-17s|%-10s|%-44s|%-16s|%s\n",
		"BoxID", "Status", "Owner", "Name", "BoxTime")
	for _, box := range bi {
		out += fmt.Sprintf("%-17s|%-10s|%-44s|%-16s|%d\n",
			box.GetId(), box.GetStatus(), box.GetOwner().String(), box.GetName(), box.GetCreatedTime())
	}
	return strings.TrimSpace(out)
}

type LockBoxInfo struct {
	Id               string         `json:"id"`
	BoxType          string         `json:"type"`
	Status           string         `json:"status"`
	Owner            sdk.AccAddress `json:"owner"`
	Name             string         `json:"name"`
	CreatedTime      int64          `json:"created_time"`
	TotalAmount      BoxToken       `json:"total_amount"`
	Description      string         `json:"description"`
	TransferDisabled bool           `json:"transfer_disabled"`
	Lock             LockBox        `json:"lock"`
}
type DepositBoxInfo struct {
	Id               string         `json:"id"`
	BoxType          string         `json:"type"`
	Status           string         `json:"status"`
	Owner            sdk.AccAddress `json:"owner"`
	Name             string         `json:"name"`
	CreatedTime      int64          `json:"created_time"`
	TotalAmount      BoxToken       `json:"total_amount"`
	Description      string         `json:"description"`
	TransferDisabled bool           `json:"transfer_disabled"`
	Deposit          DepositBox     `json:"deposit"`
}
type FutureBoxInfo struct {
	Id               string         `json:"id"`
	BoxType          string         `json:"type"`
	Status           string         `json:"status"`
	Owner            sdk.AccAddress `json:"owner"`
	Name             string         `json:"name"`
	CreatedTime      int64          `json:"created_time"`
	TotalAmount      BoxToken       `json:"total_amount"`
	Description      string         `json:"description"`
	TransferDisabled bool           `json:"transfer_disabled"`
	Future           FutureBox      `json:"future"`
}
type LockBoxInfos []LockBoxInfo
type DepositBoxInfos []DepositBoxInfo
type FutureBoxInfos []FutureBoxInfo

//nolint
func getString(Id string, Status string, Owner sdk.AccAddress, Name string, CreatedTime int64,
	TotalAmount BoxToken, Description string, TransferDisabled bool) string {
	return fmt.Sprintf(`BoxInfo:
  Id:				%s
  Status:			%s
  Owner:			%s
  Name:				%s
  TotalAmount:			%s
  CreatedTime:			%d
  Description:			%s
  TransferDisabled:		%t`,
		Id, Status, Owner.String(), Name, TotalAmount.String(),
		CreatedTime, Description, TransferDisabled)
}

//nolint
func (bi LockBoxInfo) String() string {
	str := getString(bi.Id, bi.Status, bi.Owner, bi.Name,
		bi.CreatedTime, bi.TotalAmount, bi.Description, bi.TransferDisabled)

	return fmt.Sprintf(`%s
%s`, str, bi.Lock.String())
}

//nolint
func (bi DepositBoxInfo) String() string {
	str := getString(bi.Id, bi.Status, bi.Owner, bi.Name,
		bi.CreatedTime, bi.TotalAmount, bi.Description, bi.TransferDisabled)

	return fmt.Sprintf(`%s
%s`, str, bi.Deposit.String())
}

//nolint
func (bi FutureBoxInfo) String() string {
	str := getString(bi.Id, bi.Status, bi.Owner, bi.Name,
		bi.CreatedTime, bi.TotalAmount, bi.Description, bi.TransferDisabled)

	return fmt.Sprintf(`%s
%s`, str, bi.Future.String())
}

//nolint
func (bi LockBoxInfos) String() string {
	out := fmt.Sprintf("%-17s|%-44s|%-16s|%-40s|%s\n",
		"BoxID", "Owner", "Name", "TotalAmount", "EndTime")
	for _, box := range bi {
		out += fmt.Sprintf("%-17s|%-44s|%-16s|%-40s|%s\n",
			box.Id, box.Owner.String(), box.Name, box.TotalAmount.Token.String(), time.Unix(box.Lock.EndTime, 0).String())
	}
	return strings.TrimSpace(out)
}

//nolint
func (bi DepositBoxInfos) String() string {
	out := fmt.Sprintf("%-17s|%-44s|%-16s|%-40s|%s\n",
		"BoxID", "Owner", "Name", "TotalAmount", "CreatedTime")
	for _, box := range bi {
		out += fmt.Sprintf("%-17s|%-44s|%-16s|%-40s|%d\n",
			box.Id, box.Owner.String(), box.Name, box.TotalAmount.Token.String(), box.CreatedTime)
	}
	return strings.TrimSpace(out)
}

//nolint
func (bi FutureBoxInfos) String() string {
	out := fmt.Sprintf("%-17s|%-44s|%-16s|%-40s|%s\n",
		"BoxID", "Owner", "Name", "TotalAmount", "CreatedTime")
	for _, box := range bi {
		out += fmt.Sprintf("%-17s|%-44s|%-16s|%-40s|%d\n",
			box.Id, box.Owner.String(), box.Name, box.TotalAmount.Token.String(), box.CreatedTime)
	}
	return strings.TrimSpace(out)
}

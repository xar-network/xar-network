package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func GetBoxTags(id string, boxType string, sender sdk.AccAddress) sdk.Tags {
	return sdk.NewTags(
		Category, boxType,
		BoxID, id,
		Sender, sender.String(),
	)
}

type FutureBox struct {
	Injects         []AddressInject `json:"injects"`
	TimeLine        []int64         `json:"time"`
	Receivers       [][]string      `json:"receivers"`
	TotalWithdrawal sdk.Int         `json:"total_withdrawal"`
}

//nolint
func (bi FutureBox) String() string {
	return fmt.Sprintf(`FutureInfo:
  Injects:			%s
  TimeLine:			%d
  Receivers:			%s
  TotalWithdrawal:			%s`,
		bi.Injects, bi.TimeLine, bi.Receivers, bi.TotalWithdrawal.String())
}

func GetCliContext(cdc *codec.Codec) (authtxb.TxBuilder, context.CLIContext, auth.Account, error) {
	txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().
		WithCodec(cdc).
		WithAccountDecoder(cdc)
	from := cliCtx.GetFromAddress()
	account, err := cliCtx.GetAccount(from)

	return txBldr, cliCtx, account, err
}

func GetBoxByID(cdc *codec.Codec, cliCtx context.CLIContext, id string) (Box, error) {
	var boxInfo Box
	// Query the box
	res, err := boxqueriers.QueryBoxByID(id, cliCtx)
	if err != nil {
		return nil, err
	}
	cdc.MustUnmarshalJSON(res, &boxInfo)
	return boxInfo, nil
}

func BoxOwnerCheck(cdc *codec.Codec, cliCtx context.CLIContext, sender auth.Account, id string) (Box, error) {
	boxInfo, err := GetBoxByID(cdc, cliCtx, id)
	if err != nil {
		return nil, err
	}
	if !sender.GetAddress().Equals(boxInfo.GetOwner()) {
		return nil, errors.Errorf(errors.ErrOwnerMismatch(id))
	}
	return boxInfo, nil
}

func GetCoinDecimal(cdc *codec.Codec, cliCtx context.CLIContext, coin sdk.Coin) (uint, error) {
	if coin.Denom == Agard {
		return AgardDecimals, nil
	}
	issueInfo, err := issueclientutils.GetIssueByID(cdc, cliCtx, coin.Denom)
	if err != nil {
		return 0, err
	}
	return issueInfo.GetDecimals(), nil
}
func GetBoxInfo(box BoxInfo) fmt.Stringer {
	switch box.BoxType {
	case Lock:
		var clientBox boxclienttype.LockBoxInfo
		StructCopy(&clientBox, &box)
		return clientBox
	case Deposit:
		return processDepositBoxInfo(box)
	case Future:
		var clientBox boxclienttype.FutureBoxInfo
		StructCopy(&clientBox, &box)
		return clientBox
	default:
		return box
	}
}
func GetBoxParams(params Params, boxType string) fmt.Stringer {
	switch boxType {
	case Lock:
		var clientParams LockBoxParams
		StructCopy(&clientParams, &params)
		return clientParams
	case Deposit:
		var clientParams DepositBoxParams
		StructCopy(&clientParams, &params)
		return clientParams
	case Future:
		var clientParams FutureBoxParams
		StructCopy(&clientParams, &params)
		return clientParams
	default:
		return params
	}
}
func processDepositBoxInfo(box BoxInfo) boxclienttype.DepositBoxInfo {
	var clientBox boxclienttype.DepositBoxInfo
	StructCopy(&clientBox, &box)
	return clientBox
}
func GetBoxList(boxs BoxInfos, boxType string) fmt.Stringer {
	switch boxType {
	case Lock:
		var boxInfos = make(boxclienttype.LockBoxInfos, 0, len(boxs))
		for _, box := range boxs {
			var clientBox boxclienttype.LockBoxInfo
			StructCopy(&clientBox, &box)
			boxInfos = append(boxInfos, clientBox)
		}

		return boxInfos
	case Deposit:
		var boxInfos = make(boxclienttype.DepositBoxInfos, 0, len(boxs))
		for _, box := range boxs {
			boxInfos = append(boxInfos, processDepositBoxInfo(box))
		}

		return boxInfos
	case Future:
		var boxInfos = make(boxclienttype.FutureBoxInfos, 0, len(boxs))
		for _, box := range boxs {
			var clientBox boxclienttype.FutureBoxInfo
			StructCopy(&clientBox, &box)
			boxInfos = append(boxInfos, clientBox)
		}

		return boxInfos
	}
	return boxs
}
func deepFields(faceType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < faceType.NumField(); i++ {
		v := faceType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, deepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}
	return fields
}

func StructCopy(destPtr interface{}, srcPtr interface{}) {
	srcv := reflect.ValueOf(srcPtr)
	dstv := reflect.ValueOf(destPtr)
	srct := reflect.TypeOf(srcPtr)
	dstt := reflect.TypeOf(destPtr)
	if srct.Kind() != reflect.Ptr || dstt.Kind() != reflect.Ptr ||
		srct.Elem().Kind() == reflect.Ptr || dstt.Elem().Kind() == reflect.Ptr {
		panic("Fatal error:type of parameters must be Ptr of value")
	}
	if srcv.IsNil() || dstv.IsNil() {
		panic("Fatal error:value of parameters should not be nil")
	}
	srcV := srcv.Elem()
	dstV := dstv.Elem()
	srcfields := deepFields(reflect.ValueOf(srcPtr).Elem().Type())
	for _, v := range srcfields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return
}
func GetInjectMsg(cdc *codec.Codec, cliCtx context.CLIContext, account auth.Account,
	id string, amountStr string, operation string, cli bool) (sdk.Msg, error) {

	if err := boxutils.CheckId(id); err != nil {
		return nil, errors.Errorf(err)
	}
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		return nil, errors.Errorf(errors.ErrAmountNotValid(amountStr))
	}

	boxInfo, err := GetBoxByID(cdc, cliCtx, id)
	if err != nil {
		return nil, err
	}
	if BoxInjecting != boxInfo.GetStatus() {
		return nil, errors.Errorf(errors.ErrNotAllowedOperation(boxInfo.GetStatus()))
	}
	if cli {
		decimal, err := GetCoinDecimal(cdc, cliCtx, boxInfo.GetTotalAmount().Token)
		if err != nil {
			return nil, err
		}
		amount = boxutils.MulDecimals(boxutils.ParseCoin(boxInfo.GetTotalAmount().Token.Denom, amount), decimal)
	}
	var msg sdk.Msg
	switch operation {
	case Inject:
		if err = checkAmountByInject(amount, boxInfo); err != nil {
			return nil, err
		}
		msg = msgs.NewMsgBoxInject(id, account.GetAddress(),
			sdk.NewCoin(boxInfo.GetTotalAmount().Token.Denom, amount))
	case Cancel:
		if err = checkAmountByCancel(amount, boxInfo, account); err != nil {
			return nil, err
		}
		msg = msgs.NewMsgBoxInjectCancel(id, account.GetAddress(),
			sdk.NewCoin(boxInfo.GetTotalAmount().Token.Denom, amount))
	default:
		return nil, errors.ErrNotSupportOperation()
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, errors.Errorf(err)
	}
	return msg, nil
}
func GetInterestMsg(cdc *codec.Codec, cliCtx context.CLIContext, account auth.Account,
	id string, amountStr string, operation string, cli bool) (sdk.Msg, error) {

	if err := boxutils.CheckId(id); err != nil {
		return nil, errors.Errorf(err)
	}
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		return nil, errors.Errorf(errors.ErrAmountNotValid(amountStr))
	}
	box, err := GetBoxByID(cdc, cliCtx, id)
	if err != nil {
		return nil, err
	}
	if box.GetBoxType() != Deposit {
		return nil, errors.Errorf(errors.ErrNotSupportOperation())
	}
	if box.GetStatus() != BoxCreated {
		return nil, errors.Errorf(errors.ErrNotSupportOperation())
	}
	if cli {
		decimal, err := GetCoinDecimal(cdc, cliCtx, box.GetDeposit().Interest.Token)
		if err != nil {
			return nil, err
		}
		amount = boxutils.MulDecimals(boxutils.ParseCoin(box.GetDeposit().Interest.Token.Denom, amount), decimal)
	}
	var msg sdk.Msg
	switch operation {
	case Cancel:
		flag := true
		for i, v := range box.GetDeposit().InterestInjects {
			if v.Address.Equals(account.GetAddress()) {
				if box.GetDeposit().InterestInjects[i].Amount.GTE(amount) {
					flag = false
					break
				}
			}
		}
		if flag {
			return nil, errors.ErrNotEnoughAmount()
		}
		msg = msgs.NewMsgBoxInterestCancel(id, account.GetAddress(), sdk.NewCoin(box.GetDeposit().Interest.Token.Denom, amount))
	case Inject:
		if box.GetDeposit().InterestInjects != nil {
			totalInterest := sdk.ZeroInt()
			for _, v := range box.GetDeposit().InterestInjects {
				if v.Address.Equals(account.GetAddress()) {
					totalInterest = totalInterest.Add(v.Amount)
				}
			}
			if totalInterest.Add(amount).GT(box.GetDeposit().Interest.Token.Amount) {
				return nil, errors.Errorf(errors.ErrInterestInjectNotValid(
					sdk.NewCoin(box.GetDeposit().Interest.Token.Denom, amount)))
			}
		}
		msg = msgs.NewMsgBoxInterestInject(id, account.GetAddress(), sdk.NewCoin(box.GetDeposit().Interest.Token.Denom, amount))
	default:
		return nil, errors.ErrUnknownOperation()
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, errors.Errorf(err)
	}
	return msg, nil
}
func GetWithdrawMsg(cdc *codec.Codec, cliCtx context.CLIContext, account auth.Account, id string) (sdk.Msg, error) {
	if account.GetCoins().AmountOf(id).IsZero() {
		return nil, errors.Errorf(errors.ErrNotEnoughAmount())
	}
	boxInfo, err := GetBoxByID(cdc, cliCtx, id)
	if err != nil {
		return nil, err
	}
	switch boxInfo.GetBoxType() {
	case Deposit:
		if BoxFinished != boxInfo.GetStatus() {
			return nil, errors.Errorf(errors.ErrNotAllowedOperation(boxInfo.GetStatus()))
		}
	case Future:
		if BoxCreated == boxInfo.GetStatus() {
			return nil, errors.Errorf(errors.ErrNotAllowedOperation(boxInfo.GetStatus()))
		}
		seq := boxutils.GetSeqFromFutureBoxSeq(id)
		if boxInfo.GetFuture().TimeLine[seq-1] > time.Now().Unix() {
			return nil, errors.Errorf(errors.ErrNotAllowedOperation(BoxUndue))
		}
	default:
		return nil, errors.Errorf(errors.ErrNotAllowedOperation(boxInfo.GetStatus()))
	}
	msg := msgs.NewMsgBoxWithdraw(id, account.GetAddress())
	if err := msg.ValidateBasic(); err != nil {
		return nil, errors.Errorf(err)
	}
	return msg, nil
}
func checkAmountByCancel(amount sdk.Int, boxInfo Box, account auth.Account) error {
	switch boxInfo.GetBoxType() {
	case Deposit:
		if !amount.Mod(boxInfo.GetDeposit().Price).IsZero() {
			return errors.ErrAmountNotValid(amount.String())
		}
		if account.GetCoins().AmountOf(boxInfo.GetId()).LT(amount.Quo(boxInfo.GetDeposit().Price)) {
			return errors.Errorf(errors.ErrNotEnoughAmount())
		}
	case Future:
		if boxInfo.GetFuture().Injects == nil {
			return errors.Errorf(errors.ErrNotEnoughAmount())
		}
		for _, v := range boxInfo.GetFuture().Injects {
			if v.Address.Equals(account.GetAddress()) {
				if v.Amount.GTE(amount) {
					return nil
				}
			}
		}
		return errors.Errorf(errors.ErrNotEnoughAmount())
	default:
		return errors.Errorf(errors.ErrNotSupportOperation())
	}
	return nil
}
func checkAmountByInject(amount sdk.Int, boxInfo Box) error {
	switch boxInfo.GetBoxType() {
	case Deposit:
		if !amount.Mod(boxInfo.GetDeposit().Price).IsZero() {
			return errors.ErrAmountNotValid(amount.String())
		}
		if amount.Add(boxInfo.GetDeposit().TotalInject).GT(boxInfo.GetTotalAmount().Token.Amount) {
			return errors.Errorf(errors.ErrNotEnoughAmount())
		}
	case Future:
		total := sdk.ZeroInt()
		if boxInfo.GetFuture().Injects != nil {
			for _, v := range boxInfo.GetFuture().Injects {
				total = total.Add(v.Amount)
			}
		}
		if amount.Add(total).GT(boxInfo.GetTotalAmount().Token.Amount) {
			return errors.Errorf(errors.ErrNotEnoughAmount())
		}
	default:
		return errors.Errorf(errors.ErrNotSupportOperation())
	}
	return nil
}

func MulDecimals(coin sdk.Coin, decimals uint) sdk.Int {
	if coin.Denom == Agard {
		return coin.Amount
	}
	return issueutils.MulDecimals(coin.Amount, decimals)
}
func ParseCoin(denom string, amount sdk.Int) sdk.Coin {
	if denom == Agard {
		denom = Gard
	}
	coin, _ := sdk.ParseCoin(fmt.Sprintf("%s%s", amount, denom))
	return coin
}
func CalcInterest(perCoupon sdk.Dec, share sdk.Int, interest BoxToken) sdk.Int {
	dec := perCoupon.MulInt(share)
	decimals := interest.Decimals
	if interest.Token.Denom == Agard {
		decimals = GardDecimals
	}
	dec = GetMaxPrecision(dec, decimals)
	return dec.MulInt(issueutils.GetDecimalsInt(decimals)).TruncateInt()
}

func IsId(id string) bool {
	return strings.HasPrefix(id, IDPreStr)
}

func CheckId(id string) sdk.Error {
	if !IsId(id) {
		return errors.ErrBoxID(id)
	}
	return nil
}

func CalcInterestRate(totalAmount sdk.Int, price sdk.Int, interest sdk.Coin, decimals uint) sdk.Dec {
	totalCoupon := totalAmount.Quo(price)
	perCoupon := sdk.NewDecFromBigInt(interest.Amount.BigInt()).QuoInt(totalCoupon)
	if interest.Denom == Agard {
		decimals = GardDecimals
	}
	return quoMaxPrecisionByDecimal(perCoupon, decimals)
}

func quoMaxPrecisionByDecimal(dec sdk.Dec, decimals uint) sdk.Dec {
	dec = dec.QuoInt(issueutils.GetDecimalsInt(decimals))
	dec = GetMaxPrecision(dec, decimals)
	return dec
}

func GetBoxTypeByValue(value string) string {
	value = strings.ReplaceAll(value, IDPreStr, "")
	for k, v := range BoxType {
		if strings.HasPrefix(value, v) {
			return k
		}
	}
	return ""
}
func GetCoinDenomByFutureBoxSeq(id string, seq int) string {
	return fmt.Sprintf("%s%02d", id, seq)
}
func GetIdFromBoxSeqID(idSeq string) string {
	if len(idSeq) > IdLength {
		return idSeq[:IdLength]
	}
	return idSeq
}
func GetSeqFromFutureBoxSeq(boxSeqStr string) int {
	seqStr := boxSeqStr[len(boxSeqStr)-2:]
	seq, _ := strconv.Atoi(seqStr)
	return seq
}
func GetMaxPrecision(dec sdk.Dec, decimals uint) sdk.Dec {
	precision := MaxPrecision
	if decimals < MaxPrecision {
		precision = decimals
	}
	decStr := dec.String()
	len := strings.Index(decStr, ".") + int(precision)
	str := decStr[0 : len+1]
	dec, _ = sdk.NewDecFromStr(str)
	return dec
}

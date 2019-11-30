package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type MsgInterface interface{ sdk.Msg }

func validateError(cases []struct {
	valid bool
	tx    MsgInterface
}, t *testing.T) {
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, fmt.Sprintf("transaction [no: %d] [%v] failed but was supposed to be valid", i, tc.tx))
		} else {
			require.NotNil(t, err, fmt.Sprintf("transaction [no: %d] [%v] is valid but is supposed to have an error", i, tc.tx))
		}
	}
}

// Tests

func TestMsgIssueToken(t *testing.T) {
	var (
		name                 = "Zap"
		originalSymbol       = "ZAP"
		symbol               = "zap"
		total          int64 = 1
		owner                = sdk.AccAddress([]byte("me"))
		msg                  = NewMsgIssueToken(owner, name, symbol, originalSymbol, total, false)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "issue_token")
}

func TestMsgIssueTokenValidation(t *testing.T) {
	var (
		name                 = "Zap"
		originalSymbol       = "ZAP"
		symbol               = "zap"
		total          int64 = 1
		totalInvalid   int64 = 0
		acc                  = sdk.AccAddress([]byte("me"))
		name2                = "a"
		total2         int64 = 2
		acc2                 = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgIssueToken(acc, name, symbol, originalSymbol, total, false)},
		{true, NewMsgIssueToken(acc, name, symbol, originalSymbol, total, false)},
		{false, NewMsgIssueToken(acc, name, symbol, originalSymbol, totalInvalid, false)},
		{true, NewMsgIssueToken(acc2, name2, symbol, originalSymbol, total2, false)},
		{true, NewMsgIssueToken(acc2, name2, symbol, originalSymbol, total, false)},
		{true, NewMsgIssueToken(acc, name2, symbol, originalSymbol, total2, false)},
		{false, NewMsgIssueToken(nil, name, symbol, originalSymbol, total2, false)},
		{false, NewMsgIssueToken(acc2, "", symbol, originalSymbol, total2, false)},
		{false, NewMsgIssueToken(acc2, name, symbol, originalSymbol, totalInvalid, false)},
	}

	validateError(cases, t)
}

func TestMsgIssueTokenGetSignBytes(t *testing.T) {
	var (
		name                 = "Zap"
		originalSymbol       = "ZAP"
		symbol               = "zap"
		total          int64 = 1
		owner                = sdk.AccAddress([]byte("me"))
		msg                  = NewMsgIssueToken(owner, name, symbol, originalSymbol, total, false)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/IssueToken",` +
		`"value":{` +
		`"mintable":false,` +
		`"name":"Zap",` +
		`"original_symbol":"ZAP",` +
		`"source_address":"cosmos1d4js690r9j",` +
		`"symbol":"` + symbol + `",` +
		`"total_supply":"1"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgMintCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgMintCoins(amount, symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "mint_coins")
}

func TestMsgMintCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 10
		symbol        = "ZAP-001"
		symbol2       = "MNT-500"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgMintCoins(amount, symbol, owner)},
		{true, NewMsgMintCoins(amount, symbol2, owner2)},
		{false, NewMsgMintCoins(-1, symbol, owner)},
		{false, NewMsgMintCoins(0, symbol, owner)},
		{true, NewMsgMintCoins(1, symbol, owner)},
		{false, NewMsgMintCoins(amount, symbol, nil)},
		{false, NewMsgMintCoins(amount, "", owner)},
	}

	validateError(cases, t)
}

func TestMsgMintCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgMintCoins(amount, symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/MintCoins","value":{` +
		`"amount":"10",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"ZAP-001"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgBurnCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgBurnCoins(amount, symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "burn_coins")
}

func TestMsgBurnCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 20
		symbol        = "ZAP-001"
		symbol2       = "BRN-002"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgBurnCoins(amount, symbol, owner)},
		{true, NewMsgBurnCoins(amount, symbol2, owner2)},
		{false, NewMsgBurnCoins(-1, symbol, owner)},
		{false, NewMsgBurnCoins(0, symbol, owner)},
		{true, NewMsgBurnCoins(1, symbol, owner)},
		{false, NewMsgBurnCoins(amount, symbol, nil)},
		{false, NewMsgBurnCoins(amount, "", owner)},
	}

	validateError(cases, t)
}

func TestMsgBurnCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 100
		symbol       = "ZAP-999"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgBurnCoins(amount, symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/BurnCoins","value":{` +
		`"amount":"100",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"ZAP-999"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgFreezeCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgFreezeCoins(amount, symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "freeze_coins")
}

func TestMsgFreezeCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 15
		symbol        = "ZAP-001"
		symbol2       = "FRZ-112"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgFreezeCoins(amount, symbol, owner)},
		{true, NewMsgFreezeCoins(amount, symbol2, owner2)},
		{false, NewMsgFreezeCoins(-1, symbol, owner)},
		{false, NewMsgFreezeCoins(0, symbol, owner)},
		{true, NewMsgFreezeCoins(1, symbol, owner)},
		{false, NewMsgFreezeCoins(amount, symbol, nil)},
		{false, NewMsgFreezeCoins(amount, "", owner)},
	}

	validateError(cases, t)
}

func TestMsgFreezeCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 100
		symbol       = "FRZ-999"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgFreezeCoins(amount, symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/FreezeCoins","value":{` +
		`"amount":"100",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"FRZ-999"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgUnfreezeCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "UFZ-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgUnfreezeCoins(amount, symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "unfreeze_coins")
}

func TestMsgUnfreezeCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 15
		symbol        = "ZAP-001"
		symbol2       = "UFZ-130"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgUnfreezeCoins(amount, symbol, owner)},
		{true, NewMsgUnfreezeCoins(amount, symbol2, owner2)},
		{false, NewMsgUnfreezeCoins(-1, symbol, owner)},
		{false, NewMsgUnfreezeCoins(0, symbol, owner)},
		{true, NewMsgUnfreezeCoins(1, symbol, owner)},
		{false, NewMsgUnfreezeCoins(amount, symbol, nil)},
		{false, NewMsgUnfreezeCoins(amount, "", owner)},
	}

	validateError(cases, t)
}

func TestMsgUnfreezeCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 100
		symbol       = "UFZ-999"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgUnfreezeCoins(amount, symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"assetmanagement/UnfreezeCoins","value":{` +
		`"amount":"100",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"UFZ-999"}}`

	require.Equal(t, expected, string(actual))
}

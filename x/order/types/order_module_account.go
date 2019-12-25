package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

// Currently not used in code. made for a demonstration purpose.
// The main idea behind such an approach is to guarantee that if a coin is frozen - there is no chance to loose it from module account.
// What the problem actualy is?
// If we save frost coins as an object in a KvStore we do not guarantee that the related coins from module account won't be spent (e.g. as a result of a bug)
// But if we store frozen coins in a separate variable at the account, it won't be possible to call any transaction with frozen coins before unfreeze.
// So we attempt to provide safety with such an approach.
// The problem is that we must replace an original moduleAccount created for the OrderModule with our custom struct, which is not very clean.
type OrderModuleAccountI interface {
	exported.ModuleAccountI
	UnfreezeCoins(coinsToUnfreeze sdk.Coins) error
	FreezeCoins(coinsToFreeze sdk.Coins) error
}

type OrderModuleAccount struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins" yaml:"frozen_coins"`
	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	Name          string         `json:"name"`
	Permissions   []string       `json:"permissions"`
}

func OrderModuleFromModuleAccount(account exported.ModuleAccountI) *OrderModuleAccount {
	return &OrderModuleAccount{
		account.GetAddress(),
		account.GetCoins(),
		sdk.Coins{},
		account.GetPubKey(),
		account.GetAccountNumber(),
		account.GetSequence(),
		account.GetName(),
		account.GetPermissions(),
	}
}

// GetAddress - Implements sdk.Account.
func (o OrderModuleAccount) GetAddress() sdk.AccAddress {
	return o.Address
}

// SetAddress - Implements sdk.Account.
func (o *OrderModuleAccount) SetAddress(addr sdk.AccAddress) error {
	if len(o.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	o.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (o OrderModuleAccount) GetPubKey() crypto.PubKey {
	return o.PubKey
}

// SetPubKey - Implements sdk.Account.
func (o *OrderModuleAccount) SetPubKey(pubKey crypto.PubKey) error {
	o.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (o OrderModuleAccount) GetCoins() sdk.Coins {
	return o.Coins
}

// SetCoins - Implements sdk.Account.
func (o *OrderModuleAccount) SetCoins(coins sdk.Coins) error {
	o.Coins = coins

	return nil
}

func (o *OrderModuleAccount) FreezeCoins(coinsToFreeze sdk.Coins) error {
	for _, freezingCoin := range coinsToFreeze {
		coinAmt := o.Coins.AmountOf(freezingCoin.Denom)
		if coinAmt.LT(freezingCoin.Amount) {
			return errors.New(fmt.Sprintf("cannot freeze coin with denom %s: coin amount is %v, while it is requested to freeze %v", freezingCoin.Denom, coinAmt, freezingCoin.Amount))
		}

		//remove coin from available and add to frozen
		frostCoins := sdk.Coins{freezingCoin}
		o.Coins = o.Coins.Sub(frostCoins)
		o.FrozenCoins = o.FrozenCoins.Add(frostCoins)
	}

	return nil
}

func (o *OrderModuleAccount) UnfreezeCoins(coinsToUnfreeze sdk.Coins) error {
	for _, unfreezingCoin := range coinsToUnfreeze {
		coinAmt := o.FrozenCoins.AmountOf(unfreezingCoin.Denom)
		if coinAmt.LT(unfreezingCoin.Amount) {
			return errors.New(fmt.Sprintf("cannot unfreeze coin with denom %s: coin amount is %v, while it is requested to unfreeze %v", unfreezingCoin.Denom, coinAmt, unfreezingCoin.Amount))
		}

		//remove coin from available and add to frozen
		cns := sdk.Coins{unfreezingCoin}
		o.Coins = o.Coins.Sub(cns)
		o.FrozenCoins = o.FrozenCoins.Add(cns)
	}

	return nil
}

// if it is not allowed to place freezeCoins to a separate variable
// we can use something like this to evade panic falls
func (o OrderModuleAccount) checkCoins(newCoins sdk.Coins) bool {
	avalibleCoins := o.GetAvailableCoins()
	for _, newCoin := range newCoins {
		for _, oldCoin := range o.Coins {
			if oldCoin.Denom == newCoin.Denom {
				avalibleDecline := avalibleCoins.AmountOf(oldCoin.Denom)
				if !opIsAvailible(oldCoin, newCoin, avalibleDecline) {
					return false
				}
			}
		}
	}

	return true
}

func opIsAvailible(oldCoin, newCoin sdk.Coin, availableDeclineAmt sdk.Int) bool {
	if oldCoin.Denom != newCoin.Denom {
		return false
	}

	oldCoinAmt := oldCoin.Amount
	newCoinAmt := newCoin.Amount
	amtGrowth := newCoinAmt.Sub(oldCoinAmt)

	if amtGrowth.IsPositive() {
		return true
	}

	amtDecline := amtGrowth.Neg() // make positive

	return amtDecline.LTE(availableDeclineAmt)
}

// returns true only if all
func (o OrderModuleAccount) IsAdditionalCoins(newCoins sdk.Coins) bool {
	for _, newCoin := range newCoins {
		for _, oldCoin := range o.Coins {
			if newCoin.Denom == oldCoin.Denom {
				if newCoin.Amount.LT(oldCoin.Amount) {
					return false
				}
			}
		}
	}
	return true
}

// should be replaced with a variable?
func (o OrderModuleAccount) GetAvailableCoins() sdk.Coins {
	return o.Coins.Sub(o.FrozenCoins)
}

// GetAccountNumber - Implements Account
func (o OrderModuleAccount) GetAccountNumber() uint64 {
	return o.AccountNumber
}

// SetAccountNumber - Implements Account
func (o *OrderModuleAccount) SetAccountNumber(accNumber uint64) error {
	o.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (o OrderModuleAccount) GetSequence() uint64 {
	return o.Sequence
}

// SetSequence - Implements sdk.Account.
func (o *OrderModuleAccount) SetSequence(seq uint64) error {
	o.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (o OrderModuleAccount) SpendableCoins(_ time.Time) sdk.Coins {
	return o.GetCoins()
}

// Validate checks for errors on the account fields
func (o OrderModuleAccount) Validate() error {
	if o.PubKey != nil && o.Address != nil &&
		!bytes.Equal(o.PubKey.Address().Bytes(), o.Address.Bytes()) {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}

func (o OrderModuleAccount) String() string {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (o OrderModuleAccount) GetName() string {
	return o.Name
}

func (o OrderModuleAccount) GetPermissions() []string {
	return o.Permissions
}

func (o OrderModuleAccount) HasPermission(permission string) bool {
	for _, perm := range o.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

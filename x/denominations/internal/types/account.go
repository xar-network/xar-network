package types

import (
	"encoding/json"
	"errors"
	"fmt"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	_ authexported.GenesisAccount = (*FreezeAccount)(nil)
	_ exported.Account            = (*FreezeAccount)(nil)
)

func init() {
	authtypes.RegisterAccountTypeCodec(&FreezeAccount{}, "denominations/FreezeAccount")
}

// FreezeAccount is customised to allow temporary freezing of coins to exclude them from transactions
type FreezeAccount struct {
	*authtypes.BaseAccount
	FrozenCoins sdk.Coins `json:"frozen" yaml:"frozen"`
}

func NewFreezeAccount(ba *authtypes.BaseAccount, frozenCoins sdk.Coins) *FreezeAccount {

	return &FreezeAccount{
		BaseAccount: ba,
		FrozenCoins: frozenCoins,
	}
}

// GetFrozenCoins retrieves frozen coins from account
func (acc FreezeAccount) GetFrozenCoins() sdk.Coins {
	return acc.FrozenCoins
}

// SetFrozenCoins sets frozen coins for account
func (acc *FreezeAccount) SetFrozenCoins(frozen sdk.Coins) error {
	acc.FrozenCoins = frozen
	return nil
}

// Validate checks for errors on the account fields
func (acc FreezeAccount) Validate() error {
	if !acc.FrozenCoins.IsValid() {
		return errors.New("invalid coins")

	}
	return acc.BaseAccount.Validate()
}

func AreAnyCoinsZero(coins *sdk.Coins) bool {
	for _, coin := range *coins {
		if sdk.NewInt(0).Equal(coin.Amount) {
			return true
		}
	}
	return false
}

// FreezeCoins freezes unfrozen coins for account according to input
func (acc *FreezeAccount) FreezeCoins(coinsToFreeze sdk.Coins) error {
	// Have enough coins to freeze?
	if coinsToFreeze == nil || coinsToFreeze.Empty() || coinsToFreeze.IsAnyNegative() || AreAnyCoinsZero(&coinsToFreeze) {
		return sdk.ErrInvalidCoins("No coins chosen to freeze")
	}

	currentCoins := acc.GetCoins()
	if currentCoins == nil || currentCoins.IsAllLT(coinsToFreeze) {
		return sdk.ErrInvalidCoins("Not enough coins to freeze")
	}

	// Freeze coins
	if newBalance, isNegative := currentCoins.SafeSub(coinsToFreeze); !isNegative {
		if err := acc.SetCoins(newBalance); err != nil {
			return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set coins: %s", err))
		}
	} else {
		return sdk.ErrInternal("failed to subtract coins for freezing")
	}

	frozen := acc.GetFrozenCoins()
	if frozen == nil {
		frozen = coinsToFreeze
	} else {
		frozen = frozen.Add(coinsToFreeze)
	}

	if err := acc.SetFrozenCoins(frozen); err != nil {
		return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set frozen coins: %s", err))
	}

	return nil
}

// UnfreezeCoins unfreezes frozen coins for account according to input
func (acc *FreezeAccount) UnfreezeCoins(coinsToUnfreeze sdk.Coins) error {
	// Have enough coins to unfreeze?
	if coinsToUnfreeze == nil || coinsToUnfreeze.Empty() || coinsToUnfreeze.IsAnyNegative() {
		return sdk.ErrInvalidCoins("No coins chosen to unfreeze")
	}

	currentlyFrozen := acc.GetFrozenCoins()
	if currentlyFrozen == nil || currentlyFrozen.IsAllLT(coinsToUnfreeze) {
		return sdk.ErrInvalidCoins("Not enough coins to unfreeze")
	}

	// Unfreeze coins
	currentCoins := acc.GetCoins()
	if currentCoins == nil {
		currentCoins = coinsToUnfreeze
	} else {
		currentCoins = currentCoins.Add(coinsToUnfreeze)
	}

	if newFrozenBalance, isNegative := currentlyFrozen.SafeSub(coinsToUnfreeze); !isNegative {
		if err := acc.SetFrozenCoins(newFrozenBalance); err != nil {
			return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set frozen coins: %s", err))
		}
	} else {
		return sdk.ErrInternal("failed to subtract coins for unfreezing")
	}

	if err := acc.SetCoins(currentCoins); err != nil {
		return sdk.ErrInvalidCoins(fmt.Sprintf("failed to set coins: %s", err))
	}

	return nil
}

type freezeAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	FrozenCoins   sdk.Coins      `json:"frozen_coins" yaml:"frozen_coins"`
}

func (acc FreezeAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of a ModuleAccount.
func (acc FreezeAccount) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(freezeAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		PubKey:        "",
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		FrozenCoins:   acc.FrozenCoins,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

// MarshalJSON returns the JSON representation of a ModuleAccount.
func (acc FreezeAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(freezeAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		PubKey:        "",
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
		FrozenCoins:   acc.FrozenCoins,
	})
}

// UnmarshalJSON unmarshals raw JSON bytes into a ModuleAccount.
func (acc FreezeAccount) UnmarshalJSON(bz []byte) error {
	var alias freezeAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	acc.BaseAccount = authtypes.NewBaseAccount(alias.Address, alias.Coins, nil, alias.AccountNumber, alias.Sequence)
	acc.FrozenCoins = alias.FrozenCoins

	return nil
}

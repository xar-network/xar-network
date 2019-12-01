package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	fundsStoreKey sdk.StoreKey // key for the keystore that contains the pairs account -> deposited funds
	bankKeeper    bank.Keeper
	cdc           *codec.Codec
}

func NewKeeper(fundsStoreKey sdk.StoreKey, bankKeeper bank.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		fundsStoreKey: fundsStoreKey,
		bankKeeper:    bankKeeper,
		cdc:           cdc,
	}
}

// DepositFundFromAddress allows to take the given amount from the account balance and store it into the pool
func (k Keeper) DepositFundFromAddress(ctx sdk.Context, account sdk.AccAddress, amount sdk.Coin) sdk.Error {

	// remove the coins from the sender
	_, err := k.bankKeeper.SubtractCoins(ctx, account, []sdk.Coin{amount})
	if err != nil {
		return err
	}

	// get the pool
	pool, addError := sdk.AccAddressFromHex("00000000000000000000")
	if addError != nil {
		return sdk.ErrInternal(addError.Error())
	}

	// add the funds to the pool
	_, err = k.bankKeeper.AddCoins(ctx, pool, []sdk.Coin{amount})
	if err != nil {
		return err
	}

	// get any existing funds for the user
	coin, fundsError := k.GetAccountFunds(ctx, account)
	if fundsError != nil {
		coin = sdk.NewCoin(amount.Denom, sdk.NewInt(0))
	}

	if coin.Denom != amount.Denom {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid coin, expected %s got %s", coin.Denom, amount.Denom))
	}

	// add the coin amount and set the denom
	coin.Amount = coin.Amount.Add(amount.Amount)
	coin.Denom = amount.Denom

	// set the deposit
	bytes, marshalError := k.cdc.MarshalBinaryBare(coin)
	if marshalError != nil {
		return sdk.ErrInternal(marshalError.Error())
	}

	store := ctx.KVStore(k.fundsStoreKey)
	store.Set(account, bytes)

	return err
}

// WithdrawFundToAddress allows to take the specified amount from the pool and store into the given account balance
func (k Keeper) WithdrawFundToAddress(ctx sdk.Context, amount sdk.Coin, account sdk.AccAddress) sdk.Error {

	existingAmount, err := k.GetAccountFunds(ctx, account)
	if err != nil {
		return err
	}

	// check for valid denom
	if existingAmount.Denom != amount.Denom {
		return sdk.ErrInvalidCoins("specified address has no funds with given coin denom")
	}

	// check for valid amount
	if existingAmount.Amount.LT(amount.Amount) {
		return sdk.ErrInsufficientCoins("specified address has not enough funds to withdraw")
	}

	// update the funds status
	existingAmount.Amount = existingAmount.Amount.Sub(amount.Amount)
	existingAmount.Denom = amount.Denom

	store := ctx.KVStore(k.fundsStoreKey)
	store.Set(account, k.cdc.MustMarshalBinaryBare(existingAmount))

	// get the pool
	pool, addressErr := sdk.AccAddressFromHex("00000000000000000000")
	if addressErr != nil {
		return sdk.ErrInternal(addressErr.Error())
	}

	// remove the coins from the pool
	_, err = k.bankKeeper.SubtractCoins(ctx, pool, []sdk.Coin{amount})
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}

	// add the funds to the user balance
	_, err = k.bankKeeper.AddCoins(ctx, account, []sdk.Coin{amount})
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}

	return nil
}

// GetAccountFunds returns the funds of the given account
func (k Keeper) GetAccountFunds(ctx sdk.Context, account sdk.AccAddress) (sdk.Coin, sdk.Error) {

	// default funds have the given bond denom and 0
	var amount sdk.Coin

	// read the saved funds
	store := ctx.KVStore(k.fundsStoreKey)
	value := store.Get(account)

	// if nil return 0
	if value == nil {
		return amount, sdk.ErrInsufficientCoins("address has no funds")
	}

	// unwrap the value
	err := k.cdc.UnmarshalBinaryBare(value, &amount)
	if err != nil {
		return amount, sdk.ErrInternal(err.Error())
	}

	return amount, nil
}

// DistributeReward distributes the given reward between all the funders
func (k Keeper) DistributeReward(ctx sdk.Context) error {

	// find the total funds for the given reward coin
	funds, err := k.GetTotalFunds(ctx)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.fundsStoreKey)

	// iterate over all the funds coins
	for _, fundCoin := range funds {

		// divide the reward
		investorsIterator := store.Iterator(nil, nil)
		for ; investorsIterator.Valid(); investorsIterator.Next() {

			// get the fund
			var fund sdk.Coin
			k.cdc.MustUnmarshalBinaryBare(investorsIterator.Value(), &fund)

			// update only the same coins, not others
			if fund.Denom == fundCoin.Denom {

				// compute the dividend
				dividend := fund.Amount.Quo(fundCoin.Amount)

				// add it to the fund
				fund.Amount = fund.Amount.Add(dividend)

				// save it inside the store
				store.Set(investorsIterator.Key(), k.cdc.MustMarshalBinaryBare(fund))
			}
		}

	}

	return nil
}

func (k Keeper) GetTotalFunds(ctx sdk.Context) (sdk.Coins, sdk.Error) {
	var values sdk.Coins

	// compute all the funds sum
	funds := make(map[string]sdk.Int)

	store := ctx.KVStore(k.fundsStoreKey)
	sumIterator := store.Iterator(nil, nil)
	for ; sumIterator.Valid(); sumIterator.Next() {

		// get the fund
		var fund sdk.Coin
		err := k.cdc.UnmarshalBinaryBare(sumIterator.Value(), &fund)
		if err != nil {
			return values, sdk.ErrInternal(err.Error())
		}

		// sum up the total
		amount := fund.Amount
		if value, ok := funds[fund.Denom]; ok {
			amount = value.Add(amount)
		}

		funds[fund.Denom] = amount
	}

	for k, v := range funds {
		values = append(values, sdk.NewCoin(k, v))
	}

	return values, nil
}

package keeper

import (
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/oracle/internal/types"
)

// Keeper struct for oracle module
type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	codespace sdk.CodespaceType
}

// NewKeeper returns a new keeper for the oracle modle
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		codespace: codespace,
	}
}

// AddOracle adds an Oracle to the store
func (k Keeper) AddOracle(ctx sdk.Context, address string) {

	oracles := k.GetOracles(ctx)
	oracles = append(oracles, types.Oracle{OracleAddress: address})
	store := ctx.KVStore(k.storeKey)
	store.Set(
		[]byte(types.OraclePrefix), k.cdc.MustMarshalBinaryBare(oracles),
	)
}

// AddAsset adds an asset to the store
func (k Keeper) AddAsset(
	ctx sdk.Context,
	assetCode string,
	desc string,
) {
	assets := k.GetAssets(ctx)
	assets = append(assets, types.Asset{AssetCode: assetCode, Description: desc})
	store := ctx.KVStore(k.storeKey)
	store.Set(
		[]byte(types.AssetPrefix), k.cdc.MustMarshalBinaryBare(assets),
	)
}

// SetPrice updates the posted price for a specific oracle
func (k Keeper) SetPrice(
	ctx sdk.Context,
	oracle sdk.AccAddress,
	assetCode string,
	price sdk.Dec,
	expiry sdk.Int) (types.PostedPrice, sdk.Error) {
	// If the expiry is greater than or equal to the current blockheight, we consider the price valid
	if expiry.GTE(sdk.NewInt(ctx.BlockHeight())) {
		store := ctx.KVStore(k.storeKey)
		prices := k.GetRawPrices(ctx, assetCode)
		var index int
		found := false
		for i := range prices {
			if prices[i].OracleAddress == oracle.String() {
				index = i
				found = true
				break
			}
		}
		// set the price for that particular oracle
		if found {
			prices[index] = types.PostedPrice{AssetCode: assetCode, OracleAddress: oracle.String(), Price: price, Expiry: expiry}
		} else {
			prices = append(prices, types.PostedPrice{
				AssetCode: assetCode, OracleAddress: oracle.String(),
				Price: price, Expiry: expiry,
			})
			index = len(prices) - 1
		}

		store.Set(
			[]byte(types.RawPriceFeedPrefix+assetCode), k.cdc.MustMarshalBinaryBare(prices),
		)
		return prices[index], nil
	}
	return types.PostedPrice{}, types.ErrExpired(k.codespace)

}

// SetCurrentPrices updates the price of an asset to the median of all valid oracle inputs
func (k Keeper) SetCurrentPrices(ctx sdk.Context) sdk.Error {
	assets := k.GetAssets(ctx)
	for _, v := range assets {
		assetCode := v.AssetCode
		prices := k.GetRawPrices(ctx, assetCode)
		var notExpiredPrices []types.CurrentPrice
		// filter out expired prices
		for _, v := range prices {
			if v.Expiry.GTE(sdk.NewInt(ctx.BlockHeight())) {
				notExpiredPrices = append(notExpiredPrices, types.CurrentPrice{
					AssetCode: v.AssetCode,
					Price:     v.Price,
					Expiry:    v.Expiry,
				})
			}
		}
		l := len(notExpiredPrices)
		var medianPrice sdk.Dec
		var expiry sdk.Int
		// TODO make threshold for acceptance (ie. require 51% of oracles to have posted valid prices
		if l == 0 {
			// Error if there are no valid prices in the raw oracle

			//return types.ErrNoValidPrice(k.codespace)
			medianPrice = sdk.NewDec(0)
			expiry = sdk.NewInt(0)
		} else if l == 1 {

			// Return immediately if there's only one price
			medianPrice = notExpiredPrices[0].Price
			expiry = notExpiredPrices[0].Expiry
		} else {
			// sort the prices
			sort.Slice(notExpiredPrices, func(i, j int) bool {
				return notExpiredPrices[i].Price.LT(notExpiredPrices[j].Price)
			})
			// If there's an even number of prices
			if l%2 == 0 {
				// TODO make sure this is safe.
				// Since it's a price and not a blance, division with precision loss is OK.
				price1 := notExpiredPrices[l/2-1].Price
				price2 := notExpiredPrices[l/2].Price
				sum := price1.Add(price2)
				divsor, _ := sdk.NewDecFromStr("2")
				medianPrice = sum.Quo(divsor)
				// TODO Check if safe, makes sense
				// Takes the average of the two expiries rounded down to the nearest Int.
				expiry = notExpiredPrices[l/2-1].Expiry.Add(notExpiredPrices[l/2].Expiry).Quo(sdk.NewInt(2))
			} else {
				// integer division, so we'll get an integer back, rounded down
				medianPrice = notExpiredPrices[l/2].Price
				expiry = notExpiredPrices[l/2].Expiry
			}
		}

		store := ctx.KVStore(k.storeKey)
		oldPrice := k.GetCurrentPrice(ctx, assetCode)

		// Only update if there is a price or expiry change, no need to update after every block
		if !oldPrice.Price.Equal(medianPrice) || !oldPrice.Expiry.Equal(expiry) {

			currentPrice := types.CurrentPrice{
				AssetCode: assetCode,
				Price:     medianPrice,
				Expiry:    expiry,
			}

			store.Set(
				[]byte(types.CurrentPricePrefix+assetCode), k.cdc.MustMarshalBinaryBare(currentPrice),
			)
		}
	}

	return nil
}

// GetOracles returns the oracles in the oracle store
func (k Keeper) GetOracles(ctx sdk.Context) []types.Oracle {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.OraclePrefix))
	var oracles []types.Oracle
	k.cdc.MustUnmarshalBinaryBare(bz, &oracles)
	return oracles
}

// GetAssets returns the assets in the oracle store
func (k Keeper) GetAssets(ctx sdk.Context) []types.Asset {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.AssetPrefix))
	var assets []types.Asset
	k.cdc.MustUnmarshalBinaryBare(bz, &assets)
	return assets
}

// GetAsset returns the asset if it is in the oracle system
func (k Keeper) GetAsset(ctx sdk.Context, assetCode string) (types.Asset, bool) {
	assets := k.GetAssets(ctx)

	for i := range assets {
		if assets[i].AssetCode == assetCode {
			return assets[i], true
		}
	}
	return types.Asset{}, false

}

// GetOracle returns the oracle address as a string if it is in the oracle store
func (k Keeper) GetOracle(ctx sdk.Context, oracle string) (types.Oracle, bool) {
	oracles := k.GetOracles(ctx)

	for i := range oracles {
		if oracles[i].OracleAddress == oracle {
			return oracles[i], true
		}
	}
	return types.Oracle{}, false

}

// GetCurrentPrice fetches the current median price of all oracles for a specific asset
func (k Keeper) GetCurrentPrice(ctx sdk.Context, assetCode string) types.CurrentPrice {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.CurrentPricePrefix + assetCode))
	// TODO panic or return error if not found
	var price types.CurrentPrice
	k.cdc.MustUnmarshalBinaryBare(bz, &price)
	return price
}

// GetRawPrices fetches the set of all prices posted by oracles for an asset
func (k Keeper) GetRawPrices(ctx sdk.Context, assetCode string) []types.PostedPrice {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.RawPriceFeedPrefix + assetCode))
	var prices []types.PostedPrice
	k.cdc.MustUnmarshalBinaryBare(bz, &prices)
	return prices
}

// ValidatePostPrice makes sure the person posting the price is an oracle
func (k Keeper) ValidatePostPrice(ctx sdk.Context, msg types.MsgPostPrice) sdk.Error {
	// TODO implement this

	_, assetFound := k.GetAsset(ctx, msg.AssetCode)
	if !assetFound {
		return types.ErrInvalidAsset(k.codespace)
	}
	_, oracleFound := k.GetOracle(ctx, msg.From.String())
	if !oracleFound {
		return types.ErrInvalidOracle(k.codespace)
	}

	return nil
}

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/xar-network/xar-network/x/uniswap/internal/types"
)

type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramspace params.Subspace
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramspace params.Subspace) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (keeper Keeper) SetPoolConfig(ctx sdk.Context, config types.PoolConfig) {
	keeper.paramspace.SetParamSet(ctx, &config)
}

func (keeper Keeper) GetPoolConfig(ctx sdk.Context) types.PoolConfig {
	config := types.PoolConfig{}
	keeper.paramspace.GetParamSet(ctx, &config)
	return config
}

func (keeper Keeper) AddLiquidity(ctx sdk.Context, coin sdk.Coin, token sdk.Coin) sdk.Error {
	pool := types.Pool{}
	if keeper.HasPool(ctx, token.Denom) == false {
		pool = types.NewPool(sdk.NewCoin(coin.Denom, sdk.NewInt(0)), sdk.NewCoin(token.Denom, sdk.NewInt(0)))
	} else {
		var err sdk.Error
		pool, err = keeper.GetPool(ctx, token.Denom)
		if err != nil {
			return err
		}
	}

	pool.BalanceCoin = pool.BalanceCoin.Add(coin)
	pool.BalanceToken = pool.BalanceToken.Add(token)

	err := keeper.SetPool(ctx, pool)
	if err != nil {
		return err
	}

	return nil
}

func (keeper Keeper) Swap(ctx sdk.Context, asset sdk.Coin, targetDenom string) (sdk.Coin, sdk.Error) {
	if asset.Denom == targetDenom {
		return sdk.Coin{}, sdk.ErrInternal("Can't swap identical token")
	}

	config := keeper.GetPoolConfig(ctx)
	if config.CoinDenom == asset.Denom {
		coin, err := keeper.swapFromCoin(ctx, asset, targetDenom)
		if err != nil {
			return sdk.Coin{}, err
		}
		return coin, nil
	} else {
		if targetDenom == config.CoinDenom {
			coin, err := keeper.swapToCoin(ctx, asset)
			if err != nil {
				return sdk.Coin{}, err
			}
			return coin, nil
		} else {
			intermediate, err := keeper.swapToCoin(ctx, asset)
			if err != nil {
				return sdk.Coin{}, err
			}
			coin, err := keeper.swapFromCoin(ctx, intermediate, targetDenom)
			if err != nil {
				return sdk.Coin{}, err
			}
			return coin, nil
		}
	}
}

func (keeper Keeper) swapToCoin(ctx sdk.Context, asset sdk.Coin) (sdk.Coin, sdk.Error) {
	config := keeper.GetPoolConfig(ctx)
	if config.CoinDenom == asset.Denom {
		return sdk.Coin{}, sdk.ErrInternal("Can't swap coin to coin")
	}

	pool, err := keeper.GetPool(ctx, asset.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}

	pool.BalanceToken = pool.BalanceToken.Add(asset)
	ratio := sdk.NewDecFromInt(asset.Amount).Quo(sdk.NewDecFromInt(pool.BalanceToken.Amount))

	result := sdk.NewCoin(config.CoinDenom, ratio.Mul(sdk.NewDecFromInt(pool.BalanceCoin.Amount)).RoundInt())
	pool.BalanceCoin = pool.BalanceCoin.Sub(result)

	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return sdk.Coin{}, err
	}

	return result, nil
}

func (keeper Keeper) swapFromCoin(ctx sdk.Context, coin sdk.Coin, tokenDenom string) (sdk.Coin, sdk.Error) {
	config := keeper.GetPoolConfig(ctx)
	if config.CoinDenom != coin.Denom {
		return sdk.Coin{}, sdk.ErrInternal("Invalid coin denom")
	}

	pool, err := keeper.GetPool(ctx, tokenDenom)
	if err != nil {
		return sdk.Coin{}, err
	}

	pool.BalanceCoin = pool.BalanceCoin.Add(coin)
	ratio := sdk.NewDecFromInt(coin.Amount).Quo(sdk.NewDecFromInt(pool.BalanceCoin.Amount))

	result := sdk.NewCoin(tokenDenom, ratio.Mul(sdk.NewDecFromInt(pool.BalanceToken.Amount)).RoundInt())
	pool.BalanceToken = pool.BalanceToken.Sub(result)

	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return sdk.Coin{}, err
	}

	return result, err
}

func (keeper Keeper) SetPool(ctx sdk.Context, pool types.Pool) sdk.Error {
	config := keeper.GetPoolConfig(ctx)
	if len(config.CoinDenom) == 0 {
		return sdk.ErrInternal("Pool config not set")
	}

	if config.CoinDenom != pool.BalanceCoin.Denom {
		return sdk.ErrInternal(fmt.Sprintf("Not coin denom (expected: %s, actual: %s)", config.CoinDenom, pool.BalanceCoin.Denom))
	}

	if pool.BalanceCoin.Denom == pool.BalanceToken.Denom {
		return sdk.ErrInternal("Coin's denom and Token's denom should not be equal")
	}

	if pool.BalanceCoin.Amount.LTE(sdk.NewInt(0)) || pool.BalanceToken.Amount.LTE(sdk.NewInt(0)) {
		return sdk.ErrInternal(fmt.Sprintf("Pool has 0 coin or token (coin: %s, token: %s)", pool.BalanceCoin.String(), pool.BalanceToken.String()))
	}

	store := ctx.KVStore(keeper.storeKey)
	key := []byte(config.CoinDenom + "-" + pool.BalanceToken.Denom)

	bz, gerr := keeper.cdc.MarshalBinaryBare(pool)
	if gerr != nil {
		return sdk.ErrInternal(gerr.Error())
	}

	store.Set(key, bz)

	return nil
}

func (keeper Keeper) GetPool(ctx sdk.Context, tokenDenom string) (types.Pool, sdk.Error) {
	config := keeper.GetPoolConfig(ctx)
	if len(config.CoinDenom) == 0 {
		return types.Pool{}, sdk.ErrInternal("Pool config not set")
	}

	store := ctx.KVStore(keeper.storeKey)
	key := []byte(config.CoinDenom + "-" + tokenDenom)

	bz := store.Get(key)
	if len(bz) == 0 {
		return types.Pool{}, sdk.ErrInternal("Unkown token pool")
	}

	pool := types.Pool{}
	gerr := keeper.cdc.UnmarshalBinaryBare(bz, &pool)
	if gerr != nil {
		return types.Pool{}, sdk.ErrInternal(gerr.Error())
	}

	return pool, nil
}

func (keeper Keeper) HasPool(ctx sdk.Context, tokenDenom string) bool {
	config := keeper.GetPoolConfig(ctx)
	if len(config.CoinDenom) == 0 {
		panic(fmt.Errorf("Pool config not set"))
	}

	store := ctx.KVStore(keeper.storeKey)
	key := []byte(config.CoinDenom + "-" + tokenDenom)

	bz := store.Get(key)
	if len(bz) == 0 {
		return false
	}

	return true
}

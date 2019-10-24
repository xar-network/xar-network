package uniswap

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
)

func TestSwapFromCoin(t *testing.T) {
	chain := setupChain()

	chain.sk.SetPoolConfig(chain.ctx, PoolConfig{
		CoinDenom: "uatom",
	})
	config := chain.sk.GetPoolConfig(chain.ctx)
	require.Equal(t, "uatom", config.CoinDenom)

	err := chain.sk.AddLiquidity(chain.ctx, sdk.NewCoin("uatom", sdk.NewInt(1000)), sdk.NewCoin("udai", sdk.NewInt(1000)))
	require.Nil(t, err)

	pool, err := chain.sk.GetPool(chain.ctx, "udai")
	require.Nil(t, err)
	require.Equal(t, "1000", pool.BalanceCoin.Amount.String())
	require.Equal(t, "1000", pool.BalanceToken.Amount.String())

	asset, _, err := chain.sk.Swap(chain.ctx, sdk.NewCoin("uatom", sdk.NewInt(500)), "udai")
	require.Nil(t, err)
	require.Equal(t, "udai", asset.Denom)
	require.Equal(t, sdk.NewInt(333).String(), asset.Amount.String())

	pool, err = chain.sk.GetPool(chain.ctx, "udai")
	require.Nil(t, err)
	require.Equal(t, "1500", pool.BalanceCoin.Amount.String())
	require.Equal(t, "667", pool.BalanceToken.Amount.String())
}

func TestSwapToCoin(t *testing.T) {
	chain := setupChain()

	chain.sk.SetPoolConfig(chain.ctx, PoolConfig{
		CoinDenom: "uatom",
	})
	config := chain.sk.GetPoolConfig(chain.ctx)
	require.Equal(t, "uatom", config.CoinDenom)

	err := chain.sk.AddLiquidity(chain.ctx, sdk.NewCoin("uatom", sdk.NewInt(1000)), sdk.NewCoin("udai", sdk.NewInt(1000)))
	require.Nil(t, err)

	pool, err := chain.sk.GetPool(chain.ctx, "udai")
	require.Nil(t, err)
	require.Equal(t, "1000", pool.BalanceCoin.Amount.String())
	require.Equal(t, "1000", pool.BalanceToken.Amount.String())

	asset, _, err := chain.sk.Swap(chain.ctx, sdk.NewCoin("udai", sdk.NewInt(500)), "uatom")
	require.Nil(t, err)
	require.Equal(t, "uatom", asset.Denom)
	require.Equal(t, sdk.NewInt(333).String(), asset.Amount.String())

	pool, err = chain.sk.GetPool(chain.ctx, "udai")
	require.Nil(t, err)
	require.Equal(t, "667", pool.BalanceCoin.Amount.String())
	require.Equal(t, "1500", pool.BalanceToken.Amount.String())
}

func TestSwapFromTokenToToken(t *testing.T) {
	chain := setupChain()

	chain.sk.SetPoolConfig(chain.ctx, PoolConfig{
		CoinDenom: "uatom",
	})
	config := chain.sk.GetPoolConfig(chain.ctx)
	require.Equal(t, "uatom", config.CoinDenom)

	err := chain.sk.AddLiquidity(chain.ctx, sdk.NewCoin("uatom", sdk.NewInt(1000)), sdk.NewCoin("udai", sdk.NewInt(1000)))
	require.Nil(t, err)

	pool, err := chain.sk.GetPool(chain.ctx, "udai")
	require.Nil(t, err)
	require.Equal(t, "1000", pool.BalanceCoin.Amount.String())
	require.Equal(t, "1000", pool.BalanceToken.Amount.String())

	err = chain.sk.AddLiquidity(chain.ctx, sdk.NewCoin("uatom", sdk.NewInt(1000)), sdk.NewCoin("ubatom", sdk.NewInt(1000)))
	require.Nil(t, err)

	pool, err = chain.sk.GetPool(chain.ctx, "ubatom")
	require.Nil(t, err)
	require.Equal(t, "1000", pool.BalanceCoin.Amount.String())
	require.Equal(t, "1000", pool.BalanceToken.Amount.String())

	asset, _, err := chain.sk.Swap(chain.ctx, sdk.NewCoin("ubatom", sdk.NewInt(500)), "udai")
	require.Nil(t, err)
	require.Equal(t, "udai", asset.Denom)
	require.Equal(t, sdk.NewInt(250).String(), asset.Amount.String())

	pool, err = chain.sk.GetPool(chain.ctx, "udai")
	require.Nil(t, err)
	require.Equal(t, "1333", pool.BalanceCoin.Amount.String())
	require.Equal(t, "750", pool.BalanceToken.Amount.String())
}

type testChain struct {
	cdc *codec.Codec
	ctx sdk.Context
	ak  auth.AccountKeeper
	bk  bank.Keeper
	pk  params.Keeper
	sk  Keeper
}

func setupChain() *testChain {
	db := dbm.NewMemDB()

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)

	accountKey := sdk.NewKVStoreKey("account")
	swapKey := sdk.NewKVStoreKey("SwapKey")
	keyParams := sdk.NewKVStoreKey("subspace")
	tkeyParams := sdk.NewTransientStoreKey("transient_subspace")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(swapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(
		cdc,
		accountKey,
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	bk := bank.NewBaseKeeper(
		ak,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	sk := NewKeeper(cdc, swapKey, pk.Subspace(DefaultParamspace), bk)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test"}, false, log.NewNopLogger())

	return &testChain{cdc: cdc, ctx: ctx, ak: ak, bk: bk, pk: pk, sk: sk}
}

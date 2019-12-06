/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/x/denominations/internal/keeper"
	"github.com/xar-network/xar-network/x/denominations/internal/types"

	cstore "github.com/cosmos/cosmos-sdk/store"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func TestKeeperCoverage(t *testing.T) {

	cdc := MakeTestCodec()

	logger := log.NewNopLogger()

	var (
		keyParams  = sdk.NewKVStoreKey(params.StoreKey)
		keyAcc     = sdk.NewKVStoreKey(auth.StoreKey)
		keySupply  = sdk.NewKVStoreKey(supply.StoreKey)
		keyDenom   = sdk.NewKVStoreKey(types.StoreKey)
		tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	)

	db := dbm.NewMemDB()
	ms := cstore.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDenom, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "xar-chain"}, true, logger)

	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter, supply.Burner},
	}
	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, make(map[string]bool))
	sk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)

	addr := sdk.AccAddress([]byte("someName"))
	acc := ak.NewAccountWithAddress(ctx, addr)
	ak.SetAccount(ctx, acc)

	naddr := sdk.AccAddress([]byte("nominee"))
	nacc := ak.NewAccountWithAddress(ctx, naddr)
	ak.SetAccount(ctx, nacc)

	addrerr := sdk.AccAddress([]byte("error"))
	acc = ak.NewAccountWithAddress(ctx, addrerr)
	ak.SetAccount(ctx, acc)

	dk := keeper.NewKeeper(keyDenom, cdc, ak, sk, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace)
	params := types.NewParams([]string{"cosmos1dehk66twv4js5dq8xr"})
	dk.SetParams(ctx, params)

	sk.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

	res := dk.BurnCoins(ctx, addrerr, sdk.NewInt(100), "uftm")
	require.Equal(t, false, res.IsOK())

	token := types.NewToken(
		"Max Mintable", "max",
		"MAX",
		sdk.NewInt(100),
		addr,
		true,
	)

	// Issue a new token
	res = dk.IssueToken(ctx, naddr, addr, *token)
	require.Equal(t, true, res.IsOK())

	token = types.NewToken(
		"Fantom", "uftm",
		"FTM",
		sdk.NewInt(3175000000000000),
		addr,
		true,
	)

	// Issue a new token
	res = dk.IssueToken(ctx, naddr, addr, *token)
	require.Equal(t, true, res.IsOK())

	// Try to issue again and fail
	res = dk.IssueToken(ctx, naddr, addr, *token)
	require.Equal(t, false, res.IsOK())

	// Issue a new token with total supply exceeding max supply
	token = types.NewToken(
		"Over max", "ovr",
		"ovr",
		sdk.NewInt(3175000000000000),
		addr,
		true,
	)
	res = dk.IssueToken(ctx, naddr, addr, *token)
	require.Equal(t, true, res.IsOK())

	// Issue a new token with total supply exceeding max supply
	token = types.NewToken(
		"Unmintable", "unm",
		"UNM",
		sdk.NewInt(3175000000000000),
		addrerr,
		false,
	)
	res = dk.IssueToken(ctx, naddr, addrerr, *token)
	require.Equal(t, true, res.IsOK())

	// Try to mint unmintable coin
	res = dk.MintCoins(ctx, addr, sdk.NewInt(1), "unm")
	require.Equal(t, false, res.IsOK())

	// Try to mint over max supply
	res = dk.MintCoins(ctx, addr, sdk.NewInt(101), "max")
	require.Equal(t, false, res.IsOK())

	// Try a normal mint
	res = dk.MintCoins(ctx, addr, sdk.NewInt(2), "uftm")
	require.Equal(t, true, res.IsOK())

	// Burn over max supply
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("BurnCoins should have panicked!")
			}
		}()
		// This function should cause a panic
		dk.BurnCoins(ctx, addr, sdk.NewInt(3), "uftm")
	}()

	// Try a normal burn
	res = dk.BurnCoins(ctx, addr, sdk.NewInt(1), "uftm")
	require.Equal(t, true, res.IsOK())

	// Freeze coins the address doesn't have
	res = dk.FreezeCoins(ctx, addr, addrerr, sdk.NewInt(1), "uftm")
	require.Equal(t, false, res.IsOK())

	// Freeze coins the address has
	res = dk.FreezeCoins(ctx, addr, addr, sdk.NewInt(1), "uftm")
	require.Equal(t, true, res.IsOK())

	// Unfreeze coins the address has
	res = dk.UnfreezeCoins(ctx, addr, addr, sdk.NewInt(1), "uftm")
	require.Equal(t, true, res.IsOK())
	acc = ak.GetAccount(ctx, addrerr)

	t.Logf("%s", sk.GetSupply(ctx).String())
}

func MakeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()

	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)

	return
}

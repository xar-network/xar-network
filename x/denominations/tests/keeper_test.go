package tests

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/x/denominations"
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
		tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	)

	db := dbm.NewMemDB()
	ms := cstore.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
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

	addrerr := sdk.AccAddress([]byte("error"))
	acc = ak.NewAccountWithAddress(ctx, addrerr)
	ak.SetAccount(ctx, acc)

	dk := denominations.NewKeeper(ak, sk, pk.Subspace(denominations.DefaultParamspace), denominations.DefaultCodespace)

	sk.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	dk.SetParams(ctx, types.NewParams("cosmos1wdhk6e2wv9kk2j88d92"))

	msg := types.NewMsgMint(addr, sdk.NewCoins(sdk.NewCoin("unew", sdk.NewInt(100))))
	msgb := types.NewMsgBurn(addr, sdk.NewCoins(sdk.NewCoin("unew", sdk.NewInt(100))))
	res := dk.Mint(ctx, msg)
	t.Errorf("%s", sk.GetSupply(ctx))
	require.Equal(t, true, res.IsOK())
	res = dk.Burn(ctx, msgb)
	t.Errorf("%s", sk.GetSupply(ctx))
	require.Equal(t, true, res.IsOK())

	msgbe := types.NewMsgBurn(addrerr, sdk.NewCoins(sdk.NewCoin("unew", sdk.NewInt(100))))
	res = dk.Burn(ctx, msgbe)
	t.Errorf("%s", sk.GetSupply(ctx))
	require.Equal(t, false, res.IsOK())
}

func MakeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()

	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)

	return
}

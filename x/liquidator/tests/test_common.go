package liquidator

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/xar-network/xar-network/x/auction"
	"github.com/xar-network/xar-network/x/csdt"
	"github.com/xar-network/xar-network/x/oracle"
)

// Avoid cluttering test cases with long function name
func i(in int64) sdk.Int                    { return sdk.NewInt(in) }
func c(denom string, amount int64) sdk.Coin { return sdk.NewInt64Coin(denom, amount) }
func cs(coins ...sdk.Coin) sdk.Coins        { return sdk.NewCoins(coins...) }

type keepers struct {
	paramsKeeper     params.Keeper
	accountKeeper    auth.AccountKeeper
	bankKeeper       bank.Keeper
	oracleKeeper  oracle.Keeper
	auctionKeeper    auction.Keeper
	csdtKeeper        csdt.Keeper
	liquidatorKeeper Keeper
}

func setupTestKeepers() (sdk.Context, keepers) {

	// Setup in memory database
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyPriceFeed := sdk.NewKVStoreKey(oracle.StoreKey)
	keyCSDT := sdk.NewKVStoreKey("csdt")
	keyAuction := sdk.NewKVStoreKey("auction")
	keyLiquidator := sdk.NewKVStoreKey("liquidator")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyPriceFeed, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCSDT, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAuction, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyLiquidator, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	// Create Codec
	cdc := makeTestCodec()

	// Create Keepers
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(
		cdc,
		keyAcc,
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	oracleKeeper := oracle.NewKeeper(keyPriceFeed, cdc, oracle.DefaultCodespace)
	csdtKeeper := csdt.NewKeeper(
		cdc,
		keyCSDT,
		paramsKeeper.Subspace("csdtSubspace"),
		oracleKeeper,
		bankKeeper,
	)
	auctionKeeper := auction.NewKeeper(cdc, csdtKeeper, keyAuction) // Note: csdt keeper stands in for bank keeper
	liquidatorKeeper := NewKeeper(
		cdc,
		keyLiquidator,
		paramsKeeper.Subspace("liquidatorSubspace"),
		csdtKeeper,
		auctionKeeper,
		csdtKeeper,
	) // Note: csdt keeper stands in for bank keeper

	// Create context
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "testchain"}, false, log.NewNopLogger())

	return ctx, keepers{
		paramsKeeper,
		accountKeeper,
		bankKeeper,
		oracleKeeper,
		auctionKeeper,
		csdtKeeper,
		liquidatorKeeper,
	}
}

func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	oracle.RegisterCodec(cdc)
	auction.RegisterCodec(cdc)
	csdt.RegisterCodec(cdc)
	RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

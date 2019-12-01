package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/xar-network/xar-network/x/csdt"

	"github.com/xar-network/xar-network/x/auction"
	"github.com/xar-network/xar-network/x/liquidator/internal/types"
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
	oracleKeeper     oracle.Keeper
	auctionKeeper    auction.Keeper
	csdtKeeper       csdt.Keeper
	liquidatorKeeper Keeper
}

func setupTestKeepers() (sdk.Context, keepers) {

	// Setup in memory database
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyPriceFeed := sdk.NewKVStoreKey(oracle.StoreKey)
	keyCSDT := sdk.NewKVStoreKey(csdt.StoreKey)
	keyAuction := sdk.NewKVStoreKey(auction.StoreKey)
	keyLiquidator := sdk.NewKVStoreKey(types.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyPriceFeed, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCSDT, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAuction, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyLiquidator, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
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
	blacklistedAddrs := make(map[string]bool)
	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		csdt.ModuleName: {supply.Minter, supply.Burner},
	}

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	oracleKeeper := oracle.NewKeeper(keyPriceFeed, cdc, paramsKeeper.Subspace(oracle.DefaultParamspace), oracle.DefaultCodespace)
	auctionKeeper := auction.NewKeeper(cdc, bankKeeper, keyAuction, paramsKeeper.Subspace(auction.DefaultParamspace)) // Note: csdt keeper stands in for bank keeper
	csdtKeeper := csdt.NewKeeper(
		cdc,
		keyCSDT,
		paramsKeeper.Subspace(csdt.DefaultParamspace),
		oracleKeeper,
		bankKeeper,
		supplyKeeper,
	)
	liquidatorKeeper := NewKeeper(
		cdc,
		keyLiquidator,
		paramsKeeper.Subspace(types.DefaultParamspace),
		csdtKeeper,
		auctionKeeper,
		csdtKeeper,
	) // Note: csdt keeper stands in for bank keeper

	// Create context
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "testchain"}, false, log.NewNopLogger())
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

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
	types.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)
	return cdc
}

func defaultParams() types.LiquidatorParams {
	return types.LiquidatorParams{
		DebtAuctionSize: sdk.NewInt(1000),
		CollateralParams: []types.CollateralParams{
			{
				Denom:       "btc",
				AuctionSize: sdk.NewInt(1),
			},
		},
	}
}

func csdtDefaultGenesis() csdt.GenesisState {
	return csdt.GenesisState{
		csdt.Params{
			GlobalDebtLimit: sdk.NewCoins(sdk.NewCoin(csdt.StableDenom, sdk.NewInt(1000000))),
			CollateralParams: csdt.CollateralParams{
				{
					Denom:            "btc",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(csdt.StableDenom, sdk.NewInt(500000))),
				},
			},
		},
		sdk.ZeroInt(),
		csdt.CSDTs{},
	}
}

func oracleGenesis(address string) oracle.GenesisState {
	ap := oracle.Params{
		Assets: []oracle.Asset{
			oracle.Asset{AssetCode: "btc", BaseAsset: "btc", QuoteAsset: "usd"},
		},
		Nominees: []string{address},
	}
	return oracle.GenesisState{
		Params: ap,
		PostedPrices: []oracle.PostedPrice{
			oracle.PostedPrice{
				AssetCode:     "btc",
				OracleAddress: sdk.AccAddress([]byte("someName")),
				Price:         sdk.MustNewDecFromStr("8000.00"),
				Expiry:        time.Now().Add(time.Hour * 1),
			},
		},
	}
}

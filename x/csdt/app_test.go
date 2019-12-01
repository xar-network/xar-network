package csdt

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
	"github.com/xar-network/xar-network/x/oracle"
)

func TestApp_CreateModifyDeleteCSDT(t *testing.T) {
	// Setup
	mapp, keeper := setUpMockAppWithoutGenesis()
	genAccs, addrs, _, privKeys := mock.CreateGenAccounts(1, cs(c("uftm", 100)))
	testAddr := addrs[0]
	testPrivKey := privKeys[0]
	mock.SetGenesis(mapp, genAccs)
	// setup oracle, TODO can this be shortened a bit?
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	keeper.SetParams(ctx, types.DefaultParams())
	keeper.SetGlobalDebt(ctx, sdk.NewInt(1000000000))
	keeper.GetSupply().SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	oracleParams := oracle.DefaultParams()
	oracleParams.Assets = oracle.Assets{
		oracle.Asset{
			AssetCode:  "uftm",
			BaseAsset:  "uftm",
			QuoteAsset: StableDenom,
			Oracles: oracle.Oracles{
				oracle.Oracle{
					Address: addrs[0],
				},
			},
		},
	}
	oracleParams.Nominees = []string{addrs[0].String()}

	keeper.GetOracle().SetParams(ctx, oracleParams)
	_, _ = keeper.GetOracle().SetPrice(
		ctx, addrs[0], "uftm",
		sdk.MustNewDecFromStr("1.00"),
		time.Now().Add(time.Hour*1))
	_ = keeper.GetOracle().SetCurrentPrices(ctx)
	mapp.EndBlock(abci.RequestEndBlock{})
	mapp.Commit()

	// Create CSDT
	msgs := []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(10), i(5))}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{0}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(types.StableDenom, 5), c("uftm", 90)))

	// Modify CSDT
	msgs = []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(40), i(5))}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{1}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(types.StableDenom, 10), c("uftm", 50)))

	// Delete CSDT
	msgs = []sdk.Msg{types.NewMsgCreateOrModifyCSDT(testAddr, "uftm", i(-50), i(-10))}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{2}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c("uftm", 100)))
}

// Avoid cluttering test cases with long function name
func i(in int64) sdk.Int                    { return sdk.NewInt(in) }
func d(str string) sdk.Dec                  { return sdk.MustNewDecFromStr(str) }
func c(denom string, amount int64) sdk.Coin { return sdk.NewInt64Coin(denom, amount) }
func cs(coins ...sdk.Coin) sdk.Coins        { return sdk.NewCoins(coins...) }

func setUpMockAppWithoutGenesis() (*mock.App, Keeper) {
	// Create uninitialized mock app
	mapp := mock.NewApp()

	// Register codecs
	types.RegisterCodec(mapp.Cdc)
	supply.RegisterCodec(mapp.Cdc)

	// Create keepers
	keyCSDT := sdk.NewKVStoreKey(types.StoreKey)
	keyOracle := sdk.NewKVStoreKey(oracle.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter, supply.Burner},
	}

	oracleKeeper := oracle.NewKeeper(keyOracle, mapp.Cdc, mapp.ParamsKeeper.Subspace(oracle.DefaultParamspace), oracle.DefaultCodespace)
	bankKeeper := bank.NewBaseKeeper(mapp.AccountKeeper, mapp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})
	supplyKeeper := supply.NewKeeper(mapp.Cdc, keySupply, mapp.AccountKeeper, bankKeeper, maccPerms)
	csdtKeeper := NewKeeper(mapp.Cdc, keyCSDT, mapp.ParamsKeeper.Subspace(types.DefaultParamspace), oracleKeeper, bankKeeper, supplyKeeper)

	// Register routes
	mapp.Router().AddRoute("csdt", NewHandler(csdtKeeper))
	// Mount and load the stores
	err := mapp.CompleteSetup(keyOracle, keyCSDT, keySupply)
	if err != nil {
		panic("mock app setup failed")
	}

	return mapp, csdtKeeper
}

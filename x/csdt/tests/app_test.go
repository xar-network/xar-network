package csdt

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
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
	keeper.GetOracle().AddAsset(ctx, "uftm", "uftm test")
	_, _ = keeper.GetOracle().SetPrice(
		ctx, sdk.AccAddress{}, "uftm",
		sdk.MustNewDecFromStr("1.00"),
		sdk.NewInt(10))
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

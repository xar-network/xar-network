package issue_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/staking"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/xar-network/xar-network/x/issue"
	"github.com/xar-network/xar-network/x/issue/internal/keeper"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

var (
	ReceiverCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("receiverCoins")))
	TransferAccAddr      sdk.AccAddress
	SenderAccAddr        sdk.AccAddress

	IssueParams = types.IssueParams{
		Name:               "testCoin",
		Symbol:             "TEST",
		TotalSupply:        sdk.NewInt(10000),
		BurnOwnerDisabled:  false,
		BurnHolderDisabled: false,
		BurnFromDisabled:   false,
		MintingFinished:    false}

	CoinIssueInfo = types.CoinIssueInfo{
		Owner:              SenderAccAddr,
		Issuer:             SenderAccAddr,
		Name:               "testCoin",
		Symbol:             "TEST",
		TotalSupply:        sdk.NewInt(10000),
		BurnOwnerDisabled:  false,
		BurnHolderDisabled: false,
		BurnFromDisabled:   false,
		MintingFinished:    false}
)

// Wrapper struct
type MockHooks struct {
	keeper issue.Keeper
}

// Create new issue hooks
func NewMockHooks(ik issue.Keeper) MockHooks { return MockHooks{ik} }

func (hooks MockHooks) CanSend(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (bool, sdk.Error) {
	return true, nil
}

func (hooks MockHooks) CheckMustMemoAddress(ctx sdk.Context, toAddr sdk.AccAddress, memo string) (bool, sdk.Error) {
	return false, nil
}

// initialize the mock application for this module
func getMockApp(t *testing.T, genState issue.GenesisState, genAccs []exported.Account) (
	mapp *mock.App, k keeper.Keeper, sk staking.Keeper, addrs []sdk.AccAddress,
	pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {
	mapp = mock.NewApp()
	types.RegisterCodec(mapp.Cdc)
	keyIssue := sdk.NewKVStoreKey(types.StoreKey)

	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	pk := mapp.ParamsKeeper
	ck := bank.NewBaseKeeper(mapp.AccountKeeper, mapp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, make(map[string]bool))

	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter, supply.Burner},
	}

	supk := supply.NewKeeper(mapp.Cdc, keySupply, mapp.AccountKeeper, ck, maccPerms)
	ik := issue.NewKeeper(keyIssue, pk.Subspace("testissue"), ck, sk, types.DefaultCodespace, auth.FeeCollectorName)

	mapp.Router().AddRoute(types.RouterKey, issue.NewHandler(ik))
	mapp.QueryRouter().AddRoute(types.QuerierRoute, keeper.NewQuerier(ik))
	//mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, ik, sk, genState))

	require.NoError(t, mapp.CompleteSetup(keyIssue))

	valTokens := sdk.TokensFromConsensusPower(1000000000000)
	if len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(2,
			sdk.NewCoins(sdk.NewCoin("uftm", valTokens)))
	}
	SenderAccAddr = genAccs[0].GetAddress()
	TransferAccAddr = genAccs[1].GetAddress()

	CoinIssueInfo.Owner = SenderAccAddr
	CoinIssueInfo.Issuer = SenderAccAddr

	mock.SetGenesis(mapp, genAccs)

	return mapp, ik, sk, addrs, pubKeys, privKeys
}
func getInitChainer(mapp *mock.App, keeper keeper.Keeper, stakingKeeper staking.Keeper, genState issue.GenesisState) sdk.InitChainer {

	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {

		mapp.InitChainer(ctx, req)

		//validators, err := staking.InitGenesis(ctx, stakingKeeper, stakingGenesis)
		//if err != nil {
		//	panic(err)
		//}
		if genState.IsEmpty() {
			issue.InitGenesis(ctx, keeper, issue.DefaultGenesisState())
		} else {
			issue.InitGenesis(ctx, keeper, genState)
		}
		return abci.ResponseInitChain{
			//Validators: validators,
		}
	}
}

package tests

import (
	"testing"
	"time"

	"github.com/hashgard/hashgard/x/box/msgs"

	keeper2 "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	"github.com/hashgard/hashgard/x/box/utils"

	"github.com/hashgard/hashgard/x/box/params"

	"github.com/hashgard/hashgard/x/box/types"

	"github.com/cosmos/cosmos-sdk/x/staking"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"

	"github.com/hashgard/hashgard/x/box"
	"github.com/hashgard/hashgard/x/box/keeper"
	issueutils "github.com/hashgard/hashgard/x/issue/utils"
)

var (
	Receiver          = "Receiver"
	TransferAccAddr   = sdk.AccAddress(crypto.AddressHash([]byte("transferAddress")))
	SenderAccAddr     sdk.AccAddress
	TestTokenDecimals = uint(18)

	newBoxInfo = types.BoxInfo{
		Name:             "testBox",
		BoxType:          types.Lock,
		Description:      "{}",
		TransferDisabled: true,
		TotalAmount: types.BoxToken{
			Token: sdk.NewCoin(
				"text",
				issueutils.MulDecimals(sdk.NewInt(10000), TestTokenDecimals)),
			Decimals: TestTokenDecimals},
	}
)

func GetLockBoxInfo() *params.BoxLockParams {
	box := &params.BoxLockParams{}

	box.Name = newBoxInfo.Name
	box.TotalAmount = newBoxInfo.TotalAmount
	box.Lock = types.LockBox{EndTime: time.Now().Add(time.Duration(5) * time.Second).Unix()}
	return box
}
func GetDepositBoxInfo() *params.BoxDepositParams {
	box := &params.BoxDepositParams{}

	box.Name = newBoxInfo.Name
	box.TotalAmount = newBoxInfo.TotalAmount
	box.Deposit = types.DepositBox{
		StartTime:     time.Now().Add(time.Duration(10) * time.Second).Unix(),
		EstablishTime: time.Now().Add(time.Duration(20) * time.Second).Unix(),
		MaturityTime:  time.Now().Add(time.Duration(30) * time.Second).Unix(),
		BottomLine:    issueutils.MulDecimals(sdk.NewInt(200), TestTokenDecimals),
		Price:         issueutils.MulDecimals(sdk.NewInt(100), TestTokenDecimals),
		Interest: types.BoxToken{
			Token: sdk.NewCoin(
				"interest",
				issueutils.MulDecimals(sdk.NewInt(1000), TestTokenDecimals)),
			Decimals: TestTokenDecimals}}
	box.Deposit.PerCoupon = utils.CalcInterestRate(box.TotalAmount.Token.Amount, box.Deposit.Price,
		box.Deposit.Interest.Token, box.Deposit.Interest.Decimals)
	return box
}
func GetFutureBoxInfo() *params.BoxFutureParams {
	box := &params.BoxFutureParams{}
	box.Name = newBoxInfo.Name
	box.TotalAmount = newBoxInfo.TotalAmount
	box.TotalAmount.Token.Amount = issueutils.MulDecimals(sdk.NewInt(2000), TestTokenDecimals)
	box.Future.TimeLine = []int64{
		time.Now().Add(time.Duration(20) * time.Second).Unix(),
		time.Now().Add(time.Duration(21) * time.Second).Unix(),
		time.Now().Add(time.Duration(22) * time.Second).Unix()}
	box.Future.Receivers = [][]string{
		{sdk.AccAddress(crypto.AddressHash([]byte(Receiver + "1"))).String(),
			issueutils.MulDecimals(sdk.NewInt(100), TestTokenDecimals).String(),
			issueutils.MulDecimals(sdk.NewInt(200), TestTokenDecimals).String(),
			issueutils.MulDecimals(sdk.NewInt(300), TestTokenDecimals).String()},

		{sdk.AccAddress(crypto.AddressHash([]byte(Receiver + "2"))).String(),
			issueutils.MulDecimals(sdk.NewInt(200), TestTokenDecimals).String(),
			issueutils.MulDecimals(sdk.NewInt(300), TestTokenDecimals).String(),
			issueutils.MulDecimals(sdk.NewInt(200), TestTokenDecimals).String()},

		{sdk.AccAddress(crypto.AddressHash([]byte(Receiver + "3"))).String(),
			issueutils.MulDecimals(sdk.NewInt(100), TestTokenDecimals).String(),
			issueutils.MulDecimals(sdk.NewInt(400), TestTokenDecimals).String(),
			issueutils.MulDecimals(sdk.NewInt(200), TestTokenDecimals).String()}}
	return box
}

// gov and staking endblocker
func getEndBlocker(keeper keeper.Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		tags := box.EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{
			Tags: tags,
		}
	}
}

// initialize the mock application for this module
func getMockApp(t *testing.T, genState box.GenesisState, genAccs []auth.Account) (
	mapp *mock.App, keeper keeper.Keeper, sk staking.Keeper, addrs []sdk.AccAddress,
	pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {
	mapp = mock.NewApp()
	msgs.RegisterCodec(mapp.Cdc)
	keyBox := sdk.NewKVStoreKey(types.StoreKey)
	//keyIssue := sdk.NewKVStoreKey(issue.StoreKey)

	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)

	pk := mapp.ParamsKeeper
	ck := bank.NewBaseKeeper(mapp.AccountKeeper, mapp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)
	//ik := issue.NewKeeper(mapp.Cdc, keyIssue, pk, pk.Subspace("testIssue"), ck, issue.DefaultCodespace)

	ik := NewIssueKeeper()
	fck := keeper2.DummyFeeCollectionKeeper{}

	keeper = box.NewKeeper(mapp.Cdc, keyBox, pk, pk.Subspace("testBox"), &ck, ik, fck, types.DefaultCodespace)
	sk = staking.NewKeeper(mapp.Cdc, keyStaking, tkeyStaking, ck, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)

	ck.SetHooks(NewMockHooks(keeper))

	mapp.Router().AddRoute(types.RouterKey, box.NewHandler(keeper))
	mapp.QueryRouter().AddRoute(types.QuerierRoute, box.NewQuerier(keeper))
	mapp.SetEndBlocker(getEndBlocker(keeper))
	mapp.SetInitChainer(getInitChainer(mapp, keeper, sk, genState))

	require.NoError(t, mapp.CompleteSetup(keyBox))

	valTokens := sdk.TokensFromTendermintPower(10000000)
	if len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(1,
			sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens)))
	}

	SenderAccAddr = genAccs[0].GetAddress()
	mock.SetGenesis(mapp, genAccs)

	return mapp, keeper, sk, addrs, pubKeys, privKeys
}
func getInitChainer(mapp *mock.App, keeper keeper.Keeper, stakingKeeper staking.Keeper, genState box.GenesisState) sdk.InitChainer {

	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {

		mapp.InitChainer(ctx, req)

		stakingGenesis := staking.DefaultGenesisState()
		tokens := sdk.TokensFromTendermintPower(100000)
		stakingGenesis.Pool.NotBondedTokens = tokens

		//validators, err := staking.InitGenesis(ctx, stakingKeeper, stakingGenesis)
		//if err != nil {
		//	panic(err)
		//}
		if genState.IsEmpty() {
			box.InitGenesis(ctx, keeper, box.DefaultGenesisState())
		} else {
			box.InitGenesis(ctx, keeper, genState)
		}
		return abci.ResponseInitChain{
			//Validators: validators,
		}
	}
}

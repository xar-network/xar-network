// nolint:deadcode unused
package keeper

import (
	"os"
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/stretchr/testify/require"

	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/x/interest/internal/types"
)

type testInput struct {
	ctx        sdk.Context
	cdc        *codec.Codec
	mintKeeper Keeper
}

func newTestInput(t *testing.T) testInput {
	db := dbm.NewMemDB()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyInterest := sdk.NewKVStoreKey(types.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyInterest, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))

	paramsKeeper := params.NewKeeper(types.ModuleCdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(types.ModuleCdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, make(map[string]bool))
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          {supply.Minter},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(types.ModuleCdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

	inflationKeeper := NewKeeper(types.ModuleCdc, keyInterest, paramsKeeper.Subspace(types.DefaultParamspace), supplyKeeper, auth.FeeCollectorName)

	// set module accounts
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	minterAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Minter)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner)

	inflationKeeper.SetState(ctx, types.NewInterestState("x0jpy", "0.05", "x2chf", "0.10", "x2eur", "0.01"))

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, minterAcc)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	supplyKeeper.SetModuleAccount(ctx, bondPool)

	return testInput{ctx, types.ModuleCdc, inflationKeeper}
}
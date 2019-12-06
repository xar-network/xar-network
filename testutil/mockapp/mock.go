/*

Copyright 2019 All in Bits, Inc
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

package mockapp

import (
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/xar-network/xar-network/app"
	"github.com/xar-network/xar-network/execution"
	"github.com/xar-network/xar-network/types"
	"github.com/xar-network/xar-network/x/market"
	"github.com/xar-network/xar-network/x/order"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type nopWriter struct{}

func (w nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type MockApp struct {
	Cdc             *codec.Codec
	Mq              types.Backend
	Ctx             sdk.Context
	SupplyKeeper    supply.Keeper
	MarketKeeper    market.Keeper
	OrderKeeper     order.Keeper
	BankKeeper      bank.Keeper
	ExecutionKeeper execution.Keeper
}

type Option func(t *testing.T, app *MockApp)

func New(t *testing.T, options ...Option) *MockApp {
	appDB := dbm.NewMemDB()
	mkDataDB := dbm.NewMemDB()
	dex := app.NewXarApp(log.NewNopLogger(), appDB, mkDataDB, nil, true, 0)

	genesisState := app.ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(dex.Codec(), genesisState)
	if err != nil {
		return nil
	}
	dex.InitChain(abci.RequestInitChain{
		AppStateBytes: stateBytes,
	})
	ctx := dex.BaseApp.NewContext(false, abci.Header{ChainID: "unit-test-chain", Height: 1, Time: time.Unix(1558332092, 0)})

	mock := &MockApp{
		Cdc:             dex.Codec(),
		Mq:              dex.MQ(),
		Ctx:             ctx,
		SupplyKeeper:    dex.SupplyKeeper(),
		MarketKeeper:    dex.MarketKeeper(),
		OrderKeeper:     dex.OrderKeeper(),
		BankKeeper:      dex.BankKeeper(),
		ExecutionKeeper: dex.ExecKeeper(),
	}

	for _, opt := range options {
		opt(t, mock)
	}

	return mock
}

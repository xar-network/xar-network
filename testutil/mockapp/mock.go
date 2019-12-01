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
)

type nopWriter struct{}

func (w nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type MockApp struct {
	Cdc             *codec.Codec
	Mq              types.Backend
	Ctx             sdk.Context
	MarketKeeper    market.Keeper
	OrderKeeper     order.Keeper
	BankKeeper      bank.Keeper
	ExecutionKeeper execution.Keeper
}

type Option func(t *testing.T, app *MockApp)

func New(t *testing.T, options ...Option) *MockApp {
	appDB := dbm.NewMemDB()
	mkDataDB := dbm.NewMemDB()
	dex := app.NewXarApp(log.NewNopLogger(), appDB, mkDataDB, &nopWriter{})
	dex.InitChain(abci.RequestInitChain{
		AppStateBytes: []byte("{}"),
	})
	ctx := dex.BaseApp.NewContext(false, abci.Header{ChainID: "unit-test-chain", Height: 1, Time: time.Unix(1558332092, 0)})

	mock := &MockApp{
		Cdc:             dex.Cdc,
		Mq:              dex.Mq,
		Ctx:             ctx,
		MarketKeeper:    dex.MarketKeeper,
		OrderKeeper:     dex.OrderKeeper,
		BankKeeper:      dex.BankKeeper,
		ExecutionKeeper: dex.ExecKeeper,
	}

	for _, opt := range options {
		opt(t, mock)
	}

	return mock
}

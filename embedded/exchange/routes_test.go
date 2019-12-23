package exchange

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/xar-network/xar-network/embedded/order"
	"github.com/xar-network/xar-network/testutil"
	"github.com/xar-network/xar-network/testutil/testflags"
	"github.com/xar-network/xar-network/types"
	"github.com/xar-network/xar-network/types/store"
)

func TestGetOrderHandler(t *testing.T) {
	testflags.UnitTest(t)
	cdc := codec.New()
	db := dbm.NewMemDB()
	k := order.NewKeeper(db, cdc)
	ctx := context.NewCLIContext()

	id := store.NewEntityID(0)
	genOwner := testutil.RandAddr()
	for i := 0; i < 110; i++ {
		id = id.Inc()
		var owner sdk.AccAddress
		if i%2 == 0 {
			owner = genOwner
		}

		var market store.EntityID
		if i%2 == 0 {
			market = store.NewEntityID(2)
		} else {
			market = store.NewEntityID(1)
		}

		require.NoError(t, k.OnEvent(types.OrderCreated{
			MarketID: market,
			ID:       id,
			Owner:    owner,
			CreatedTime: time.Unix(int64(i), 0),
		}))
	}

	req, err := http.NewRequest("GET", "/orders/50/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getOrderHandler(ctx, cdc))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	assert.Equal(t, rr.Body.String(), "Test OK = 50")
}

package market

import (
	"fmt"

	"github.com/xar-network/xar-network/types/errs"
	"github.com/xar-network/xar-network/types/store"
	"github.com/xar-network/xar-network/x/market/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

const (
	seqKey = "seq"
	valKey = "val"
)

type IteratorCB func(mkt types.Market) bool

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSubspace subspace.Subspace
}

func NewKeeper(sk sdk.StoreKey, cdc *codec.Codec, paramstore subspace.Subspace) Keeper {
	return Keeper{
		storeKey:      sk,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
	}
}

func (k Keeper) Get(ctx sdk.Context, id store.EntityID) (types.Market, sdk.Error) {
	params := k.GetParams(ctx)
	market := params.Markets[id.Uint64()]
	if (market == types.Market{}) {
		return types.Market{}, errs.ErrNotFound("not found")
	}
	return market, nil
}

func (k Keeper) Pair(ctx sdk.Context, id store.EntityID) (string, sdk.Error) {
	mkt, err := k.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", mkt.BaseAssetDenom, mkt.QuoteAssetDenom), nil
}

func (k Keeper) Has(ctx sdk.Context, id store.EntityID) bool {
	_, err := k.Get(ctx, id)
	//err == nil could have side effects, should check the error type
	return err == nil
}

func (k Keeper) Iterator(ctx sdk.Context, cb IteratorCB) {
	params := k.GetParams(ctx)
	for _, mkt := range params.Markets {
		if !cb(mkt) {
			break
		}
	}
}

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized market message type: %T", msg)).Result()
	}
}

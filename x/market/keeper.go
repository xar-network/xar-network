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
	codespace     sdk.CodespaceType
}

func NewKeeper(sk sdk.StoreKey, cdc *codec.Codec, paramstore subspace.Subspace, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:      sk,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
		codespace:     codespace,
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

func (k Keeper) CreateMarket(ctx sdk.Context, msg types.MsgCreateMarket) sdk.Result {
	if !k.IsNominee(ctx, msg.Nominee.String()) {
		return sdk.ErrInternal(fmt.Sprintf("not a nominee: '%s'", msg.Nominee.String())).Result()
	}
	params := k.GetParams(ctx)
	id := uint64(len(params.Markets))
	market := types.NewMarket(store.NewEntityID(id).Inc(), msg.BaseAsset, msg.QuoteAsset)
	params.Markets = append(params.Markets, market)
	k.SetParams(ctx, params)

	return sdk.Result{Events: ctx.EventManager().Events()}
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

func (k Keeper) IsNominee(ctx sdk.Context, nominee string) bool {
	params := k.GetParams(ctx)
	nominees := params.Nominees
	for _, v := range nominees {
		if v == nominee {
			return true
		}
	}
	return false
}

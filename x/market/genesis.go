package market

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/types/store"
	"github.com/xar-network/xar-network/x/market/types"
)

type GenesisState struct {
	Markets types.Markets
}

func NewGenesisState(markets types.Markets) GenesisState {
	return GenesisState{Markets: markets}
}

func ValidateGenesis(data GenesisState) error {
	currentId := store.ZeroEntityID
	for _, market := range data.Markets {
		currentId = currentId.Inc()
		if !currentId.Equals(market.ID) {
			return errors.New("Invalid Market: ID must monotonically increase.")
		}
		if market.BaseAssetDenom == "" {
			return errors.New("Invalid Market: Must specify a non-zero base asset denom.")
		}
		if market.QuoteAssetDenom == "" {
			return errors.New("Invalid Market: Must specify a non-zero quote asset denom.")
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Markets: types.Markets{
			{
				ID:              store.NewEntityID(1),
				BaseAssetDenom:  "uftm",
				QuoteAssetDenom: "uzar",
			},
			{
				ID:              store.NewEntityID(2),
				BaseAssetDenom:  "ueur",
				QuoteAssetDenom: "uzar",
			},
			{
				ID:              store.NewEntityID(3),
				BaseAssetDenom:  "uusd",
				QuoteAssetDenom: "uzar",
			},
			{
				ID:              store.NewEntityID(4),
				BaseAssetDenom:  "ubtc",
				QuoteAssetDenom: "uzar",
			},
		},
	}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetParams(ctx, types.NewParams(data.Markets))
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	return GenesisState{Markets: params.Markets}
}

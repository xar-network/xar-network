package market

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/xar-network/xar-network/types/store"
	"github.com/xar-network/xar-network/x/market/types"
)

type GenesisState struct {
	Markets []types.Market
}

func NewGenesisState(markets []types.Market) GenesisState {
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
		Markets: []types.Market{
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

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, asset := range data.Markets {
		keeper.Inject(ctx, asset)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var markets []types.Market
	k.Iterator(ctx, func(asset types.Market) bool {
		markets = append(markets, asset)
		return true
	})
	return GenesisState{Markets: markets}
}

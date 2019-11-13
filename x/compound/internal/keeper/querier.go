package keeper

// query endpoints supported by the moneymarket Querier
const (
	QueryMoneyMarkets   = "moneymarket"
	QueryMarketPosition   = "marketposition"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryMoneyMarkets:
			return queryMoneyMarkets(ctx, path[1:], req, keeper)
		case QueryMarketPosition:
			return queryGetMarketPosition(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown moneymarket query endpoint")
		}
	}
}

func queryMoneyMarkets(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	moneymarket := keeper.GetMarketInfo(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, moneymarket)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryGetMarketPosition(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	accAddress,err := sdk.AccAddressFromBech32(path[0])
	marketposition := keeper.GetMarketPosition(ctx, accAddress)
	fmt.Println(marketposition)
	res, err := codec.MarshalJSONIndent(keeper.cdc, marketposition)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

// query endpoints supported by the assetmanagement Querier
const (
	QuerySymbols = "symbols"
	QueryToken   = "token"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryToken:
			return queryToken(ctx, path[1:], req, keeper)
		case QuerySymbols:
			return querySymbols(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown assetmanagement query endpoint")
		}
	}
}

// nolint: unparam
func queryToken(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	searchToken := path[0]
	token, err := keeper.GetToken(ctx, searchToken)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, token)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func insertInto(s string, interval int, sep rune) string {
	var buffer bytes.Buffer
	before := interval - 1
	last := len(s) - 1
	for i, char := range s {
		buffer.WriteRune(char)
		if i%interval == before && i != last {
			buffer.WriteRune(sep)
		}
	}
	return buffer.String()
}

// nolint: unparam
func querySymbols(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var symbolList types.QueryResultSymbol

	iterator := keeper.GetTokensIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		symbolList = append(symbolList, string(iterator.Key()))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, symbolList)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

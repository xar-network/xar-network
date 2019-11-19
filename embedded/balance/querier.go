package balance

import (
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/xar-network/xar-network/types/errs"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryGet = "get"
)

func NewQuerier(bk bank.Keeper, sk supply.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGet(ctx, bk, sk, req.Data)
		default:
			return nil, sdk.ErrUnknownRequest("unknown balance request")
		}
	}
}

func queryGet(ctx sdk.Context, bk bank.Keeper, sk supply.Keeper, reqB []byte) ([]byte, sdk.Error) {
	var req GetQueryRequest
	err := amino.UnmarshalBinaryBare(reqB, &req)
	if err != nil {
		return nil, errs.ErrUnmarshalFailure("failed to unmarshal get query request")
	}

	res := GetQueryResponse{
		Balances: make([]GetQueryResponseBalance, 0),
	}
	balances := sk.GetSupply(ctx).GetTotal()
	for _, coin := range balances {
		bal := bk.GetCoins(ctx, req.Address).AmountOf(coin.Denom)
		if !bal.IsZero() {
			res.Balances = append(res.Balances, GetQueryResponseBalance{
				Denom:  coin.Denom,
				Liquid: sdk.NewUintFromString(bal.String()),
				AtRisk: sdk.ZeroUint(),
			})
		}
	}

	b, err := codec.MarshalJSONIndent(codec.New(), res)
	if err != nil {
		return nil, errs.ErrMarshalFailure("could not marshal result")
	}
	return b, nil
}

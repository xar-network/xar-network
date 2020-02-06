package util

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

type Genesis struct {
	Params        types.Params
	GlobalDebt    sdk.Int
	CSDTs         types.CSDTs
	TotalBorrows  types.CoinUs
	TotalSupplies types.CoinUs
	TotalReserves types.CoinUs
}

func GetGenesis(k keeper.Keeper, ctx sdk.Context) Genesis {
	params := k.GetParams(ctx)
	csdts := types.CSDTs{}
	borrows := types.CoinUs{}
	supplies := types.CoinUs{}
	reserves := types.CoinUs{}
	for _, param := range params.CollateralParams {
		l, err := k.GetCSDTs(ctx, param.Denom, sdk.Dec{})
		if err != nil {
			panic(err)
		} else {
			csdts = append(csdts, l...)
		}

		borrow, ok := k.GetTotalBorrows(ctx, param.Denom)
		if !ok {
			panic(fmt.Sprintf("Failed to retrieve total borrows for '%s'", param.Denom))
		}
		borrows = append(borrows, types.CoinU{
			Denom:  param.Denom,
			Amount: borrow,
		})

		supply, ok := k.GetTotalCash(ctx, param.Denom)
		if !ok {
			panic(fmt.Sprintf("Failed to retrieve total cash/supply for '%s'", param.Denom))
		}
		supplies = append(supplies, types.CoinU{
			Denom:  param.Denom,
			Amount: supply,
		})

		reserve, ok := k.GetTotalReserve(ctx, param.Denom)
		if !ok {
			panic(fmt.Sprintf("Failed to retrieve total reserve for '%s'", param.Denom))
		}
		reserves = append(reserves, types.CoinU{
			Denom:  param.Denom,
			Amount: reserve,
		})
	}
	debt := k.GetGlobalDebt(ctx)

	return Genesis{
		Params:        params,
		GlobalDebt:    debt,
		CSDTs:         csdts,
		TotalBorrows:  borrows,
		TotalSupplies: supplies,
		TotalReserves: reserves,
	}
}

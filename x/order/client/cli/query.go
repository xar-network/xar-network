package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/xar-network/xar-network/x/order/types"
)

func GetCmdListOrders(queryRoute string, cdc *codec.Codec) *cobra.Command {
	out := &cobra.Command{
		Use:   "list",
		Short: "lists orders with the specified filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)
			res, height, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/list", queryRoute), nil)
			if err != nil {
				return err
			}
			ctx = ctx.WithHeight(height)

			var out types.ListQueryResult
			cdc.MustUnmarshalJSON(res, &out)
			return ctx.PrintOutput(out)
		},
	}
	return out
}

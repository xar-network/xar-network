package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// GetCmd_GetCsdt queries the latest info about a particular csdt
func GetCmd_GetCsdt(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "csdt [ownerAddress] [collateralType]",
		Short: "get info about a csdt",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Prepare params for querier
			ownerAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			collateralType := args[1] // TODO validation?
			bz, err := cdc.MarshalJSON(types.QueryCsdtsParams{
				Owner:           ownerAddress,
				CollateralDenom: collateralType,
			})
			if err != nil {
				return err
			}

			// Query
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetCsdts)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				fmt.Printf("error when getting csdt info - %s", err)
				fmt.Printf("could not get current csdt info - %s %s \n", string(ownerAddress), string(collateralType))
				return err
			}
			cliCtx = cliCtx.WithHeight(height)

			// Decode and print results
			var csdts types.CSDTs
			cdc.MustUnmarshalJSON(res, &csdts)
			if len(csdts) != 1 {
				panic("Unexpected number of CSDTs returned from querier. This shouldn't happen.")
			}
			return cliCtx.PrintOutput(csdts[0])
		},
	}
}

func GetCmd_GetCsdts(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "csdts [collateralType]",
		Short: "get info about many csdts",
		Long:  "Get all CSDTs or specify a collateral type to get only CSDTs with that collateral type.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Prepare params for querier
			bz, err := cdc.MarshalJSON(types.QueryCsdtsParams{CollateralDenom: args[0]}) // denom="" returns all CSDTs // TODO will this fail if there are no args?
			if err != nil {
				return err
			}

			// Query
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetCsdts)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)

			// Decode and print results
			var out types.CSDTs
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmd_GetUnderCollateralizedCsdts(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bad-csdts [collateralType] [price]",
		Short: "get under collateralized CSDTs",
		Long:  "Get all CSDTS of a particular collateral type that will be under collateralized at the specified price. Pass in the current price to get currently under collateralized CSDTs.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Prepare params for querier
			price, errSdk := sdk.NewDecFromStr(args[1])
			if errSdk != nil {
				return fmt.Errorf(errSdk.Error()) // TODO check this returns useful output
			}
			bz, err := cdc.MarshalJSON(types.QueryCsdtsParams{
				CollateralDenom:       args[0],
				UnderCollateralizedAt: price,
			})
			if err != nil {
				return err
			}

			// Query
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetCsdts)
			res, height, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)

			// Decode and print results
			var out types.CSDTs
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmd_GetParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "get the csdt module parameters",
		Long:  "Get the current global csdt module parameters.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Query
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGetParams)
			res, height, err := cliCtx.QueryWithData(route, nil) // TODO use cliCtx.QueryStore?
			if err != nil {
				return err
			}
			cliCtx = cliCtx.WithHeight(height)

			// Decode and print results
			var out types.Params
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

package cli

import (
	"bufio"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	tmtime "github.com/tendermint/tendermint/types/time"
	"github.com/xar-network/xar-network/x/oracle/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Oracle transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		client.PostCommands(
			GetCmdPostPrice(cdc),
			getCmdAddOracle(cdc),
			getCmdSetOracles(cdc),
		)...,
	)

	return cmd
}

// GetCmdPostPrice cli command for posting prices.
func GetCmdPostPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "postprice [from_key_or_address] [assetCode] [price] [expiry]",
		Short: "post the latest price for a particular asset",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			price, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}
			expiryInt, ok := sdk.NewIntFromString(args[2])
			if !ok {
				fmt.Printf("invalid expiry - %s \n", args[2])
				return nil
			}
			expiry := tmtime.Canonical(time.Unix(expiryInt.Int64(), 0))
			msg := types.NewMsgPostPrice(cliCtx.GetFromAddress(), args[1], price, expiry)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdAddOracle(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "add-oracle [nominee_key] [denom] [oracle_address]",
		Example: "xarcli oracle add-oracle nominee xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu",
		Short:   "Create a new oracle",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			oracleAddr, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddOracle(cliCtx.GetFromAddress(), args[1], oracleAddr)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdSetOracles(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "set-oracles [nominee_key] [denom] [oracle_addresses]",
		Example: "xarcli oracle add-oracle nominee xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu,xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu,xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu",
		Short:   "Sets a list of oracles for a denom",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			oracles, err := types.ParseOracles(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetOracles(cliCtx.GetFromAddress(), args[1], oracles)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

package cli

import (
	"bufio"

	"github.com/xar-network/xar-network/x/authority/internal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/spf13/cobra"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	authorityCmds := &cobra.Command{
		Use:                "authority",
		Short:              "Authority transactions subcommands",
		DisableFlagParsing: false,
		RunE:               client.ValidateCmd,
	}

	authorityCmds.AddCommand(
		client.PostCommands(
			getCmdCreateIssuer(cdc),
			getCmdCreateOracle(cdc),
			getCmdCreateMarket(cdc),
			getCmdDestroyIssuer(cdc),
			getCmdSetSupply(cdc),
		)...,
	)

	return authorityCmds
}

func getCmdCreateIssuer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "create-issuer [authority_key_or_address] [issuer_address] [denominations]",
		Example: "xarcli authority create-issuer masterkey xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu x2eur,x0jpy",
		Short:   "Create a new issuer",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			issuerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denoms, err := types.ParseDenominations(args[2])
			if err != nil {
				return err
			}

			msg := types.MsgCreateIssuer{
				Issuer:        issuerAddr,
				Denominations: denoms,
				Authority:     cliCtx.GetFromAddress(),
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdCreateOracle(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "create-oracle [authority_key_or_address] [oracle_address]",
		Example: "xarcli authority create-oracle masterkey xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu",
		Short:   "Create a new oracle",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			oracleAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgCreateOracle{
				Oracle:    oracleAddr,
				Authority: cliCtx.GetFromAddress(),
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdCreateMarket(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "create-market [authority_key_or_address] [base_asset_denom] [quote_asset_denom]",
		Example: "xarcli authority create-market masterkey uftm ucsdt",
		Short:   "Create a new market",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			msg := types.MsgCreateMarket{
				BaseAsset:  args[1],
				QuoteAsset: args[2],
				Authority:  cliCtx.GetFromAddress(),
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdSetSupply(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "set-supply [authority_key_or_address] [supply]",
		Example: "xarcli authority create-market masterkey 100000uftm",
		Short:   "Set the supply for a denomination",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgSetSupply{
				Supply:    coins,
				Authority: cliCtx.GetFromAddress(),
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func getCmdDestroyIssuer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "destroy-issuer [authority_key_or_address] [issuer_address]",
		Example: "xarcli authority destory-issuer masterkey xar17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu",
		Short:   "Delete an issuer",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			issuerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgDestroyIssuer{
				Issuer:    issuerAddr,
				Authority: cliCtx.GetFromAddress(),
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

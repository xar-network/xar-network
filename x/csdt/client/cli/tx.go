package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	csdtTxCmd := &cobra.Command{
		Use:   "csdt",
		Short: "CSDT transactions subcommands",
	}

	csdtTxCmd.AddCommand(
		GetCmdModifyCsdt(cdc),
	)

	return csdtTxCmd
}

// GetCmdModifyCsdt cli command for creating and modifying csdts.
func GetCmdModifyCsdt(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modifycsdt [from_key_or_addres] [collateralType] [collateralChange] [debtChange]",
		Short: "create or modify a csdt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			collateralChange, ok := sdk.NewIntFromString(args[2])
			if !ok {
				fmt.Printf("invalid collateral amount - %s \n", string(args[2]))
				return nil
			}
			debtChange, ok := sdk.NewIntFromString(args[3])
			if !ok {
				fmt.Printf("invalid debt amount - %s \n", string(args[3]))
				return nil
			}
			msg := types.NewMsgCreateOrModifyCSDT(cliCtx.GetFromAddress(), args[1], collateralChange, debtChange)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

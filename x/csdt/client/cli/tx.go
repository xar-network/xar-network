/*

Copyright 2016 All in Bits, Inc
Copyright 2019 Kava Labs, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package cli

import (
	"bufio"
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
		Use:   "modifycsdt [from_key_or_addres] [collateralDenom] [collateralChange] [debtChange]",
		Short: "create or modify a csdt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

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

// GetCmdDepositCollateral cli command for depositing collateral.
func GetCmdDepositCollateral(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-collateral [from_key_or_addres] [collateralDenom] [collateralChange]",
		Short: "deposit collateral to csdt",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			collateralChange, ok := sdk.NewIntFromString(args[2])
			if !ok || collateralChange.IsZero() || collateralChange.IsNegative() {
				fmt.Printf("invalid collateral amount - %s \n", string(args[2]))
				return nil
			}
			denom := args[1]
			if len(denom) == 0 {
				fmt.Printf("invalid denom - %s \n", string(args[1]))
				return nil
			}
			msg := types.NewMsgDepositCollateral(cliCtx.GetFromAddress(), denom, collateralChange)
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

// GetCmdWithdrawCollateral cli command for withdrawing collateral.
func GetCmdWithdrawCollateral(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-collateral [from_key_or_addres] [collateralDenom] [collateralChange]",
		Short: "withdraw collateral from csdt",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			collateralChange, ok := sdk.NewIntFromString(args[2])
			if !ok || collateralChange.IsZero() || collateralChange.IsNegative() {
				fmt.Printf("invalid collateral amount - %s \n", string(args[2]))
				return nil
			}
			denom := args[1]
			if len(denom) == 0 {
				fmt.Printf("invalid denom - %s \n", string(args[1]))
				return nil
			}
			msg := types.NewMsgWithdrawCollateral(cliCtx.GetFromAddress(), denom, collateralChange)
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

// GetCmdSettleDebt cli command for settling debt.
func GetCmdDepositDebt(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settle-debt [from_key_or_addres] [collateralDenom] [debtDenom] [debtChange]",
		Short: "settle debt with csdt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			debtChange, ok := sdk.NewIntFromString(args[3])
			if !ok || debtChange.IsZero() || debtChange.IsNegative() {
				fmt.Printf("invalid debt amount - %s \n", string(args[3]))
				return nil
			}
			collateralDenom := args[1]
			if len(collateralDenom) == 0 {
				fmt.Printf("invalid collateral denom - %s \n", string(args[1]))
				return nil
			}
			debtDenom := args[2]
			if len(debtDenom) == 0 {
				fmt.Printf("invalid debt denom - %s \n", string(args[2]))
				return nil
			}
			msg := types.NewMsgSettleDebt(cliCtx.GetFromAddress(), collateralDenom, debtDenom, debtChange)
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

// GetCmdWithdrawDebt cli command for settling debt.
func GetCmdWithdrawDebt(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-debt [from_key_or_addres] [collateralDenom] [debtDenom] [debtChange]",
		Short: "withdraw debt from csdt",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)

			debtChange, ok := sdk.NewIntFromString(args[3])
			if !ok || debtChange.IsZero() || debtChange.IsNegative() {
				fmt.Printf("invalid debt amount - %s \n", string(args[3]))
				return nil
			}
			collateralDenom := args[1]
			if len(collateralDenom) == 0 {
				fmt.Printf("invalid collateral denom - %s \n", string(args[1]))
				return nil
			}
			debtDenom := args[2]
			if len(debtDenom) == 0 {
				fmt.Printf("invalid debt denom - %s \n", string(args[2]))
				return nil
			}
			msg := types.NewMsgWithdrawDebt(cliCtx.GetFromAddress(), collateralDenom, debtDenom, debtChange)
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

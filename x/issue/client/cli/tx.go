package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xar-network/xar-network/x/issue/client/rest"
	"github.com/xar-network/xar-network/x/issue/internal/types"
)

// GetIssueCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	issueCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Issue transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	issueCmd.AddCommand(
		IssueCreateCmd(cdc),
		IssueTransferOwnershipCmd(cdc),
		IssueDescriptionCmd(cdc),
		IssueMintCmd(cdc),
		IssueDisableFeatureCmd(cdc),
		IssueFreezeCmd(cdc),
		IssueUnFreeCmd(cdc),
		IssueBurnCmd(cdc),
		IssueBurnFromCmd(cdc),
		IssueSendFromCmd(cdc),
		IssueApproveCmd(cdc),
		IssueIncreaseApprovalCmd(cdc),
		IssueDecreaseApprovalCmd(cdc),
	)

	return issueCmd
}

// IssueCreateCmd implements the issue a coin transaction command.
func IssueCreateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [from_key_or_address] [name] [symbol] [total_supply]",
		Args:    cobra.ExactArgs(4),
		Short:   "Issue a new token",
		Example: "$ xarcli issue create coin_key Coin CN 1",
		RunE: func(cmd *cobra.Command, args []string) error {
			totalSupply, ok := sdk.NewIntFromString(args[3])

			if !ok {
				return fmt.Errorf("Total supply %s not a valid int, please input a valid total supply", args[2])
			}

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			coinIssueInfo := types.IssueParams{
				Name:               args[1],
				Symbol:             strings.ToUpper(args[2]),
				BurnOwnerDisabled:  viper.GetBool(flagBurnOwnerDisabled),
				BurnHolderDisabled: viper.GetBool(flagBurnHolderDisabled),
				BurnFromDisabled:   viper.GetBool(flagBurnFromDisabled),
				MintingFinished:    viper.GetBool(flagMintingFinished),
				FreezeDisabled:     viper.GetBool(flagFreezeDisabled),
				TotalSupply:        totalSupply,
				Decimals:           uint(viper.GetInt(flagDecimals)),
			}
			coinIssueInfo.TotalSupply = types.MulDecimals(coinIssueInfo.TotalSupply, coinIssueInfo.Decimals)
			msg := types.NewMsgIssue(cliCtx.GetFromAddress(), &coinIssueInfo)

			validateErr := msg.ValidateBasic()

			if validateErr != nil {
				return types.Errorf(validateErr)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Uint(flagDecimals, types.CoinDecimalsMaxValue, "Decimals of the token")
	cmd.Flags().Bool(flagBurnOwnerDisabled, false, "Disable token owner burn")
	cmd.Flags().Bool(flagBurnHolderDisabled, false, "Disable token holder burn")
	cmd.Flags().Bool(flagBurnFromDisabled, false, "Disable token owner burn from any holder")
	cmd.Flags().Bool(flagMintingFinished, false, "Token owner can not mint")
	cmd.Flags().Bool(flagFreezeDisabled, false, "Token holder can transfer the token")

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// IssueTransferOwnershipCmd implements transfer a coin owner ship transaction command.
func IssueTransferOwnershipCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer-ownership [from_key_or_address] [issue_id] [to_address]",
		Args:    cobra.ExactArgs(3),
		Short:   "Transfer ownership a token",
		Long:    "Token owner transfer the ownership to new account",
		Example: "$ xarcli issue transfer-ownership coin_key coin174876e800 xard1vf7pnhwh5v4lmdp59dms2andn2hhperghppkxc",
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[1]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			_, err = types.IssueOwnerCheck(cliCtx, cliCtx.GetFromAddress(), issueID)
			if err != nil {
				return err
			}
			msg := types.NewMsgIssueTransferOwnership(issueID, cliCtx.GetFromAddress(), to)

			validateErr := msg.ValidateBasic()

			if validateErr != nil {
				return types.Errorf(validateErr)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// IssueDescriptionCmd implements issue a coin transaction command.
func IssueDescriptionCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe [from_key_or_address] [issue_id] [description_file]",
		Args:    cobra.ExactArgs(3),
		Short:   "Add description to a token",
		Long:    "Owner can add a description of the token. The description needs to be in json format. You can customize preferences or use recommended templates.",
		Example: "$ xarcli issue describe coin_key coin174876e800 path/description.json --from foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[0]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			contents, err := ioutil.ReadFile(args[2])
			if err != nil {
				return err
			}
			buffer := bytes.Buffer{}
			err = json.Compact(&buffer, contents)
			if err != nil {
				return types.ErrCoinDescriptionNotValid()
			}
			contents = buffer.Bytes()

			_, err = types.IssueOwnerCheck(cliCtx, cliCtx.GetFromAddress(), issueID)
			if err != nil {
				return err
			}
			msg := types.NewMsgIssueDescription(issueID, cliCtx.GetFromAddress(), contents)

			validateErr := msg.ValidateBasic()

			if validateErr != nil {
				return types.Errorf(validateErr)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// IssueMintCmd implements mint a coinIssue transaction command.
func IssueMintCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mint [from_key_or_address] [issue_id] [to] [amount]",
		Args:    cobra.ExactArgs(4),
		Short:   "Mint tokens",
		Long:    "Token owner can mint the token to an address",
		Example: "$ xarcli issue mint coin_key coin174876e800 xard1vf7pnhwh5v4lmdp59dms2andn2hhperghppkxc 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[1]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			amount, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return types.Errorf(types.ErrAmountNotValid())
			}
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			issueInfo, err := types.IssueOwnerCheck(cliCtx, cliCtx.GetFromAddress(), issueID)
			if err != nil {
				return err
			}

			if issueInfo.IsMintingFinished() {
				return types.Errorf(types.ErrCanNotMint())
			}

			amount = types.MulDecimals(amount, issueInfo.GetDecimals())
			msg := types.NewMsgIssueMint(issueID, cliCtx.GetFromAddress(), to, amount, issueInfo.GetDecimals())

			validateErr := msg.ValidateBasic()
			if validateErr != nil {
				return types.Errorf(validateErr)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

// IssueDisableFeatureCmd implements disable feature a coinIssue transaction command.
func IssueDisableFeatureCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable [from_key_or_address] [issue_id] [feature]",
		Args:  cobra.ExactArgs(3),
		Short: "Disable a feature from the token",
		Long: fmt.Sprintf("Token Owner disabled the features:\n"+
			"%s:Token owner can burn the token\n"+
			"%s:Token holder can burn the token\n"+
			"%s:Token owner can burn the token from any holder\n"+
			"%s:Token owner can freeze in and out the token from any address\n"+
			"%s:Token owner can mint the token", types.BurnOwner, types.BurnHolder, types.BurnFrom, types.Freeze, types.Minting),
		Example: fmt.Sprintf("$ xarcli issue disable coin_key coin174876e800 %s\n"+
			"$ xarcli issue disable coin_key coin174876e800 %s\n"+
			"$ xarcli issue disable coin_key coin174876e800 %s\n"+
			"$ xarcli issue disable coin_key coin174876e800 %s\n"+
			"$ xarcli issue disable coin_key coin174876e800 %s",
			types.BurnOwner, types.BurnHolder, types.BurnFrom, types.Freeze, types.Minting),

		RunE: func(cmd *cobra.Command, args []string) error {
			feature := args[2]

			_, ok := types.Features[feature]
			if !ok {
				return types.Errorf(types.ErrUnknownFeatures())
			}

			issueID := args[1]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			_, err := types.IssueOwnerCheck(cliCtx, cliCtx.GetFromAddress(), issueID)
			if err != nil {
				return err
			}

			msg := types.NewMsgIssueDisableFeature(issueID, cliCtx.GetFromAddress(), feature)
			validateErr := msg.ValidateBasic()
			if validateErr != nil {
				return types.Errorf(validateErr)
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

// IssueFreezeCmd implements freeze a token transaction command.
func IssueFreezeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze [from_key_or_address] [freeze_type] [issue_id] [address]",
		Args:  cobra.ExactArgs(4),
		Short: "Freeze transfers from an address",
		Long: fmt.Sprintf("Token owner freeze the transfer from an address:\n\n"+
			"%s:The address can not transfer in\n"+
			"%s:The address can not transfer out\n"+
			"%s:The address not can transfer in or out\n\n", types.FreezeIn, types.FreezeOut, types.FreezeInAndOut),
		Example: "$ xarcli issue freeze coin_key in coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n\n" +
			"$ xarcli issue freeze coin_key out coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n\n" +
			"$ xarcli issue freeze coin_key in-out coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueFreeze(cdc, args, true)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

// IssueUnFreeCmd implements un freeze  a token transaction command.
func IssueUnFreeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfreeze [from_key_or_address] [freeze_type] [issue_id] [address]",
		Args:  cobra.ExactArgs(4),
		Short: "UnFreeze transfers from an address",
		Long: fmt.Sprintf("Token owner unFreeze the transfer from a address:\n\n"+
			"%s:The address can transfer in\n"+
			"%s:The address can transfer out\n"+
			"%s:The address can transfer in and out", types.FreezeIn, types.FreezeOut, types.FreezeInAndOut),
		Example: "$ xarcli issue unfreeze coin_key in coin174876e800 xar15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n\n" +
			"$ xarcli issue unfreeze coin_key out coin174876e800 xar15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n\n" +
			"$ xarcli issue unfreeze coin_key in-out coin174876e800 xar15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueFreeze(cdc, args, false)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

func issueFreeze(cdc *codec.Codec, args []string, freeze bool) error {
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

	msg, err := rest.GetIssueFreezeMsg(cliCtx, cliCtx.GetFromAddress(), args[1], args[2], args[3], freeze)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

// IssueBurnCmd implements burn a coinIssue transaction command.
func IssueBurnCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "burn [from_key_or_address] [issue_id] [amount]",
		Args:    cobra.ExactArgs(3),
		Short:   "Token holder can burn the token",
		Long:    "Token holder or the Owner burns the token he holds (the Owner can burn if 'burning_owner_disabled' is false, the holder can burn if 'burning_holder_disabled' is false)",
		Example: "$ xarcli issue burn coin_key coin174876e800 88888",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueBurnFrom(cdc, args, types.BurnHolder)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

// IssueBurnFromCmd implements burn a coinIssue transaction command.
func IssueBurnFromCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "burn-from [from_key_or_address] [issue_id] [from_address] [amount]",
		Args:    cobra.ExactArgs(4),
		Short:   "Token owner burns the token",
		Long:    "Token Owner burns the token from any holder (the Owner can burn if 'burning_any_disabled' is false)",
		Example: "$ xarcli issue burn-from coin_key coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueBurnFrom(cdc, args, types.BurnFrom)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

func issueBurnFrom(cdc *codec.Codec, args []string, burnFromType string) error {
	issueID := args[1]
	if err := types.CheckIssueId(issueID); err != nil {
		return types.Errorf(err)
	}
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

	amountStr := ""
	//burn sender
	accAddress := cliCtx.GetFromAddress()

	if types.BurnFrom == burnFromType {
		acc, err := sdk.AccAddressFromBech32(args[2])
		accAddress = acc
		if err != nil {
			return err
		}
		amountStr = args[3]
	} else {
		amountStr = args[2]
	}
	amount, ok := sdk.NewIntFromString(amountStr)
	if !ok {
		return types.Errorf(types.ErrAmountNotValid())
	}
	msg, err := rest.GetBurnMsg(cliCtx, cliCtx.GetFromAddress(), accAddress, issueID, amount, burnFromType, true)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

// IssueSendFromCmd implements send from a token transaction command.
func IssueSendFromCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send-from [from_key_or_address] [issue_id] [from_address] [to_address] [amount]",
		Args:    cobra.ExactArgs(5),
		Short:   "Send tokens from one address to another",
		Long:    "Send tokens from one address to another by allowance",
		Example: "$ xarcli issue send-from coin_key coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n xard1vud9ptwagudgq7yht53cwuf8qfmgkd0qcej0ah 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			issueID := args[1]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			fromAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			toAddress, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return types.Errorf(types.ErrAmountNotValid())
			}

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			if err := rest.CheckAllowance(cliCtx, issueID, fromAddress, cliCtx.GetFromAddress(), amount); err != nil {
				return err
			}

			if err = rest.CheckFreeze(cliCtx, issueID, fromAddress, toAddress); err != nil {
				return err
			}

			issueInfo, err := rest.GetIssueByID(cliCtx, issueID)
			if err != nil {
				return err
			}
			amount = types.MulDecimals(amount, issueInfo.GetDecimals())

			msg := types.NewMsgIssueSendFrom(issueID, cliCtx.GetFromAddress(), fromAddress, toAddress, amount)

			validateErr := msg.ValidateBasic()
			if validateErr != nil {
				return types.Errorf(validateErr)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

// IssueApproveCmd implements approve a token transaction command.
func IssueApproveCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approve [from_key_or_address] [issue-id] [address] [amount]",
		Args:    cobra.ExactArgs(4),
		Short:   "Approve tokens on behalf of sender",
		Long:    "Approve the passed address to spend the specified amount of tokens on behalf of sender",
		Example: "$ xarcli issue approve coin_key coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueApprove(cdc, args, types.Approve)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

// IssueIncreaseApprovalCmd implements increase approval a token transaction command.
func IssueIncreaseApprovalCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "increase-approval [from_key_or_address] [issue_id] [address] [amount]",
		Args:    cobra.ExactArgs(4),
		Short:   "Increase approval to spend tokens on behalf of sender",
		Long:    "Increase approval to spend the specified amount of tokens on behalf of sender",
		Example: "$ xarcli issue increase-approval coin_key coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueApprove(cdc, args, types.IncreaseApproval)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}

// IssueDecreaseApprovalCmd implements decrease approval a token transaction command.
func IssueDecreaseApprovalCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrease-approval [from_key_or_address] [issue_id] [address] [amount]",
		Args:    cobra.ExactArgs(4),
		Short:   "Decrease approval to spend tokens on behalf of sender",
		Long:    "Decrease approval to spend the specified amount of tokens on behalf of sender",
		Example: "$ xarcli issue decrease-approval coin_key coin174876e800 xard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n 100",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issueApprove(cdc, args, types.DecreaseApproval)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	return cmd
}
func issueApprove(cdc *codec.Codec, args []string, approveType string) error {
	issueID := args[1]
	accAddress, err := sdk.AccAddressFromBech32(args[2])
	if err != nil {
		return err
	}
	amount, ok := sdk.NewIntFromString(args[3])
	if !ok {
		return types.Errorf(types.ErrAmountNotValid())
	}
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

	msg, err := rest.GetIssueApproveMsg(cliCtx, issueID, cliCtx.GetFromAddress(), accAddress, approveType, amount, true)
	if err != nil {
		return err
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}

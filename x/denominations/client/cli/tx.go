package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/xar-network/xar-network/x/denominations/rand"

	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txRootCmd := &cobra.Command{
		Use:                        "token",
		Short:                      "Asset Management transaction sub-commands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txRootCmd.AddCommand(client.PostCommands(
		GetCmdIssueToken(cdc),
		GetCmdMintCoins(cdc),
		GetCmdBurnCoins(cdc),
		GetCmdFreezeCoins(cdc),
		GetCmdUnfreezeCoins(cdc),
	)...)

	return txRootCmd
}

func flagError2String(name string, err error) string {
	return fmt.Sprintf("unable to find '%s' flag: %v", name, err)
}

func fetchStringFlag(cmd *cobra.Command, flagName string) string {
	flag, err := cmd.Flags().GetString(flagName)
	if err != nil {
		panic(flagError2String(flagName, err))
	}

	return flag
}

func fetchInt64Flag(cmd *cobra.Command, flagName string) int64 {
	flag, err := cmd.Flags().GetInt64(flagName)
	if err != nil {
		panic(flagError2String(flagName, err))
	}

	return flag
}

func fetchBoolFlag(cmd *cobra.Command, flagName string) bool {
	flag, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		panic(flagError2String(flagName, err))
	}

	return flag
}

func setupRequiredFlag(cmd *cobra.Command, name string) {
	err := cmd.MarkFlagRequired(name)
	if err != nil {
		panic(fmt.Sprintf("failed to setup '%s' flag: %s", name, err))
	}
}

func setupBoolFlag(cmd *cobra.Command, name string, shorthand string, value bool, usage string, required bool) {
	cmd.Flags().BoolP(name, shorthand, value, usage)
	if required {
		setupRequiredFlag(cmd, name)
	}
}

func setupStringFlag(cmd *cobra.Command, name string, shorthand string, value string, usage string, required bool) {
	cmd.Flags().StringP(name, shorthand, value, usage)
	if required {
		setupRequiredFlag(cmd, name)
	}
}

func setupInt64Flag(cmd *cobra.Command, name string, shorthand string, value int64, usage string, required bool) {
	cmd.Flags().Int64P(name, shorthand, value, usage)
	if required {
		setupRequiredFlag(cmd, name)
	}
}

func getAccountAddress(cliCtx client.CLIContext) sdk.AccAddress {
	from := cliCtx.GetFromName()
	address := cliCtx.GetFromAddress()
	fmt.Printf("token account: %s / %v", from, address)

	return address
}

// GetCmdIssueToken is the CLI command for sending a IssueToken transaction
func GetCmdIssueToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: `issue --token-name [name] --total-supply [amount]
			--symbol [ABC] --mintable --from [account]`,
		Short: "create a new asset",
		// Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			address := getAccountAddress(cliCtx)

			name := fetchStringFlag(cmd, "token-name")
			originalSymbol := fetchStringFlag(cmd, "symbol")
			symbol := strings.ToLower(rand.GenerateNewSymbol(originalSymbol))
			totalSupply := sdk.NewInt(fetchInt64Flag(cmd, "total-supply"))
			maxSupply := sdk.NewInt(fetchInt64Flag(cmd, "max-supply"))
			mintable := fetchBoolFlag(cmd, "mintable")
			log.Debugf("token is mintable? %t", mintable)

			msg := types.NewMsgIssueToken(address, name, symbol, originalSymbol, totalSupply, maxSupply, mintable)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupBoolFlag(cmd, "mintable", "", false, "is the new token mintable", false)
	setupStringFlag(cmd, "token-name", "", "", "the name of the new token", true)
	setupInt64Flag(cmd, "total-supply", "", -1,
		"what is the total supply for the new token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC, for the new token", true)

	return cmd
}

func printIssuedSymbol(response *sdk.TxResponse, cliCtx context.CLIContext) {
	if response != nil {
		txHash := response.TxHash
		tmTags := []string{fmt.Sprintf("tx.hash='%s'", txHash)}

		txs, err := utils.QueryTxsByEvents(cliCtx, tmTags, 1, 1)
		if err != nil {
			log.Errorf("Failed to find new transaction for hash: %s because: %s", txHash, err)
		} else if txs != nil {
			if txs.Count == 0 || len(txs.Txs[0].Logs) == 0 {
				log.Errorf("Failed to find new transaction for hash: %s", txHash)
			} else {
				log.Debugf("Found new transaction: %s", txs.Txs)
				logResponse := txs.Txs[0].Logs[0]
				if logResponse.Success == false {
					log.Errorf("Transaction failed: %v", logResponse)
				} else {
					newSymbol := logResponse.Log
					log.Infof("Symbol issued: %s", newSymbol)
				}
			}
		}
	} else {
		log.Errorf("No response")
	}
}

// GetCmdMintCoins is the CLI command for sending a MintCoins transaction
func GetCmdMintCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `mint --amount [amount] --symbol [ABC-123]`,
		Short: "mint more coins for the specified token",
		// Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// if err := cliCtx.EnsureAccountExists(); err != nil {
			// 	return err
			// }

			address, symbol, amount := getCommonParameters(cliCtx, cmd)

			msg := types.NewMsgMintCoins(amount, symbol, address)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to mint for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}

func getCommonParameters(cliCtx client.CLIContext, cmd *cobra.Command) (sdk.AccAddress, string, sdk.Int) {
	// find given account
	address := getAccountAddress(cliCtx)
	symbol := fetchStringFlag(cmd, "symbol")
	amount := sdk.NewInt(fetchInt64Flag(cmd, "amount"))
	return address, symbol, amount
}

// GetCmdBurnCoins is the CLI command for sending a BurnCoins transaction
func GetCmdBurnCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `burn --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "destroy the given amount of token/coins, reducing the total supply",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cliCtx, cmd)

			msg := types.NewMsgBurnCoins(amount, symbol, address)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to burn for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}

// GetCmdFreezeCoins is the CLI command for sending a FreezeCoins transaction
func GetCmdFreezeCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `freeze [freeze_address] --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "move specified amount of token/coins into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cliCtx, cmd)
			freezeAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgFreezeCoins(amount, symbol, address, freezeAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to freeze for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}

// GetCmdUnfreezeCoins is the CLI command for sending a FreezeCoins transaction
func GetCmdUnfreezeCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `unfreeze [freeze_address] --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "move specified amount of token into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cliCtx, cmd)
			freezeAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnfreezeCoins(amount, symbol, address, freezeAddress)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	setupInt64Flag(cmd, "amount", "", -1,
		"what is the total amount of coins to unfreeze for the given token", true)
	setupStringFlag(cmd, "symbol", "", "",
		"what is the shorthand symbol, eg ABC-123, for the existing token", true)

	return cmd
}

package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xar-network/xar-network/x/denominations/rand"

	"github.com/xar-network/xar-network/x/denominations/internal/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			cliCtx.BroadcastMode = flags.BroadcastBlock // wait in order to query later

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			address := getAccountAddress(cliCtx)

			name := fetchStringFlag(cmd, "token-name")
			originalSymbol := fetchStringFlag(cmd, "symbol")
			symbol := strings.ToLower(rand.GenerateNewSymbol(originalSymbol))
			totalSupply := fetchInt64Flag(cmd, "total-supply")
			mintable := fetchBoolFlag(cmd, "mintable")
			log.Debugf("token is mintable? %t", mintable)

			msg := types.NewMsgIssueToken(address, name, symbol, originalSymbol, totalSupply, mintable)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			response, result := GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

			printIssuedSymbol(response, cliCtx)

			return result
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

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

func getCommonParameters(cliCtx client.CLIContext, cmd *cobra.Command) (sdk.AccAddress, string, int64) {
	// find given account
	address := getAccountAddress(cliCtx)
	symbol := fetchStringFlag(cmd, "symbol")
	amount := fetchInt64Flag(cmd, "amount")
	return address, symbol, amount
}

// GetCmdBurnCoins is the CLI command for sending a BurnCoins transaction
func GetCmdBurnCoins(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   `burn --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "destroy the given amount of token/coins, reducing the total supply",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

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
		Use:   `freeze --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "move specified amount of token/coins into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cliCtx, cmd)

			msg := types.NewMsgFreezeCoins(amount, symbol, address)
			err := msg.ValidateBasic()
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
		Use:   `unfreeze --amount [amount] --symbol [ABC-123] --from [account]`,
		Short: "move specified amount of token into frozen status, preventing their sale",
		// Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			address, symbol, amount := getCommonParameters(cliCtx, cmd)

			msg := types.NewMsgUnfreezeCoins(amount, symbol, address)
			err := msg.ValidateBasic()
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

// GenerateOrBroadcastMsgs creates a StdTx given a series of messages. If
// the provided context has generate-only enabled, the tx will only be printed
// to STDOUT in a fully offline manner. Otherwise, the tx will be signed and
// broadcast.
func GenerateOrBroadcastMsgs(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, msgs []sdk.Msg) (*sdk.TxResponse, error) {
	if cliCtx.GenerateOnly {
		return nil, utils.PrintUnsignedStdTx(txBldr, cliCtx, msgs)
	}

	return CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
}

// CompleteAndBroadcastTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node
// and returns the response and/or error
func CompleteAndBroadcastTxCLI(txBldr authtypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (*sdk.TxResponse, error) {
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return nil, err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return nil, err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return nil, nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return nil, err
		}

		var json []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return nil, err
		}
	}

	passphrase, err := keys.GetPassphrase(fromName)
	if err != nil {
		return nil, err
	}

	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, passphrase, msgs)
	if err != nil {
		return nil, err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	return &res, cliCtx.PrintOutput(res)
}

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/cosmos/modules/incubator/nft"
	nftcmd "github.com/cosmos/modules/incubator/nft/client/cli"
	nftrest "github.com/cosmos/modules/incubator/nft/client/rest"
	issuecmd "github.com/zar-network/zar-network/x/issue/client/cli"
	issuerest "github.com/zar-network/zar-network/x/issue/client/rest"

	auctionclient "github.com/zar-network/zar-network/x/auction/client"
auctionrest "github.com/zar-network/zar-network/x/auction/client/rest"
cdpclient "github.com/zar-network/zar-network/x/cdp/client"
cdprest "github.com/zar-network/zar-network/x/cdp/client/rest"
liquidatorclient "github.com/zar-network/zar-network/x/liquidator/client"
liquidatorrest "github.com/zar-network/zar-network/x/liquidator/client/rest"
priceclient "github.com/zar-network/zar-network/x/pricefeed/client"
pricerest "github.com/zar-network/zar-network/x/pricefeed/client/rest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/zar-network/zar-network/app"
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("zar", "zarp")
	config.SetBech32PrefixForValidator("zva", "zvap")
	config.SetBech32PrefixForConsensusNode("zca", "zcap")
	config.Seal()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	rootCmd := &cobra.Command{
		Use:   "zarcli",
		Short: "Command line interface for interacting with zard",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		version.Cmd,
		client.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		client.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		client.LineBreak,
		nftcmd.GetQueryCmd(nft.StoreKey, cdc),
		client.LineBreak,
		issuecmd.QueryCmd(cdc),
		issuecmd.QueryIssueCmd(cdc),
		issuecmd.QueryFreezeCmd(cdc),
		issuecmd.QueryIssuesCmd(cdc),
		issuecmd.QueryParamsCmd(cdc),
		issuecmd.QueryFreezesCmd(cdc),
		issuecmd.QueryAllowanceCmd(cdc),
		issuecmd.QuerySearchIssuesCmd(cdc),
		client.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		client.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		nftcmd.GetTxCmd(nft.StoreKey, cdc),
		issuecmd.GetCmdIssueBurn(cdc),
		issuecmd.GetCmdIssueFreeze(cdc),
		issuecmd.GetCmdIssueApprove(cdc),
		issuecmd.GetCmdIssueBurnFrom(cdc),
		issuecmd.GetCmdIssueSendFrom(cdc),
		issuecmd.GetCmdIssueUnFreeze(cdc),
		issuecmd.GetCmdIssueDecreaseApproval(cdc),
		issuecmd.GetCmdIssueIncreaseApproval(cdc),
		authcmd.GetDecodeCommand(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)

	issuerest.RegisterRoutes(rs.CliCtx, rs.Mux)
	nftrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.CliCtx.Codec, nft.StoreKey)
	pricerest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.CliCtx.Codec, pricefeed.StoreKey)
auctionrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.CliCtx.Codec)
cdprest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.CliCtx.Codec)
liquidatorrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.CliCtx.Codec)

	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

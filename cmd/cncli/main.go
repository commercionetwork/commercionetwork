package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"commercio-network"

	// CommercioAUTH
	authclient "commercio-network/x/commercioauth/client"
	authrest "commercio-network/x/commercioauth/client/rest"

	// CommercioDOCS
	docsclient "commercio-network/x/commerciodocs/client"
	docsrest "commercio-network/x/commerciodocs/client/rest"

	// CommercioID
	idclient "commercio-network/x/commercioid/client"
	idrest "commercio-network/x/commercioid/client/rest"

	sdk "github.com/cosmos/cosmos-sdk/types"
	//at "github.com/cosmos/cosmos-sdk/x/auth"
	auth "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	dist "github.com/cosmos/cosmos-sdk/x/distribution/client/rest"
	gv "github.com/cosmos/cosmos-sdk/x/gov"
	gov "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	sl "github.com/cosmos/cosmos-sdk/x/slashing"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/client/rest"
	st "github.com/cosmos/cosmos-sdk/x/staking"
	staking "github.com/cosmos/cosmos-sdk/x/staking/client/rest"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	distcmd "github.com/cosmos/cosmos-sdk/x/distribution"
	distClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	govClient "github.com/cosmos/cosmos-sdk/x/gov/client"
	slashingClient "github.com/cosmos/cosmos-sdk/x/slashing/client"
	stakingClient "github.com/cosmos/cosmos-sdk/x/staking/client"
	//_ "github.com/cosmos/cosmos-sdk/client/lcd/statik"
)

const (
	storeAcc  = "acc"
	storeAUTH = "commercioauth"
	storeID   = "commercioid"
	storeDOCS = "commerciodocs"
)

var defaultCLIHome = os.ExpandEnv("$HOME/.cncli")

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	// Instantiate the codec for the command line application
	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.Seal()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	// Module clients hold cli commands (tx,query) and lcd routes
	// TODO: Make the lcd command take a list of ModuleClient
	mc := []sdk.ModuleClients{
		govClient.NewModuleClient(gv.StoreKey, cdc),
		distClient.NewModuleClient(distcmd.StoreKey, cdc),
		stakingClient.NewModuleClient(st.StoreKey, cdc),
		slashingClient.NewModuleClient(sl.StoreKey, cdc),

		// CommercioAUTH
		authclient.NewModuleClient(storeAUTH, cdc),

		// CommercioID
		idclient.NewModuleClient(storeID, cdc),

		// CommercioDOCS
		docsclient.NewModuleClient(storeDOCS, cdc),
	}

	rootCmd := &cobra.Command{
		Use:   "cncli",
		Short: "Commercio.network client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Set the app version
	version.Version = app.Version

	// Build root command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(defaultCLIHome),
		queryCmd(cdc, mc),
		txCmd(cdc, mc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
		client.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with CN
	executor := cli.PrepareMainCmd(rootCmd, "CN", defaultCLIHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

// registerRoutes registers the routes from the different modules for the LCD.
func registerRoutes(rs *lcd.RestServer) {
	//rs.CliCtx = rs.CliCtx.WithAccountDecoder(rs.Cdc)
	rpc.RegisterRoutes(rs.CliCtx, rs.Mux)
	tx.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	auth.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeAcc)
	bank.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	dist.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, distcmd.StoreKey)
	staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	slashing.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	gov.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)

	// CommercioAUTH
	authrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeAUTH)

	// CommercioID
	idrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeID)

	// CommercioDOCS
	docsrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeDOCS)
}

func queryCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		tx.SearchTxCmd(cdc),
		tx.QueryTxCmd(cdc),
		client.LineBreak,
		authcmd.GetAccountCmd(storeAcc, cdc),
		//authcmd.GetAccountCmd(at.StoreKey, cdc),
	)

	for _, m := range mc {
		queryCmd.AddCommand(m.GetQueryCmd())
	}

	return queryCmd
}

func txCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		client.LineBreak,
		authcmd.GetAccountCmd(storeAcc, cdc),
		authcmd.GetMultiSignCommand(cdc),
		tx.GetBroadcastCommand(cdc),
		authcmd.GetSignCommand(cdc),
		//tx.GetBroadcastCommand(cdc),
		//authcmd.GetBroadcastCommand(cdc),

		// RECHECK THIS POINT: not sure!!!! Marco
		// tx.GetBroadcastCommand(cdc),
		//tx.GetEncodeCommand(cdc),

		client.LineBreak,
	)

	for _, m := range mc {
		txCmd.AddCommand(m.GetTxCmd())
	}

	return txCmd
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

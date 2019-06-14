package main

import (
	"commercio-network"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/version"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	cnInit "commercio-network/init"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// cnd custom flags
const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	rootCmd := &cobra.Command{
		Use:               "cnd",
		Short:             "Commercio.network app daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// Set the app version
	version.Version = app.Version

	// Build root command
	rootCmd.AddCommand(cnInit.InitCmd(ctx, cdc))
	rootCmd.AddCommand(cnInit.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(cnInit.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(cnInit.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(cnInit.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(cnInit.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "CN", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		1, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewCommercioNetworkApp(
		logger, db, traceStore, true,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	)
}

// Substitutions of old function appExporter()
func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		dapp := app.NewCommercioNetworkApp(logger, db, traceStore, false)
		err := dapp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return dapp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	dapp := app.NewCommercioNetworkApp(logger, db, traceStore, true)
	return dapp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

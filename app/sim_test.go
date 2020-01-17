package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	simapp.GetSimulatorFlags()
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

func testAndRunTxs(app *CommercioNetworkApp, config simulation.Config) []simulation.WeightedOperation {
	// TODO: fill this with module simulations messages
	ap := make(simulation.AppParams)

	paramChanges := app.sm.GenerateParamChanges(config.Seed)

	_ = paramChanges

	if config.ParamsFile != "" {
		bz, err := ioutil.ReadFile(config.ParamsFile)
		if err != nil {
			panic(err)
		}

		app.cdc.MustUnmarshalJSON(bz, &ap)
	}

	// nolint: govet
	return []simulation.WeightedOperation{}
}

func TestFullAppSimulation(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	var logger log.Logger
	config := simapp.NewConfigFromFlags()

	if simapp.FlagVerboseValue {
		logger = log.TestingLogger()
	} else {
		logger = log.NewNopLogger()
	}

	var db dbm.DB
	dir, _ := ioutil.TempDir("", "goleveldb-app-sim")
	db, _ = sdk.NewLevelDB("Simulation", dir)

	defer func() {
		db.Close()
		_ = os.RemoveAll(dir)
	}()

	commercioNetworkApp := NewCommercioNetworkApp(logger, db, nil, true, simapp.FlagPeriodValue, fauxMerkleModeOpt)
	require.Equal(t, "Commercio.network", commercioNetworkApp.Name())

	// Run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		t, os.Stdout, commercioNetworkApp.BaseApp, simapp.AppStateFn(commercioNetworkApp.cdc, commercioNetworkApp.sm),
		testAndRunTxs(commercioNetworkApp, config), commercioNetworkApp.ModuleAccountAddrs(), config,
	)

	// export state and params before the simulation error is checked
	if config.ExportStatePath != "" {
		err := ExportStateToJSON(commercioNetworkApp, config.ExportStatePath)
		require.NoError(t, err)
	}

	if config.ExportParamsPath != "" {
		err := simapp.ExportParamsToJSON(simParams, config.ExportParamsPath)
		require.NoError(t, err)
	}

	require.NoError(t, simErr)

	if config.Commit {
		// for memdb:
		// fmt.Println("Database Size", db.Stats()["database.size"])
		fmt.Println("\nGoLevelDB Stats")
		fmt.Println(db.Stats()["leveldb.stats"])
		fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
	}
}

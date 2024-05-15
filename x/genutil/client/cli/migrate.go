package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	v300 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v3.0.0"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/cometbft/cometbft/types"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	cometjson "github.com/cometbft/cometbft/libs/json"
)

var migrationMap = map[string][]extypes.MigrationCallback{
	"v3.0.0": {v300.Migrate},
}

const (
	flagGenesisTime   = "genesis-time"
	flagChainID       = "chain-id"
	flagInitialHeight = "initial-height"
)

func MigrationsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrations-list",
		Short: "Lists all the available migrations",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			var migrations []string
			for migration := range migrationMap {
				migrations = append(migrations, migration)
			}

			sort.Strings(migrations)
			for _, m := range migrations {
				fmt.Println(m)
			}

			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "Override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "Override chain_id with this flag")

	return cmd
}

func MigrateGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Please note that migrations should be only be done sequentially. As a reference, suppose we have the following versions:
- v2.2.0
- v3.0.0

If you want to migrate from version v2.0.0 to v3.0.0, you need to execute two migrations:
1. From v2.2.0 to v3.0.0
   $ %s migrate v3.0.0 ...

To see get a full list of available migrations, use the migrations-list command.

Example:
$ %s migrate v3.0.0 /path/to/genesis.json --chain-id=commercio-testnetXXXX --genesis-time=2019-04-22T17:00:00Z --initial-height 1234
`, version.AppName, version.AppName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var err error
			target := args[0]
			importGenesis := args[1]

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return err
			}

			var initialState extypes.AppMap
			if err := json.Unmarshal(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			migrations := migrationMap[target]
			if migrations == nil {
				return fmt.Errorf("unknown migration function version: %s", target)
			}

			newGenState := initialState
			for _, migration := range migrations {
				newGenState = migration(newGenState, clientCtx)
			}

			genDoc.AppState, err = json.Marshal(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime, _ := cmd.Flags().GetString(flagGenesisTime)
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return err
				}

				genDoc.GenesisTime = t
			}

			chainID, _ := cmd.Flags().GetString(flagChainID)
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			initialHeight, _ := cmd.Flags().GetInt(flagInitialHeight)
			genDoc.InitialHeight = int64(initialHeight)

			bz, err := cometjson.Marshal(genDoc)
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			fmt.Println(string(sortedBz))
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "Override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "Override chain_id with this flag")
	cmd.Flags().Int(flagInitialHeight, 0, "Override intial height with this flag")

	return cmd
}

package cli

import (
	"fmt"
	"sort"
	"time"

	v038 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v0_38"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"

	v120 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.2.0"
	v121 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.2.1"
	v130 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.3.0"
	v131 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.3.1"
	v132 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.3.2"
	v133 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.3.3"
	v134 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.3.4"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"
)

var migrationMap = map[string][]extypes.MigrationCallback{
	"v1.2.0": {v038.Migrate, v120.Migrate},
	"v1.2.1": {v121.Migrate},
	"v1.3.0": {v130.Migrate},
	"v1.3.1": {v131.Migrate},
	"v1.3.2": {v132.Migrate},
	"v1.3.3": {v133.Migrate},
	"v1.3.4": {v134.Migrate},
}

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

func MigrationsListCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrations-list",
		Short: "Lists all the available migrations",
		Args:  cobra.ExactArgs(0),
		RunE:  migrationListCmdFunc,
	}

	cmd.Flags().String(flagGenesisTime, "", "Override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "Override chain_id with this flag")

	return cmd
}

func migrationListCmdFunc(cmd *cobra.Command, args []string) error {

	var migrations []string
	for migration := range migrationMap {
		migrations = append(migrations, migration)
	}

	sort.Strings(migrations)
	for _, m := range migrations {
		fmt.Println(m)
	}

	return nil
}

func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Please note that migrations should be only be done sequentially. As a reference, suppose we have the following versions:
- v1.2.0
- v1.2.1
- v1.3.0

If you want to migrate from version v1.2.0 to v1.3.0, you need to execute two migrations:
1. From v1.2.0 to v1.2.1
   $ %s migrate v1.2.1 ...
2. From v1.2.1 to v1.3.0
   $ %s migrate v1.3.0 ...

To see get a full list of available migrations, use the migrations-list command.

Example:
$ %s migrate v1.2.0 /path/to/genesis.json --chain-id=commercio-testnetXXXX --genesis-time=2019-04-22T17:00:00Z
`, version.ServerName, version.ServerName, version.ServerName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrateGenesisCmdFunc(cmd, args, cdc)
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "Override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "Override chain_id with this flag")

	return cmd
}

func migrateGenesisCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	target := args[0]
	importGenesis := args[1]

	genDoc, err := types.GenesisDocFromFile(importGenesis)
	if err != nil {
		return err
	}

	var initialState extypes.AppMap
	cdc.MustUnmarshalJSON(genDoc.AppState, &initialState)

	migrations := migrationMap[target]
	if migrations == nil {
		return fmt.Errorf("unknown migration function version: %s", target)
	}

	newGenState := initialState
	for _, migration := range migrations {
		newGenState = migration(newGenState)
	}

	genDoc.AppState = cdc.MustMarshalJSON(newGenState)

	genesisTime := cmd.Flag(flagGenesisTime).Value.String()
	if genesisTime != "" {
		var t time.Time

		err := t.UnmarshalText([]byte(genesisTime))
		if err != nil {
			return err
		}

		genDoc.GenesisTime = t
	}

	chainID := cmd.Flag(flagChainID).Value.String()
	if chainID != "" {
		genDoc.ChainID = chainID
	}

	out, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(sdk.MustSortJSON(out)))
	return nil
}

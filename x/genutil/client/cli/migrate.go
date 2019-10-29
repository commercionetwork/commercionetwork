package cli

import (
	"fmt"
	"time"

	v038 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v0_38"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"

	v120 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.2.0"
	v130 "github.com/commercionetwork/commercionetwork/x/genutil/legacy/v1.3.0"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"
)

var migrationMap = map[string][]extypes.MigrationCallback{
	"v1.2.0": {v038.Migrate, v120.Migrate},
	"v1.2.1": {v038.Migrate, v120.Migrate},
	"v1.3.0": {v130.Migrate},
	"v1.3.1": {v130.Migrate},
}

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ %s migrate v1.2.0 /path/to/genesis.json --chain-id=commercio-testnetXXXX --genesis-time=2019-04-22T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "Override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "Override chain_id with this flag")

	return cmd
}

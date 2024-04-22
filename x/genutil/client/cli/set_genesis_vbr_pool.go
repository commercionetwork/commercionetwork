package cli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cometbft/cometbft/libs/cli"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	vbrTypes "github.com/commercionetwork/commercionetwork/x/vbr/types"
)

// SetGenesisVbrPoolAmount returns set-genesis-vbr-pool-amount cobra Command.
func SetGenesisVbrPoolAmount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-vbr-pool-amount [amount]",
		Short: "Sets the given amount as the initial VBR pool inside genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			if coins.Len() > 1 {
				return errors.New("cannot have multiple coins inside the VBR pool")
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			genDoc, err := tmtypes.GenesisDocFromFile(genFile)
			if err != nil {
				return err
			}

			/*var genState map[string]json.RawMessage
			if err = json.Unmarshal(genDoc.AppState, &genState); err != nil {
				return fmt.Errorf("error unmarshalling genesis doc %s: %s", genFile, err.Error())
			}

			// set pool amount into the app state
			var genStateVbr vbrTypes.GenesisState
			json.Unmarshal(genState[vbrTypes.ModuleName], &genStateVbr)

			genStateVbr.PoolAmount = sdk.NewDecCoinsFromCoins(coins...)

			genesisStateBzVbr, err := tmjson.Marshal(genStateVbr)
			if err != nil {
				return sdkErrors.Wrap(err, "failed to marshal genesis doc")
			}
			genState[vbrTypes.ModuleName] = genesisStateBzVbr*/

			genState, err := SetVbrPoolAmount(genDoc.AppState, coins)
			if err != nil {
				return err
			}

			appStateJSON, err := json.Marshal(genState)
			if err != nil {
				return err
			}

			// export app state
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	//cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	//cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	return cmd
}
func SetVbrPoolAmount(appState json.RawMessage, coins sdk.Coins) (map[string]json.RawMessage, error) {
	var genState map[string]json.RawMessage
	if err := json.Unmarshal(appState, &genState); err != nil {
		return genState, fmt.Errorf("error unmarshalling genesis doc for vbr: %s", err.Error())
	}

	// set pool amount into the app state
	var genStateVbr vbrTypes.GenesisState
	json.Unmarshal(genState[vbrTypes.ModuleName], &genStateVbr)
	genStateVbr.PoolAmount = sdk.NewDecCoinsFromCoins(coins...)

	genesisStateBzVbr, err := tmjson.Marshal(genStateVbr)
	if err != nil {
		return genState, errorsmod.Wrap(err, "failed to marshal genesis doc")
	}
	genState[vbrTypes.ModuleName] = genesisStateBzVbr

	return genState, nil

}

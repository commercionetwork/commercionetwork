package cli

import (
	"errors"

	"github.com/commercionetwork/commercionetwork/x/tbr"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// SetGenesisTbrPoolAmount returns add-genesis-minter cobra Command.
func SetGenesisTbrPoolAmount(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-tbr-pool-amount [amount]",
		Short: "Sets the given amount as the initial TBR pool inside genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			if coins.Len() > 1 {
				return errors.New("cannot have multiple coins inside the TBR pool")
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// set pool amount into the app state
			var genState tbr.GenesisState
			cdc.MustUnmarshalJSON(appState[tbr.ModuleName], &genState)
			genState.PoolAmount = sdk.NewDecCoins(coins)
			genState.RewardDenom = coins[0].Denom

			genesisStateBz := cdc.MustMarshalJSON(genState)
			appState[tbr.ModuleName] = genesisStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			// export app state
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	return cmd
}

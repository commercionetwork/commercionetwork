package cli

import (
	"errors"

	vbrTypes "github.com/commercionetwork/commercionetwork/x/vbr/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/commercionetwork/commercionetwork/x/vbr"
)

// SetGenesisVbrPoolAmount returns set-genesis-tbr-pool-amount cobra Command.
func SetGenesisVbrPoolAmount(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-vbr-pool-amount [amount]",
		Short: "Sets the given amount as the initial VBR pool inside genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			if coins.Len() > 1 {
				return errors.New("cannot have multiple coins inside the VBR pool")
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// set pool amount into the app state
			var genState vbr.GenesisState
			cdc.MustUnmarshalJSON(appState[vbrTypes.ModuleName], &genState)
			genState.PoolAmount = sdk.NewDecCoinsFromCoins(coins...)

			genesisStateBz := cdc.MustMarshalJSON(genState)
			appState[vbrTypes.ModuleName] = genesisStateBz

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

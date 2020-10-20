package cli

import (
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/commercionetwork/commercionetwork/x/government"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// SetGenesisTumblerAddressCmd returns set-genesis-tumbler-address cobra Command.
func SetGenesisTumblerAddressCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-tumbler-address [tumbler_address_or_key]",
		Short: "Sets the given address as the tumbler address inside genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			address, err := getAddressFromString(args[0])
			if err != nil {
				return err
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// add minter to the app state
			var genState government.GenesisState
			cdc.MustUnmarshalJSON(appState[governmentTypes.ModuleName], &genState)

			genState.TumblerAddress = address

			genesisStateBz := cdc.MustMarshalJSON(genState)
			appState[governmentTypes.ModuleName] = genesisStateBz

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

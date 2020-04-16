package cli

import (
	"encoding/json"
	"fmt"
	"github.com/commercionetwork/commercionetwork/x/memberships"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// SetGenesisGovernmentAddressCmd returns set-genesis-government-address cobra Command.
func SetGenesisGovernmentAddressCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-government-address [government_address_or_key]",
		Short: "Sets the given address as the government address inside genesis.json",
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
			var genStateGovernment government.GenesisState
			cdc.MustUnmarshalJSON(appState[government.ModuleName], &genStateGovernment)

			if !genStateGovernment.GovernmentAddress.Empty() {
				return fmt.Errorf("cannot replace existing government address")
			}

			genStateGovernment.GovernmentAddress = address

			genesisStateBzGovernment := cdc.MustMarshalJSON(genStateGovernment)
			appState[government.ModuleName] = genesisStateBzGovernment

			// set a black membership to the government address
			// add a membership to the genesis state
			var genStateMemberships memberships.GenesisState
			err = json.Unmarshal(appState[memberships.ModuleName], &genStateMemberships)
			if err != nil {
				return err
			}

			membership := memberships.NewMembership("black", address)
			genStateMemberships.Memberships, _ = genStateMemberships.Memberships.AppendIfMissing(membership)

			genesisStateBzMemberships := cdc.MustMarshalJSON(genStateMemberships)
			appState[memberships.ModuleName] = genesisStateBzMemberships

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

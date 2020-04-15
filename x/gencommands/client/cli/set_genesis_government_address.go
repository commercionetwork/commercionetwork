package cli

import (
	"encoding/json"
	"fmt"

	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/commercionetwork/commercionetwork/x/memberships"
	membershipsTypes "github.com/commercionetwork/commercionetwork/x/memberships/types"

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
		Short: "Sets the given address as the government address inside genesis.json, and assings a black membership to it",
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
			cdc.MustUnmarshalJSON(appState[governmentTypes.ModuleName], &genStateGovernment)

			if !genStateGovernment.GovernmentAddress.Empty() {
				return fmt.Errorf("cannot replace existing government address")
			}

			genStateGovernment.GovernmentAddress = address

			genesisStateBzGovernment := cdc.MustMarshalJSON(genStateGovernment)
			appState[governmentTypes.ModuleName] = genesisStateBzGovernment

			// set a black membership to the government address
			// add a membership to the genesis state
			var genStateMemberships memberships.GenesisState
			err = json.Unmarshal(appState[membershipsTypes.ModuleName], &genStateMemberships)
			if err != nil {
				return err
			}

			membership := membershipsTypes.NewMembership("black", address)
			genStateMemberships.Memberships, _ = genStateMemberships.Memberships.AppendIfMissing(membership)

			genesisStateBzMemberships := cdc.MustMarshalJSON(genStateMemberships)
			appState[membershipsTypes.ModuleName] = genesisStateBzMemberships

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

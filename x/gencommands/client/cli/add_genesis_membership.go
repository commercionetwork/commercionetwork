package cli

import (
	"encoding/json"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/memberships"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

// AddGenesisTspCmd returns add-genesis-tsp cobra Command.
func AddGenesisMembershipCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-membership [membership_type] [account_address_or_key]",
		Short: "Creates a new membership of the specified type and associates it to the given address, saving it inside the genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			address, err := getAddressFromString(args[1])
			if err != nil {
				return err
			}

			membershipType := args[0]
			if !memberships.IsMembershipTypeValid(membershipType) {
				return fmt.Errorf("invalid membership type: %s", membershipType)
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// add a membership to the genesis state
			var genState memberships.GenesisState
			err = json.Unmarshal(appState[memberships.ModuleName], &genState)
			if err != nil {
				return err
			}

			membership := memberships.NewMembership(membershipType, address)
			genState.Memberships, _ = genState.Memberships.AppendIfMissing(membership)

			// save the state
			genesisStateBz := cdc.MustMarshalJSON(genState)
			appState[memberships.ModuleName] = genesisStateBz

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

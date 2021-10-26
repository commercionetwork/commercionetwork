package cli

import (
	"encoding/json"
	"fmt"

	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
)

// SetGenesisGovernmentAddressCmd returns set-genesis-government-address cobra Command.
func SetGenesisGovernmentAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-government-address [government_address_or_key]",
		Short: "Sets the given address as the government address inside genesis.json, and assings a black membership to it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			genDoc, err := tmtypes.GenesisDocFromFile(genFile)
			if err != nil {
				return err
			}

			var genState map[string]json.RawMessage
			if err = json.Unmarshal(genDoc.AppState, &genState); err != nil {
				return fmt.Errorf("error unmarshalling genesis doc %s: %s", genFile, err.Error())
			}

			// add minter to the app state
			var genStateGovernment govTypes.GenesisState
			json.Unmarshal(genState[govTypes.ModuleName], &genStateGovernment)

			if genStateGovernment.GovernmentAddress != "" {
				return fmt.Errorf("cannot replace existing government address")
			}

			genStateGovernment.GovernmentAddress = address.String()

			genesisStateBzGovernment, err := tmjson.Marshal(genStateGovernment)
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}
			genState[govTypes.ModuleName] = genesisStateBzGovernment

			// set a black membership to the government address
			// add a membership to the genesis state
			/*var genStateMemberships commerciokyc.GenesisState
			err = json.Unmarshal(appState[commerciokycTypes.ModuleName], &genStateMemberships)
			if err != nil {
				return err
			}

			initSecondsPerYear := time.Hour * 24 * 365
			initExpirationDate := time.Now().Add(initSecondsPerYear) // It's safe becouse command is executed in one machine

			membership := commerciokycTypes.NewMembership(commerciokycTypes.MembershipTypeBlack, address, address, initExpirationDate)
			genStateMemberships.Memberships, _ = genStateMemberships.Memberships.AppendIfMissing(membership)

			genesisStateBzMemberships := cdc.MustMarshalJSON(genStateMemberships)
			appState[commerciokycTypes.ModuleName] = genesisStateBzMemberships
			*/
			genDoc.AppState, err = json.Marshal(genState)

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	//cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	//cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	return cmd
}

// getAddressFromString reads the given value as an AccAddress object, retuning an error if
// the specified value is not a valid address
/*func getAddressFromString(value string) (sdk.AccAddress, error) {
	minterAddr, err := sdk.AccAddressFromBech32(value)
	if err != nil {
		kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
		if err != nil {
			return nil, err
		}

		info, err := kb.Get(value)
		if err != nil {
			return nil, err
		}

		minterAddr = info.GetAddress()
	}

	return minterAddr, nil
}*/

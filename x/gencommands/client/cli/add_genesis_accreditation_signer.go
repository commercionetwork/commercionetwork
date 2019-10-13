package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/accreditations"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// SetGenesisGovernmentAddressCmd returns add-genesis-minter cobra Command.
func AddGenesisAccreditationSigner(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-accreditation-signer [signer_address_or_key]",
		Short: "Add a trusted accreditation signer to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			minterAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return err
				}

				minterAddr = info.GetAddress()
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// add minter to the app state
			var genState accreditations.GenesisState

			cdc.MustUnmarshalJSON(appState[accreditations.ModuleName], &genState)

			if genState.TrustedServiceProviders.Contains(minterAddr) {
				return fmt.Errorf("cannot add already existing signer %v", minterAddr)
			}

			genState.TrustedServiceProviders, _ = genState.TrustedServiceProviders.AppendIfMissing(minterAddr)

			genesisStateBz := cdc.MustMarshalJSON(genState)
			appState[accreditations.ModuleName] = genesisStateBz

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

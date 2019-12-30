package cli

import (
	"github.com/commercionetwork/commercionetwork/x/memberships"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// AddGenesisTspCmd returns add-genesis-tsp cobra Command.
func AddGenesisTspCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-tsp [tsp_address_or_key]",
		Short: "Add a trusted accreditation signer to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cm *cobra.Command, args []string) error {
			return addGenesisTspCmdFunc(cm, args, ctx, cdc)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	return cmd
}

func addGenesisTspCmdFunc(_ *cobra.Command, args []string, ctx *server.Context, cdc *codec.Codec) error {
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

	// add tsp to the app state
	var genState memberships.GenesisState
	cdc.MustUnmarshalJSON(appState[memberships.ModuleName], &genState)

	genState.TrustedServiceProviders, _ = genState.TrustedServiceProviders.AppendIfMissing(address)

	// save the app state
	genesisStateBz := cdc.MustMarshalJSON(genState)
	appState[memberships.ModuleName] = genesisStateBz

	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		return err
	}

	// Export app state
	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

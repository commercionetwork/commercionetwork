package cli

import (
	"fmt"
	pricefeedTypes "github.com/commercionetwork/commercionetwork/x/pricefeed/types"

	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetGenesisPriceCmd returns set-genesis-price cobra Command.
func SetGenesisPriceCmd(ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-price [asset_name] [value] [expiration_height] [oracle]",
		Short: "Adds a new raw price for the given asset having the specified value, expiration block height and oracle",
		Long: `Adda a new raw price having the given value and expiration block height to the list of raw prices 
for the specified asset. 
Also adds the specified address as a valid oracle and the given token name as a supported asset. 
`,
		Args: cobra.ExactArgs(4),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			value, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return fmt.Errorf("invalid price value, %s", args[1])
			}

			expiry, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid expiration height, %s", args[2])
			}

			oracle, err2 := getAddressFromString(args[3])
			if err2 != nil {
				return err
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err2 := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err2 != nil {
				return err
			}

			// create the price object
			price := pricefeedTypes.Price{AssetName: args[0], Value: value, Expiry: expiry}

			// add the price to the app state
			var genState pricefeed.GenesisState
			cdc.MustUnmarshalJSON(appState[pricefeedTypes.ModuleName], &genState)

			// save the raw price, the asset name and the oracle
			rawPrice := pricefeedTypes.OraclePrice{Oracle: oracle, Price: price, Created: sdk.ZeroInt()}
			genState.RawPrices, _ = genState.RawPrices.UpdatePriceOrAppendIfMissing(rawPrice)
			genState.Assets, _ = genState.Assets.AppendIfMissing(price.AssetName)
			genState.Oracles, _ = genState.Oracles.AppendIfMissing(oracle)

			genesisStateBz := cdc.MustMarshalJSON(genState)
			appState[pricefeedTypes.ModuleName] = genesisStateBz

			appStateJSON, err2 := cdc.MarshalJSON(appState)
			if err2 != nil {
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

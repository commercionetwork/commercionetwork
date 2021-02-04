package cli

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/commercionetwork/commercionetwork/x/vbr"
	vbrTypes "github.com/commercionetwork/commercionetwork/x/vbr/types"
)

// SetGenesisVbrPoolAmount returns set-genesis-vbr-pool-amount cobra Command.
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

// SetGenesisVbrRewardRate returns set-genesis-vbr-reward-rate cobra Command.
func SetGenesisVbrRewardRate(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-genesis-vbr-reward-rate [rate]",
		Short: "Sets the given value as the initial VBR reward rate inside genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			value, err := sdk.NewDecFromStr(args[0])
			if err != nil {
				return err
			}

			if value.IsZero() || value.IsNegative() {
				return errors.New("cannot have zero or negative value of reward rate")
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
			genState.RewardRate = value

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

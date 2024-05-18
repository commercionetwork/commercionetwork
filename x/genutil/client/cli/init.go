package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/cometbft/cometbft/types"

	"github.com/cosmos/go-bip39"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	gencli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cometbft/cometbft/libs/cli"
	cometos "github.com/cometbft/cometbft/libs/os"
	cometrand "github.com/cometbft/cometbft/libs/rand"

	cfg "github.com/cometbft/cometbft/config"
)

const (
	// FlagGovAddr defines a flag to set Government address.
	FlagGovAddr = "government-address"

	// FlagVbrPoolAmount defines a flag to initialize the vbr-pool-amount
	FlagVbrPoolAmount = "vbr-pool-amount"
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))

	return err
}

func InitCmd(mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", cometrand.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			recover, _ := cmd.Flags().GetBool(gencli.FlagRecover)
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(gencli.FlagOverwrite)

			if !overwrite && cometos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}
			appState, err := json.MarshalIndent(mbm.DefaultGenesis(cdc), "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return errors.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			// Verify government address flag
			govAddress, _ := cmd.Flags().GetString(FlagGovAddr)
			if govAddress != "" {
				address, err := sdk.AccAddressFromBech32(govAddress)
				if err != nil {
					return err
				}

				genState, err := SetGovernmentAddress(clientCtx, appState, address)
				if err != nil {
					return err
				}
				appState, err = json.Marshal(genState)
				if err != nil {
					return err
				}
			}

			// Verify vbr pool amount flag
			vbrPoolAmount, _ := cmd.Flags().GetString(FlagVbrPoolAmount)
			if vbrPoolAmount != "" {
				coins, err := sdk.ParseCoinsNormalized(vbrPoolAmount)
				if err != nil {
					return err
				}
				genState, err := SetVbrPoolAmount(appState, coins)
				if err != nil {
					return err
				}
				appState, err = json.Marshal(genState)
				if err != nil {
					return err
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState
			if _, err = genutiltypes.AppGenesisFromFile(genFile); err != nil {
				return errors.Wrap(err, "Failed to export gensis file")
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(gencli.FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(gencli.FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagGovAddr, "", "genesis government address, if left blank will need to set with set-government-address command")
	cmd.Flags().String(FlagVbrPoolAmount, "", "genesis vbr pool amount, if left blank will need to set with set-genesis-vbr-pool-amount command")

	return cmd
}

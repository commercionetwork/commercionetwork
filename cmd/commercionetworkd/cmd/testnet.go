package cmd

// DONTCOVER

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	commerciokycTypes "github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	vbrTypes "github.com/commercionetwork/commercionetwork/x/vbr/types"

	//epochsTypes "github.com/commercionetwork/commercionetwork/x/epochs/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"

	"github.com/commercionetwork/commercionetwork/app"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"

	//"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	//authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	tmconfig "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagNodeDirPrefix     = "node-dir-prefix"
	flagNumValidators     = "v"
	flagOutputDir         = "output-dir"
	flagNodeDaemonHome    = "node-daemon-home"
	flagNodeCLIHome       = "node-cli-home"
	flagStartingIPAddress = "starting-ip-address"
)

// get cmd to initialize all files for tendermint testnet and application
func testnetCmd(mbm module.BasicManager, genBalIterator bankTypes.GenesisBalancesIterator) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a Commercio testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	cnd testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			outputDir := viper.GetString(flagOutputDir)
			chainID := viper.GetString(flags.FlagChainID)
			minGasPrices := viper.GetString(server.FlagMinGasPrices)
			nodeDirPrefix := viper.GetString(flagNodeDirPrefix)
			nodeDaemonHome := viper.GetString(flagNodeDaemonHome)
			nodeCLIHome := viper.GetString(flagNodeCLIHome)
			startingIPAddress := viper.GetString(flagStartingIPAddress)
			numValidators := viper.GetInt(flagNumValidators)

			return InitTestnet(clientCtx, cmd, config, mbm, genBalIterator, outputDir, chainID,
				minGasPrices, nodeDirPrefix, nodeDaemonHome, nodeCLIHome, startingIPAddress, numValidators)
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "cnd",
		"Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeCLIHome, "cncli",
		"Home directory of the node's cli configuration")
	cmd.Flags().String(flagStartingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(
		flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(
		//server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", app.DefaultBondDenom),
		//"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
		server.FlagMinGasPrices, fmt.Sprintf(""),
		"No gas prices need setup")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")

	return cmd
}

const nodeDirPerm = 0755

// Initialize the testnet
func InitTestnet(
	clientCtx client.Context,
	cmd *cobra.Command,
	config *tmconfig.Config,
	mbm module.BasicManager,
	genBalIterator bankTypes.GenesisBalancesIterator,
	outputDir,
	chainID,
	minGasPrices,
	nodeDirPrefix,
	nodeDaemonHome,
	nodeCLIHome,
	startingIPAddress string,
	numValidators int,
) error {

	if chainID == "" {
		chainID = "chain-" + tmrand.NewRand().Str(6)
	}

	monikers := make([]string, numValidators)
	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]cryptotypes.PubKey, numValidators)

	cndConfig := srvconfig.DefaultConfig()
	cndConfig.MinGasPrices = minGasPrices

	//nolint:prealloc
	var (
		genBalances []bankTypes.Balance
		genAccounts []authTypes.GenesisAccount
		genFiles    []string
	)

	inBuf := bufio.NewReader(cmd.InOrStdin())
	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		clientDir := filepath.Join(outputDir, nodeDirName, nodeCLIHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")

		config.SetRoot(nodeDir)
		config.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		if err := os.MkdirAll(clientDir, nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		monikers = append(monikers, nodeDirName)
		config.Moniker = nodeDirName

		ip, err := getIP(i, startingIPAddress)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(config)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, config.GenesisFile())

		kb, err := keyring.New(
			sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend),
			clientDir,
			inBuf)
		if err != nil {
			return err
		}

		/*kb, err := keys.NewKeyring(
			sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend),
			clientDir,
			inBuf,
		)
		if err != nil {
			return err
		}*/

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
		if err != nil {
			return err
		}

		//keyPass := clientkeys.DefaultKeyPass
		addr, secret, err := server.GenerateSaveCoinKey(kb, nodeDirName, true, algo)
		//addr, secret, err := server.GenerateSaveCoinKeyFromPath(kb, nodeDirName, true, algo, app.FullFundraiserPath)

		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint); err != nil {
			return err
		}

		accTokens := sdk.TokensFromConsensusPower(10000000)
		accStakingTokens := sdk.TokensFromConsensusPower(10000000)
		coins := sdk.Coins{
			sdk.NewCoin(app.StableCreditDenom, accTokens),
			sdk.NewCoin(app.StakeDenom, accStakingTokens),
		}
		genBalances = append(genBalances, bankTypes.Balance{Address: addr.String(), Coins: coins.Sort()})
		genAccounts = append(genAccounts, authTypes.NewBaseAccount(addr, nil, 0, 0))

		valTokens := sdk.TokensFromConsensusPower(100)
		msg, err := stakingTypes.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(app.StakeDenom, valTokens),
			stakingTypes.NewDescription(nodeDirName, "", "", "", ""),
			stakingTypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec()),
			sdk.OneInt(),
		)
		if err != nil {
			return err
		}

		//tx := authTypes.NewStdTx([]sdk.Msg{msg}, auth.StdFee{}, []auth.StdSignature{}, memo)
		//txBldr := authTypes.NewTxBuilderFromCLI(inBuf).WithChainID(chainID).WithMemo(memo).WithKeybase(kb)

		txBuilder := clientCtx.TxConfig.NewTxBuilder()
		if err := txBuilder.SetMsgs(msg); err != nil {
			return err
		}
		txBuilder.SetMemo(memo)
		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(chainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(clientCtx.TxConfig)

		if err := tx.Sign(txFactory, nodeDirName, txBuilder, true); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		txBytes, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		/*signedTx, err := txBuilder.SignStdTx(nodeDirName, clientkeys.DefaultKeyPass, tx, false)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		txBytes, err := cdc.MarshalJSON(signedTx)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}*/

		// gather gentxs folder
		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBytes); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		// TODO: Rename config file to server.toml as it's not particular to Gaia
		// (REF: https://github.com/cosmos/cosmos-sdk/issues/4125).
		cndConfigFilePath := filepath.Join(nodeDir, "config/cnd.toml")
		srvconfig.WriteConfigFile(cndConfigFilePath, cndConfig)
	}

	if err := initGenFiles(clientCtx, mbm, chainID, genBalances, genAccounts, genFiles, numValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		clientCtx, config, chainID, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genBalIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

func initGenFiles(
	clientCtx client.Context,
	mbm module.BasicManager,
	chainID string,
	genBalances []bankTypes.Balance,
	genAccounts []authTypes.GenesisAccount,
	genFiles []string,
	numValidators int,
) error {

	appGenState := mbm.DefaultGenesis(clientCtx.JSONMarshaler)

	// set the accounts in the genesis state
	var authGenState authTypes.GenesisState
	cdc := clientCtx.JSONMarshaler
	cdc.MustUnmarshalJSON(appGenState[authTypes.ModuleName], &authGenState)

	accounts, err := authTypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = accounts
	appGenState[authTypes.ModuleName] = cdc.MustMarshalJSON(&authGenState)

	var bankGenState bankTypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[bankTypes.ModuleName], &bankGenState)

	bankGenState.Balances = bankTypes.SanitizeGenesisBalances(genBalances)
	for _, bal := range bankGenState.Balances {
		bankGenState.Supply = bankGenState.Supply.Add(bal.Coins...)
	}
	appGenState[bankTypes.ModuleName] = cdc.MustMarshalJSON(&bankGenState)

	// commercionetworkd set-genesis-government-address
	var governmentState governmentTypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[governmentTypes.ModuleName], &governmentState)
	governmentState.GovernmentAddress = genAccounts[0].GetAddress().String()
	appGenState[governmentTypes.ModuleName] = cdc.MustMarshalJSON(&governmentState)

	// set-genesis-vbr-pool-amount 1000000000ucommercio
	var vbrState vbrTypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[vbrTypes.ModuleName], &vbrState)
	tokens := sdk.TokensFromConsensusPower(1000)
	vbrState.PoolAmount = sdk.NewDecCoinsFromCoins(sdk.NewCoin(app.StakeDenom, tokens))
	vbrState.Params.EarnRate = sdk.NewDecWithPrec(5, 1)
	vbrState.Params.DistrEpochIdentifier = "minute"
	appGenState[vbrTypes.ModuleName] = cdc.MustMarshalJSON(&vbrState)

	// cnd add-genesis-tsp
	var kycState commerciokycTypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[commerciokycTypes.ModuleName], &kycState)
	tsps := ctypes.Strings{}

	kycState.TrustedServiceProviders, _ = tsps.AppendIfMissing(genAccounts[0].GetAddress().String())
	// cnd add-genesis-membership black
	currentTime := time.Now()
	expiryAt, _ := time.Parse("2006-01-02", currentTime.Format("2006-01-02"))
	expiryAt = expiryAt.Add(time.Hour * 24 * 365)

	membership := commerciokycTypes.NewMembership("black", genAccounts[0].GetAddress(), genAccounts[0].GetAddress(), expiryAt)

	kycState.Memberships = append(kycState.Memberships, &membership)
	genesisStateBz := cdc.MustMarshalJSON(&kycState)
	appGenState[commerciokycTypes.ModuleName] = genesisStateBz

	appGenStateJSON, err := json.MarshalIndent(appGenState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}
	return nil
}

func collectGenFiles(
	clientCtx client.Context, nodeConfig *tmconfig.Config, chainID string,
	nodeIDs []string, valPubKeys []cryptotypes.PubKey, numValidators int,
	outputDir, nodeDirPrefix, nodeDaemonHome string, genBalIterator bankTypes.GenesisBalancesIterator,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")
		nodeConfig.Moniker = nodeDirName

		nodeConfig.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := genutiltypes.NewInitConfig(chainID, gentxsDir, nodeID, valPubKey)

		genDoc, err := types.GenesisDocFromFile(nodeConfig.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(clientCtx.JSONMarshaler, clientCtx.TxConfig, nodeConfig, initCfg, *genDoc, genBalIterator)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := nodeConfig.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func getIP(i int, startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return calculateIP(startingIPAddr, i)
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := tmos.EnsureDir(writePath, 0700)
	if err != nil {
		return err
	}

	err = tmos.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}

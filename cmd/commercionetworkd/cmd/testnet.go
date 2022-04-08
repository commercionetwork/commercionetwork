package cmd

// DONTCOVER

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

	commerciokycTypes "github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	vbrTypes "github.com/commercionetwork/commercionetwork/x/vbr/types"

	"github.com/commercionetwork/commercionetwork/app"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/client/tx"

	tmconfig "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
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
//GenTxCmd(mbm module.BasicManager, txEncCfg client.TxEncodingConfig, genBalIterator types.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {

func testnetCmd(mbm module.BasicManager, genBalIterator banktypes.GenesisBalancesIterator) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a Commercio testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	commercionetworkd testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			minGasPrices, _ := cmd.Flags().GetString(server.FlagMinGasPrices)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)

			return InitTestnet(
				clientCtx, cmd, config, mbm, genBalIterator, outputDir, chainID, minGasPrices,
				nodeDirPrefix, nodeDaemonHome, startingIPAddress, keyringBackend, algo, numValidators,
			)

		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "commercionetwork",
		"Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeCLIHome, "cncli",
		"Home directory of the node's cli configuration")
	cmd.Flags().String(flagStartingIPAddress, "192.168.0.1",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(
		flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(
		server.FlagMinGasPrices, fmt.Sprintf("0.000006%s", app.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")

	return cmd
}

const nodeDirPerm = 0755

// Initialize the testnet
func InitTestnet(
	clientCtx client.Context,
	cmd *cobra.Command,
	nodeConfig *tmconfig.Config,
	mbm module.BasicManager,
	genBalIterator banktypes.GenesisBalancesIterator,
	outputDir,
	chainID,
	minGasPrices,
	nodeDirPrefix,
	nodeDaemonHome,
	startingIPAddress,
	keyringBackend,
	algoStr string,
	numValidators int,
) error {

	// Setup chain-id
	if chainID == "" {
		chainID = "chain-" + tmrand.NewRand().Str(6)
	}

	monikers := make([]string, numValidators)
	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]cryptotypes.PubKey, numValidators)

	// Setup app.toml
	cndConfig := srvconfig.DefaultConfig()
	cndConfig.MinGasPrices = minGasPrices
	cndConfig.API.Enable = true
	cndConfig.API.Swagger = true

	var (
		genAccounts []authtypes.GenesisAccount
		genBalances []banktypes.Balance
		genFiles    []string
	)

	inBuf := bufio.NewReader(cmd.InOrStdin())

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")

		nodeConfig.SetRoot(nodeDir)
		nodeConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		monikers = append(monikers, nodeDirName)
		nodeConfig.Moniker = nodeDirName

		ip, err := getIP(i, startingIPAddress)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(nodeConfig)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, nodeConfig.GenesisFile())

		kb, err := keyring.New(
			sdk.KeyringServiceName(),
			keyringBackend,
			nodeDir,
			inBuf,
		)
		if err != nil {
			return err
		}

		//keyPass := clientkeys.DefaultKeyPass
		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
		if err != nil {
			return err
		}

		addr, secret, err := server.GenerateSaveCoinKey(kb, nodeDirName, true, algo)
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

		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), nodeDir, cliPrint); err != nil {
			return err
		}

		accTokens := sdk.TokensFromConsensusPower(10000000)
		accStakingTokens := sdk.TokensFromConsensusPower(10000000)
		coins := sdk.Coins{
			sdk.NewCoin(app.StableCreditsDenom, accTokens),
			sdk.NewCoin(app.DefaultBondDenom, accStakingTokens),
		}

		genBalances = append(genBalances, banktypes.Balance{Address: addr.String(), Coins: coins.Sort()})
		genAccounts = append(genAccounts, authtypes.NewBaseAccount(addr, nil, 0, 0))

		valTokens := sdk.TokensFromConsensusPower(100)
		createValMsg, err := stakingtypes.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(app.DefaultBondDenom, valTokens),
			stakingtypes.NewDescription(nodeDirName, "", "", "", ""),
			stakingtypes.NewCommissionRates(sdk.OneDec(), sdk.OneDec(), sdk.OneDec()),
			sdk.OneInt(),
		)
		if err != nil {
			return err
		}
		txBuilder := clientCtx.TxConfig.NewTxBuilder()
		if err := txBuilder.SetMsgs(createValMsg); err != nil {
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
			return err
		}

		txBz, err := clientCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return err
		}

		// gather gentxs folder
		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBz); err != nil {
			//_ = os.RemoveAll(outputDir)
			return err
		}

		srvconfig.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), cndConfig)
	}

	if err := initGenFiles(clientCtx, mbm, chainID, genAccounts, genBalances, genFiles, numValidators); err != nil {
		return err
	}

	err := collectGenFiles(
		clientCtx, nodeConfig, chainID, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genBalIterator,
	)

	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

func initGenFiles(
	clientCtx client.Context, mbm module.BasicManager, chainID string,
	genAccounts []authtypes.GenesisAccount, genBalances []banktypes.Balance,
	genFiles []string, numValidators int,
) error {

	cdc := clientCtx.JSONMarshaler

	appGenState := mbm.DefaultGenesis(cdc)

	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[authtypes.ModuleName], &authGenState)
	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}
	authGenState.Accounts = accounts
	appGenState[authtypes.ModuleName] = cdc.MustMarshalJSON(&authGenState)

	// set the balances in the genesis state
	var bankGenState banktypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[banktypes.ModuleName], &bankGenState)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(genBalances)
	for _, bal := range bankGenState.Balances {
		bankGenState.Supply = bankGenState.Supply.Add(bal.Coins...)
	}
	appGenState[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankGenState)

	// setup staking denom
	var stakingState stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[stakingtypes.ModuleName], &stakingState)
	stakingState.Params.BondDenom = app.DefaultBondDenom
	appGenState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingState)

	// setup crisis constant fee denom
	var crisisState crisistypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[crisistypes.ModuleName], &crisisState)
	crisisState.ConstantFee.Denom = app.DefaultBondDenom
	appGenState[crisistypes.ModuleName] = cdc.MustMarshalJSON(&crisisState)

	var govState govtypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[govtypes.ModuleName], &govState)
	govState.DepositParams.MinDeposit[0].Denom = app.DefaultBondDenom
	appGenState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govState)

	// commercionetworkd set-genesis-government-address
	var governmentState governmentTypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[governmentTypes.ModuleName], &governmentState)
	governmentState.GovernmentAddress = genAccounts[0].GetAddress().String()
	appGenState[governmentTypes.ModuleName] = cdc.MustMarshalJSON(&governmentState)

	// set-genesis-vbr-pool-amount 1000000000ucommercio
	var vbrState vbrTypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[vbrTypes.ModuleName], &vbrState)
	tokens := sdk.TokensFromConsensusPower(1000)
	vbrState.PoolAmount = sdk.NewDecCoinsFromCoins(sdk.NewCoin(app.DefaultBondDenom, tokens))
	vbrState.Params.DistrEpochIdentifier = "minute"
	vbrState.Params.EarnRate = sdk.NewDecWithPrec(5, 2)
	appGenState[vbrTypes.ModuleName] = cdc.MustMarshalJSON(&vbrState)

	// commercionetworkd set-genesis-price ucommercio 1 100000000

	// commercionetworkd add-genesis-tsp
	var kycState commerciokycTypes.GenesisState
	cdc.MustUnmarshalJSON(appGenState[commerciokycTypes.ModuleName], &kycState)
	kycState.TrustedServiceProviders = append(kycState.TrustedServiceProviders, genAccounts[0].GetAddress().String())
	// commercionetworkd add-genesis-membership black
	membership := commerciokycTypes.NewMembership("black", genAccounts[0].GetAddress(), genAccounts[0].GetAddress(), time.Now().Add(time.Hour*24*365))
	kycState.Memberships = append(kycState.Memberships, &membership)
	cdc.MustUnmarshalJSON(appGenState[commerciokycTypes.ModuleName], &kycState)

	//appGenStateJSON, err := codec.MarshalJSONIndent(cdc, appGenState)
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
	outputDir, nodeDirPrefix, nodeDaemonHome string, genBalIterator banktypes.GenesisBalancesIterator,
) error {

	var appState json.RawMessage
	genTime := tmtime.Now()
	genFile := ""
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
	// Save genesis in common directory
	base_config := filepath.Join(outputDir, "base_config")
	nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, 0)
	nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)

	genesis_src := filepath.Join(nodeDir, "config", "genesis.json")
	genesis_dest := filepath.Join(base_config, "genesis.json")
	if err := os.MkdirAll(base_config, nodeDirPerm); err != nil {
		return err
	}

	bytesRead, err := ioutil.ReadFile(genesis_src)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(genesis_dest, bytesRead, nodeDirPerm)
	if err != nil {
		return err
	}

	genFile = nodeConfig.GenesisFile()

	gentxsDir := filepath.Join(outputDir, "gentxs")

	initCfg := genutiltypes.NewInitConfig(chainID, gentxsDir, "", nil)
	genDoc, err := types.GenesisDocFromFile(nodeConfig.GenesisFile())
	if err != nil {
		return err
	}

	_, persistentPeers, _ := genutil.CollectTxs(
		clientCtx.JSONMarshaler, clientCtx.TxConfig.TxJSONDecoder(), "", initCfg.GenTxsDir, *genDoc, genBalIterator,
	)

	writeFile("persistent.txt", filepath.Join(outputDir, "base_config"), []byte(persistentPeers))

	if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
		return err
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

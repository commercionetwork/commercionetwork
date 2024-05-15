package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	//"github.com/tendermint/spm/openapiconsole"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"

	// ------------------------------------------
	// Comet base
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	cometos "github.com/cometbft/cometbft/libs/os"

	// ------------------------------------------
	// Cosmwasm module
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	// ------------------------------------------
	// Cosmos SDK utils
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"

	// ------------------------------------------
	// Cosmos SDK modules
	//  Auth
	comosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	//  Authz
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	//  Bank
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	//  Capability
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	//  Crisis
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	//  Distribution
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	//  Evidence
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"

	//  Feegrant
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"

	//  Genesis Utilities
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	//  Governance
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	//  Params
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	//  Slashing
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	//  Staking
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	//  Upgrade
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	cometjson "github.com/cometbft/cometbft/libs/json"

	//  Vesting
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	// ------------------------------------------
	// IBC v7
	//  Transfer
	transfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	ibcaddresslimit "github.com/commercionetwork/commercionetwork/x/ibc-address-limiter"
	ibcaddresslimittypes "github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/types"

	//  Core
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	// ------------------------------------------
	// Commercio.Network
	appparams "github.com/commercionetwork/commercionetwork/app/params"
	docs "github.com/commercionetwork/commercionetwork/docs_3_0"
	"github.com/commercionetwork/commercionetwork/x/ante"

	// ------------------------------------------
	// Commercio.Network Modules
	//  Kyc
	commerciokycmodule "github.com/commercionetwork/commercionetwork/x/commerciokyc"
	commerciokycKeeper "github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
	commerciokycTypes "github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	//  Mint
	commerciomintmodule "github.com/commercionetwork/commercionetwork/x/commerciomint"
	commerciomintKeeper "github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	commerciomintTypes "github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	//  Did
	"github.com/commercionetwork/commercionetwork/x/did"
	didkeeper "github.com/commercionetwork/commercionetwork/x/did/keeper"
	didTypes "github.com/commercionetwork/commercionetwork/x/did/types"

	//  Documents
	"github.com/commercionetwork/commercionetwork/x/documents"
	documentskeeper "github.com/commercionetwork/commercionetwork/x/documents/keeper"
	documentstypes "github.com/commercionetwork/commercionetwork/x/documents/types"

	//  Epochs
	"github.com/commercionetwork/commercionetwork/x/epochs"
	epochskeeper "github.com/commercionetwork/commercionetwork/x/epochs/keeper"
	epochstypes "github.com/commercionetwork/commercionetwork/x/epochs/types"

	//  Government
	governmentmodule "github.com/commercionetwork/commercionetwork/x/government"
	governmentmodulekeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentmoduletypes "github.com/commercionetwork/commercionetwork/x/government/types"

	//  Vbr
	vbrmodule "github.com/commercionetwork/commercionetwork/x/vbr"
	vbrmodulekeeper "github.com/commercionetwork/commercionetwork/x/vbr/keeper"
	vbrmoduletypes "github.com/commercionetwork/commercionetwork/x/vbr/types"
)

const Name = "commercionetwork"

var (
	// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
	ProposalsEnabled = "true"

	// If set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificProposals = ""
	DefaultBondDenom        = "ucommercio"
	StableCreditsDenom      = "uccc"
	Bech32Prefix 			= "did:com"
)

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasmtypes.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasmtypes.EnableAllProposals
		}
		return wasmtypes.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

// ConvertToProposals maps each key to a ProposalType and returns a typed list.
// If any string is not a valid type (in this file), then return an error
func ConvertToProposals(keys []string) ([]wasmtypes.ProposalType, error) {
	valid := make(map[string]bool, len(wasmtypes.EnableAllProposals))
	for _, key := range wasmtypes.EnableAllProposals {
		valid[string(key)] = true
	}

	proposals := make([]wasmtypes.ProposalType, len(keys))
	for i, key := range keys {
		if _, ok := valid[key]; !ok {
			return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "'%s' is not a valid ProposalType", key)
		}
		proposals[i] = wasmtypes.ProposalType(key)
	}
	return proposals, nil
}

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
	)

	return govProposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},

		vbrmodule.AppModuleBasic{},
		governmentmodule.AppModuleBasic{},
		did.AppModuleBasic{},
		documents.AppModuleBasic{},
		commerciokycmodule.AppModuleBasic{},
		commerciomintmodule.AppModuleBasic{},
		wasm.AppModuleBasic{},
		epochs.AppModuleBasic{},
		ibcaddresslimit.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:       nil,
		distrtypes.ModuleName:            nil,
		stakingtypes.BondedPoolName:      {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:   {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:              {authtypes.Burner},
		vbrmoduletypes.ModuleName:        {authtypes.Minter},
		governmentmoduletypes.ModuleName: nil,
		commerciokycTypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
		commerciomintTypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		documentstypes.ModuleName:        nil,
		didTypes.ModuleName:              nil,
		ibctransfertypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		wasmtypes.ModuleName:             {authtypes.Burner},
	}
)

var (
	_ CosmosApp               = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
	_ runtime.AppI			  = (*App)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.BaseKeeper
	AuthzKeeper      authzkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        *govkeeper.Keeper
	CrisisKeeper     *crisiskeeper.Keeper
	UpgradeKeeper    *upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	WasmKeeper       wasmkeeper.Keeper
	ContractKeeper   *wasmkeeper.PermissionedKeeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper

	GovernmentKeeper    governmentmodulekeeper.Keeper
	CommercioMintKeeper commerciomintKeeper.Keeper
	CommercioKycKeeper  commerciokycKeeper.Keeper

	VbrKeeper vbrmodulekeeper.Keeper

	DidKeeper       didkeeper.Keeper
	DocumentsKeeper documentskeeper.Keeper
	// the module manager
	mm           *module.Manager
	EpochsKeeper epochskeeper.Keeper

	// simulation manager
	sm *module.SimulationManager

	AddressLimitingICS4Wrapper *ibcaddresslimit.ICS4Wrapper
	RawIcs20TransferAppModule  transfer.AppModule
	TransferStack              *ibcaddresslimit.IBCModule
}

// Remove assertNoPrefix
func NewKVStoreKeys(names ...string) map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey, len(names))
	for _, n := range names {
		keys[n] = sdk.NewKVStoreKey(n)
	}
	return keys
}

// New returns a reference to an initialized Commercionetwork.
// NewSimApp returns a reference to an initialized SimApp.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	enabledProposals []wasmtypes.ProposalType,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {

	appCodec := encodingConfig.Codec
	legacyAmino := encodingConfig.Amino
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		crisistypes.StoreKey,
		paramstypes.StoreKey,
		consensusparamtypes.StoreKey,
		ibcexported.StoreKey,
		upgradetypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey,
		feegrant.StoreKey,
		authzkeeper.StoreKey,
		wasmtypes.StoreKey,
		vbrmoduletypes.StoreKey,
		didTypes.StoreKey,
		commerciokycTypes.StoreKey,
		commerciomintTypes.StoreKey,
		governmentmoduletypes.StoreKey,
		documentstypes.StoreKey,
		epochstypes.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// Setup all modules keepers

	// add params keeper
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[consensusparamtypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String())
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	app.ScopedTransferKeeper = scopedTransferKeeper

	// -----------------------------------------
	// add keepers

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms, Bech32Prefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, BlockedAddresses(), authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[slashingtypes.StoreKey],
		app.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		keys[crisistypes.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	cfg := module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	// register the staking hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// -----------------------------------------
	// ... other modules keepers

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)
	app.WireICS20PreWasmKeeper(appCodec, bApp)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	epochsKeeper := epochskeeper.NewKeeper(appCodec, keys[epochstypes.StoreKey])
	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// Government keeper must be set before other modules keeper that depend on it
	app.GovernmentKeeper = *governmentmodulekeeper.NewKeeper(
		appCodec,
		keys[governmentmoduletypes.StoreKey],
		keys[governmentmoduletypes.MemStoreKey],
	)
	governmentModule := governmentmodule.NewAppModule(appCodec, app.GovernmentKeeper)

	// Create Vbr keeper
	app.VbrKeeper = *vbrmodulekeeper.NewKeeper(
		appCodec,
		keys[vbrmoduletypes.StoreKey],
		keys[vbrmoduletypes.MemStoreKey],
		app.DistrKeeper,
		app.BankKeeper,
		app.AccountKeeper,
		app.GovernmentKeeper,
		app.EpochsKeeper,
		app.GetSubspace(vbrmoduletypes.ModuleName),
		*app.StakingKeeper,
	)
	vbrModule := vbrmodule.NewAppModule(appCodec, app.VbrKeeper)

	// CommercioMint keeper must be set before CommercioKyc
	app.CommercioMintKeeper = *commerciomintKeeper.NewKeeper(
		appCodec,
		keys[commerciomintTypes.StoreKey],
		keys[commerciomintTypes.MemStoreKey],
		app.BankKeeper,
		app.AccountKeeper,
		app.GovernmentKeeper,
		app.GetSubspace(commerciomintTypes.ModuleName),
	)
	commercioMintModule := commerciomintmodule.NewAppModule(appCodec, app.CommercioMintKeeper)

	// Create commerciokyc keeper
	app.CommercioKycKeeper = *commerciokycKeeper.NewKeeper(
		appCodec,
		keys[commerciokycTypes.StoreKey],
		keys[commerciokycTypes.MemStoreKey],
		app.BankKeeper,
		app.GovernmentKeeper,
		app.AccountKeeper,
		app.CommercioMintKeeper,
	)
	commerciokycModule := commerciokycmodule.NewAppModule(appCodec, app.CommercioKycKeeper)

	// Create did keeper
	app.DidKeeper = *didkeeper.NewKeeper(
		appCodec,
		keys[didTypes.StoreKey],
		keys[didTypes.MemStoreKey],
	)
	didModule := did.NewAppModule(appCodec, app.DidKeeper)

	// Create documents keeper
	app.DocumentsKeeper = *documentskeeper.NewKeeper(
		appCodec,
		keys[documentstypes.StoreKey],
		keys[documentstypes.MemStoreKey],
	)
	documentsModule := documents.NewAppModule(appCodec, app.DocumentsKeeper)

	// Create epoch keeper
	app.EpochsKeeper = *epochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(
			// insert epoch hooks receivers here
			app.VbrKeeper.Hooks(),
		),
	)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, app.TransferStack)

	// Wasm keeper support
	wasmDir := filepath.Join(homePath, "data")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	supportedFeatures := "iterator,staking,stargate"

	//wasmOpts = append(wasmOpts, wasmkeeper.WithCustomIBCPortNameGenerator(wasmkeeper.HexIBCPortNameGenerator{}))

	app.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		keys[wasmtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.AddressLimitingICS4Wrapper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)
	// Pass the contract keeper to all the structs (generally ICS4Wrappers for ibc middlewares) that need it
	app.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	app.AddressLimitingICS4Wrapper.ContractKeeper = app.ContractKeeper

	// wire up x/wasm to IBC
	ibcRouter.AddRoute(wasmtypes.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper))
	app.IBCKeeper.SetRouter(ibcRouter)

	// The gov proposal types can be individually enabled
	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasmtypes.RouterKey, wasmkeeper.NewLegacyWasmProposalHandler(app.WasmKeeper, enabledProposals))
	}
	govConfig := govtypes.DefaultConfig()

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		params.NewAppModule(app.ParamsKeeper),
		app.RawIcs20TransferAppModule,
		// custoum modules
		governmentModule,
		vbrModule,
		commerciokycModule,
		commercioMintModule,
		didModule,
		documentsModule,
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		ibcaddresslimit.NewAppModule(*app.AddressLimitingICS4Wrapper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)), // always be last to make sure that it checks for all invariants and not only part of them
	)

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}

	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overrideModules)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		// Note: epochs' begin should be "real" start of epochs, we keep epochs beginblock at the beginning
		epochstypes.ModuleName,
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		govtypes.ModuleName,
		ibctransfertypes.ModuleName,
		crisistypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		wasmtypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		commerciokycTypes.ModuleName,
		commerciomintTypes.ModuleName,
		documentstypes.ModuleName,
		didTypes.ModuleName,
		governmentmoduletypes.ModuleName,
		vbrmoduletypes.ModuleName,
		ibcaddresslimittypes.ModuleName,
	)

	// TODO: check order for End Blockers
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		govtypes.ModuleName,
		ibctransfertypes.ModuleName,
		crisistypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		genutiltypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		commerciokycTypes.ModuleName,
		commerciomintTypes.ModuleName,
		documentstypes.ModuleName,
		didTypes.ModuleName,
		governmentmoduletypes.ModuleName,
		vbrmoduletypes.ModuleName,
		// Note: epochs' endblock should be "real" end of epochs, we keep epochs endblock at the end
		epochstypes.ModuleName,
		wasmtypes.ModuleName,
		ibcaddresslimittypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	// TODO: check init genesis correct order.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName, // Required if your application uses the localhost client (opens new window) to connect two different modules from the same chain
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		feegrant.ModuleName,
		governmentmoduletypes.ModuleName,
		commerciokycTypes.ModuleName,
		commerciomintTypes.ModuleName,
		vbrmoduletypes.ModuleName,
		didTypes.ModuleName,
		documentstypes.ModuleName,
		epochstypes.ModuleName,
		authz.ModuleName,
		wasmtypes.ModuleName,
		ibcaddresslimittypes.ModuleName,
	)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	cfg = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(cfg)

	// TODO: add NewSimulationManager for fuzzy testing

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(
		ante.NewAnteHandler(
			app.AccountKeeper,
			app.BankKeeper,
			app.GovernmentKeeper,
			app.CommercioMintKeeper,
			comosante.DefaultSigVerificationGasConsumer,
			encodingConfig.TxConfig.SignModeHandler(),
			DefaultBondDenom,
			StableCreditsDenom,
			app.FeeGrantKeeper,
			app.IBCKeeper,
			&wasmConfig,
			keys[wasmtypes.StoreKey],
		),
	)
	app.SetEndBlocker(app.EndBlocker)

	// Old upgrade v3.1.0 only for dev enviroment
	/*upgradeName := "v3.1.0"

	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, plan upgradetypes.Plan) {

		},
	)*/

	upgradeName := "v4.0.0"
	/*app.UpgradeKeeper.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
			// Update modules params
			// Update gov params
			app.GovKeeper.SetVotingParams(ctx, govtypes.NewVotingParams(time.Hour*24))
			app.GovKeeper.SetDepositParams(ctx, govtypes.NewDepositParams(sdk.NewCoins(sdk.NewCoin(DefaultBondDenom, sdk.NewInt(5000000000))), time.Hour*48))

			// Update wasm params
			wasmParams := wasmtypes.DefaultParams()
			wasmParams.CodeUploadAccess.Permission = wasmtypes.AccessTypeOnlyAddress
			wasmParams.CodeUploadAccess.Address = app.GovernmentKeeper.GetGovernment300Address(ctx).String()
			app.WasmKeeper.SetParams(ctx, wasmParams)

			// Update slashing params
			slashingParams := slashingtypes.NewParams(
				20000,
				sdk.NewDecFromIntWithPrec(sdk.NewInt(5), 2),
				time.Minute*10,
				sdk.NewDecFromIntWithPrec(sdk.NewInt(1), 2),
				sdk.NewDecFromIntWithPrec(sdk.NewInt(5), 2),
			)
			app.SlashingKeeper.SetParams(ctx, slashingParams)

			fromVM := make(map[string]uint64)
			for moduleName := range app.mm.Modules {
				fromVM[moduleName] = 1
			}

			delete(fromVM, "ibc") // Force delete ibc reference

			return app.mm.RunMigrations(ctx, cfg, fromVM)
		},
	)*/

	upgradeNameV420 := "v4.2.0"
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeNameV420,
		func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
			return vm, nil
		},
	)

	upgradeNameV5 := "v5.0.0"
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeNameV5,
		func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
			return vm, nil
		},
	)

	upgradeNameV51 := "v5.1.0"
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradeNameV51,
		func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
			return vm, nil
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	// Setup store loader
	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	if upgradeInfo.Name == upgradeNameV420 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgradesV4 := storetypes.StoreUpgrades{}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgradesV4))
	}

	if upgradeInfo.Name == upgradeNameV5 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgradesV5 := storetypes.StoreUpgrades{
			Added: []string{ibcaddresslimittypes.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgradesV5))
	}

	if upgradeInfo.Name == upgradeNameV51 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgradesV51 := storetypes.StoreUpgrades{
			Added: []string{ibcaddresslimittypes.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgradesV51))
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			cometos.Exit(err.Error())
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		app.NewUncachedContext(true, tmproto.Header{})
		app.CapabilityKeeper.Seal()
	}

	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))

		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper

	return app
}

// Create the IBC Transfer Stack from bottom to top:
//
// * SendPacket. Originates from the transferKeeper and and goes up the stack:
// transferKeeper.SendPacket -> ibc_address_limit.SendPacket -> channel.SendPacket
// * RecvPacket, message that originates from core IBC and goes down to app, the flow is the other way
// channel.RecvPacket -> ibc_address_limit.OnRecvPacket -> transfer.OnRecvPacket
//
// After this, the wasm keeper is required to be set on appKeepers.AddressLimitingICS4Wrapper
func (appKeepers *App) WireICS20PreWasmKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp) {

	// ChannelKeeper wrapper for address limiting SendPacket(). The wasmKeeper needs to be added after it's created
	addressLimitingParams := appKeepers.GetSubspace(ibcaddresslimittypes.ModuleName)
	addressLimitingParams = addressLimitingParams.WithKeyTable(ibcaddresslimittypes.ParamKeyTable())
	AddressLimitingICS4Wrapper := ibcaddresslimit.NewICS4Middleware(
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.AccountKeeper,
		// wasm keeper we set later.
		nil,
		&appKeepers.BankKeeper,
		addressLimitingParams,
	)
	appKeepers.AddressLimitingICS4Wrapper = &AddressLimitingICS4Wrapper

	// Create Transfer Keepers
	transferKeeper := ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		// The ICS4Wrapper is replaced by the addressLimitingICS4Wrapper instead of the channel
		appKeepers.AddressLimitingICS4Wrapper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)
	appKeepers.TransferKeeper = transferKeeper
	appKeepers.RawIcs20TransferAppModule = transfer.NewAppModule(appKeepers.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(appKeepers.TransferKeeper)

	// AddressLimiting IBC Middleware
	addressLimitingTransferModule := ibcaddresslimit.NewIBCModule(transferIBCModule, appKeepers.AddressLimitingICS4Wrapper)
	appKeepers.TransferStack = &addressLimitingTransferModule
}

func (app *App) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := cometjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Commercio's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Commercio's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
	// register app's OpenAPI routes.
	apiSvr.Router.Handle("/static/openapi.yml", http.FileServer(http.FS(docs.Docs)))
	//apiSvr.Router.HandleFunc("/", openapiconsole.Handler(Name, "/static/openapi.yml"))
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

func (app *App) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range GetMaccPerms() {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	// allow the following addresses to receive funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// Required for ibctesting
func (app *App) GetStakingKeeper() stakingkeeper.Keeper {
	return *app.StakingKeeper
}

func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper // This is a *ibckeeper.Keeper
}

func (app *App) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *App) GetTxConfig() client.TxConfig {
	return MakeEncodingConfig().TxConfig
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName)

	paramsKeeper.Subspace(governmentmoduletypes.ModuleName)
	paramsKeeper.Subspace(vbrmoduletypes.ModuleName)
	paramsKeeper.Subspace(didTypes.ModuleName)
	paramsKeeper.Subspace(documentstypes.ModuleName)
	paramsKeeper.Subspace(commerciomintTypes.ModuleName)
	paramsKeeper.Subspace(commerciokycTypes.ModuleName)
	paramsKeeper.Subspace(ibcaddresslimittypes.ModuleName)

	return paramsKeeper
}

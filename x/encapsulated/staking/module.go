package customstaking

import (
	"encoding/json"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simulation2 "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking/simulation"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func DefaultGenesis(defaultBondDenom string) staking.GenesisState {
	return staking.GenesisState{
		Params: staking.Params{
			UnbondingTime: types.DefaultUnbondingTime,
			MaxValidators: types.DefaultMaxValidators,
			MaxEntries:    types.DefaultMaxEntries,
			BondDenom:     defaultBondDenom,
		},
	}
}

type AppModuleBasic struct {
	stakingModule    staking.AppModuleBasic
	DefaultBondDenom string
}

func NewAppModuleBasic(defaultBondDenom string) AppModuleBasic {
	return AppModuleBasic{
		stakingModule:    staking.AppModuleBasic{},
		DefaultBondDenom: defaultBondDenom,
	}
}

func (am AppModuleBasic) Name() string {
	return am.stakingModule.Name()
}

func (am AppModuleBasic) RegisterCodec(arg *codec.Codec) {
	am.stakingModule.RegisterCodec(arg)
}

func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return staking.ModuleCdc.MustMarshalJSON(DefaultGenesis(am.DefaultBondDenom))
}

func (am AppModuleBasic) ValidateGenesis(arg json.RawMessage) error {
	return am.stakingModule.ValidateGenesis(arg)
}

func (am AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	am.stakingModule.RegisterRESTRoutes(ctx, r)
}

func (am AppModuleBasic) GetTxCmd(arg *codec.Codec) *cobra.Command {
	return am.stakingModule.GetTxCmd(arg)
}

func (am AppModuleBasic) GetQueryCmd(arg *codec.Codec) *cobra.Command {
	return am.stakingModule.GetQueryCmd(arg)
}

// AppModule implements an application module for the staking module.
type AppModule struct {
	AppModuleBasic

	keeper        staking.Keeper
	accountKeeper types.AccountKeeper
	supplyKeeper  types.SupplyKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper staking.Keeper, accountKeeper types.AccountKeeper, supplyKeeper types.SupplyKeeper) AppModule {

	return AppModule{
		AppModuleBasic: NewAppModuleBasic("ucommercio"),
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		supplyKeeper:   supplyKeeper,
	}
}

// Name returns the staking module's name.
func (AppModule) Name() string {
	return staking.ModuleName
}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	staking.RegisterInvariants(ir, am.keeper)
}

// Route returns the message routing key for the staking module.
func (AppModule) Route() string {
	return staking.RouterKey
}

// NewHandler returns an sdk.Handler for the staking module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the staking module's querier route name.
func (AppModule) QuerierRoute() string {
	return staking.QuerierRoute
}

// NewQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return staking.NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the staking module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState staking.GenesisState
	staking.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return staking.InitGenesis(ctx, am.keeper, am.accountKeeper, am.supplyKeeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the staking
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := staking.ExportGenesis(ctx, am.keeper)
	return staking.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the staking module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	staking.BeginBlocker(ctx, am.keeper)
}

// EndBlock returns the end blocker for the staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return staking.EndBlocker(ctx, am.keeper)
}

//____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the staking module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simulation2.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized staking param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simulation2.ParamChange {
	return simulation.ParamChanges(r)
}

// RegisterStoreDecoder registers a decoder for staking module's types
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[staking.StoreKey] = simulation.DecodeStore
}

// WeightedOperations returns the all the staking module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simulation2.WeightedOperation {
	return simulation.WeightedOperations(simState.AppParams, simState.Cdc,
		am.accountKeeper, am.keeper)
}

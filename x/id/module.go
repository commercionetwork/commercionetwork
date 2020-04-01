package id

import (
	"encoding/json"

	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"

	"github.com/commercionetwork/commercionetwork/x/id/types"

	"github.com/commercionetwork/commercionetwork/x/id/client/cli"

	"github.com/commercionetwork/commercionetwork/x/id/client/rest"
	"github.com/commercionetwork/commercionetwork/x/id/keeper"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the id module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// module genesis validation
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, types.QuerierRoute)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command {
	return nil
}

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the auth module.
type AppModuleSimulation struct{}

// RegisterStoreDecoder registers a decoder for auth module's types
func (AppModuleSimulation) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation
	keeper       keeper.Keeper
	govKeeper    governmentKeeper.Keeper
	supplyKeeper supply.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper, govKeeper governmentKeeper.Keeper, supplyKeeper supply.Keeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		keeper:              keeper,
		govKeeper:           govKeeper,
		supplyKeeper:        supplyKeeper,
	}
}

// module name
func (AppModule) Name() string {
	return types.ModuleName
}

// register invariants
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// module message route name
func (AppModule) Route() string {
	return types.ModuleName
}

// module handler
func (am AppModule) NewHandler() sdk.Handler {
	return keeper.NewHandler(am.keeper, am.govKeeper)
}

// module querier route name
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

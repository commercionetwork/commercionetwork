package custombank

import (
	"encoding/json"

	"github.com/commercionetwork/commercionetwork/x/encapsulated/bank/client/rest"
	"github.com/commercionetwork/commercionetwork/x/encapsulated/bank/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	cosmosbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the bank module.
type AppModuleBasic struct {
	cosmosbank.AppModuleBasic
}

func NewAppModuleBasic(basic cosmosbank.AppModuleBasic) AppModuleBasic {
	return AppModuleBasic{
		AppModuleBasic: basic,
	}
}

// Name returns the bank module's name.
func (amb AppModuleBasic) Name() string { return amb.AppModuleBasic.Name() }

// RegisterCodec registers the bank module's types for the given codec.
func (amb AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	amb.AppModuleBasic.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the bank
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	defaultState := DefaultGenesisState()
	bytes, _ := json.Marshal(&defaultState)
	return bytes
}

// ValidateGenesis performs genesis state validation for the bank module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := json.Unmarshal(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the bank module.
func (amb AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	amb.AppModuleBasic.RegisterRESTRoutes(ctx, rtr)
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the bank module.
func (amb AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return amb.AppModuleBasic.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the bank module.
func (amb AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return amb.AppModuleBasic.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the bank module.
type AppModule struct {
	cosmosbank.AppModule
	keeper    keeper.Keeper
	govKeeper government.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(appModule cosmosbank.AppModule, keeper keeper.Keeper, govKeeper government.Keeper) AppModule {
	return AppModule{
		AppModule: appModule,
		keeper:    keeper,
		govKeeper: govKeeper,
	}
}

// Name returns the bank module's name.
func (am AppModule) Name() string { return am.AppModule.Name() }

// RegisterInvariants registers the bank module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.AppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the bank module.
func (am AppModule) Route() string { return am.AppModule.Route() }

// NewHandler returns an sdk.Handler for the bank module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.AppModule.NewHandler(), am.keeper, am.govKeeper)
}

// QuerierRoute returns the bank module's querier route name.
func (am AppModule) QuerierRoute() string { return am.AppModule.QuerierRoute() }

// NewQuerierHandler returns the bank module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.AppModule.NewQuerierHandler(), am.keeper)
}

// InitGenesis performs genesis initialization for the bank module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	if err := json.Unmarshal(data, &genesisState); err != nil {
		panic(err)
	}

	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the bank
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)

	bytes, err := json.Marshal(&gs)
	if err != nil {
		panic(err)
	}

	return bytes
}

// BeginBlock performs a no-op.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the bank module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

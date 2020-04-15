package commerciomint

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/client/cli"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/client/rest"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const DefaultCreditsDenom = "uccc"

// AppModuleBasic defines the basic application module used by the docs module.
type AppModuleBasic struct {
	CreditsDenom string
}

func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

// module name
func (AppModuleBasic) Name() string                   { return types.ModuleName }
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { types.RegisterCodec(cdc) }

// default genesis state
func (amb AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState(DefaultCreditsDenom))
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
	rest.RegisterRoutes(ctx, rtr)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the auth module.
type AppModuleSimulation struct{}

// RegisterStoreDecoder registers a decoder for auth module's types
func (AppModuleSimulation) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {}

//____________________________________________________________________________

// AppModule implements an application module for the id module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation
	keeper       keeper.Keeper
	supplyKeeper supply.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper, sk supply.Keeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		keeper:              keeper,
		supplyKeeper:        sk,
	}
}

// module name
func (AppModule) Name() string { return types.ModuleName }

// register invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	keeper.RegisterInvariants(ir, am.keeper)
}

// module message route name
func (AppModule) Route() string { return types.RouterKey }

// module handler
func (am AppModule) NewHandler() sdk.Handler        { return keeper.NewHandler(am.keeper) }
func (AppModule) QuerierRoute() string              { return types.QuerierRoute }
func (am AppModule) NewQuerierHandler() sdk.Querier { return keeper.NewQuerier(am.keeper) }

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.supplyKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(ExportGenesis(ctx, am.keeper))
}

// module begin-block
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {}

// module end-block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.keeper.AutoLiquidatePositions(ctx)
	return []abci.ValidatorUpdate{}
}

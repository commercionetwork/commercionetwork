package vbr

import (
	"encoding/json"
	"github.com/commercionetwork/commercionetwork/x/vbr/keeper"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"

	"github.com/commercionetwork/commercionetwork/x/vbr/client/rest"

	"github.com/commercionetwork/commercionetwork/x/vbr/client/cli"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

//module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

//module codecs
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

//genesis default state
func (amb AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

//genesis validation
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc, types.ModuleName, types.QuerierRoute)
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the auth module.
type AppModuleSimulation struct{}

// RegisterStoreDecoder registers a decoder for auth module's types
func (AppModuleSimulation) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

//____________________________________________________________________________

//___________________________
// app module
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation
	keeper        keeper.Keeper
	stakingKeeper staking.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, sk staking.Keeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		keeper:              k,
		stakingKeeper:       sk,
	}
}

func (AppModule) Name() string {
	return types.ModuleName
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (AppModule) Route() string {
	return types.ModuleName
}

func (am AppModule) NewHandler() sdk.Handler {
	return keeper.NewHandler(am.keeper)
}

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
func (am AppModule) BeginBlock(ctx sdk.Context, rbb abci.RequestBeginBlock) {
	BeginBlocker(ctx, rbb, am.keeper, am.stakingKeeper)
}

// module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

package creditrisk

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/creditrisk/client/cli"
	"github.com/commercionetwork/commercionetwork/x/creditrisk/client/rest"
	"github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string                                { return types.ModuleName }
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec)              { RegisterCodec(cdc) }
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command    { return nil }
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command { return cli.GetQueryCmd(cdc) }
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.GenesisState{Pool: sdk.NewCoins()})
}

func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRESTRoutes(ctx, rtr, types.QuerierRoute)
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {}
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry)            {}
func (AppModule) Route() string                                             { return types.RouterKey }
func (AppModule) QuerierRoute() string                                      { return types.QuerierRoute }
func (am AppModule) NewQuerierHandler() sdk.Querier                         { return NewQuerier(am.keeper) }

func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) NewHandler() sdk.Handler {
	return func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
		return nil, errors.Wrap(errors.ErrUnknownRequest, "creditrisk doesn't handle messages")
	}
}

// InitGenesis initialise the state.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis serializes the state to JSON.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.GenesisState{Pool: am.keeper.GetPoolFunds(ctx)})
}

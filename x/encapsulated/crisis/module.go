package customcrisis

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func DefaultGenesis(defaultBondDenom string) crisis.GenesisState {
	return crisis.GenesisState{
		ConstantFee: sdk.NewCoin(defaultBondDenom, sdk.NewInt(1000)),
	}
}

// app module basics object
type AppModuleBasic struct {
	CrisisModule     crisis.AppModuleBasic
	DefaultBondDenom string
}

func NewAppModuleBasic(defaultBondDenom string) AppModuleBasic {
	return AppModuleBasic{
		CrisisModule:     crisis.AppModuleBasic{},
		DefaultBondDenom: defaultBondDenom,
	}
}

func (am AppModuleBasic) Name() string {
	return am.CrisisModule.Name()
}

func (am AppModuleBasic) RegisterCodec(arg *codec.Codec) {
	am.CrisisModule.RegisterCodec(arg)
}

func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return crisis.ModuleCdc.MustMarshalJSON(DefaultGenesis(am.DefaultBondDenom))
}

func (am AppModuleBasic) ValidateGenesis(arg json.RawMessage) error {
	return am.CrisisModule.ValidateGenesis(arg)
}

func (am AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	am.CrisisModule.RegisterRESTRoutes(ctx, r)
}

func (am AppModuleBasic) GetTxCmd(arg *codec.Codec) *cobra.Command {
	return am.CrisisModule.GetTxCmd(arg)
}

func (am AppModuleBasic) GetQueryCmd(arg *codec.Codec) *cobra.Command {
	return am.CrisisModule.GetQueryCmd(arg)
}

package custommint

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func DefaultGenesis(defaultBondDenom string) mint.GenesisState {
	return mint.GenesisState{
		Minter: mint.DefaultInitialMinter(),
		Params: mint.Params{
			MintDenom:           defaultBondDenom,
			InflationRateChange: sdk.NewDecWithPrec(13, 2),
			InflationMax:        sdk.NewDecWithPrec(20, 2),
			InflationMin:        sdk.NewDecWithPrec(7, 2),
			GoalBonded:          sdk.NewDecWithPrec(67, 2),
			BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
		},
	}
}

type AppModuleBasic struct {
	mintModule       mint.AppModuleBasic
	DefaultBondDenom string
}

func (am AppModuleBasic) Name() string {
	return am.mintModule.Name()
}

func (am AppModuleBasic) RegisterCodec(arg *codec.Codec) {
	am.mintModule.RegisterCodec(arg)
}

func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return mint.ModuleCdc.MustMarshalJSON(DefaultGenesis(am.DefaultBondDenom))
}

func (am AppModuleBasic) ValidateGenesis(arg json.RawMessage) error {
	return am.mintModule.ValidateGenesis(arg)
}

func (am AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	am.mintModule.RegisterRESTRoutes(ctx, r)
}

func (am AppModuleBasic) GetTxCmd(arg *codec.Codec) *cobra.Command {
	return am.mintModule.GetTxCmd(arg)
}

func (am AppModuleBasic) GetQueryCmd(arg *codec.Codec) *cobra.Command {
	return am.mintModule.GetQueryCmd(arg)
}

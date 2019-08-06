package customstaking

import (
	"encoding/json"
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

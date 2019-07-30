package customgov

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func DefaultGenesis(defaultBondName string) gov.GenesisState {
	minDepositTokens := sdk.TokensFromConsensusPower(50000)
	return gov.GenesisState{
		StartingProposalID: 1,
		DepositParams: gov.DepositParams{
			MinDeposit:       sdk.Coins{sdk.NewCoin(defaultBondName, minDepositTokens)},
			MaxDepositPeriod: gov.DefaultPeriod,
		},
		VotingParams: gov.VotingParams{
			VotingPeriod: gov.DefaultPeriod,
		},
		TallyParams: gov.TallyParams{
			Quorum:    sdk.NewDecWithPrec(334, 3),
			Threshold: sdk.NewDecWithPrec(5, 1),
			Veto:      sdk.NewDecWithPrec(334, 3),
		},
	}
}

type AppModuleBasic struct {
	govModule       gov.AppModuleBasic
	DefaultBondName string
}

func (am AppModuleBasic) Name() string {
	return am.govModule.Name()
}

func (am AppModuleBasic) RegisterCodec(arg *codec.Codec) {
	am.govModule.RegisterCodec(arg)
}

func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return gov.ModuleCdc.MustMarshalJSON(DefaultGenesis(am.DefaultBondName))
}

func (am AppModuleBasic) ValidateGenesis(arg json.RawMessage) error {
	return am.govModule.ValidateGenesis(arg)
}

func (am AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	am.govModule.RegisterRESTRoutes(ctx, r)
}

func (am AppModuleBasic) GetTxCmd(arg *codec.Codec) *cobra.Command {
	return am.govModule.GetTxCmd(arg)
}

func (am AppModuleBasic) GetQueryCmd(arg *codec.Codec) *cobra.Command {
	return am.govModule.GetQueryCmd(arg)
}

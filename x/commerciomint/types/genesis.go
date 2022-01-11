package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/durationpb"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	freezePeriod := durationpb.New(DefaultFreezePeriod).AsDuration()
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Positions:  []*Position{},
		PoolAmount: sdk.Coins{},
		Params: Params{
			ConversionRate: DefaultConversionRate,
			FreezePeriod:   freezePeriod, //TODO CONTROL CAST
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {

	// this line is used by starport scaffolding # genesis/types/validate

	if err := gs.Params.Validate(); err != nil {
		return err
	}

	//  TODO validate PoolAmount

	for _, position := range gs.Positions {
		if err := position.Validate(); err != nil {
			return err
		}
	}
	return nil
}

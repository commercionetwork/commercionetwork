package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/durationpb"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	freezePeriod := durationpb.New(DefaultFreezePeriod).AsDuration()
	return &GenesisState{
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

	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}

	if ok := gs.PoolAmount.IsValid(); !ok {
		return errors.New("invalid pool amount")
	}

	for _, position := range gs.Positions {
		if err := position.Validate(); err != nil {
			return fmt.Errorf("invalid position: %s", err)
		}
	}
	return nil
}

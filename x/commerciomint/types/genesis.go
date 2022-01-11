package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	DefaultCreditsDenom               = "uccc"
	DefaultFreezePeriod time.Duration = time.Hour * 24 * 7 * 3
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	collateralRate := sdk.NewDec(1)
	freezePeriod := durationpb.New(DefaultFreezePeriod).AsDuration()
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Positions:  []*Position{},
		PoolAmount: sdk.Coins{},
		Params: Params{
			CollateralRate: collateralRate,
			FreezePeriod:   &freezePeriod, //TODO CONTROL CAST
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	// this line is used by starport scaffolding # genesis/types/validate
	if gs.Params.CollateralRate.IsZero() {
		return fmt.Errorf("conversion rate cannot be zero")
	}
	if gs.Params.CollateralRate.IsNegative() {
		return fmt.Errorf("conversion rate must be positive")
	}

	if gs.Params.FreezePeriod.Seconds() < 0 {
		return fmt.Errorf("freeze period cannot be lower than zero")
	}
	for _, position := range gs.Positions {
		err := position.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

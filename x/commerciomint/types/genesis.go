package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCreditsDenom               = "uccc"
	DefaultFreezePeriod time.Duration = time.Hour * 24 * 7 * 3
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	collateralRate := sdk.DecProto{
		Dec: sdk.NewDec(1),
	}

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Positions:      []*Position{},
		PoolAmount:     []*sdk.Coin{},
		CollateralRate: &collateralRate,
		FreezePeriod:   DefaultFreezePeriod.String(), //TODO CONTROL CAST
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	// this line is used by starport scaffolding # genesis/types/validate
	if gs.CollateralRate.Dec.IsZero() {
		return fmt.Errorf("conversion rate cannot be zero")
	}
	if gs.CollateralRate.Dec.IsNegative() {
		return fmt.Errorf("conversion rate must be positive")
	}

	freezePeriod, err := time.ParseDuration(gs.FreezePeriod)
	if freezePeriod.Seconds() < 0 || err != nil {
		return fmt.Errorf("freeze period cannot be lower than zero or have a incorrect form %s", err.Error())
	}
	for _, position := range gs.Positions {
		err := position.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

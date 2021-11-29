package types

import (
	"errors"
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AutomaticWithdraw: true,
		Params: Params{
			DistrEpochIdentifier: "day",
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if gs.PoolAmount == nil || gs.PoolAmount.Empty() {
		return errors.New("validator block reward pool cannot be empty")
	}

	if !gs.PoolAmount.IsValid() {
		return fmt.Errorf("invalid validator block reward pool: %s", gs.PoolAmount.String())
	}

	if gs.Params.DistrEpochIdentifier == "" {
		return errors.New("epoch identifier should NOT be empty")
	}

	return ValidateRewardRate(gs.RewardRate)
}

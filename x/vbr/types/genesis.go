package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefautGenesis is not valid

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: NewParams(EpochDay, sdk.NewDecWithPrec(5, 1)),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if gs.PoolAmount == nil || gs.PoolAmount.Empty() {
		return errors.New("validator block reward pool cannot be empty")
	}

	// this check seems redundant by construction of sdk.DecCoins
	if !gs.PoolAmount.IsValid() {
		return fmt.Errorf("invalid validator block reward pool: %s", gs.PoolAmount.String())
	}

	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}

	return nil
}

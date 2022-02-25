package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// unused constant

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefautGenesis is not valid

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{
			DistrEpochIdentifier: EpochDay,
			EarnRate:             sdk.NewDecWithPrec(5, 1),
		},
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

	if gs.Params.DistrEpochIdentifier == "" {
		return errors.New("epoch identifier should NOT be empty")
	}

	if gs.Params.EarnRate.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("EarnRate: %d must be positive", gs.Params.EarnRate))
	}

	return nil
}

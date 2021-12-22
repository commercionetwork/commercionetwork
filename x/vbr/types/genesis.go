package types

import (
	"errors"
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{
			DistrEpochIdentifier: /*EpochMinute*/EpochDay,
			EarnRate: sdk.NewDecWithPrec(050,2),
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

	if gs.Params.EarnRate.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("EarnRate: %d must be positive", gs.Params.EarnRate))
	}

	return nil
}

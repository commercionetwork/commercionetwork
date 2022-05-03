package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if _, err := sdk.AccAddressFromBech32(gs.GovernmentAddress); err != nil {
		return err
	}
	return nil
}

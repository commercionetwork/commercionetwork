package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default

	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in Document
	if gs.GovernmentAddress == "" {
		return fmt.Errorf("government address cannot be empty. Use the set-genesis-government-address command to set one")
	}

	return nil
}

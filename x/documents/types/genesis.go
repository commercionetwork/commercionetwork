package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Documents: []*Document{},
		Receipts:  []*DocumentReceipt{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in Document
	DocumentIdMap := make(map[string]bool)

	for _, elem := range gs.Documents {
		if _, ok := DocumentIdMap[elem.UUID]; ok {
			return fmt.Errorf("duplicated id for Document")
		}
		DocumentIdMap[elem.UUID] = true
	}

	// Check for duplicated ID in Document
	ReceiptIdMap := make(map[string]bool)

	for _, elem := range gs.Receipts {
		if _, ok := ReceiptIdMap[elem.UUID]; ok {
			return fmt.Errorf("duplicated id for Receipt")
		}
		ReceiptIdMap[elem.UUID] = true
	}

	return nil
}

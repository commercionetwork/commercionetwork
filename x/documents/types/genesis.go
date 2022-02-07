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

	// Check for duplicated ID in documents
	DocumentIdMap := make(map[string]struct{})

	for _, document := range gs.Documents {
		if _, ok := DocumentIdMap[document.UUID]; ok {
			return fmt.Errorf("duplicated id %s for document", document.UUID)
		}
		if err := document.Validate(); err != nil {
			return fmt.Errorf("document with UUID %s is invalid: %e", document.UUID, err)
		}
		DocumentIdMap[document.UUID] = struct{}{}
	}

	// Check for duplicated ID in receipts
	ReceiptIdMap := make(map[string]struct{})

	for _, receipt := range gs.Receipts {
		if _, found := DocumentIdMap[receipt.DocumentUUID]; !found {
			return fmt.Errorf("could not find corresponding document for %s", receipt.UUID)
		}
		if _, ok := ReceiptIdMap[receipt.UUID]; ok {
			return fmt.Errorf("duplicated id %s for receipt", receipt.UUID)
		}
		if err := receipt.Validate(); err != nil {
			return fmt.Errorf("receipt with UUID %s is invalid: %e", receipt.UUID, err)
		}
		ReceiptIdMap[receipt.UUID] = struct{}{}
	}

	return nil
}

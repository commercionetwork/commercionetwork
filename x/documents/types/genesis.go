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

	// Support map for duplicated ID in documents
	// It will be used later also for accepting only receipts that concern valid documents
	duplicateDocumentIdMap := make(map[string]struct{})

	for _, document := range gs.Documents {
		// Check for duplicate document UUID
		if _, duplicated := duplicateDocumentIdMap[document.UUID]; duplicated {
			return fmt.Errorf("duplicated id %s for document", document.UUID)
		}
		if err := document.Validate(); err != nil {
			return fmt.Errorf("document with UUID %s is invalid: %e", document.UUID, err)
		}
		duplicateDocumentIdMap[document.UUID] = struct{}{}
	}

	// Support map for duplicated ID in receipts
	duplicateReceiptIdMap := make(map[string]struct{})

	for _, receipt := range gs.Receipts {
		// Check if the receipt concers an invalid document
		if _, found := duplicateDocumentIdMap[receipt.DocumentUUID]; !found {
			return fmt.Errorf("could not find corresponding document for %s", receipt.UUID)
		}
		// Check for duplicate receipt UUID
		if _, duplicated := duplicateReceiptIdMap[receipt.UUID]; duplicated {
			return fmt.Errorf("duplicated id %s for receipt", receipt.UUID)
		}
		if err := receipt.Validate(); err != nil {
			return fmt.Errorf("receipt with UUID %s is invalid: %e", receipt.UUID, err)
		}
		duplicateReceiptIdMap[receipt.UUID] = struct{}{}
	}

	return nil
}

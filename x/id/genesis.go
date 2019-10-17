package id

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - id genesis state
type GenesisState struct {
	DidDocuments []DidDocument `json:"did_documents"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets ids information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, didDocument := range data.DidDocuments {
		if err := keeper.SaveDidDocument(ctx, didDocument); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	identities, err := keeper.GetDidDocuments(ctx)
	if err != nil {
		panic(err)
	}
	return GenesisState{
		DidDocuments: identities,
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}

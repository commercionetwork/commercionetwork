package docs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	Documents                      Documents        `json:"documents"`
	Receipts                       DocumentReceipts `json:"receipts"`
	SupportedMetadataSchemes       MetadataSchemes  `json:"supported_metadata_schemes"`
	TrustedMetadataSchemaProposers []sdk.AccAddress `json:"trusted_metadata_schema_proposers"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, doc := range data.Documents {
		if err := keeper.SaveDocument(ctx, doc); err != nil {
			panic(err)
		}
	}

	for _, receipt := range data.Receipts {
		if err := keeper.SaveReceipt(ctx, receipt); err != nil {
			panic(err)
		}
	}

	for _, schema := range data.SupportedMetadataSchemes {
		keeper.AddSupportedMetadataScheme(ctx, schema)
	}

	for _, proposer := range data.TrustedMetadataSchemaProposers {
		keeper.AddTrustedSchemaProposer(ctx, proposer)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Documents:                      exportDocuments(ctx, keeper),
		Receipts:                       exportReceipts(ctx, keeper),
		SupportedMetadataSchemes:       exportMetadataSchemes(ctx, keeper),
		TrustedMetadataSchemaProposers: exportTrustedSchemaProviders(ctx, keeper),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
